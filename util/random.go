package util

import (
	"math/rand"
	"time"
)

var letters = []rune("abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ0123456789")

type RandUtil struct {
	r *rand.Rand
}

// NewRandUtil 创建一个随机工具类：
// - 不传 seed → 使用 time.Now().UnixNano()
// - 传入 seed → 可重复生成随机序列（用于测试）
func NewRandUtil(seed ...int64) *RandUtil {
	var src rand.Source
	if len(seed) > 0 {
		src = rand.NewSource(seed[0])
	} else {
		src = rand.NewSource(time.Now().UnixNano())
	}
	return &RandUtil{r: rand.New(src)}
}

// Int 返回 [min, max) 的随机整数
func (ru *RandUtil) Int(min, max int) int {
	if max <= min {
		return min
	}
	return ru.r.Intn(max-min) + min
}

// String 返回指定长度的随机字符串
func (ru *RandUtil) String(n int) string {
	b := make([]rune, n)
	for i := range b {
		b[i] = letters[ru.r.Intn(len(letters))]
	}
	return string(b)
}
