FROM golang:1.17
RUN mkdir /code
WORKDIR /code
COPY .  .
ENV GO111MODULE=on
RUN GOARCH=amd64 GOOS=linux CGO_ENABLED=0 go build -o fc-kaniko-builder

FROM gcr.io/kaniko-project/executor:debug
COPY --from=0 /code/fc-kaniko-builder /kaniko/fc-kaniko-builder
RUN mkdir -p /kaniko/workspace
ENTRYPOINT ["/kaniko/fc-kaniko-builder"]

# 如果在镜像中直接完成 dockerhub 或者 acr 的 auth, 其中 auth 的值为 $user:$pwd base64 后的值
# {
#   "auths": {
#     "registry.cn-hangzhou.aliyuncs.com": {
#       "auth": "xxx"
#     }
#   }
# }
# COPY config.json /kaniko/.docker/config.json