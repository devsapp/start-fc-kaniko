# start-fc-kaniko 帮助文档

<p align="center" class="flex justify-center">
    <a href="https://www.serverless-devs.com" class="ml-1">
    <img src="http://editor.devsapp.cn/icon?package=start-fc-kaniko&type=packageType">
  </a>
  <a href="http://www.devsapp.cn/details.html?name=start-fc-kaniko" class="ml-1">
    <img src="http://editor.devsapp.cn/icon?package=start-fc-kaniko&type=packageVersion">
  </a>
  <a href="http://www.devsapp.cn/details.html?name=start-fc-kaniko" class="ml-1">
    <img src="http://editor.devsapp.cn/icon?package=start-fc-kaniko&type=packageDownload">
  </a>
</p>

<description>

>  ***快速部署一个Dockfile编译和push镜像的应用到阿里云函数计算***

</description>

<table>
</table>

<codepre id="codepre">
</codepre>

<deploy>

## 部署 & 体验

<appcenter>

-  :fire:  通过 [Serverless 应用中心](https://fcnext.console.aliyun.com/applications/create?template=start-fc-kaniko) ，
[![Deploy with Severless Devs](https://img.alicdn.com/imgextra/i1/O1CN01w5RFbX1v45s8TIXPz_!!6000000006118-55-tps-95-28.svg)](https://fcnext.console.aliyun.com/applications/create?template=start-fc-kaniko)  该应用。 

</appcenter>

- 通过 [Serverless Devs Cli](https://www.serverless-devs.com/serverless-devs/install) 进行部署：
    - [安装 Serverless Devs Cli 开发者工具](https://www.serverless-devs.com/serverless-devs/install) ，并进行[授权信息配置](https://www.serverless-devs.com/fc/config) ；
    - 初始化项目：`s init start-fc-kaniko -d start-fc-kaniko`   
    - 进入项目，并进行项目部署：`cd start-fc-kaniko && s deploy -y`

</deploy>

<appdetail id="flushContent">

# 应用详情

项目部署完成, 直接使用生成的域名发起 HTTP 调用就可以实现对指定 git 仓库的 Dockfile 工程实现镜像的 build 和 push

```bash
$ curl -v  http://builder.fc-kaniko.123456789.cn-hangzhou.fc.devsapp.net  -d '{"url":"https://github.com/rsonghuster/TestKitBackend", "registry":"registry.cn-hangzhou.aliyuncs.com", "usr":"my-usr", "pwd":"my-pwd", "image":"registry.cn-hangzhou.aliyuncs.com/rsong/test:v1"}'
```

其中：

- **url** : git url,  如果是私有 github, 可以把 token 一起传过来，public 示例值为 `https://github.com/rsonghuster/TestKitBackend`, private 示例值为 `https://oauth2:access_token@github.com/username/xxx.git`

- **registry** : 镜像仓库，比如 dockerhub 或者 阿里云容器镜像服务 ACR，默认值为 dockerhub 的 `https://index.docker.io/v1/`， ACR 的示例值：`registry.cn-hangzhou.aliyuncs.com`

- **usr** : 对应镜像仓库的账户名

- **pwd** : 对应镜像仓库的账户密码

- **image** : 完整的镜像名字

## 二次开发

如果您想进行二次开发， 对 code 目录进行开发， 编译新的镜像作为函数 customContainerConfig 的 image 参数即可

</appdetail>

<devgroup>

## 开发者社区

您如果有关于错误的反馈或者未来的期待，您可以在 [Serverless Devs repo Issues](https://github.com/serverless-devs/serverless-devs/issues) 中进行反馈和交流。如果您想要加入我们的讨论组或者了解 FC 组件的最新动态，您可以通过以下渠道进行：

<p align="center">

| <img src="https://serverless-article-picture.oss-cn-hangzhou.aliyuncs.com/1635407298906_20211028074819117230.png" width="130px" > | <img src="https://serverless-article-picture.oss-cn-hangzhou.aliyuncs.com/1635407044136_20211028074404326599.png" width="130px" > | <img src="https://serverless-article-picture.oss-cn-hangzhou.aliyuncs.com/1635407252200_20211028074732517533.png" width="130px" > |
|--- | --- | --- |
| <center>微信公众号：`serverless`</center> | <center>微信小助手：`xiaojiangwh`</center> | <center>钉钉交流群：`33947367`</center> | 

</p>

</devgroup>