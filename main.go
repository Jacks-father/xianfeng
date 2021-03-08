package main

import (
	"fmt"
	"XianfengChain03/chain"
	"github.com/boltdb/bolt"
)

const DBFILE = "xianfeng03.db"

/**
 * 项目的主入口
 */
func main() {

	fmt.Println("hello world")

	engine, err := bolt.Open(DBFILE, 0600, nil)
	if err != nil {
		panic(err.Error())
	}

	blockChain := chain.NewBlockChain(engine)
	//创世区块
	blockChain.CreateGenesis([]byte("hello world"))
	//新增一个区块
	err = blockChain.AddNewBlock([]byte("hello"))
	if err != nil {
		fmt.Println(err.Error())
		return
	}

	//获取最新区块
	//lastBlock, err := blockChain.GetLastBlock()
	//if err != nil {
	//	fmt.Println(err.Error())
	//	return
	//}
	//fmt.Println(lastBlock)

	allBlocks, err := blockChain.GetAllBlocks()
	if err != nil {
		fmt.Println(err.Error())
		return
	}
	for _, block := range allBlocks{
		fmt.Println(block)
	}

}
