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
	fmt.Println("当前共有区块个数：",len(blockchain.Blocks))
	fmt.Println(blockchain.Blocks[0])
	fmt.Println(blockchain.Blocks[1])
	fmt.Println(blockchain.Blocks[2])

}
