package util

import "database/sql"

func NullStringToString(ns sql.NullString) string {
	if ns.Valid {
		return ns.String
	}
	return "" // 默认值
}

func StringToNullString(s string) sql.NullString {
	return sql.NullString{
		String: s,
		Valid:  true,
	}
}
