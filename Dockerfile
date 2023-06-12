# Build the Go API
FROM golang:1.19.4 AS go_builder
ARG APP_NAME=todo-service
ARG APP_PATH=/app

WORKDIR ${APP_PATH}

ADD . ${APP_PATH}
RUN go mod download
RUN CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build -ldflags "-w" -a -o ${APP_PATH}/${APP_NAME} cmd/main.go

#CREATE PRODUCTION IMAGE
FROM alpine:3.16

ARG APP_NAME=todo-service
ARG APP_PATH=/app
ARG APP_PORT=80
ARG BUILD_DATE

LABEL build-date=$BUILD_DATE

WORKDIR ${APP_PATH}
COPY --from=go_builder /${APP_PATH}/${APP_NAME} /${APP_PATH}/${APP_NAME}
COPY --from=go_builder /${APP_PATH}/conf ${APP_PATH}/conf

RUN apk --no-cache add ca-certificates curl
RUN mkdir -p ${APP_PATH}/conf
RUN chmod +x ${APP_PATH}/${APP_NAME}
RUN ln -s ${APP_PATH}/${APP_NAME} /bin/${APP_NAME}

EXPOSE ${APP_PORT}

ENTRYPOINT [ "todo-service" ]