package chain

/**
 * 迭代器的接口标准声明:
		① 判断是否还有数据
		② 取出下一个数据
 */
type ChainIterator interface {
	HasNext() bool
	Next() Block
}
