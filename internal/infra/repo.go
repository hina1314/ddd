package infra

import (
	"database/sql"
	"errors"

	"github.com/jackc/pgx/v5/pgconn"
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
	var pgErr *pgconn.PgError
	if errors.As(err, &pgErr) && pgErr.Code == "23505" {
		return true
	}
	return false
}
