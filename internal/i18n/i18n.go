package i18n

import "fmt"

// Language represents supported languages
type Language string

const (
	Chinese Language = "zh"
	English Language = "en"
)

// Translator handles translations
type Translator struct {
	lang         Language
	translations map[Language]map[string]string
}

// NewTranslator creates a new translator with the specified language
func NewTranslator(lang Language) *Translator {
	return &Translator{
		lang:         lang,
		translations: translations,
	}
}

// T translates a key to the current language
func (t *Translator) T(key string) string {
	if langMap, ok := t.translations[t.lang]; ok {
		if translation, ok := langMap[key]; ok {
			return translation
		}
	}
	// Fallback to Chinese if key not found
	if langMap, ok := t.translations[Chinese]; ok {
		if translation, ok := langMap[key]; ok {
			return translation
		}
	}
	return fmt.Sprintf("[MISSING: %s]", key)
}

// SetLanguage changes the current language
func (t *Translator) SetLanguage(lang Language) {
	t.lang = lang
}

// GetLanguage returns the current language
func (t *Translator) GetLanguage() Language {
	return t.lang
}
