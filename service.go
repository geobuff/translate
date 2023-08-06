package main

import (
	"context"
	"fmt"

	"cloud.google.com/go/translate"
	"golang.org/x/text/language"
)

type IService interface {
	translateText(targetLanguage, text string) (string, error)
}

type Service struct{}

func NewService() *Service {
	return &Service{}
}

func (s *Service) translateText(targetLanguage, text string) (string, error) {
	lang, err := language.Parse(targetLanguage)
	if err != nil {
		return "", err
	}

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
