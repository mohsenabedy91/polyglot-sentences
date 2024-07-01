package translation_test

import (
	"encoding/json"
	"github.com/stretchr/testify/require"
	"os"
	"testing"

	"github.com/mohsenabedy91/polyglot-sentences/internal/core/config"
	"github.com/mohsenabedy91/polyglot-sentences/pkg/translation"
	"github.com/nicksnyder/go-i18n/v2/i18n"
	"golang.org/x/text/language"
)

func setup(t *testing.T) {
	err := os.MkdirAll("testdata/locales", os.ModePerm)
	require.NoError(t, err)

	file, err := os.Create("testdata/locales/en.json")
	require.NoError(t, err)

	write, err := file.WriteString(`{"hello_world": "Hello, World!"}`)
	require.NoError(t, err)
	require.Greater(t, write, 0)

	err = file.Close()
	require.NoError(t, err)
}

func teardown(t *testing.T) {
	err := os.RemoveAll("testdata")
	require.NoError(t, err)
}

func TestInitialize(t *testing.T) {
	setup(t)
	defer teardown(t)

	conf := config.App{
		PathLocale: "testdata/locales",
	}

	translation.Initialize(conf)

	expectedMessage := "Hello, World!"
	localizer := i18n.NewLocalizer(translation.Bundle, "en")
	message, err := localizer.Localize(&i18n.LocalizeConfig{
		MessageID: "hello_world",
	})
	require.NoError(t, err)
	require.Equal(t, expectedMessage, message)
}

func TestInitialize_InvalidJSONFile(t *testing.T) {
	setup(t)
	defer teardown(t)

	invalidJSONFile, err := os.Create("testdata/locales/invalid.json")
	require.NoError(t, err)

	write, err := invalidJSONFile.WriteString(`{"invalid_json": `)
	require.NoError(t, err)
	require.Greater(t, write, 0)

	err = invalidJSONFile.Close()
	require.NoError(t, err)

	conf := config.App{
		PathLocale: "testdata/locales",
	}

	translation.Initialize(conf)

	localizer := i18n.NewLocalizer(translation.Bundle, "en")
	message, err := localizer.Localize(&i18n.LocalizeConfig{
		MessageID: "invalid_json",
	})
	require.Error(t, err)
	require.Equal(t, "", message)
}

func TestGetLocalizer(t *testing.T) {
	setup(t)
	defer teardown(t)

	conf := config.App{
		PathLocale: "testdata/locales",
		Locale:     "en",
	}
	trans := translation.NewTranslation(conf)

	enMessage := &i18n.Message{
		ID:    "test_message",
		Other: "This is an English message",
	}
	esMessage := &i18n.Message{
		ID:    "test_message",
		Other: "Este es un mensaje en español",
	}

	err := translation.Bundle.AddMessages(language.English, enMessage)
	require.NoError(t, err)

	err = translation.Bundle.AddMessages(language.Spanish, esMessage)
	require.NoError(t, err)

	localizer := trans.GetLocalizer("")
	require.NotNil(t, localizer)
	message, err := localizer.Localize(&i18n.LocalizeConfig{
		MessageID: "test_message",
	})
	require.NoError(t, err)
	require.Equal(t, "This is an English message", message)

	localizer = trans.GetLocalizer("es")
	require.NotNil(t, localizer)
	message, err = localizer.Localize(&i18n.LocalizeConfig{
		MessageID: "test_message",
	})
	require.NoError(t, err)
	require.Equal(t, "Este es un mensaje en español", message)

	localizer = trans.GetLocalizer("invalid-lang")
	require.NotNil(t, localizer)
	message, err = localizer.Localize(&i18n.LocalizeConfig{
		MessageID: "test_message",
	})
	require.NoError(t, err)
	require.Equal(t, "This is an English message", message)
}

func TestLang(t *testing.T) {
	setup(t)
	defer teardown(t)

	langEnglish := "en"

	conf := config.App{
		PathLocale: "testdata/locales",
		Locale:     langEnglish,
	}
	trans := translation.NewTranslation(conf)

	translation.Bundle = i18n.NewBundle(language.English)
	translation.Bundle.RegisterUnmarshalFunc("json", json.Unmarshal)

	enMessage := &i18n.Message{
		ID:    "test_message",
		Other: "This is a test message",
	}
	err := translation.Bundle.AddMessages(language.English, enMessage)
	require.NoError(t, err)

	translation.AcceptLanguage = i18n.NewLocalizer(translation.Bundle, langEnglish)

	message := trans.Lang("test_message", nil, nil)
	require.Equal(t, "This is a test message", message)

	message = trans.Lang("non_existent_key", nil, nil)
	require.Equal(t, "non_existent_key", message)

	message = trans.Lang("test_message", nil, &langEnglish)
	require.Equal(t, "This is a test message", message)
}
