# 运维相关配置信息

```yaml

service_name: hgms-layout       #服务名称
service_level: S1               #服务等级(部署前与上级沟通填写)
service_brief: xxx              #服务描述(自行填写)
service_owner: xxx              #服务owner(自行填写)

program_language: go            #编程语言
program_language_version: go1.17 #依赖语言版本

git_addr: ssh://git@dev.int.hypergryph.com:7999/pubtech/hgms-layout.git #git地址
build_cmd: make build           #构建命令
build_dir: ./bin                #构建输出目录
run_cmd: ./hgms-layout          #启动命令

first_deploy_branch: master     #首次发布分支(自行填写)

APOLLO_APPID: hgms-layout       #apollo appid
APOLLO_NS: bootstrap.yaml       #apollo namespace

if_enable_istio: true           #是否注入sidecar
service_protocol:               #服务协议
    http: 8080
    grpc: 8081
domian_conf:                    #域名配置[公网](内网域名&公司白名单域名 默认自动生成)(部署前与上级沟通填写)
    domain1: {path: /api, network: 公网, 环境: pre,prod}

```