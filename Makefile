build:
	env GOOS=linux CGO_ENABLED=0 go build -ldflags="-s -w" -o bin/solarpal ./cmd/http/.

deploy_prod: build
	serverless deploy --stage prod 

run: build
	bin/solarpal
