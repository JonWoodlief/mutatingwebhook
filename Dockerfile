FROM golang:latest as builder

WORKDIR /app
COPY . .
RUN go mod tidy
RUN GOOS=linux GOARCH=amd64 go build

FROM alpine:latest  

RUN apk --no-cache add ca-certificates

WORKDIR /root/
COPY --from=builder /app/webhookdemo .
EXPOSE 443

CMD ["./webhookdemo"] 
