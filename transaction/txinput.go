package transaction

/**
 * 定义交易输入结构体
 */
type TxInput struct {
	TxId      [32]byte //标识引用自哪笔交易
	Vout      int      //引用自哪个交易输出
	ScriptSig []byte   //解锁脚本
}
