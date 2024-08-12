TEXT = "Building off of our last successful Causal Inference and Experimentation Summit, we held another week-long internal conference this year to learn from our stunning colleagues. We brought together speakers from across the business to learn about methodological developments and innovative applications."


run:
	# go run . "en" "vi" $(TEXT)
	go run . "en" "vi" $(TEXT) | mpv --speed=2.0 --quiet -

build-linux:
	GOOS=linux GOARCH=amd64 go build -o tool

# build windows binary on linux
build-win:
	GOOS=windows GOARCH=amd64 go build -o tool.exe
