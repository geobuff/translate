package main

import (
	"context"
	"fmt"
	"strings"

	"cloud.google.com/go/translate"
	"golang.org/x/text/language"
)

type IService interface {
	translateInnerText(text string) (string, error)
}

type TranslationService struct {
	lang language.Tag
}

func NewTranslationService(targetLanguage string) (*TranslationService, error) {
	lang, err := language.Parse(targetLanguage)
	if err != nil {
		return nil, err
	}

	return &TranslationService{
		lang,
	}, nil
}

func (s *TranslationService) translateInnerText(text string) (string, error) {
	if !strings.Contains(text, `"`) {
		return text, nil
	}

	sub := text[(strings.Index(text, `"`) + 1):]
	value := sub[:strings.Index(sub, `"`)]
	translated, err := translateText(s.lang, value)
	if err != nil {
		return "", err
	}

	result := fmt.Sprintf(`%s"%s"%s`, text[:strings.Index(text, `"`)], translated, sub[strings.Index(sub, `"`)+1:])
	return result, err
}

var translateText = func(lang language.Tag, text string) (string, error) {
	ctx := context.Background()
	client, err := translate.NewClient(ctx)
	if err != nil {
		return "", err
	}
	defer client.Close()

	resp, err := client.Translate(ctx, []string{text}, lang, nil)
	if err != nil {
		return "", err
	}

	if len(resp) == 0 {
		return "", fmt.Errorf("translate returned empty response to text: %s", text)
	}

	return resp[0].Text, nil
}
