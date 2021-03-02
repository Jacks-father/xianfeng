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

	block := chain.CreateBlock(0, [32]byte{0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0}, nil)

	fmt.Println("新区块：",block)
}
