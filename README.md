# Tp0t OJ

延续老版本的Tp0t OJ，前端继承并优化了视觉效果，后端从JAVA切换为Golang，定位也从一个面向校内的持续性训练平台，转向于一个比CTFd更好用的CTF比赛平台

### 我们的特性

| 特性                                         | Tp0t OJ | CTFd |
| -------------------------------------------- | :-----: | :--: |
| K8s集群 -> 一键部署题目环境                  | :heavy_check_mark: | :x: |
| 题目定义与实例分离 -> 原生支持动态Flag       | :heavy_check_mark: | :x: |
| 预构建单执行文件 -> 无其他依赖，平台一键启动 | :heavy_check_mark: | :x: |
| 提供WriteUp的上传和下载通道                  | :heavy_check_mark: | :x: |

## 开发指南

### 前端

+ **请在`app`目录下打开vscode**
+ **npm相关命令，请在`app`目录下运行**
+ 总之请将`app`做为工作目录

### 后端

+ **使用Goland将`server`作为工作目录**

### 接口相关

+ GraphQL的schema文件定义在`server/src/resources/schema`目录下

- 请求成功返回message 为空字符串
