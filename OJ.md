# TESTING

- [ ] disable 题目和用户后排行榜的正常显示问题

# TODO

- [x] Announcement页面逻辑修改，卡片间距离存在问题，对于多行内容，仅显示1~2行，剩下的使用点击卡片，卡片展开的形式呈现

- [x] 题目1、2、3血功能未加上 @@@@@

- [x] Challenge美观调整：Challenges卡片过大，需要按比例缩小

- [x] 分值溢出问题

- [x] Announcement invalued date时间问题

- [ ] 各种依赖包的版本更新和接口修复

- [x] 头像判断：新增多分支头像，优先邮箱匹配gravatar国内代理 `https://cravatar.cn/avatar/`（路径抽取到配置），返回404则使用vue默认头像

- [x] Challenge功能调整：对于需要启动pod的动态题目，提供开启环境的按钮

- [x] Admin功能调整：Admin可以查看User的Submit（解题列表），将Submit查看页面由admin的user页，转移到对应User的profile页

- [x] Disable后该user仍在排行榜上的问题

- [x] 创建题目的时候的逻辑问题：创建题目的时候Disable应该为不可用默认值，修改题目的时候才是可用按钮

- [ ] session过期逻辑：在本地存在cookies但是任意页面返回forbidden的时候，应该弹出提示框重新登陆 @@

- [x] 创建题目成功后自动刷新页面

- [x] 被forbiden的账户不能登录，且登录时显示账户被封禁

- [x] 被disable的用户仍然在排行榜上的问题

- [x] 创建的题目外部链接，会附带一个空链接

- [ ] 缺少上传dockerfile或者配置包的页面

- [ ] 缺少目前服务器集群镜像管理的页面

- [ ] 缺少目前服务器集群机器状态的页面

- [ ] 附加外部链接格式统一解析（优先级不高）

- [x] session持久化存储，目前F5 cookie会掉

- [x] 由于submit查询转移到了userprofile，所以admin的user管理页面点击头像应该能跳转到user profile

- [x] ~~管理员上传home.html自定义页面~~

- [x] admin不参与排行榜

- [x] challenge展示的时候加上所属类型

- [ ] 后端start replic stop replic

- [ ] 单用户全局唯一replic

```c
//maybe :为空则username匹配`https://avatar.sourcegcdn.com/gh/<github name>`
```

# RoadMap

- [ ] 第一次点击题目的时间记录

- [ ] 第一次点击题目时间和完成题目时间的数据导出功能

- [ ] 前后端联动，每个用户同时只能启动一个Pod（全局）

- [ ] 开赛时间、比赛结束时间、比赛页面
