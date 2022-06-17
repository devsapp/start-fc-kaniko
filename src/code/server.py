# coding=utf-8

from flask import Flask, request
import subprocess
import json
import sys
import traceback
import logging
import base64

logging.getLogger('werkzeug').setLevel(logging.ERROR)

app = Flask(__name__)

REQUEST_ID_HEADER = 'x-fc-request-id'


@app.route("/", methods=['GET', 'POST'])
def build_image():
    rid = request.headers.get(REQUEST_ID_HEADER)
    print('FC Invoke Start RequestId: {}'.format(rid))
    data = request.stream.read()
    try:
        evt = json.loads(data)
        print(evt)
        url = evt['url']
        # 如果是私有 github, 可以把 token 一起传过来， 或者 token 作为函数的环境变量
        # https://codeantenna.com/a/PuZENneALu
        # 这里示例只考虑 public
        subprocess.check_call(
            "git clone {} /tmp/workspace".format(url), shell=True)

        image = evt['image']  # 镜像名称

        # 当然， 如果是您个人账号构建自己的镜像,也可以直接在镜像中完成 /kaniko/.docker/config.json 的创建, 下面的代码就不需要了

        # registry 默认是 dockerhub
        # 比如是阿里云的 acr， 示例 registry.cn-hangzhou.aliyuncs.com
        registry = evt.get('registry', 'https://index.docker.io/v1/')
        usr = evt['usr']  # dockerhub 或者 acr 账户名
        pwd = evt['pwd']  # dockerhub 或者 acr 账户密码
        auth = str(base64.b64encode(
            "{}:{}".format(usr, pwd).encode()), "utf-8")
        content = {
            "auths": {
                registry: {
                    "auth": auth
                }
            }
        }
        with open('/kaniko/.docker/config.json', 'w+') as f:
            f.write(json.dumps(content))

        subprocess.check_call("cat /kaniko/.docker/config.json", shell=True)
        # 镜像的编译和 push
        subprocess.check_call("ls -lh /tmp/workspace", shell=True)
        subprocess.check_call('executor --force=true --cache=false --use-new-run=true  \
            --dockerfile /tmp/workspace/Dockerfile --context /tmp/workspace \
            --destination {}'.format(image), shell=True)

    except Exception as e:
        exc_info = sys.exc_info()
        trace = traceback.format_tb(exc_info[2])
        errRet = {
            "message": str(e),
            "stack": trace
        }
        print("FC Invoke End RequestId: " + rid +
              ", Error: Unhandled function error")
        print(errRet)
        return errRet, 404, [("x-fc-status", "404")]

    print('FC Invoke End RequestId: {}'.format(rid))
    return "OK"


if __name__ == "__main__":
    app.run(host='0.0.0.0', port=9000)
