# MEEGO 米格会议管理插件后端

## 部署说明

### 前置准备

项目需要同时使用到飞书和飞书项目的开放平台能力，在服务端程序启动的时候，需要将应用凭证作为参数。如果你还没有创建过飞书开发平台应用和飞书项目插件应用，你需要在运行程序之前准备好这些内容。

**程序启动的必要参数如下**

| 参数 | 说明|
| --- | --- |
| lark_app_id | 飞书应用凭证 app_id |
| lark_app_secret | 飞书应用凭证 app_secret |
| meego_plugin_id | 飞书项目插件凭证 app_id |
| meego_plugin_secret | 飞书项目插件凭证 app_secret |

本章节提及到的所有 `yourhost.com` 都需要替换为你想要部署的域名地址


#### 创建&&配置飞书应用

1. 飞书项目空间管理员或开发者参考[如何开发企业自建应用](https://www.feishu.cn/hc/zh-CN/articles/360049067916)？会议管理需要创建一个飞书应用。

应用能力需要选择包含机器人

2. 飞书项目空间管理员或开发者申请妙计、日历、会议等权限。具体可参考：[申请权限](https://open.feishu.cn/document/ukTMukTMukTM/uMTNz4yM1MjLzUzM)。

可以将如下权限直接导入使用

```
{
  "scopes": {
    "tenant": [
      "calendar:calendar",
      "calendar:calendar:readonly",
      "contact:user.email:readonly",
      "contact:user.employee:readonly",
      "contact:user.employee_id:readonly",
      "contact:user.phone:readonly",
      "im:message.group_msg",
      "im:message:readonly",
      "minutes:minutes:readonly",
      "vc:meeting.all_meeting:readonly",
      "vc:meeting:readonly",
      "vc:record:readonly"
    ],
    "user": [
      "calendar:calendar",
      "calendar:calendar:readonly",
      "contact:user.email:readonly",
      "contact:user.employee:readonly",
      "contact:user.employee_id:readonly",
      "contact:user.phone:readonly",
      "im:message:readonly",
      "minutes:minutes:readonly",
      "vc:meeting:readonly",
      "vc:record:readonly"
    ]
  }
}
```
![alt text](https://raw.githubusercontent.com/jimu5/meego_meeting_plugin/main/docs/img/server/1.png)

3. 在事件与回调中配置好订阅方式和回调事件

订阅方式中的请求域名请替换成将要部署的域名 `https://yourhost.com/api/v1/lark/webhook/event`


![alt text](https://raw.githubusercontent.com/jimu5/meego_meeting_plugin/main/docs/img/server/2.png)

4. 在安全设置中配置重定向 url

配置的 url 为(需要自行替换域名) `https://yourhost.com/api/v1/meego/lark/auth`

![alt text](https://raw.githubusercontent.com/jimu5/meego_meeting_plugin/main/docs/img/server/3.png)

5. 记录你的凭证信息, 作为后续服务端程序的启动参数

![alt text](https://raw.githubusercontent.com/jimu5/meego_meeting_plugin/main/docs/img/server/4.png)

#### 创建飞书项目插件

由于本插件暂时不需要额外的权限设置，创建过程详见 [飞书项目开发者中心](https://project.feishu.cn/b/helpcenter/1p8d7djs/359lzbgu)。

创建好之后, 记录下飞书项目插件的凭证信息

![alt text](https://raw.githubusercontent.com/jimu5/meego_meeting_plugin/main/docs/img/server/5.png)

### 从二进制文件运行服务端程序

1. 下载或者从源码编译二进制文件

下载链接: TODO

从源码编译: 详见从源码编辑章节

2. 执行服务端二进制文件

```
./meego_meeting_plugin -lark_app_id=飞书开放app_id -lark_app_secret=飞书开放app_secret -meego_plugin_id=飞书项目插件id -meego_plugin_secret=飞书项目插件plugin_secret
```

启动参数说明

| 参数 | 说明|
| --- | --- |
| lark_app_id | 飞书应用凭证 app_id |
| lark_app_secret | 飞书应用凭证 app_secret |
| meego_plugin_id | 飞书项目插件凭证 app_id |
| meego_plugin_secret | 飞书项目插件凭证 app_secret |

3. 执行命令之后，如果没有意外的话，后端服务将运行在 7999 端口上

> TODO: 后续会将支持自定义配置启动端口

4. 配置反向代理等

配置好域名, 并将域名配置到飞书开放平台的配置上。

### 从源码构建程序

1. 安装好 golang 1.19

2. 拉取本项目源码到本地

```
git clone https://github.com/jimu5/meego_meeting_plugin.git
```

3. 编译

默认情况下会编译出当前平台程序, 如有其他诉求, 自行修改 `go build` 指令

```
cd server && go mod tidy && go build -ldflags="-s -w" -o meego_meeting_plugin .
```

上面命令执行完毕之后, 编译产物会出现在 `./server/` 目录下, 文件名是`meego_meeting_plugin`

