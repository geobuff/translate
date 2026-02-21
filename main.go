package main

import (
	"bufio"
	"context"
	"flag"
	"fmt"
	"io"
	"os"
)

// run reads lines from r, translates each line via service, and writes the
// results to w. It is extracted from main so it can be tested with in-memory
// readers and writers.
func run(service *TranslationService, r io.Reader, w io.Writer) error {
	scanner := bufio.NewScanner(r)
	buf := make([]byte, 0, 64*1024)
	scanner.Buffer(buf, 1024*1024)
	for scanner.Scan() {
		text := scanner.Text()
		translated, err := service.translateInnerText(text)
		if err != nil {
			return err
		}
		fmt.Fprintln(w, translated)
	}
	return scanner.Err()
}

func main() {
	input := flag.String("input", "", "file to be translated")
	language := flag.String("language", "", "target language")
	flag.Parse()

	if *input == "" {
		fmt.Fprintf(os.Stderr, "please specify input file (--input=)\n")
		os.Exit(1)
	}

	if *language == "" {
		fmt.Fprintf(os.Stderr, "please specify target language (--language=)\n")
		os.Exit(1)
	}

	ctx := context.Background()
	service, err := NewTranslationService(ctx, *language)
	if err != nil {
		fmt.Fprintf(os.Stderr, "failed to create translation service: %v\n", err)
		os.Exit(1)
	}
	defer service.Close()

	file, err := os.Open(*input)
	if err != nil {
		fmt.Fprintf(os.Stderr, "failed to open input file: %v\n", err)
		os.Exit(1)
	}
	defer file.Close()

	if err := run(service, file, os.Stdout); err != nil {
		fmt.Fprintf(os.Stderr, "translation failed: %v\n", err)
		os.Exit(1)
	}
}
