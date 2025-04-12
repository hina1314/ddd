package repository

import (
	"errors"
	"github.com/lib/pq"
	er "study/util/errors"
)

// 检查是否为唯一键冲突错误（数据库特定实现）
//func isDuplicateKeyError(err error) bool {
//	return strings.Contains(err.Error(), "UNIQUE constraint failed") ||
//		strings.Contains(err.Error(), "Duplicate entry")
//}

// IsNotFoundError 判断是否为"未找到"错误
func IsNotFoundError(err error) bool {
	var domainErr *er.DomainError
	if errors.As(err, &domainErr) {
		return domainErr.Code == er.ErrUserNotFound
	}
	return false
}

func isDuplicateKeyError(err error) bool {
	var pqErr *pq.Error
	if errors.As(err, &pqErr) {
		switch pqErr.Code.Name() {
		case "unique_violation":
			return true
		}
	}
	return false
}
