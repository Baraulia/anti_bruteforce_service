# Собираем в гошке
FROM golang:1.21.6 as build

ENV BIN_FILE /opt/ab-service
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
        -o ${BIN_FILE} cmd/*

# На выходе тонкий образ
FROM alpine:3.9

LABEL SERVICE="anti brutforce service"
LABEL MAINTAINERS="baraulia@yandex.ru"

ENV BIN_FILE "/opt/ab-service"
COPY --from=build ${BIN_FILE} ${BIN_FILE}

ENV CONFIG_FILE /etc/ab-service/config.yaml

COPY ./configs/config.yaml ${CONFIG_FILE}
EXPOSE 8085

CMD ${BIN_FILE} -config ${CONFIG_FILE}
