# Собираем в гошке
FROM golang:1.21 as build

ENV BIN_FILE /opt/calendar/sender
ENV CODE_DIR /go/src/

WORKDIR ${CODE_DIR}

# Кэшируем слои с модулями
COPY go.mod .
COPY go.sum .
RUN go mod download

COPY . ${CODE_DIR}

# Собираем статический бинарник Go (без зависимостей на Си API),
# иначе он не будет работать в alpine образе.
ARG LDFLAGS
RUN CGO_ENABLED=0 go build \
        -ldflags "$LDFLAGS" \
        -o ${BIN_FILE} cmd/calendar_sender/*

# На выходе тонкий образ
FROM alpine:3.9

LABEL ORGANIZATION="OTUS Online Education"
LABEL SERVICE="calendar"
LABEL MAINTAINERS="student@otus.ru"

ENV BIN_FILE "/opt/calendar/sender"
COPY --from=build ${BIN_FILE} ${BIN_FILE}

ENV CONFIG_FILE /etc/calendar/sender.yaml
COPY ./configs/sender.yaml ${CONFIG_FILE}

CMD ${BIN_FILE} --config ${CONFIG_FILE}
