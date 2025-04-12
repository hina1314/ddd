package i18n

import (
	"encoding/json"
	"fmt"
	"gopkg.in/yaml.v2"
	"os"
	"path/filepath"
	"strings"
	"sync"
)

// Translator 是一个多语言翻译器接口
type Translator interface {
	// T 翻译一个键
	T(key string, locale string, params map[string]interface{}) string

	// LoadTranslations 加载翻译数据
	LoadTranslations(dir string) error

	// SetFallbackLocale 设置fallback语言
	SetFallbackLocale(locale string)
}

// FileTranslator 实现了基于文件的翻译器
type FileTranslator struct {
	mu             sync.RWMutex
	translations   map[string]map[string]string // locale -> key -> message
	fallbackLocale string
}

// NewFileTranslator 创建一个新的文件翻译器
func NewFileTranslator(fallbackLocale string) *FileTranslator {
	return &FileTranslator{
		translations:   make(map[string]map[string]string),
		fallbackLocale: fallbackLocale,
	}
}

// LoadTranslations 从目录加载所有翻译文件
func (t *FileTranslator) LoadTranslations(dir string) error {
	t.mu.Lock()
	defer t.mu.Unlock()

	// 遍历目录中的所有文件
	err := filepath.Walk(dir, func(path string, info os.FileInfo, err error) error {
		if err != nil {
			return err
		}

		// 跳过目录
		if info.IsDir() {
			return nil
		}

		// 只处理yaml和json文件
		ext := strings.ToLower(filepath.Ext(path))
		if ext != ".yaml" && ext != ".yml" && ext != ".json" {
			return nil
		}

		// 获取locale，假设文件名格式为 locale.ext，如 en.yml, zh.json
		base := filepath.Base(path)
		locale := strings.TrimSuffix(base, ext)

		// 初始化map
		if _, ok := t.translations[locale]; !ok {
			t.translations[locale] = make(map[string]string)
		}

		// 读取文件内容
		data, err := os.ReadFile(path)
		if err != nil {
			return fmt.Errorf("error reading translation file %s: %w", path, err)
		}

		// 解析文件内容
		var messages map[string]string

		if ext == ".json" {
			if err := json.Unmarshal(data, &messages); err != nil {
				return fmt.Errorf("error parsing JSON translation file %s: %w", path, err)
			}
		} else {
			if err := yaml.Unmarshal(data, &messages); err != nil {
				return fmt.Errorf("error parsing YAML translation file %s: %w", path, err)
			}
		}

		// 添加到翻译映射
		for key, message := range messages {
			t.translations[locale][key] = message
		}

		return nil
	})

	return err
}

// T 根据键和语言环境翻译消息
func (t *FileTranslator) T(key string, locale string, params map[string]interface{}) string {
	t.mu.RLock()
	defer t.mu.RUnlock()

	// 尝试获取指定语言的翻译
	if localeMessages, ok := t.translations[locale]; ok {
		if message, ok := localeMessages[key]; ok {
			return t.interpolate(message, params)
		}
	}

	// 尝试获取fallback语言的翻译
	if locale != t.fallbackLocale {
		if fallbackMessages, ok := t.translations[t.fallbackLocale]; ok {
			if message, ok := fallbackMessages[key]; ok {
				return t.interpolate(message, params)
			}
		}
	}

	// 没有找到翻译，返回键本身
	return key
}

// SetFallbackLocale 设置fallback语言
func (t *FileTranslator) SetFallbackLocale(locale string) {
	t.mu.Lock()
	t.fallbackLocale = locale
	t.mu.Unlock()
}

// interpolate 替换消息中的参数占位符
func (t *FileTranslator) interpolate(message string, params map[string]interface{}) string {
	if params == nil {
		return message
	}

	result := message
	for key, value := range params {
		placeholder := fmt.Sprintf("{%s}", key)
		result = strings.Replace(result, placeholder, fmt.Sprintf("%v", value), -1)
	}
	return result
}
