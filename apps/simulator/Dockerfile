FROM golang:1.26-alpine AS builder
WORKDIR /app
COPY go.work.sum ./
COPY shared/ ./shared/
COPY apps/simulator/ ./apps/simulator/
RUN printf 'go 1.26.1\n\nuse (\n\t./shared\n\t./apps/simulator\n)\n' > go.work
RUN cd apps/simulator && go build -o /bin/simulator ./cmd/api/

FROM alpine:3.19
RUN apk add --no-cache ca-certificates tzdata
COPY --from=builder /bin/simulator /simulator
EXPOSE 8090
CMD ["/simulator"]
