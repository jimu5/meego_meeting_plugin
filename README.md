# Meego 米格会议管理插件

## 前言

你是否会遇到此类问题：围绕一个项目/需求要开很多会，会议的信息目前散落在飞书的日程、会议、妙计中，查询成本高企，无法沉淀到一起管理。米格会议管理插件通过手动和自动关联的形式将会议相关的信息高效汇总到一起。

节点展示
![alt text](https://raw.githubusercontent.com/jimu5/meego_meeting_plugin/main/docs/img/1.png)
详情页展示
![alt text](https://raw.githubusercontent.com/jimu5/meego_meeting_plugin/main/docs/img/2.png)

## 整体解决方案

![alt text](https://raw.githubusercontent.com/jimu5/meego_meeting_plugin/main/docs/img/3.png)

## 功能说明

### 手动关联

在节点表单或者详情页标签页中搜搜关联日程，关联后即可把会议名称、组织者、描述、时间、妙计、状态、参与者数量等信息汇总到一起。其中会议名称点击可跳转回飞书日历，方便查看更多的信息。

**搜索日程名称进行关联**

如下图所示，搜索并关联
![alt text](https://raw.githubusercontent.com/jimu5/meego_meeting_plugin/main/docs/img/4.jpeg)

**重复日程处理**

重复日程不需要手动挨个关联，插件会自动识别并将重复的日程所有关联上，后续也可以一起取消关联

![alt text](https://raw.githubusercontent.com/jimu5/meego_meeting_plugin/main/docs/img/5.png)

### 自动关联

打开自动关联日程开关后，将当前工作项实例群添加为日程参与者并将日程分享到群内时，系统会自动将日程与当前实例关联。

![alt text](https://raw.githubusercontent.com/jimu5/meego_meeting_plugin/main/docs/img/6.png)

## 未来规划

1. **权限隔离：用户只能看到自己有权限的会议信息**
2. **部分功能支持 **AI** 增强**
   1. **列表页展示 AI 会议纪要**
   2. **通过 AI 生成代办，并转成节点工作项任务 or 工作项**实例**，在 **meego** 实现会议资产的闭环管理**
   3. **会前会议智能创建：根据节点排期、实例角色等信息自动 book 会议**
3. **支持重度的会议管理模式：部分客户会使用工作项来管理会议，通过预制工作项的方式管理会议。插件支持预置会议管理工作项的重会议管理模式，关联日程后自动创建会议实例并将会议信息同步到各类字段信息中。**

## 目录说明

### server

plugin 服务端文件

### fe

plugin 前端文件

## 部署说明

详见 server 和 fe 文件夹下的 README
