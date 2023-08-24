# Building stage
FROM golang:1.21-alpine as builder
WORKDIR /app
COPY . .
RUN CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build -a -o app cmd/main.go

#Run
FROM alpine:latest  
COPY --from=builder /app/app /app/
CMD ["/app/app", "/data/input.txt"]