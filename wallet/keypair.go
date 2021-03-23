package wallet

import (
	"crypto/ecdsa"
	"crypto/elliptic"
	"crypto/rand"
)

/**
 * 定义Wallet结构体, 存放私钥和公钥密钥对
 */
type KeyPair struct {
	Priv *ecdsa.PrivateKey
	Pub  []byte
}

/**
 * 生成一对私钥和公钥的秘钥对，返回秘钥对结构体指针和错误信息
 */
func NewKeyPair() (*KeyPair, error) {
	curve := elliptic.P256()
	pri, err := ecdsa.GenerateKey(curve, rand.Reader)
	if err != nil {
		return nil, err
	}
	pub := elliptic.Marshal(curve, pri.X, pri.Y)
	keyPair := &KeyPair{
		Priv: pri,
		Pub:  pub,
	}
	return keyPair, nil
}
