FROM golang:1.23.2-alpine3.20 AS build
WORKDIR /app
ADD go.mod go.sum ./
RUN go mod download
ADD  main.go ./
RUN GOARCH=amd64 GOOS=linux go build && GOARCH=wasm GOOS=js go build -o web/app.wasm

FROM scratch
WORKDIR /app
COPY --from=build /app/web/app.wasm /app/web/app.wasm
COPY --from=build /app/flush-log /app/flush-log
CMD ["/app/flush-log"]
