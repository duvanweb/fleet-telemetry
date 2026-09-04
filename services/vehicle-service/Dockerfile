FROM golang:1.26-alpine AS builder
WORKDIR /app
COPY go.work go.work.sum ./
COPY shared/ ./shared/
COPY services/vehicle-service/ ./services/vehicle-service/
RUN cd services/vehicle-service && go build -o /bin/vehicle-service ./cmd/api/

FROM alpine:3.19
RUN apk add --no-cache ca-certificates tzdata
COPY --from=builder /bin/vehicle-service /vehicle-service
EXPOSE 8081
CMD ["/vehicle-service"]
