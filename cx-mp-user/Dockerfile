FROM golang:latest as builder

COPY cx-mp-user /app/cx-mp-user
COPY conf/user.toml /app/conf/user.toml
WORKDIR /app

FROM alpine:latest

RUN mkdir /app
WORKDIR /app

COPY --from=builder /app /app
RUN echo "Asia/Shanghai" >  /etc/timezone

EXPOSE 8000

CMD ./cx-mp-user -config ./conf/user.toml
