# CLAUDE.md

This file provides guidance to Claude Code (claude.ai/code) when working with code in this repository.

## Overview

A single-binary Go CLI (`tool`) that wraps Google's public (unofficial) translate and TTS endpoints. It translates text, synthesizes speech, or does both, writing results to stdout.

## Commands

```bash
# Run against the sample TEXT in the Makefile, piping audio to mpv at 2x
make run

# Cross-compile
make build-linux   # -> ./tool     (GOOS=linux GOARCH=amd64)
make build-win     # -> ./tool.exe (GOOS=windows GOARCH=amd64)

# Run directly: tool [mode] [source] [target] [text]
go run . trans en vi "hello world"
go run . tts en vi "hello world" > out.mp3
go run . trans+tts en vi "hello world"   # also writes out.mp3

# Tests (live in ./engines)
go test ./engines/
go test -run TestSplitText1 ./engines/   # single test
```

Modes: `tts` (speech only, uses target lang), `trans` (text translation), `trans+tts` (translate then speak; also writes `out.mp3`). `main.go` requires exactly 4 positional args or it prints usage.

## Architecture

- `main.go` — arg parsing and mode dispatch only. All logic lives in `package engines`.
- `engines/` — one file per concern:
  - `google_translate.go` — `GoogleTranslate`: builds the `translate_a/single` URL, then parses the deeply-nested untyped JSON array response in `extractTranslation` (walks `jsonData[0]` and takes element `[0]` of each sub-array).
  - `google_tts.go` — `TTS` (sequential) and `TTSConcurrent` (preferred; used by `main.go`). Both chunk text via `SplitText` because the TTS endpoint caps input length, then concatenate the returned MP3 bytes into one stream.
  - `text_utils.go` — `SplitText(text, limit)`: splits long text into <=`limit` chunks, preferring to break at sentence ends (`.`/`?`), then clause breaks (`,`/`;`), then spaces. Invariant enforced at the end: the summed chunk lengths must equal the input length.
  - `request.go` — HTTP layer. `ApiRequest` (single GET) and `ApiRequests` (concurrent GETs via goroutines + channel, results re-sorted by original index to preserve order — important because TTS chunks must stay in sequence).
  - `invariant.go` — `assert(cond, msg)`: panics on failure. Used throughout for precondition/invariant checks.

## Conventions

- Error handling is deliberately fail-fast: functions `panic` (directly or via `assert`) rather than returning errors up the stack. Preserve this style when editing existing engine code.
- The `textLimit` for TTS chunking is `200` (`google_tts.go`).
- Depends on Google's undocumented endpoints (`translate.googleapis.com`, `translate.google.com/translate_tts`) with hardcoded client params (`gtx`, `tw-ob`); response shapes can change without notice, so `extractTranslation` is the fragile spot.
