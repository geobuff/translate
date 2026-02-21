CLI tool to translate TypeScript i18n files into other languages using the Google Cloud Translation API.

## Prerequisites

- **Go 1.23+**
- A Google Cloud project with the [Cloud Translation API](https://cloud.google.com/translate/docs/setup) enabled
- [Application Default Credentials](https://cloud.google.com/docs/authentication/provide-credentials-adc) configured locally
- Input file must match the expected TypeScript i18n format (see `example.ts`)
- Language flag must be a valid ISO-639 code from the [supported list](https://cloud.google.com/translate/docs/languages)

## Architecture

| File | Role |
|---|---|
| `main.go` | CLI entry point; parses flags, opens the input file, calls `run()` |
| `translation_service.go` | `TranslationService` struct; wraps the GCP client and per-line translation logic |

`run(service, r, w)` is extracted from `main` for testability: it reads lines from any `io.Reader` and writes results to any `io.Writer`.

## Implementation details

- **Line-by-line processing** — the file is read with a `bufio.Scanner`; each line is processed independently.
- **Single quoted value per line** — `translateInnerText` assumes each line contains at most one double-quoted value, matching the expected TypeScript i18n file format. Lines without a quoted value (braces, blank lines) pass through unchanged.
- **Client reuse** — a single `*translate.Client` is created once in `NewTranslationService` and reused for every translation call; call `Close()` when done.
- **Mock-friendly design** — `translateText` and `newTranslateClient` are package-level variables so tests can substitute fakes without real GCP credentials.

## Install

```
git clone https://github.com/geobuff/translate
go build
```

## Run

```
./translate --input=example.ts --language=es
```

Output to file:

```
./translate --input=example.ts --language=es > output.ts
```

## Input format

The input file should be a TypeScript module exporting a nested object of string values, e.g.:

```typescript
export default {
  validations: {
    sizeRequired: "Please choose an available size.",
  },
};
```

## Testing

No real GCP credentials are needed to run the tests:

```
go test ./...
```
