# 概述

这个是米格会议插件的前端工程文件

## 调试项目

1. 在项目目录下运行：`npm start`
2. 然后访问「飞书项目·开发者后台」首页，在页面左下角启用插件的本地调试模式。
3. 最后打开「飞书项目」，预览插件的效果。

## 发布你的插件

### 修改配置

1. 修改 plugin.config.json 文件中 `pluginID` 和 `pluginSecret` 配置
2. 修改 `src/constants/index.ts`文件中的 `PLUGIN_ID`和 `PLUGIN_SECRET`配置
3. 修改文件 ` fe/src/models/api/index.ts` 中的 `CUSTOM_API_PREFIX `, `REDRECT_URL`, `FEISHU_APP_ID`配置
4. 在前端目录下执行命令 `lpm config set pluginSecret {你的secret}`

### 发布产物

1. 在终端上运行 `npm run release` 命令来构建产物并上传。
2. 打开「飞书项目·开发者后台」相应插件详情页。
3. 左侧导航切换到「插件功能」tab，添加对应功能构成并完善配置。
4. 左侧导航切换到「插件发布」tab。
5. 点击「创建版本」填写相关信息，在「产物版本』中选择对应版本产物并提交。
6. 回到「插件发布」页面，会出现一条新增的版本记录，点击该记录的「申请发布」按钮。
7. 恭喜，你现在已经成功发布了一个插件，可以回到「飞书项目」插件市场去安装并尽情使用啦！
