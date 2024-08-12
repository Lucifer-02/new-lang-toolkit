package main

import (
	"os"
	"tool/engines"
)

func main() {

	if len(os.Args) != 5 {
		//help
		os.Stdout.Write([]byte("Usage: tool [mode] [source] [target] [text]\n"))
		os.Stdout.Write([]byte("Modes:\n"))
		os.Stdout.Write([]byte("tts: Text to speech\n"))
		os.Stdout.Write([]byte("trans: Translate text\n"))
		os.Stdout.Write([]byte("trans+tts: Translate text and convert to speech"))
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
	default:
		panic("Invalid mode")
	}
}
