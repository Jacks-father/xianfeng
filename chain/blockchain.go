package chain

import (
	"github.com/boltdb/bolt"
	"errors"
)

const BLOCKS = "blocks"
const LASTHASH = "lastHash"

/**
 * 定义区块链这个结构体，用于存储产生的区块（内存中)
 */
type BlockChain struct {
	//Blocks []Block
	Engine *bolt.DB
}

func NewBlockChain(db *bolt.DB) BlockChain {
	return BlockChain{db}
}

/**
 * 创建一个区块链实例，该实例携带一个创世区块
 */
func (chain *BlockChain) CreateGenesis(genesisData []byte) {
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
			}
			//创世区块已经存在了，不需要再写入了

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
	var lastBlock Block
	var err error
	engine.View(func(tx *bolt.Tx) error {
		bucket := tx.Bucket([]byte(BLOCKS))
		if bucket == nil {
			err = errors.New("区块链数据库操作失败，请重试!")
			return err
		}
		lastHash := bucket.Get([]byte(LASTHASH))
		lastBlockData := bucket.Get(lastHash)
		//2、拿到最后一个区块的数据，进行反序列化，得到最后一个区块结构体
		lastBlock, err = Deserialize(lastBlockData)
		if err != nil {
			err = errors.New("反序列化区块发生错误, 请重试!")
			return err
		}
		return nil
	})

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
		return nil
	})
	return err
}

/**
 * 获取最新的最后的一个区块
 */
func (chain BlockChain) GetLastBlock() (Block, error) {
	engine := chain.Engine
	var err error
	var lastBlock Block
	engine.View(func(tx *bolt.Tx) error {
		bucket := tx.Bucket([]byte(BLOCKS))
		if bucket == nil {
			err = errors.New("区块数据库操作失败,请重试!")
			return err
		}
		//获取最后的区块的hash
		lastHash := bucket.Get([]byte(LASTHASH))
		//根据最后的区块的hash获取最后的区块
		lastBlockBytes := bucket.Get(lastHash)
		lastBlock, err = Deserialize(lastBlockBytes)
		if err != nil {
			return err
		}
		return nil
	})
	return lastBlock, err
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
		//首先把最新区块取出来
		lastHash := bucket.Get([]byte(LASTHASH))
		var currentHash []byte
		currentHash = lastHash
		for {
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
			if currentBlock.Height == 0{
				break
			}
			currentHash = currentBlock.PreHash[:]
		}
		return nil
	})
	return blocks, errs
}
