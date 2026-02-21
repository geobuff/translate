package main

import (
	"bytes"
	"context"
	"errors"
	"io"
	"strings"
	"testing"

	"cloud.google.com/go/translate"
	"golang.org/x/text/language"
)

// errReader returns an error on the first Read call to simulate a mid-stream
// I/O failure and exercise the scanner.Err() path in run().
type errReader struct{}

func (errReader) Read(p []byte) (int, error) {
	return 0, errors.New("read error")
}

func newTestService(t *testing.T) *TranslationService {
	t.Helper()
	savedNewClient := newTranslateClient
	t.Cleanup(func() { newTranslateClient = savedNewClient })
	newTranslateClient = func(ctx context.Context) (*translate.Client, error) {
		return nil, nil
	}
	svc, err := NewTranslationService(context.Background(), "es")
	if err != nil {
		t.Fatalf("newTestService: %v", err)
	}
	return svc
}

func TestRun_EmptyInput(t *testing.T) {
	svc := newTestService(t)
	var buf bytes.Buffer
	if err := run(svc, strings.NewReader(""), &buf); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if buf.Len() != 0 {
		t.Errorf("expected no output; got %q", buf.String())
	}
}

func TestRun_PassThrough(t *testing.T) {
	savedTranslateText := translateText
	defer func() { translateText = savedTranslateText }()
	translateText = func(_ *translate.Client, _ language.Tag, text string) (string, error) {
		t.Fatal("translateText should not be called for structural lines")
		return "", nil
	}

	svc := newTestService(t)
	input := "export default {\n  validations: {\n};\n"
	var buf bytes.Buffer
	if err := run(svc, strings.NewReader(input), &buf); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	want := "export default {\n  validations: {\n};\n"
	if buf.String() != want {
		t.Errorf("expected %q; got %q", want, buf.String())
	}
}

func TestRun_Translated(t *testing.T) {
	savedTranslateText := translateText
	defer func() { translateText = savedTranslateText }()
	translateText = func(_ *translate.Client, _ language.Tag, text string) (string, error) {
		return "HOLA", nil
	}

	svc := newTestService(t)
	input := "greeting: \"Hello\",\n"
	var buf bytes.Buffer
	if err := run(svc, strings.NewReader(input), &buf); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	want := "greeting: \"HOLA\",\n"
	if buf.String() != want {
		t.Errorf("expected %q; got %q", want, buf.String())
	}
}

func TestRun_ScannerError(t *testing.T) {
	svc := newTestService(t)
	err := run(svc, errReader{}, &bytes.Buffer{})
	if err == nil {
		t.Fatal("expected error from scanner; got nil")
	}
}

func TestRun_TranslateInnerTextError(t *testing.T) {
	savedTranslateText := translateText
	defer func() { translateText = savedTranslateText }()
	translateText = func(_ *translate.Client, _ language.Tag, text string) (string, error) {
		return "", errors.New("api error")
	}

	svc := newTestService(t)
	input := "key: \"value\",\n"
	err := run(svc, strings.NewReader(input), io.Discard)
	if err == nil {
		t.Fatal("expected error from translateInnerText; got nil")
	}
	if err.Error() != "api error" {
		t.Errorf("expected 'api error'; got %q", err.Error())
	}
}
