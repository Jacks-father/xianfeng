package chain

import (
	"time"
	"bytes"
	"XianfengChain03/utils"
	"crypto/sha256"
)

const VERSION = 0x00

/**
 * 区块数据结构的定义
 */
type Block struct {
	Height  int64 // 高度
	Version int64
	PreHash [32]byte
	Hash    [32]byte //区块hash
	//默克尔根
	Timestamp int64
	//Difficulty int64
	Nonce int64
	Data  []byte //区块体
}

/**
 * 该方法用于计算区块的hash值
 */
func (block *Block) SetHash() {
	heightByte, _ := utils.Int2Byte(block.Height)
	versionByte, _ := utils.Int2Byte(block.Version)
	timeByte, _ := utils.Int2Byte(block.Timestamp)
	nonceByte, _ := utils.Int2Byte(block.Nonce)
	bk := bytes.Join([][]byte{heightByte, versionByte, block.PreHash[:], timeByte, nonceByte, block.Data}, []byte{})
	hash := sha256.Sum256(bk)
	block.Hash = hash
}

/**
 * 创建一个新的区块的函数
 */
func CreateBlock(height int64, prevHash [32]byte, data []byte) Block {
	block := Block{}
	block.Height = height + 1
	block.PreHash = prevHash
	block.Version = VERSION
	block.Timestamp = time.Now().Unix()
	block.Data = data

	block.SetHash() //计算hash
	return block
}

/**
 * 封装用于生成创世区块的函数, 该函数只生成创世区块
 */
func CreateGenesisBlock(data []byte) Block {
	genesis := Block{}
	genesis.Height = 0
	genesis.PreHash = [32]byte{0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0}
	genesis.Version = VERSION
	genesis.Timestamp = time.Now().Unix()
	genesis.Data = data
	genesis.SetHash() //计算hash值
	return genesis
}
