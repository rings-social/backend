FROM golang:1.20-alpine3.18 AS builder
RUN apk add --no-cache make
WORKDIR /app
COPY ./ /app
RUN make build

FROM alpine:3.18
COPY --from=builder /app/build/rings-backend /usr/bin/rings-backend
ENTRYPOINT ["/usr/bin/rings-backend"]