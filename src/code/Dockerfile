FROM gcr.io/kaniko-project/executor:debug AS kaniko

FROM alpine:3.14.2

RUN apk --no-cache add ca-certificates

RUN mkdir /workspace
RUN mkdir -p /kaniko/.docker/
RUN touch /kaniko/.docker/config.json
COPY --from=kaniko /kaniko/*  /kaniko/
COPY --from=kaniko /busybox/* /busybox/

ENV HOME /root
ENV USER root
ENV PATH $PATH:/usr/local/bin:/kaniko:/busybox
ENV SSL_CERT_DIR=/kaniko/ssl/certs
ENV DOCKER_CONFIG /kaniko/.docker/
ENV DOCKER_CREDENTIAL_GCR_CONFIG /kaniko/.config/gcloud/docker_credential_gcr_config.json

# 如果在镜像中直接完成 dockerhub 或者 acr 的 auth, 其中 auth 的值为 $user:$pwd base64 后的值
# {
#   "auths": {
#     "registry.cn-hangzhou.aliyuncs.com": {
#       "auth": "xxx"
#     }
#   }
# }
# COPY config.json /kaniko/.docker/config.json

RUN apk --update add \
  bash \
  curl \
  git \
  jq \
  python3 \
  py3-pip

# Create app directory
RUN mkdir /kaniko/App
WORKDIR /kaniko/App

RUN pip install flask -i https://mirrors.aliyun.com/pypi/simple/

# Bundle app source
COPY server.py server.py

EXPOSE 9000
ENV FLASK_APP=server
ENTRYPOINT ["python3", "-u",  "server.py"]