FROM golang:1.23-alpine as build

WORKDIR /app

RUN go install github.com/a-h/templ/cmd/templ@latest &&\
    apk add make npm nodejs &&\
    npm install -g tailwindcss 

COPY go.mod go.sum ./
RUN go mod download

COPY . .
RUN make build

FROM alpine:latest as run

WORKDIR /app
COPY --from=build /app/bin/www ./
COPY --from=build /app/static ./static

ENTRYPOINT ["/app/www"]
