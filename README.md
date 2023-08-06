CLI tool to translate our UI language file into other languages using the Google Translate API.

## Prerequisites

- Must be authenticated locally with a Google Cloud service account with translate API access. Follow the steps [here](https://cloud.google.com/translate/docs/setup) to get setup.
- Input language file must match expected format. See `example.ts` for reference.
- Language flag must match an ISO-639 code in the [supported list](https://cloud.google.com/translate/docs/languages).

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
