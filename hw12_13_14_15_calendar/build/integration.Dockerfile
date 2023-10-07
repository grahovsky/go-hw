FROM golang:1.21

COPY . /tests
WORKDIR /tests

ENV CGO_ENABLED=0
RUN go mod tidy

ENTRYPOINT [ "go", "test", "-tags=integration", "./tests/...", "-v"]