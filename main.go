package main

import (
	"XianfengChain03/client"
	"github.com/boltdb/bolt"
	"XianfengChain03/chain"
)

const DBFILE = "xianfeng03.db"

/**
 * 项目的主入口
 */
func main() {

	engine, err := bolt.Open(DBFILE, 0600, nil)
	if err != nil {
		panic(err.Error())
	}
	blockChain := chain.NewBlockChain(engine)

	cli := client.Client{blockChain}
	cli.Run()

}
