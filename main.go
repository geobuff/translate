package main

import (
	"bufio"
	"errors"
	"flag"
	"fmt"
	"os"
)

func main() {
	input := flag.String("input", "", "file to be translated")
	language := flag.String("language", "", "target language")
	flag.Parse()

	if *input == "" {
		panic(errors.New("please specify input file (--input=)"))
	}

	if *language == "" {
		panic(errors.New("please specify target language (--language=)"))
	}

	file, err := os.Open(*input)
	if err != nil {
		panic(err)
	}
	defer file.Close()

	service, err := NewTranslationService(*language)
	if err != nil {
		panic(err)
	}

	scanner := bufio.NewScanner(file)
	buf := make([]byte, 0, 64*1024)
	scanner.Buffer(buf, 1024*1024)
	for scanner.Scan() {
		text := scanner.Text()
		translated, err := service.translateInnerText(text)
		if err != nil {
			panic(err)
		}
		fmt.Println(translated)
	}
}
