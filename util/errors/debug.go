package errors

import (
	"encoding/json"
	"fmt"
	"runtime"
	"strings"
)

// DebugMode 控制是否启用错误详细信息
var DebugMode bool = false

// SetDebugMode 设置调试模式
func SetDebugMode(debug bool) {
	DebugMode = debug
}

// ErrorTrace 表示错误追踪信息
type ErrorTrace struct {
	Message  string      `json:"message"`
	Code     ErrorCode   `json:"code"`
	Location string      `json:"location,omitempty"`
	Params   interface{} `json:"params,omitempty"`
	Cause    *ErrorTrace `json:"cause,omitempty"`
}

// GetErrorTrace 获取完整的错误链
func GetErrorTrace(err error) *ErrorTrace {
	if err == nil {
		return nil
	}

	var trace *ErrorTrace

	if de, ok := err.(*DomainError); ok {
		// 记录文件和行号
		_, file, line, _ := runtime.Caller(1)
		paths := strings.Split(file, "/")
		file = paths[len(paths)-1]
		location := fmt.Sprintf("%s:%d", file, line)

		trace = &ErrorTrace{
			Message:  de.Message,
			Code:     de.Code,
			Location: location,
			Params:   de.Params,
		}

		// 递归获取原因链
		if de.Cause != nil {
			trace.Cause = GetErrorTrace(de.Cause)
		}
	} else {
		// 普通错误
		trace = &ErrorTrace{
			Message: err.Error(),
		}
	}

	return trace
}

// AddDebugInfo 在HTTP响应中添加调试信息
func AddDebugInfo(response map[string]interface{}, err error) {
	if DebugMode && err != nil {
		response["debug"] = GetErrorTrace(err)
	}
}

// ErrorStack 获取完整错误堆栈字符串
func ErrorStack(err error) string {
	trace := GetErrorTrace(err)
	if trace == nil {
		return ""
	}

	bytes, _ := json.MarshalIndent(trace, "", "  ")
	return string(bytes)
}
