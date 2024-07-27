FROM golang:1-alpine as builder

WORKDIR /app

COPY ./cmd /app/cmd
COPY ./pkg /app/pkg
COPY ./public /app/public
COPY ./templates /app/templates
COPY ./go.mod ./go.sum /app/

RUN ls -lah /app

RUN go build -o ./release/cli ./cmd/cli/main.go
RUN chmod +x ./release/cli

FROM alpine

WORKDIR /app

RUN apk add --no-cache --update bind-tools

COPY --from=builder /app/public /app/public
COPY --from=builder /app/templates /app/templates
COPY --from=builder /app/release/cli /app/cli

CMD [ "/app/cli" ]

