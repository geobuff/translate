package main

import (
	"context"
	"errors"
	"testing"

	"cloud.google.com/go/translate"
	"golang.org/x/text/language"
)

func TestNewTranslationService(t *testing.T) {
	savedNewClient := newTranslateClient
	defer func() { newTranslateClient = savedNewClient }()

	t.Run("invalid language code", func(t *testing.T) {
		_, err := NewTranslationService(context.Background(), "blah")
		if err == nil {
			t.Error("expected error on invalid language code; got nil")
		}
	})

	t.Run("client creation failure", func(t *testing.T) {
		newTranslateClient = func(ctx context.Context) (*translate.Client, error) {
			return nil, errors.New("dial failed")
		}
		_, err := NewTranslationService(context.Background(), "es")
		if err == nil {
			t.Error("expected error on client creation failure; got nil")
		}
	})

	t.Run("valid language code", func(t *testing.T) {
		newTranslateClient = func(ctx context.Context) (*translate.Client, error) {
			return nil, nil
		}
		svc, err := NewTranslationService(context.Background(), "es")
		if err != nil {
			t.Fatalf("unexpected error: %v", err)
		}
		if svc == nil {
			t.Error("expected non-nil service")
		}
	})
}

func TestTranslateInnerText(t *testing.T) {
	savedTranslateText := translateText
	defer func() { translateText = savedTranslateText }()

	savedNewClient := newTranslateClient
	defer func() { newTranslateClient = savedNewClient }()
	newTranslateClient = func(ctx context.Context) (*translate.Client, error) {
		return nil, nil
	}

	tt := []struct {
		name          string
		input         string
		translateText func(_ *translate.Client, _ language.Tag, text string) (string, error)
		expectedText  string
		expectedErr   string
	}{
		{
			name:          "root tag",
			input:         "export default {",
			translateText: savedTranslateText,
			expectedText:  "export default {",
			expectedErr:   "",
		},
		{
			name:          "group tag",
			input:         "  validations: {",
			translateText: savedTranslateText,
			expectedText:  "  validations: {",
			expectedErr:   "",
		},
		{
			name:          "empty string",
			input:         "",
			translateText: savedTranslateText,
			expectedText:  "",
			expectedErr:   "",
		},
		{
			name:          "whitespace only",
			input:         "   ",
			translateText: savedTranslateText,
			expectedText:  "   ",
			expectedErr:   "",
		},
		{
			name:          "closing brace",
			input:         "};",
			translateText: savedTranslateText,
			expectedText:  "};",
			expectedErr:   "",
		},
		{
			name:          "unmatched opening quote",
			input:         `key: "no closing quote`,
			translateText: savedTranslateText,
			expectedText:  `key: "no closing quote`,
			expectedErr:   "",
		},
		{
			name:  "empty quoted value",
			input: `key: "",`,
			translateText: func(_ *translate.Client, _ language.Tag, text string) (string, error) {
				return "", nil
			},
			expectedText: `key: "",`,
			expectedErr:  "",
		},
		{
			name:  "inner text; error on translate text",
			input: `sizeRequired: "Please choose an available size.",`,
			translateText: func(_ *translate.Client, _ language.Tag, text string) (string, error) {
				return "", errors.New("testing")
			},
			expectedText: "",
			expectedErr:  "testing",
		},
		{
			name:  "inner text; success",
			input: `sizeRequired: "Please choose an available size.",`,
			translateText: func(_ *translate.Client, _ language.Tag, text string) (string, error) {
				return "TRANSLATED_TEXT", nil
			},
			expectedText: `sizeRequired: "TRANSLATED_TEXT",`,
			expectedErr:  "",
		},
	}

	service, err := NewTranslationService(context.Background(), "es")
	if err != nil {
		t.Fatalf("failed to create service: %v", err)
	}

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
