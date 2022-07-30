package errs

import (
	"errors"
)

// 查询缓存不存在应该返回这个错误
var CacheMiss = errors.New("cache miss")

// 数据为nil
var DataIsNil = errors.New("data is nil")

var _ error = QueryErr{}

type QueryErr struct {
	mainErr error // 主要错误
	errs    map[string]error
}

func (qe QueryErr) Error() string {
	return qe.mainErr.Error()
}

// 获取指定key的error
func (qe QueryErr) GetError(key string) error {
	return qe.errs[key]
}

func NewQueryErr(errs map[string]error) error {
	if len(errs) == 0 {
		return nil
	}
	qe := QueryErr{
		errs: errs,
	}

	getErrorLevel := func(err error) int {
		switch err {
		case CacheMiss:
			return 1
		case DataIsNil:
			return 2
		}
		return 0
	}
	for _, err := range errs {
		if err == nil {
			continue
		}
		if qe.mainErr == nil || getErrorLevel(err) < getErrorLevel(qe.mainErr) {
			qe.mainErr = err
		}
	}
	if qe.mainErr == nil {
		return nil
	}
	return qe
}

// 获取指定key的error
func GetKeyError(err error, key string) error {
	qe, ok := err.(QueryErr)
	if !ok {
		return nil
	}
	return qe.GetError(key)
}
