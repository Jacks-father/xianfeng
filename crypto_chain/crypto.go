package crypto_chain

import (
	"crypto/ecdsa"
	"crypto/elliptic"
	"crypto/rand"
)

/**
 * 随机生成一个私钥
 */
func NewKey(curve elliptic.Curve) (*ecdsa.PrivateKey, error) {
	//ECC：椭圆曲线加密算法 elliptic curve cryptic
	//椭圆曲线数字签名算法
	//elliptic curve digital signature algorithm
	return ecdsa.GenerateKey(curve, rand.Reader)
}

/**
 * 根据私钥返回对应的公钥
 */
func GetPub(curve elliptic.Curve, pri *ecdsa.PrivateKey) []byte {
	return elliptic.Marshal(curve, pri.X, pri.Y)
}
