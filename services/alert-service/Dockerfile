FROM golang:1.26-alpine AS builder
WORKDIR /app
COPY go.work go.work.sum ./
COPY shared/ ./shared/
COPY services/alert-service/ ./services/alert-service/
RUN cd services/alert-service && go build -o /bin/alert-service ./cmd/api/

FROM alpine:3.19
RUN apk add --no-cache ca-certificates tzdata
COPY --from=builder /bin/alert-service /alert-service
EXPOSE 8083
CMD ["/alert-service"]
