package engines

func assert(condition bool, message string) {
	if !condition {
		panic(message)
	}
}
