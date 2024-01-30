FROM golang:latest as builder

WORKDIR /app
COPY . .
RUN go mod tidy
RUN go build

FROM alpine:latest  

RUN apk --no-cache add ca-certificates

WORKDIR /root/
COPY --from=builder /app/webhookdemo .
EXPOSE 443

CMD ["./webhookdemo"] 
