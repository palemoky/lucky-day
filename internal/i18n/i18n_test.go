package i18n

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestNewTranslator(t *testing.T) {
	tests := []struct {
		name     string
		lang     Language
		expected Language
	}{
		{
			name:     "创建中文翻译器",
			lang:     Chinese,
			expected: Chinese,
		},
		{
			name:     "创建英文翻译器",
			lang:     English,
			expected: English,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			translator := NewTranslator(tt.lang)

			require.NotNil(t, translator)
			assert.Equal(t, tt.expected, translator.GetLanguage())
		})
	}
}

func TestTranslator_T(t *testing.T) {
	tests := []struct {
		name     string
		lang     Language
		key      string
		expected string
	}{
		{
			name:     "中文翻译 - 应用标题",
			lang:     Chinese,
			key:      "app.title",
			expected: "幸运抽奖",
		},
		{
			name:     "英文翻译 - 应用标题",
			lang:     English,
			key:      "app.title",
			expected: "Lucky Draw",
		},
		{
			name:     "中文翻译 - 语言选择",
			lang:     Chinese,
			key:      "lang.chinese",
			expected: "中文",
		},
		{
			name:     "英文翻译 - 语言选择",
			lang:     English,
			key:      "lang.english",
			expected: "English",
		},
		{
			name:     "不存在的键 - 返回MISSING标记",
			lang:     Chinese,
			key:      "nonexistent.key",
			expected: "[MISSING: nonexistent.key]",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			translator := NewTranslator(tt.lang)
			result := translator.T(tt.key)

			assert.Equal(t, tt.expected, result)
		})
	}
}

func TestTranslator_SetLanguage(t *testing.T) {
	tests := []struct {
		name        string
		initialLang Language
		newLang     Language
		testKey     string
		expected    string
	}{
		{
			name:        "从中文切换到英文",
			initialLang: Chinese,
			newLang:     English,
			testKey:     "app.title",
			expected:    "Lucky Draw",
		},
		{
			name:        "从英文切换到中文",
			initialLang: English,
			newLang:     Chinese,
			testKey:     "app.title",
			expected:    "幸运抽奖",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			translator := NewTranslator(tt.initialLang)

			// 切换语言
			translator.SetLanguage(tt.newLang)

			// 验证语言已切换
			assert.Equal(t, tt.newLang, translator.GetLanguage())

			// 验证翻译使用新语言
			result := translator.T(tt.testKey)
			assert.Equal(t, tt.expected, result)
		})
	}
}

func TestTranslations_Completeness(t *testing.T) {
	// 验证中英文翻译键的完整性

	// 获取所有中文键
	chineseKeys := make(map[string]bool)
	for key := range translations[Chinese] {
		chineseKeys[key] = true
	}

	// 验证所有中文键在英文中都存在
	for key := range chineseKeys {
		t.Run("英文翻译存在_"+key, func(t *testing.T) {
			englishTranslation := translations[English][key]
			assert.NotEmpty(t, englishTranslation, "键 %s 在英文翻译中缺失", key)
		})
	}

	// 验证所有英文键在中文中都存在
	for key := range translations[English] {
		t.Run("中文翻译存在_"+key, func(t *testing.T) {
			chineseTranslation := translations[Chinese][key]
			assert.NotEmpty(t, chineseTranslation, "键 %s 在中文翻译中缺失", key)
		})
	}
}

func TestTranslations_KeyCategories(t *testing.T) {
	// 验证关键类别的翻译都存在
	categories := []struct {
		name string
		keys []string
	}{
		{
			name: "应用相关",
			keys: []string{"app.title", "app.exit", "app.error"},
		},
		{
			name: "语言选择",
			keys: []string{"lang.select", "lang.chinese", "lang.english"},
		},
		{
			name: "模式选择",
			keys: []string{"mode.select", "mode.excel", "mode.qr", "mode.db"},
		},
		{
			name: "奖品相关",
			keys: []string{"prize.title", "prize.remaining", "prize.total"},
		},
		{
			name: "抽奖相关",
			keys: []string{"draw.title", "draw.instruction"},
		},
		{
			name: "中奖者相关",
			keys: []string{"winner.title", "winner.name", "winner.prize"},
		},
	}

	for _, cat := range categories {
		t.Run(cat.name, func(t *testing.T) {
			translatorCN := NewTranslator(Chinese)
			translatorEN := NewTranslator(English)

			for _, key := range cat.keys {
				// 验证中文翻译
				cnTranslation := translatorCN.T(key)
				assert.NotEqual(t, key, cnTranslation, "中文翻译缺失: %s", key)

				// 验证英文翻译
				enTranslation := translatorEN.T(key)
				assert.NotEqual(t, key, enTranslation, "英文翻译缺失: %s", key)
			}
		})
	}
}
