package chain

import (
	"time"
	"XianfengChain03/consensus"
	"encoding/gob"
	"bytes"
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
func (block *Block) SetHash() {
	heightByte, _ := utils.Int2Byte(block.Height)
	versionByte, _ := utils.Int2Byte(block.Version)
	timeByte, _ := utils.Int2Byte(block.Timestamp)
	nonceByte, _ := utils.Int2Byte(block.Nonce)
	bk := bytes.Join([][]byte{heightByte, versionByte, block.PreHash[:], timeByte, nonceByte, block.Data}, []byte{})
	hash := sha256.Sum256(bk)
	block.Hash = hash
}
*/

/**
 * 区块的序列化，序列化为[]byte数据类型
 */
func (block *Block) Serialize() ([]byte, error) {
	buff := new(bytes.Buffer)
	encoder := gob.NewEncoder(buff)
	err := encoder.Encode(&block)
	return buff.Bytes(), err
}

/**
 * 区块的反序列化操作,传入[]byte,返回Block结构体或者error
 */
func Deserialize(data []byte) (Block, error) {
	var block Block
	decoder := gob.NewDecoder(bytes.NewReader(data))
	err := decoder.Decode(&block)
	return block, err
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

	//尝试给nonce值赋值
	//共识机制: PoW、PoS、....
	//确定选用pow实现共识机制

	proof := consensus.NewProofWork(block)
	hash, nonce := proof.SearchNonce()
	block.Nonce = nonce

	block.Hash = hash

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

	proof := consensus.NewProofWork(genesis)
	hash, nonce := proof.SearchNonce()
	genesis.Hash = hash
	genesis.Nonce = nonce

	return genesis
}

/**
 * 该方法是实现BlockInterface的GetHeight方法
 */
func (block Block) GetHeight() int64 {
	return block.Height
}

/**
 * 该方法是实现BlockInterface的GetVersion方法
 */
func (block Block) GetVersion() int64 {
	return block.Version
}

func (block Block) GetTimeStamp() int64 {
	return block.Timestamp
}

func (block Block) GetPreHash() [32]byte {
	return block.PreHash
}

func (block Block) GetData() []byte {
	return block.Data
}
