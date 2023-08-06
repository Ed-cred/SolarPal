FROM golang:1.20.5

WORKDIR D:\SolarPal

COPY go.mod go.sum ./
RUN go mod download && go mod verify

COPY . .
RUN go build -o bin/solarpal ./cmd/http/.


CMD ["app"]
