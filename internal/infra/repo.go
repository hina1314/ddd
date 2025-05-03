package infra

import (
	"database/sql"
	"errors"
	"github.com/lib/pq"
)

// 检查是否为唯一键冲突错误（数据库特定实现）
//func isDuplicateKeyError(err error) bool {
//	return strings.Contains(err.Error(), "UNIQUE constraint failed") ||
//		strings.Contains(err.Error(), "Duplicate entry")
//}

// IsNotFoundError 判断是否为"未找到"错误
func IsNotFoundError(err error) bool {
	return errors.Is(err, sql.ErrNoRows)
}

func IsDuplicateKeyError(err error) bool {
	var pqErr *pq.Error
	if errors.As(err, &pqErr) && pqErr.Code == "23505" {
		return true
	}
	return false
}
