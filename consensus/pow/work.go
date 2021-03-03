package pow

import (
	"XianfengChain03/chain"
	"fmt"
	"XianfengChain03/utils"
	"bytes"
	"crypto/sha256"
	"math/big"
)

// 256位二进制
//思路：给一个大整数，初始值为1，根据自己需要的难度进行左移位，左移的位数是256-0的个数

//000000010000000....000

const DIFFICULTY = 10 //初始难度为10，即大整数的开头有10个零

/**
 * 工作量证明
 */
type ProofWork struct {
	Block  chain.Block
	Target *big.Int
}

/**
 * 实现共识机制接口的方法
 */
func (work ProofWork) SearchNonce() ([32]byte, int64) {
	fmt.Println("这里是PoW的方法的代码实现过程")
	//block -> nonce
	// block哈希 < 系统提供的某个目标值
	//1 给定一个non值，计算带有non的区块哈希
	var nonce int64
	nonce = 0
	for {
		hash := CalculateBlockHash(work.Block, nonce)
		//2 系统给定的值
		target := work.Target
		//3 拿1和2比较
		result := bytes.Compare(hash[:], target.Bytes())
		//4 判断结果，区块哈希<给定值，返回non;
		if result == -1 {
			return hash, nonce
		}
		//否则，non自增
		nonce++
	}
}

/**
 * 根据当前的区块和当前的non值，计算区块的哈希值
 */
func CalculateBlockHash(block chain.Block, nonce int64) [32]byte {
	heightByte, _ := utils.Int2Byte(block.Height)
	versionByte, _ := utils.Int2Byte(block.Version)
	timeByte, _ := utils.Int2Byte(block.Timestamp)
	nonceByte, _ := utils.Int2Byte(nonce)
	bk := bytes.Join([][]byte{heightByte,
		versionByte,
		block.PreHash[:],
		timeByte,
		nonceByte,
		block.Data},
		[]byte{})
	return sha256.Sum256(bk)
}
