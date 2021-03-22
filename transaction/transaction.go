package transaction

import (
	"XianfengChain03/utils"
	"crypto/sha256"
)

const REWARD = 50

/**
 * 定义交易的结构体
 */
type Transaction struct {
	TxHash  [32]byte   //交易的唯一标识
	Inputs  []TxInput  //交易的交易输入
	Outputs []TxOutput //交易的交易输出
}

/**
 * 该函数用于构建一个coinbase交易，返回一个交易的结构体实例
 */
func NewCoinbaseTx(address string) (*Transaction, error) {
	//txInput 为空

	//构造txOutput
	txOutput := TxOutput{
		Value:     REWARD,
		ScriptPub: []byte(address),
	}
	tx := Transaction{
		Inputs:  []TxInput{},
		Outputs: []TxOutput{txOutput},
	}

	//序列化
	txBytes, err := utils.GobEncode(tx)
	if err != nil {
		return nil, err
	}
	//交易哈希计算，并赋值给TxHash字段
	tx.TxHash = sha256.Sum256(txBytes)
	return &tx, nil
}

/**
 * 构建一笔新的交易
 */
func NewTransaction(spent []UTXO, from string, to string, value float64) (*Transaction, error) {
	//交易输入的容器切片
	txInputs := make([]TxInput, 0)
	var inputAmount float64
	for _, utxo := range spent {
		inputAmount += utxo.Value
		input := TxInput{
			TxId:      utxo.TxId,
			Vout:      utxo.Vout,
			ScriptSig: []byte(from),
		}
		//把构建好的交易输入放入到交易输入容器中
		txInputs = append(txInputs, input)
	}

	//交易输出的容器切片
	//A->B 10 至多会产生两个交易输出
	txOutputs := make([]TxOutput, 0)

	//第一个交易输出：对应转账接收者的输出
	txOutput0 := TxOutput{
		Value:     value,
		ScriptPub: []byte(to),
	}
	txOutputs = append(txOutputs, txOutput0)

	//还有可能产生找零的一个输出：交易发起者给的钱比要转账的钱多
	if inputAmount-value > 0 { //需要找零给转账发起人
		txOutput1 := TxOutput{
			Value:     inputAmount - value,
			ScriptPub: []byte(from),
		}
		txOutputs = append(txOutputs, txOutput1)
	}

	//构建交易
	tx := Transaction{
		Inputs:  txInputs,
		Outputs: txOutputs,
	}
	//序列化
	txBytes, err := utils.GobEncode(tx)
	if err != nil {
		return nil, err
	}
	tx.TxHash = sha256.Sum256(txBytes)

	return &tx, nil
}
