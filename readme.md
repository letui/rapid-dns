## Rapid-DNS
### 功能介绍
1. A记录解析
2. Web控制台
3. 支持扩展域名
4. 多账号注册

### 技术
#### 前端 Vue ElementUI Axios
#### 后端 gin boltdb

### 使用安装

    cd rapid-dns
    go build
    ./rapid-dns

### 环境变量说明
    rapid_dns_port=53
    rapid_web_port=8053
    rapid_domain_list=what,names,you,like
环境变量用来自定义启动的一些配置，比如http的监听端口，DNS的服务监听端口以及域名的支持列表
。系统默认自带.api .test .prod .web .ui .db等域名后缀

- 默认账号密码 rapid/rapid （重启后自动重置密码，以防止密码丢失）
- 系统启动后会生成rapid.dat，用来保存域名的注册信息，如需要清空可以删除文件后重启
- 使用DNS服务，请记得修改系统DNS ServerIP，即启动该程序的主机就是DNS Server

## 截图说明

- 首页
![title 首页]（https://raw.githubusercontent.com/letui/rapid-ui/master/public/images/homepage.png）