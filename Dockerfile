FROM golang:1.20.4 as builder
ARG BIN_LABELS
ENV ENV_BIN_LABELS=${BIN_LABELS}
ENV GOPROXY https://goproxy.cn/
WORKDIR /data
COPY . /data
RUN make build-lux BIN_LABELS=${ENV_BIN_LABELS}

FROM docker.osisbim.com/deploy/alpine_base:latest
ARG BIN_LABELS
ENV ENV_BIN_LABELS=${BIN_LABELS}
COPY --from=builder /data/${BIN_LABELS} /data/${BIN_LABELS}
COPY --from=builder /usr/share/zoneinfo/Asia/Shanghai /etc/localtime
EXPOSE 8000
WORKDIR /data
CMD [ "sh", "-c","/data/${ENV_BIN_LABELS}" ]