
> 注：当前项目为 Serverless Devs 应用，由于应用中会存在需要初始化才可运行的变量（例如应用部署地区、函数名等等），所以**不推荐**直接 Clone 本仓库到本地进行部署或直接复制 s.yaml 使用，**强烈推荐**通过 `s init ${模版名称}` 的方法或应用中心进行初始化，详情可参考[部署 & 体验](#部署--体验) 。

# start-fc-kaniko 帮助文档

<description>

快速部署 build&push 镜像的应用到阿里云函数计算

</description>


## 资源准备

使用该项目，您需要有开通以下服务并拥有对应权限：

<service>



| 服务/业务 |  权限  | 相关文档 |
| --- |  --- | --- |
| 函数计算 |  AliyunFCFullAccess | [帮助文档](https://help.aliyun.com/product/2508973.html) [计费文档](https://help.aliyun.com/document_detail/2512928.html) |

</service>

<remark>



</remark>

<disclaimers>



</disclaimers>

## 部署 & 体验

<appcenter>
   
- :fire: 通过 [云原生应用开发平台 CAP](https://functionai.console.aliyun.com/template-detail?template=start-fc-kaniko) ，[![Deploy with Severless Devs](https://img.alicdn.com/imgextra/i1/O1CN01w5RFbX1v45s8TIXPz_!!6000000006118-55-tps-95-28.svg)](https://functionai.console.aliyun.com/template-detail?template=start-fc-kaniko) 该应用。
   
</appcenter>
<deploy>
    
   
</deploy>

## 案例介绍

<appdetail id="flushContent">

快速部署一个Dockfile编译和push镜像的应用到阿里云函数计算

</appdetail>







## 使用流程

<usedetail id="flushContent">

# 应用详情

项目部署完成, 直接调用函数实现对指定 git 仓库 (或者zip包的 url）的 Dockfile 工程实现镜像的 build 和 push

> 调用函数的 payload 请参考 evt-sample 示例, event.json 中的 credentials 中的 ak 用于获取阿里云 ACR/ACREE 仓库的临时账号和密码， 该 ak 一定有 ACR/ACREE GetAuthorizationToken 的权限

```bash
# 1. git repo ==> image ==> ACR
$ s invoke --invocation-type async  -f evt-sample/git-acr.json

# 2. devsapp registry ==> image ==> ACR
$ s invoke --invocation-type async  -f evt-sample/registry-acr.json

# 3. git repo ==> image ==> ACREE
$ s invoke --invocation-type async  -f evt-sample/git-acree.json

# 4. devsapp registry ==> image ==> ACREE
$ s invoke --invocation-type async  -f evt-sample/registry-acree.json
```

> **注意**: 最好将您的函数进行 region 化部署， 比如您使用的 url 是 github， 您的这个函数最好部署到新加坡或其他海外 region; 如果是 gitee, 函数最好部署到上海等国内 region; 如果是 gitlab，取决于您的 gitlab 是在国内还是海外。

</usedetail>

## 二次开发指南

<development id="flushContent">

如果您想进行二次开发， 对 code 目录进行开发， 编译新的镜像作为函数 customContainerConfig 的 image 参数即可

</development>






