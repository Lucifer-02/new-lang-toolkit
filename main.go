package main

import (
	"os"
	"tool/engines"
)

func main() {
	source := os.Args[1]
	target := os.Args[2]
	text := os.Args[3]

	translation := engines.GoogleTranslate(text, source, target)
	// fmt.Println(translation)
	// engines.TTSConcurrent(translation, "vi")
	audio := engines.TTSConcurrent(translation, "vi")
	os.Stdout.Write(audio)

}
