build:
	 go build -o bin/solarpal ./cmd/http/.

pancakes: build
	bin/solarpal
