package wallet

import (
	"github.com/boltdb/bolt"
	"XianfengChain03/utils"
)

const KEYSTORE = "keystore"

/**
 * 定义钱包结构体，管理地址和密钥对
 */
type Wallet struct {
	// key       value
	// add1      KeyPair1
	// add2      KeyPair2
	// add3      KeyPair3
	// add4      KeyPair4
	// ...
	Address map[string]*KeyPair

	Engine *bolt.DB
}

/**
 * 定义创建新地址的方法，该方法属于Wallet的功能，由Wallet进行调用
 */
func (wallet *Wallet) CreateNewAddress() (string, error) {
	//1、获取秘钥对
	keyPair, err := NewKeyPair()
	if err != nil {
		return "", err
	}
	//2、将生成的keyPair中的公钥传递给NewAddress
	address, err := NewAddress(keyPair.Pub)
	if err != nil {
		return "", err
	}
	//放入到内存中
	wallet.Address[address] = keyPair
	//把新生成的地址和对应的KeyPair写入到DB中
	err = wallet.SaveMem2DB()
	return address, err
}

/**
 * 保存内存中的地址和密钥对信息到DB文件中
 */
func (wallet *Wallet) SaveMem2DB() error {
	var err error
	wallet.Engine.Update(func(tx *bolt.Tx) error {
		bucket := tx.Bucket([]byte(KEYSTORE))
		if bucket == nil {
			bucket, err = tx.CreateBucket([]byte(KEYSTORE))
			if err != nil {
				return err
			}
		}
		//把内存中的地址和对应的秘钥对信息存入到db中
		for key, value := range wallet.Address {
			keybytes := bucket.Get([]byte(key))
			if len(keybytes) == 0 {
				keyPairBytes, err := utils.GobEncode(value)
				if err != nil {
					return err
				}
				bucket.Put([]byte(key), keyPairBytes)
			}
		}
		return nil
	})
	return err
}
