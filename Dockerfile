# COMPILE
FROM golang:alpine AS builder

ARG VERSION=dev

COPY . /app

RUN cd /app \
    && go build -ldflags "-s -w -X github.com/axllent/adguard-home-bg/cmd.Version=${VERSION}" -o adguard-home-bg . \ 
    && apk add --no-cache upx && upx -9 adguard-home-bg

# BUILD IMAGE
FROM alpine


EXPOSE 8080

RUN apk add --no-cache tzdata

COPY --from=builder /app/adguard-home-bg /usr/bin/adguard-home-bg

USER nobody

CMD ["adguard-home-bg", "server"]
