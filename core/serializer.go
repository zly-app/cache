package core

// 序列化器
type ISerializer interface {
	// 序列化
	MarshalBytes(a interface{}) ([]byte, error)
	// 反序列化
	UnmarshalBytes(data []byte, a interface{}) error
}
