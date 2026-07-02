package main

import (
	"fmt"
	"os"

	"github.com/Lucifer-02/new-lang-toolkit/engines"
)

func main() {
	// The engines are deliberately fail-fast (they panic). Recover at this
	// boundary so the user sees a clean error instead of a stack trace.
	defer func() {
		if r := recover(); r != nil {
			fmt.Fprintf(os.Stderr, "tool: %v\n", r)
			os.Exit(1)
		}
	}()

	if len(os.Args) != 5 {
		//help
		os.Stdout.Write([]byte("Usage: tool [mode] [source] [target] [text]\n"))
		os.Stdout.Write([]byte("Modes:\n"))
		os.Stdout.Write([]byte("\ttts: Text to speech\n"))
		os.Stdout.Write([]byte("\ttrans: Translate text\n"))
		os.Stdout.Write([]byte("\ttrans+tts: Translate text and convert to speech"))
		return
	}

	mode := os.Args[1]
	source := os.Args[2]
	target := os.Args[3]
	text := os.Args[4]

	switch mode {
	case "tts":
		audio := engines.TTSConcurrent(text, target)
		os.Stdout.Write(audio)
	case "trans":
		translation := engines.GoogleTranslate(text, source, target)
		os.Stdout.Write([]byte(translation))
	case "trans+tts":
		translation := engines.GoogleTranslate(text, source, target)
		audio := engines.TTSConcurrent(translation, target)
		os.Stdout.Write(audio)
		if err := os.WriteFile("out.mp3", audio, 0o644); err != nil {
			panic(err)
		}
	default:
		panic("Invalid mode")
	}
}
