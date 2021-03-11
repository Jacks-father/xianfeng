package chain

import (
	"github.com/boltdb/bolt"
	"errors"
	"math/big"
)

const BLOCKS = "blocks"
const LASTHASH = "lastHash"

/**
 * 定义区块链这个结构体，用于存储产生的区块（内存中)
 */
type BlockChain struct {
	//Blocks []Block
	Engine            *bolt.DB
	LastBlock         Block    //最新的区块
	IteratorBlockHash [32]byte //迭代到的区块哈希值
}

func NewBlockChain(db *bolt.DB) BlockChain {
	//增加为lastblock赋值的逻辑
	var lastBlock Block
	db.Update(func(tx *bolt.Tx) error {
		bucket := tx.Bucket([]byte(BLOCKS))
		if bucket == nil {
			bucket, _ = tx.CreateBucket([]byte(BLOCKS))
		}
		lastHash := bucket.Get([]byte(LASTHASH))
		if len(lastHash) == 0 {
			return nil
		}
		lastBlockBytes := bucket.Get(lastHash)
		lastBlock, _ = Deserialize(lastBlockBytes)
		return nil
	})
	return BlockChain{
		Engine:            db,
		LastBlock:         lastBlock,
		IteratorBlockHash: lastBlock.Hash,
	}
}

/**
 * 创建一个区块链实例，该实例携带一个创世区块
 */
func (chain *BlockChain) CreateGenesis(genesisData []byte) {
	//先看chain.LastBlock是否为空
	hashBig := new(big.Int)
	hashBig.SetBytes(chain.LastBlock.Hash[:])
	if hashBig.Cmp(big.NewInt(0)) > 0 {
		return
	}

	engine := chain.Engine
	//读一遍bucket，查看是否有数据
	engine.Update(func(tx *bolt.Tx) error { //
		bucket := tx.Bucket([]byte(BLOCKS))
		if bucket == nil {
			bucket, _ = tx.CreateBucket([]byte(BLOCKS))
		}
		if bucket != nil {
			lastHash := bucket.Get([]byte(LASTHASH))
			if len(lastHash) == 0 {
				genesis := CreateGenesisBlock(genesisData)
				genSerBytes, _ := genesis.Serialize()
				//存创世区块
				bucket.Put(genesis.Hash[:], genSerBytes)
				//更新最新区块的标志 lastHash -> 最新区块hash
				bucket.Put([]byte(LASTHASH), genesis.Hash[:])
				chain.LastBlock = genesis
				chain.IteratorBlockHash = genesis.Hash
			}
		}
		return nil
	})
}

/**
 * 新增一个区块
 */
func (chain *BlockChain) AddNewBlock(data []byte) error {
	//1、从db中找到最后一个区块数据
	engine := chain.Engine
	//2、获取到最新的区块
	lastBlock := chain.LastBlock

	//3、得到最后一个区块的各种属性，并利用这些属性生成新区块
	newBlock := CreateBlock(lastBlock.Height, lastBlock.Hash, data)
	newBlockByte, err := newBlock.Serialize()
	if err != nil {
		return err
	}
	//4、更新db文件，将新生成的区块写入到db中，同时更新最新区块的指向标记
	engine.Update(func(tx *bolt.Tx) error {
		bucket := tx.Bucket([]byte(BLOCKS))
		if bucket == nil {
			err = errors.New("区块链数据库操作失败，请重试!")
			return err
		}
		//将最新的区块数据存到db中
		bucket.Put(newBlock.Hash[:], newBlockByte)
		//更新最新区块的指向标记
		bucket.Put([]byte(LASTHASH), newBlock.Hash[:])

		//更新blockChain对象的LastBlock结构体实例
		chain.LastBlock = newBlock
		chain.IteratorBlockHash = newBlock.Hash
		return nil
	})
	return err
}

/**
 * 获取最新的最后的一个区块
 */
func (chain *BlockChain) GetLastBlock() (Block) {
	return chain.LastBlock
}

/**
 * 获取所有的区块
 */
func (chain BlockChain) GetAllBlocks() ([]Block, error) {
	engine := chain.Engine
	var errs error
	blocks := make([]Block, 0)
	engine.View(func(tx *bolt.Tx) error {
		bucket := tx.Bucket([]byte(BLOCKS))
		if bucket == nil {
			errs = errors.New("区块数据库操作是吧，请重试！")
			return errs
		}

		var currentHash []byte
		//直接从倒数第二个区块进行遍历
		currentHash = bucket.Get([]byte(LASTHASH))
		for { //倒数第一个区块区块开始遍历
			//根据区块hash拿[]byte类型的区块数据
			currentBlockBytes := bucket.Get(currentHash)
			//[]byte类型的区块数据 反序列化为  struct类型
			currentBlock, err := Deserialize(currentBlockBytes)
			if err != nil {
				errs = err
				break
			}
			blocks = append(blocks, currentBlock)
			//终止循环的逻辑
			if currentBlock.Height == 0 {
				break
			}
			//创世区块的hash值
			currentHash = currentBlock.PreHash[:]
		}
		return nil
	})
	return blocks, errs
}

/**
 * 该方法用于实现ChainIterator迭代器接口的方法，用于判断是否还有区块
 */
func (chain *BlockChain) HasNext() bool {
	//是否还有前一个区块
	//思路：先知道当前在哪个区块，根据当前的区块去判断是否还有下一个区块
	engine := chain.Engine
	var hasNext bool
	engine.View(func(tx *bolt.Tx) error {
		currentBlockHash := chain.IteratorBlockHash
		bucket := tx.Bucket([]byte(BLOCKS))
		if bucket == nil {
			return errors.New("区块数据文件操作失败,请重试")
		}
		currentBlockBytes := bucket.Get(currentBlockHash[:])
		currentBlock, err := Deserialize(currentBlockBytes)
		if err != nil {
			return err
		}

		hashBig := big.NewInt(0)
		hashBig = hashBig.SetBytes(currentBlock.Hash[:])
		if hashBig.Cmp(big.NewInt(0)) == 1 { //区块hash有值
			hasNext = true
		} else {
			hasNext = false
		}

		return nil
	})
	return hasNext
}

/**
 * 该方法用于实现ChainIterator迭代器接口的方法，用于取出下一个区块
 */
func (chain *BlockChain) Next() Block {
	engine := chain.Engine
	var currentBlock Block
	engine.View(func(tx *bolt.Tx) error {
		bucket := tx.Bucket([]byte(BLOCKS))
		if bucket == nil {
			return errors.New("区块数据文件操作失败,请重试！")
		}
		currentBlockBytes := bucket.Get(chain.IteratorBlockHash[:])
		currentBlock, _ = Deserialize(currentBlockBytes)
		chain.IteratorBlockHash = currentBlock.PreHash //赋值iteratorBlock，
		return nil
	})
	return currentBlock
}
