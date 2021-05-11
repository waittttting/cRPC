FROM golang:latest as builder

COPY cx-mp-gateway /app/cx-mp-gateway
COPY conf/gateway.toml /app/conf/gateway.toml
WORKDIR /app

FROM alpine:latest

RUN mkdir /app
WORKDIR /app

COPY --from=builder /app /app
RUN echo "Asia/Shanghai" >  /etc/timezone

EXPOSE 8000

CMD ./cx-mp-gateway -config ./conf/gateway.toml
