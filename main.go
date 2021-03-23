package main

import (
	"XianfengChain03/client"
	"github.com/boltdb/bolt"
	"XianfengChain03/chain"
)

const DBFILE = "xianfeng03.db"
//blocks:  hash为key blockBytes为value
//keystore：addr为key KeyPairBytes为value

/**
 * 项目的主入口
 */
func main() {

	engine, err := bolt.Open(DBFILE, 0600, nil)
	if err != nil {
		panic(err.Error())
	}

	blockChain, err := chain.NewBlockChain(engine)
	if err != nil {
		panic(err.Error())
		return
	}
	cli := client.Client{blockChain}
	cli.Run()

}
