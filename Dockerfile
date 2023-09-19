FROM golang:1.17 AS builder
COPY . /src
WORKDIR /src
RUN GOPROXY=https://goproxy.cn make build

FROM debian:stable-slim
RUN apt-get update && apt-get install -y --no-install-recommends \
		ca-certificates  \
        netbase \
        && rm -rf /var/lib/apt/lists/ \
        && apt-get autoremove -y && apt-get autoclean -y
COPY --from=builder /src/bin /app
WORKDIR /app

EXPOSE 8080
EXPOSE 8081
EXPOSE 9090
VOLUME /data/conf

CMD ["./hgms-layout", "-conf", "/data/conf"]