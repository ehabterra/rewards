FROM golang:1.20-alpine AS builder
COPY . /app/
WORKDIR /app
RUN apk update && apk add make
RUN make build

FROM alpine
COPY --from=builder /app/bin /app
RUN apk add findutils
WORKDIR /app
EXPOSE 8000
ENTRYPOINT [ "./server" ]