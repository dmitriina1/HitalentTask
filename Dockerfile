FROM golang:1.25 AS builder
WORKDIR /app
COPY . .
RUN go mod tidy && go build -o app main.go

FROM golang:1.25
WORKDIR /app
COPY --from=builder /app/app /app/app
COPY ./migrations /app/migrations
COPY ./openapi.json /app/openapi.json
COPY .env /app/.env
EXPOSE 8080
CMD ["/app/app"]