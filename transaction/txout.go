package transaction

/**
 * 定义交易输出结构体
 */
type TxOutput struct {
	Value     float64 //锁定的币的数量
	ScriptPub []byte  //锁定的脚本，锁
}
