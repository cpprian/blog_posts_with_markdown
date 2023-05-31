FROM golang:1.20.4-alpine3.18 AS builder

RUN adduser -D -g '' appuser

WORKDIR /opt/app

COPY . . ./

RUN go mod download && \
    go mod verify

COPY . .

RUN CGO_ENABLED=0 GOOS=linux go build -ldflags="-w -s" -a -installsuffix cgo -o /go/bin/blogposts-users ./cmd/app


FROM alpine:3.17.3

LABEL langugage="golang"

COPY --from=builder /etc/passwd /etc/passwd

COPY --from=builder --chown=appuser:1000 /go/bin/blogposts-users /blogposts-users

USER appuser

ENTRYPOINT ["./blogposts-users"]