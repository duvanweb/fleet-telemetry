FROM golang:1.26-alpine AS builder
WORKDIR /app
COPY go.work go.work.sum ./
COPY shared/ ./shared/
COPY services/telemetry-service/ ./services/telemetry-service/
RUN cd services/telemetry-service && go build -o /bin/telemetry-service ./cmd/api/

FROM alpine:3.19
RUN apk add --no-cache ca-certificates tzdata
COPY --from=builder /bin/telemetry-service /telemetry-service
EXPOSE 8082
CMD ["/telemetry-service"]
