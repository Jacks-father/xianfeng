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

	gensis := chain.CreateGenesisBlock([]byte("hello world"))
	fmt.Println("新区块：", gensis)
}
