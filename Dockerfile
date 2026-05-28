FROM golang:alpine AS builder

LABEL stage=gobuilder

ENV CGO_ENABLED 0


RUN apk update --no-cache && apk add --no-cache tzdata

WORKDIR /build

ADD go.mod .
ADD go.sum .
RUN go mod download
COPY . .

RUN go build -ldflags="-s -w" -o /app/deeplx-translategemma .


FROM alpine

COPY --from=builder /etc/ssl/certs/ca-certificates.crt /etc/ssl/certs/ca-certificates.crt
COPY --from=builder /usr/share/zoneinfo/Asia/Shanghai /usr/share/zoneinfo/Asia/Shanghai
ENV TZ Asia/Shanghai

WORKDIR /app
COPY --from=builder /app/deeplx-translategemma /app/deeplx-translategemma
COPY etc/deeplx-api.yaml.example /app/deeplx-api.yaml.example
COPY docker-entrypoint.sh /app/docker-entrypoint.sh
RUN chmod +x /app/docker-entrypoint.sh

ENTRYPOINT ["/app/docker-entrypoint.sh"]
CMD ["./deeplx-translategemma", "-f", "/app/deeplx-api.yaml"]
