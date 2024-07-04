package translation

import (
	"github.com/nicksnyder/go-i18n/v2/i18n"
	"github.com/stretchr/testify/mock"
)

type MockTranslator struct {
	mock.Mock
}

func (r *MockTranslator) GetLocalizer(lang string) *i18n.Localizer {
	return r.Called(lang).Get(0).(*i18n.Localizer)
}

func (r *MockTranslator) Lang(key string, args map[string]interface{}, lang *string) string {
	return r.Called(key, args, lang).String(0)
}
