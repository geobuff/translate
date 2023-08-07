package main

import (
	"errors"
	"testing"

	"golang.org/x/text/language"
)

func TestNewTranslationService(t *testing.T) {
	_, err := NewTranslationService("blah")
	if err == nil {
		t.Error("expected error on invalid language code; got nil")
	}
}

func TestTranslateInnerText(t *testing.T) {
	savedTranslateText := translateText
	defer func() {
		translateText = savedTranslateText
	}()

	tt := []struct {
		name          string
		input         string
		translateText func(lang language.Tag, text string) (string, error)
		expectedText  string
		expectedErr   string
	}{
		{
			name:          "root tag",
			input:         "export default {",
			translateText: translateText,
			expectedText:  "export default {",
			expectedErr:   "",
		},
		{
			name:          "group tag",
			input:         "  validations: {",
			translateText: translateText,
			expectedText:  "  validations: {",
			expectedErr:   "",
		},
		{
			name:          "inner text; error on translate text",
			input:         "sizeRequired: \"Please choose an available size.\",",
			translateText: func(lang language.Tag, text string) (string, error) { return "", errors.New("testing") },
			expectedText:  "",
			expectedErr:   "testing",
		},
		{
			name:          "inner text; success",
			translateText: func(lang language.Tag, text string) (string, error) { return "TRANSLATED_TEXT", nil },
			input:         "sizeRequired: \"Please choose an available size.\",",
			expectedText:  "sizeRequired: \"TRANSLATED_TEXT\",",
			expectedErr:   "",
		},
	}

	service, _ := NewTranslationService("es")
	for _, tc := range tt {
		t.Run(tc.name, func(t *testing.T) {
			translateText = tc.translateText
			text, err := service.translateInnerText(tc.input)

			if tc.expectedText != text {
				t.Errorf("expected %s; got %s", tc.expectedText, text)
			}

			if err != nil && err.Error() != tc.expectedErr {
				t.Errorf("expected %s; got %s", tc.expectedErr, err.Error())
			}
		})
	}
}
