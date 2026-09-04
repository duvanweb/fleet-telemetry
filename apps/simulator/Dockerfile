FROM golang:1.26-alpine AS builder
WORKDIR /app
COPY go.work go.work.sum ./
COPY shared/ ./shared/
COPY apps/simulator/ ./apps/simulator/
RUN cd apps/simulator && go build -o /bin/simulator ./cmd/api/

FROM alpine:3.19
RUN apk add --no-cache ca-certificates tzdata
COPY --from=builder /bin/simulator /simulator
EXPOSE 8090
CMD ["/simulator"]
