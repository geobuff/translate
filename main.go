package main

import (
	"bufio"
	"context"
	"fmt"
	"os"
	"strings"

	"cloud.google.com/go/translate"
	"golang.org/x/text/language"
)

func main() {
	if len(os.Args) < 2 {
		fmt.Println("missing file path and target language command line arg")
		return
	}

	path := os.Args[1]
	file, err := os.Open(path)
	if err != nil {
		fmt.Println(err)
		return
	}
	defer file.Close()

	targetLanguage := os.Args[2]

	scanner := bufio.NewScanner(file)
	buf := make([]byte, 0, 64*1024)
	scanner.Buffer(buf, 1024*1024)

	for scanner.Scan() {
		text := scanner.Text()
		if strings.Contains(text, `"`) {
			sub := text[(strings.Index(text, `"`) + 1):]
			value := sub[:strings.Index(sub, `"`)]
			translated, err := translateText(targetLanguage, value)
			if err != nil {
				panic(err)
			}

			fmt.Println(fmt.Sprintf(`%s"%s"%s`, text[:strings.Index(text, `"`)], translated, sub[strings.Index(sub, `"`)+1:]))
		} else {
			fmt.Println(text)
		}
	}
}

func translateText(targetLanguage, text string) (string, error) {
	lang, err := language.Parse(targetLanguage)
	if err != nil {
		return "", fmt.Errorf("language.Parse: %v", err)
	}

	ctx := context.Background()
	client, err := translate.NewClient(ctx)
	if err != nil {
		panic(err)
	}
	defer client.Close()

	resp, err := client.Translate(ctx, []string{text}, lang, nil)
	if err != nil {
		return "", fmt.Errorf("Translate: %v", err)
	}

	if len(resp) == 0 {
		return "", fmt.Errorf("Translate returned empty response to text: %s", text)
	}

	return resp[0].Text, nil
}
