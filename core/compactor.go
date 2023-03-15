package core

// 压缩器
type ICompactor interface {
	// 压缩
	CompressBytes(in []byte) (out []byte, err error)
	// 解压缩
	UnCompressBytes(in []byte) (out []byte, err error)
}
