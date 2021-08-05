package models

import (
	"ch35/goblog/pkg/types"
)

// BaseModel 模型基类
type BaseModel struct {
	ID uint64
}

// GetStringID 获取 ID 的字符串格式
func (a BaseModel) GetStringID() string {
	return types.Uint64ToString(a.ID)
}