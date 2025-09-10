FROM golang:1.25.1

WORKDIR /usr/
RUN go install github.com/air-verse/air@v1.62.0


# pre-copy/cache go.mod for pre-downloading dependencies and only redownloading them in subsequent builds if they change
COPY go.mod go.sum ./
RUN go mod download && go mod verify

COPY . .

CMD ["air"]
