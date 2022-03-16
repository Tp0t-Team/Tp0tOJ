# Tp0t OJ

延续老版本的Tp0t OJ，前端继承并优化了视觉效果，后端从JAVA切换为Golang，定位也从一个面向校内的持续性训练平台，转向于一个比CTFd更好用的CTF比赛平台

### 我们的特性

| 特性                       | Tp0t OJ            | CTFd |
| ------------------------ |:------------------:|:----:|
| K8s集群 -> 一键部署题目环境        | :heavy_check_mark: | :x:  |
| 题目定义与实例分离 -> 原生支持动态Flag  | :heavy_check_mark: | :x:  |
| 预构建单执行文件 -> 无其他依赖，平台一键启动 | :heavy_check_mark: | :x:  |
| 提供WriteUp的上传和下载通道        | :heavy_check_mark: | :x:  |

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

## 平台部署指南

## 平台使用指南

需要使用k8s集群的基本就是PWN题和WEB题

[Pwn_demo1_](https://github.com/Tp0t-Team/Tp0tOJ_demos/tree/main/pwn1)



### 镜像编译

在镜像编译页面上传镜像，需要上传包含Dock而file的tar包，**注意tar包没有额外目录层级，需要直接能获取到Dockerfile**

对于PWN题，推荐使用`xinetd`作为守护进程

#### Dockerfile示例

> - **注意对于所有需要执行的文件，附加执行权限，否则镜像会build成功，但是用户申请创建实例的时候会失败**
> 
> - 对于singleton的题目（所有选手共用一个容器实例），请严格注意权限管控

```dockerfile
FROM ubuntu:20.04
RUN sed -i "s/http:\/\/archive.ubuntu.com/http:\/\/mirrors.ustc.edu.cn/g" /etc/apt/sources.list
RUN apt-get update
RUN apt-get -y upgrade
RUN apt-get install -y apt-utils lib32z1 xinetd
RUN useradd -u 8888 -m pwn
COPY share/libunicorn.so.1 /usr/local/lib/libunicorn.so.1
RUN chmod 755 /usr/local/lib/libunicorn.so.1
RUN ldconfig
COPY share/easiestpwn /home/pwn/easiestpwn
RUN chmod 755 /home/pwn/easiestpwn
RUN rm /etc/xinetd.d/*
COPY xinetd /etc/xinetd.d/xinetd
COPY entrypoint.sh /home/pwn/entrypoint.sh
ENTRYPOINT ["/home/pwn/entrypoint.sh"]
EXPOSE 8888

```



#### `entrypoint.sh`示例

对于动态flag的题目，必须具备`entrypoint.sh`文件，因为平台会将生成的随机flag通过`entrypoint.sh`的第一个参数的形式传入

```shell
#!/bin/sh
echo $1 > flag #动态FLAG的必须行，用于平台将生成的FLAG写入，也可自行调整写入位置
/usr/sbin/xinetd -dontfork #启动守护进程
```



#### 守护进程`xinetd`示例

```
service pwn 
{
    disable = no
    type        = UNLISTED
    wait        = no
    server      = /bin/sh
    # replace helloworld to your program
    server_args = -c cd${IFS}/home/pwn;exec${IFS}./easiestpwn
    socket_type = stream
    protocol    = tcp
    user        = pwn 
    port        = 8888
    # bind        = 0.0.0.0
    # safety options
    flags       = REUSE
    per_source	= 10 # the maximum instances of this service per source IP address
    rlimit_cpu	= 1 # the maximum number of CPU seconds that the service may use
    #rlimit_as  = 1024M # the Address Space resource limit for the service
    #access_times = 2:00-9:00 12:00-24:00
    nice        = 18
}

```



#### 打包示例

打包可以使用gz、xz压缩，但是推荐使用tar直接打包，打包不能包含当前文件夹

举个例子，如果Dockerfile等文件都放在`/home/pwn/example`下，需要运行以下指令打包

```shell
cd /home/pwn/example
tar -cvf ../example.tar ./*
```



### 题目部署

yaml文件用于上一键导入题目信息以及配置题目所属的K8S节点，需要使用节点部署的题只能通过导入config的形式添加

```yaml
name: pwn1 #challenge页面显示的题目名称
category: PWN #题目类型
score:
  baseScore: 1000 #基本分数，动态积分和奖励根据基本分数的一定比例计算
  dynamic: true #是否开启动态积分
flag: 
  value: flag{test} #对于开启动态flag的题目value字段没用，flag由平台随机生成
  dynamic: true #是否开启动态积分
description: "this is a test pwn"
externalLink: [http://cloud.lordcasser.com/s/DkxtK] #题目附件链接
singleton: false #false时题目所有用户都是同一个环境/flag，true时每人一个环境
nodeConfig:
  - name: "pwn1" #节点名称，要求[a~z 0~9]且必须有一个字母
    image: "pwn1" #题目使用镜像的名称，需要提前上传好所使用的镜像
    servicePorts:
      - name: http #不用更改
        protocol: TCP #不用更改
        external: 8888 #指定docker对外暴露的端口
        internal: 8888 #指定docker容器内端口
        pod: 0 # 集群pod对外的端口，也是用户实际访问的端口，0为自动分配
```
