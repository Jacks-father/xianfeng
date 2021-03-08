package db

import (
	"github.com/boltdb/bolt"
	"XianfengChain03/chain"
	"fmt"
)

//自定义的存储引擎结构体,用于实现对区块数据的读、写操作
type DBEngine struct {
	DB *bolt.DB
}

//区块存到文件里（写）
func (engine DBEngine) SaveBlock2DB(block chain.Block) {
	fmt.Println("在该方法中将区块存到db中去")

}

//从文件中恢复区块（读）
func (engine DBEngine) GetBlockFromDB(hash [32]byte) chain.Block {
	//bolt.DB存储区块的方式:
	// key -> value
	// hash -> 序列化后的[]byte
	fmt.Println("该方法中从db中获取特定的区块")


}

