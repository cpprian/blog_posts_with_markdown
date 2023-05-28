FROM golang:1.20.4-alpine3.18 AS builder

RUN adduser -D -g '' appuser

WORKDIR /opt/web

COPY go.mod go.sum ./

RUN go mod download && \
    go mod verify

COPY . .

RUN CGO_ENABLED=0 GOOS=linux go build -ldflags="-w -s" -a -installsuffix cgo -o /go/bin/blogposts-website ./cmd/web


FROM alpine:3.17.3
LABEL language="golang"
LABEL org.opencontainers.image.source https://github.com/cpprian/blog_posts_with_markdown

COPY --from=builder /etc/passwd /etc/passwd

WORKDIR /opt/web

COPY --from=builder --chown=appuser:1000 /go/bin/blogposts-website blogposts-website

COPY ./ui ui

ENTRYPOINT ["./blogposts-website"]