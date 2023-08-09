build:
	 go build -o bin/solarpal ./cmd/http/.

run: build
	bin/solarpal
