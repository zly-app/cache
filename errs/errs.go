package errs

import (
	"errors"
)

// 查询缓存不存在应该返回这个错误
var CacheMiss = errors.New("cache miss")

// 数据为nil
var DataIsNil = errors.New("data is nil")
