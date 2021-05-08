FROM golang:latest as builder

COPY cx-config-center /app/cx-config-center
COPY conf/ccs.toml /app/conf/ccs.toml
WORKDIR /app

FROM alpine:latest

RUN mkdir /app
WORKDIR /app

COPY --from=builder /app /app
RUN echo "Asia/Shanghai" >  /etc/timezone

EXPOSE 8000

CMD ./cx-config-center -config ./conf/ccs.toml