package utils

import (
	"bytes"
	"encoding/binary"
	"encoding/gob"
	"encoding/json"
)

/**
 * 将一个int64数值类型，转换为[]byte类型
 */
func Int2Byte(num int64) ([]byte, error) {
	buff := new(bytes.Buffer)
	err := binary.Write(buff, binary.BigEndian, num)
	if err != nil {
		return nil, err
	}
	return buff.Bytes(), nil
}

/**
 * 公共的gob序列化功能
 */
func GobEncode(entity interface{}) ([]byte, error) {
	buff := new(bytes.Buffer)
	encoder := gob.NewEncoder(buff)
	err := encoder.Encode(entity)
	return buff.Bytes(), err
}

/**
 * 公共的gob反序列化功能
 */
func GobDecode(data []byte, entity interface{}) (interface{}, error) {
	decoder := gob.NewDecoder(bytes.NewReader(data))
	err := decoder.Decode(entity)
	return entity, err
}

// [zhangsan,lisi]
func JSONString2Slice(data string) ([]string, error) {
	var slice []string
	err := json.Unmarshal([]byte(data), &slice)
	return slice, err
}

//JSONArray：[10.12, 5.38] ->[]float64  [10.12 5.38]
func JSONFloat2Slice(data string) ([]float64, error) {
	var slice []float64
	err := json.Unmarshal([]byte(data), &slice)
	return slice, err
}
