FROM node:alpine AS web-builder
WORKDIR /build

COPY web .
RUN npm i && npm run build

FROM golang:alpine AS builder

ARG CLI_VERSION
WORKDIR /build

COPY . .
COPY --from=web-builder /build/build web/build

RUN go build -ldflags "-s -w -X localapps/constants.Version=$CLI_VERSION -X localapps/constants.runningInContainer=true"

FROM alpine
COPY --from=builder /build/localapps /bin

RUN mkdir /root/.config
ENTRYPOINT ["localapps", "up"]
