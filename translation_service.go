package main

import (
	"context"
	"fmt"
	"strings"

	"cloud.google.com/go/translate"
	"golang.org/x/text/language"
)

// newTranslateClient is a package-level variable so tests can mock out client
// creation without real GCP credentials.
var newTranslateClient = func(ctx context.Context) (*translate.Client, error) {
	return translate.NewClient(ctx)
}

// translateText is a package-level variable so tests can mock out translation
// calls. The client is passed in to enable reuse across calls.
var translateText = func(client *translate.Client, lang language.Tag, text string) (string, error) {
	resp, err := client.Translate(context.Background(), []string{text}, lang, nil)
	if err != nil {
		return "", fmt.Errorf("translate: %w", err)
	}

	if len(resp) == 0 {
		return "", fmt.Errorf("translate returned empty response to text: %s", text)
	}

	return resp[0].Text, nil
}

type TranslationService struct {
	client *translate.Client
	lang   language.Tag
}

func NewTranslationService(ctx context.Context, targetLanguage string) (*TranslationService, error) {
	lang, err := language.Parse(targetLanguage)
	if err != nil {
		return nil, err
	}

	client, err := newTranslateClient(ctx)
	if err != nil {
		return nil, fmt.Errorf("create translate client: %w", err)
	}

	return &TranslationService{
		client: client,
		lang:   lang,
	}, nil
}

// Close releases the underlying GCP client connection. Safe to call on a nil receiver.
func (s *TranslationService) Close() error {
	if s == nil || s.client == nil {
		return nil
	}
	return s.client.Close()
}

// translateInnerText extracts the quoted value from a single line, translates
// it, and returns the line with the translated value substituted in. Lines
// without a quoted value are returned unchanged.
//
// The implementation assumes at most one double-quoted value per line, matching
// the expected TypeScript i18n file format.
func (s *TranslationService) translateInnerText(text string) (string, error) {
	if !strings.Contains(text, `"`) {
		return text, nil
	}

	openIdx := strings.Index(text, `"`)
	sub := text[openIdx+1:]
	closeIdx := strings.Index(sub, `"`)
	if closeIdx < 0 {
		// Unmatched opening quote — return line unchanged to avoid a panic on
		// malformed input.
		return text, nil
	}

	value := sub[:closeIdx]
	translated, err := translateText(s.client, s.lang, value)
	if err != nil {
		return "", err
	}

	result := fmt.Sprintf(`%s"%s"%s`, text[:openIdx], translated, sub[closeIdx+1:])
	return result, nil
}
