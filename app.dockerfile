FROM golang:alpine AS build
WORKDIR /inventory
COPY go.mod go.sum ./
RUN go mod download
COPY .env .
COPY cmd cmd
COPY internal internal
COPY pkg pkg
RUN go build -o /cmd/app ./cmd

FROM alpine:3.20
WORKDIR /usr/bin
COPY --from=build /cmd/app .
CMD [ "app" ]