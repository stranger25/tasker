FROM golang:latest AS builder

WORKDIR /app
COPY . .
RUN go mod tidy
RUN GO111MODULE="on" CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build -o app ./cmd
EXPOSE 9090
ENV THIRD_PARTY_URL=http://third-party.com
CMD ["./app"]