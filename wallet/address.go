package wallet

import (
	"XianfengChain03/utils"
	"bytes"
)

const VERSION = 0x00

/**
 * 生成一个新地址
 */
func NewAddress(pub []byte) (string, error) {

	//3、sha256(pub）
	hashPub := utils.Hash256(pub)

	//4、ripemd160
	ripemdPub := utils.RipeMd160(hashPub[:])

	//5、拼接version版本号
	versionPub := append([]byte{VERSION}, ripemdPub...)

	//6、对ripemd160 双hash
	firstRipemd := utils.Hash256(versionPub)
	secondRipemd := utils.Hash256(firstRipemd[:])

	//7、截取前4个字节作为校验位
	checkBytes := secondRipemd[:4]

	//8、步骤5versionPub与校验位进行拼接
	origAddress := append(versionPub, checkBytes...)

	//9、base58编码
	address := utils.Encode(origAddress)
	return address, nil
}

/**
 *  定义该函数，用于判断和校验给定的一个字符串是否是符合地址的计算规范
 */
func IsAddressValid(addr string) bool {
	//1、base58反编码
	reverseAddr := utils.Decode(addr)

	//2、取反编码以后的后4个字节作为校验位
	check := reverseAddr[len(reverseAddr)-4:]

	//3、获取到versionPub
	versionPub := reverseAddr[:len(reverseAddr)-4]

	//4、对versionPub双hash，重新计算校验位
	firstHash := utils.Hash256(versionPub)
	secondHash := utils.Hash256(firstHash)

	reCheck := secondHash[:4]

	//5, 比较 check 和 reCheck
	return bytes.Compare(check, reCheck) == 0
}
