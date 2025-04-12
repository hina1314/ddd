package i18n

import (
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"
	"strings"
	"sync"

	"gopkg.in/yaml.v3"
)

// FileTranslator implements a file-based Translator for JSON and YAML translation files.
type FileTranslator struct {
	translations   *sync.Map // map[string]map[string]string (locale -> key -> message)
	fallbackLocale string
	once           sync.Once
}

// NewFileTranslator creates a new FileTranslator with the specified fallback locale.
func NewFileTranslator(fallbackLocale string) *FileTranslator {
	return &FileTranslator{
		translations:   &sync.Map{},
		fallbackLocale: fallbackLocale,
	}
}

// LoadTranslations loads translation files (JSON or YAML) from the specified directory.
// Files should be named as `<locale>.json` or `<locale>.yaml` (e.g., en.json, zh.yaml).
func (t *FileTranslator) LoadTranslations(dir string) error {
	var loadErr error
	t.once.Do(func() {
		err := filepath.Walk(dir, func(path string, info os.FileInfo, err error) error {
			if err != nil {
				return fmt.Errorf("failed to access path %s: %w", path, err)
			}

			// Skip directories
			if info.IsDir() {
				return nil
			}

			// Only process .json, .yaml, or .yml files
			ext := strings.ToLower(filepath.Ext(path))
			if ext != ".json" && ext != ".yaml" && ext != ".yml" {
				return nil
			}

			// Extract locale from filename (e.g., en.json -> en)
			base := filepath.Base(path)
			locale := strings.TrimSuffix(base, ext)
			if locale == "" {
				return fmt.Errorf("invalid translation file name: %s", path)
			}

			// Read file content
			data, err := os.ReadFile(path)
			if err != nil {
				return fmt.Errorf("failed to read translation file %s: %w", path, err)
			}

			// Parse file content
			var messages map[string]string
			switch ext {
			case ".json":
				if err := json.Unmarshal(data, &messages); err != nil {
					return fmt.Errorf("failed to parse JSON file %s: %w", path, err)
				}
			case ".yaml", ".yml":
				if err := yaml.Unmarshal(data, &messages); err != nil {
					return fmt.Errorf("failed to parse YAML file %s: %w", path, err)
				}
			}

			// Validate messages
			if len(messages) == 0 {
				return fmt.Errorf("no translations found in file %s", path)
			}

			// Store translations atomically
			t.translations.Store(locale, messages)
			return nil
		})

		if err != nil {
			loadErr = fmt.Errorf("failed to load translations from %s: %w", dir, err)
		}
	})

	return loadErr
}

// T translates a key for the given locale, falling back to the fallback locale or key if no translation is found.
func (t *FileTranslator) T(key string, locale string, params map[string]interface{}) string {
	// Try the specified locale
	if messages, ok := t.translations.Load(locale); ok {
		if msgMap, ok := messages.(map[string]string); ok {
			if message, ok := msgMap[key]; ok {
				return t.interpolate(message, params)
			}
		}
	}

	// Try the fallback locale if different
	if locale != t.fallbackLocale {
		if messages, ok := t.translations.Load(t.fallbackLocale); ok {
			if msgMap, ok := messages.(map[string]string); ok {
				if message, ok := msgMap[key]; ok {
					return t.interpolate(message, params)
				}
			}
		}
	}

	// Return the key if no translation is found
	return key
}

// SetFallbackLocale sets the fallback locale for translations.
func (t *FileTranslator) SetFallbackLocale(locale string) {
	t.fallbackLocale = locale
}

// interpolate replaces placeholders in the message with parameter values.
func (t *FileTranslator) interpolate(message string, params map[string]interface{}) string {
	if params == nil {
		return message
	}

	result := message
	for key, value := range params {
		// Prevent injection by ensuring placeholders are safe
		if !isValidParamKey(key) {
			continue
		}
		placeholder := fmt.Sprintf("{%s}", key)
		result = strings.ReplaceAll(result, placeholder, fmt.Sprintf("%v", value))
	}
	return result
}

// isValidParamKey checks if a parameter key is safe to use in interpolation.
func isValidParamKey(key string) bool {
	// Basic validation: alphanumeric and underscore only
	for _, r := range key {
		if !((r >= 'a' && r <= 'z') || (r >= 'A' && r <= 'Z') || (r >= '0' && r <= '9') || r == '_') {
			return false
		}
	}
	return len(key) > 0
}
