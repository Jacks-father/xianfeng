package pos

import (
	"XianfengChain03/chain"
	"fmt"
)

type ProofStock struct {
	Block chain.Block
}

func (stock ProofStock) SearchNonce()([32]byte,int64){
	fmt.Println("我是新来的，这是我写的共识机制的PoS的实现方法")
	return [32]byte{},0
}



