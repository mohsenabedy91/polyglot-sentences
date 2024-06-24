package translation

import (
	"encoding/json"
	"fmt"
	"github.com/mohsenabedy91/polyglot-sentences/internal/core/config"
	"github.com/nicksnyder/go-i18n/v2/i18n"
	"golang.org/x/text/language"
	"os"
	"path/filepath"
)

var (
	AcceptLanguage *i18n.Localizer
	Bundle         *i18n.Bundle
)

func Initialize(cfg config.App) {
	Bundle = i18n.NewBundle(language.English)
	Bundle.RegisterUnmarshalFunc("json", json.Unmarshal)

	root := cfg.PathLocale
	err := filepath.Walk(root, func(path string, info os.FileInfo, err error) error {
		if err != nil {
			fmt.Println("Error:", err)
			return err
		}

		// Check if the current path is a JSON file
		if !info.IsDir() && filepath.Ext(path) == ".json" {
			if _, err = Bundle.LoadMessageFile(path); err != nil {
				fmt.Printf("Failed to load message file %s: %v\n", path, err)
			}
		}

		return nil
	})

	if err != nil {
		fmt.Println("Error:", err)
		return
	}
}

type Translation struct {
	cfg config.App
}

// Translator is an interface that defines the methods needed to translate messages.
type Translator interface {
	GetLocalizer(lang string) *i18n.Localizer
	Lang(key string, args map[string]interface{}, lang *string) string
}

func NewTranslation(cfg config.App) *Translation {
	Initialize(cfg)
	return &Translation{
		cfg: cfg,
	}
}

// GetLocalizer initializes the localizer with the desired language.
func (r *Translation) GetLocalizer(lang string) *i18n.Localizer {
	if lang == "" {
		lang = r.cfg.Locale
	}
	tag, err := language.Parse(lang)
	if err != nil {
		fmt.Println("Failed to parse language tag:", err)
		tag = language.English
	}

	AcceptLanguage = i18n.NewLocalizer(Bundle, tag.String())

	return AcceptLanguage
}

// Lang is a helper function that translates a message.
func (r *Translation) Lang(key string, args map[string]interface{}, lang *string) string {
	cfg := &i18n.LocalizeConfig{
		MessageID:    key,
		TemplateData: args,
	}

	if lang != nil {
		AcceptLanguage = r.GetLocalizer(*lang)
	}

	message, err := AcceptLanguage.Localize(cfg)
	if err != nil {
		defaultLang := i18n.NewLocalizer(Bundle, r.cfg.FallbackLocale)
		message, err = defaultLang.Localize(cfg)
		if err != nil {
			return key
		}
	}

	return message
}
