package main

import (
	"bufio"
	"errors"
	"flag"
	"fmt"
	"os"
	"strings"
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

	service := NewService()
	scanner := bufio.NewScanner(file)
	buf := make([]byte, 0, 64*1024)
	scanner.Buffer(buf, 1024*1024)
	for scanner.Scan() {
		text := scanner.Text()
		if strings.Contains(text, `"`) {
			sub := text[(strings.Index(text, `"`) + 1):]
			value := sub[:strings.Index(sub, `"`)]
			translated, err := service.translateText(*language, value)
			if err != nil {
				panic(err)
			}

			fmt.Println(fmt.Sprintf(`%s"%s"%s`, text[:strings.Index(text, `"`)], translated, sub[strings.Index(sub, `"`)+1:]))
		} else {
			fmt.Println(text)
		}
	}
}
