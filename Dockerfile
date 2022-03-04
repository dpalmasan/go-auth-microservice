FROM golang:alpine as builder

RUN apk --update add --no-cache ca-certificates
RUN update-ca-certificates

# Build project
WORKDIR /go/src/github.com/dpalmasan/go-auth-microservice
COPY . .
#RUN dep ensure
RUN CGO_ENABLED=0 GOOS=linux go build -a -installsuffix cgo -o  bin/auth-service ./cmd/server

FROM alpine:latest

# RUN addgroup -S auth-service-acc && adduser -S -g auth-service-acc auth-service-acc
# USER auth-service-acc
WORKDIR /app/
COPY --from=builder /etc/ssl/certs/ca-certificates.crt /etc/ssl/certs/
COPY --from=builder /go/src/github.com/dpalmasan/go-auth-microservice/bin/auth-service .
COPY --from=builder /go/src/github.com/dpalmasan/go-auth-microservice/cert cert
CMD ["./auth-service"]