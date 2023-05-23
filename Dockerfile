FROM golang:1.20.4-alpine3.18 AS builder
RUN mkdir /build
ADD . /build
WORKDIR /build
RUN go get sorare-mu
RUN go build ./cmd/server/server.go

FROM alpine
RUN adduser -S -D -H -h /app appuser
USER appuser
COPY --from=builder /build/server /app/
COPY ext/ /app/ext
COPY static/ /app/static
COPY web/templates/ /app/web/templates
WORKDIR /app
CMD ["./server"]
