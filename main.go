package main

import (
	"fmt"
	"XianfengChain03/chain"
)

/**
 * 项目的主入口
 */
func main() {
	fmt.Println("hello world")

	blockchain := chain.CreateChainWithGenesis([]byte("Hello World"))

	blockchain.AddNewBlock([]byte("block1"))
	blockchain.AddNewBlock([]byte("block2"))
	fmt.Println("当前共有区块个数：", len(blockchain.Blocks))
	block0 := blockchain.Blocks[0]
	block0SerBytes, err := block0.Serialize()
	if err != nil {
		fmt.Println("序列化区块0出现错误")
		return
	}

	deBlock0, err := chain.Deserialize(block0SerBytes)
	if err != nil {
		fmt.Println("反序列化区块0出现错误，程序已停止。",err.Error())
		return
	}
	fmt.Println(string(deBlock0.Data))
}
