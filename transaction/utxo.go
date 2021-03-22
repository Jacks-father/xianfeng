package transaction

/**
 * 未花费的交易输出结构体
 */
type UTXO struct {
	TxId [32]byte //表明该可花的钱在哪笔交易上
	Vout int    //表明该可花的钱在该交易的哪个交易输出上
	//Value float64//可花的钱的数目
	//Owner string//该笔钱属于谁
	TxOutput //用集成TxOutput的方式来表示utxo上有多少可用的钱和该笔钱属于谁
}
