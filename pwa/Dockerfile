FROM golang:1.23.2-alpine3.20 AS build
WORKDIR /app
ADD go.mod go.sum ./
RUN go mod download
ADD main/* ./main/
ADD flush/* ./flush/
RUN GOARCH=wasm GOOS=js go build -o web/app.wasm ./main/main.go && \
GOARCH=amd64 GOOS=linux go build -o flush-log ./main/main.go
RUN printf 'http://localhost:6789' > /app/web/apiurl

FROM scratch
WORKDIR /app
COPY --from=build /app/web/apiurl /app/web/apiurl
COPY --from=build /app/web/app.wasm /app/web/app.wasm
COPY --from=build /app/flush-log /app/flush-log
CMD ["/app/flush-log"]
