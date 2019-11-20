FROM golang:alpine AS builder
RUN apk add git gcc g++

ENV GOOS=linux
ENV GOARCH=amd64

WORKDIR /app
COPY . .
RUN /app/build.sh
RUN chmod -R 777 /app

FROM alpine:latest AS production
WORKDIR /app
COPY --from=builder /app/plugins /app/plugins
COPY --from=builder /app/index.html /app/index.html
COPY --from=builder /app/favicon.ico /app/favicon.ico
COPY --from=builder /app/dasher /app/dasher

CMD ["/app/dasher"]