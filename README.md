# prometheus-webhook-dingtalk

###### Generating DingTalk notification from Prometheus AlertManager WebHooks.
#### Version description
```
1、定义了告警和恢复模板输出相关内容；
2、时间更改为本地时区；
3、支持多群组按级别告警；
4、支持@所有人和@多个人；
5、解决timonwong/prometheus-webhook-dingtalk项目有时告警发不出来情况！
6、安全，稳定，娇小
7、过滤多余的label
8、对于频繁告警的场景做了优化
9、支持阿里云电话告警(可按级别和时间段,恢复不产生电话告警)
10, 支持自定义告警功能
11, 支持端口检测功能
```

---

#### Alarm example:
``` 
=====【告警】=====
【告警简述】：内存使用率超过95%  
【告警时间】：2020-07-29 21:53:30 
【告警级别】：warning 
【告警信息】：内存使用超过 95% (当前值95.49). 
【告警主机】：110.133.160.246 
【主机标签】：kehu-app  
--------------------------------
@所有人 

=====【恢复】=====
【告警简述】：内存使用率超过95%  
【告警时间】：2020-07-29 21:53:30 
【恢复时间】：2020-07-29 21:55:30
【告警级别】：warning 
【告警信息】：内存使用超过 95% (当前值95.49). 
【恢复主机】：110.133.160.246  
【主机标签】：kehu-app 
--------------------------------
@所有人 
```
---
#### Install description
```
目前上传的程序支持主流的linux amd64 系统, 本人用的是centos 7！
```
##### 安装步骤
1. 解压包
    ```
   tar -zxvf dingdingalert-v2.2-linux-amd64.tar.gz  -C /usr/local/
    ```
    
2. 修改配置
    ```
    vim /usr/local/dingdingalert/alert.conf
#  release
mode: release
# 指定端口
port: 18089

# 下面这个参数很重要, 用于自动更新值班人和识别是否工作日!  
# 开启后会获取URL里面的key( 如/cambodia 会对因到 dingTalk_config/url_path ) , 自动更换dingTalk_config/second_mobile的电话号码
# URL内容结构: {"/cambodia":"17727901925","/database":"13424251847","/server":"17727901925","WorkMk":"true"}
# project: http://127.0.0.1/dingtalkserver/defaults

# 值班时间段0~24（0~23:59）
alert_time: 18~9
# 电话告警lable, 在prometheus rule定义
alert_level: severity=critical
# 告警每个手机号码收到来电的间隔(分钟), 防止告警轰炸! 阿里云默认限制:1次/分,15次/时,30次/天(这个可以通过添加白名单来开放限制)!
alert_interval: 30
# 发送值班信息时间, 秒、分、时、日、月、周, 用不上可以注释
duty_cron: "8 30 16 * * *"


# 3分钟检测一次
cron_telnet:
  - address: 127.0.0.1:9093
  - address: 127.0.0.1:9090

# 支持自定义告警, 方便在脚本中执行任务失败时, 给群通知
#  curl 'http://192.168.10.218:18089/customs' -H 'Content-Type: application/json' -d '{"message": "## <font color=#ff0000> test 192.168.10.5  时间同步异常! </font>"}'


dingTalk_config:
  # 每个告警url不能一样, 否则启动不成功, 可以不限添加URL
  - url_path: /cambodia
    group_url: https://oapi.dingtalk.com/robot/send?access_token=3edc6581908
    secret: SEC9840ba7d408c9c7f56f55a765d700ae7a965ddd1315309ab8d790d39e3
    at_all: true
    # 项目负责人电话, 负责上班期的所有告警处理主要负责人!
    main_mobile: 138888888888,132666666666
    # 值班负责人电话, 负责值班期间(alert_time), 和project(WorkMk)所有告警处理主要负责人!
    second_mobile: 138888888888,132666666666

  - url_path: /server
    group_url: https://oapi.dingtalk.com/robot/send?access_token=3edc6581908207b4d
    secret: SEC9840ba7d408c9c7f56f55a765d700ae7a965ddd1315309ab8d790d39
    at_all: true
    # 项目负责人电话
    main_mobile: 138888888888,132666666666
    # 值班负责人电话
    second_mobile: 138888888888

  - url_path: /database
    group_url: https://oapi.dingtalk.com/robot/send?access_token=e0455151c0883f433
    secret: SEC9b235b39e7e53e5bc6e9a5b3eb14063d75664ee07f9e413f35f
    at_all: false
    # 项目负责人电话
    main_mobile: 138888888888,132666666666
    # 值班负责人电话
    second_mobile: 

aly_config:
  # 是否开启阿里云电话告警: 0为关闭; 1为开启值班电话; 2为开启全天电话,正常班直接电话给 main_mobile
  open_aly_dx: 2
  # 地域
  region_id: cn-hangzhou
  # 主账号AccessKey的ID
  access_key_id: 
  # 主账号密钥
  access_secret: 
  # 已购买的固定号码,为空则用公共池的号码!
  called_show_number:
  # 文本转语音(TTS)模板ID
  tts_code:

log:
  level: "debug"
  filename: "/usr/local/dingtalkalert/logs.log"
  max_size: 20
  max_age: 30
  max_backups: 7

    ```
    
3. 启动 dingtalk
    ```
    nohup /usr/local/dingtalk/dingdingalert &
    ```
     centos 7 systemd 启动配置
    ```
    tee /etc/systemd/system/dingdingalert.service <<- 'EOF'
    [Unit]
    Description=Prometheus-alertmanager-dingtalk
    Documentation=https://github.com/smarterallen/prometheus-webhook-dingtalk/
    After=network.target

    [Service]
    Type=simple
    User=root
    ExecStart=/usr/local/dingdingalert/dingdingalert
    Restart=on-failure

    [Install]
    WantedBy=multi-user.target
    EOF
    
    systemctl enable dingdingalert; systemctl start dingdingalert
    ```
---
    
#### other config

1. prometheus.yml
    ```
    # Recommended version: 2.0+
      - targets: ['192.168.10.50:9100','192.168.10.51:9100']
        labels:
          env: 'Linux'
          svc: 'kehu-app'
    # labels 下的 key 随便定义，有告警时会把 value 都放到告警信息中的【主机标签】：
    ```

2. rules.yml
    ```
      - alert: "内存使用率过高"
        expr: round(100- node_memory_MemFree_bytes/node_memory_MemTotal_bytes*100) > 90
        for: 1m
        labels:
          severity: warning
          service: DB
        annotations:
          summary: "内存使用率过高"
          description: "当前使用率{{ $value }}%"

    ```

3. alertmanager.yml
    ```
    global:
      resolve_timeout: 5m

    route:
      group_by: ['alertname']
      group_wait: 31s
      group_interval: 3m
      repeat_interval: 30m
      receiver: 'serverAlert'
      routes:
      - receiver: 'DBAlert'
        match_re:
          service: DB.*|DB|UAT-DB|UAT-DB.*

    receivers:
    - name: 'serverAlert'
      webhook_configs:
      - url: 'http://localhost:18089/alert0'
      - url: 'http://localhost:18089/alert2'
    - name: 'DBAlert'
      webhook_configs:
      - url: 'http://localhost:18089/alert1'

    inhibit_rules:
      - source_match:
          severity: 'critical'
        target_match:
          severity: 'warning'
        equal: ['instance']

    ```
    
   4. 其他重要注意点:
    ```
    阿里云操作:   语音服务>>国内语音单呼>>语音通知>>文本转语音模板>>添加模板>>模板内容为: 服务器告警: ${description}
    
    程序内部会把 rule.yml 中的 - alert: "内存使用率过高", 赋值给${description} 进行电话告警
    
    以下情况阿里云可能会不进行通话:
    请勿在变量中添加特殊符号,如: # / : - %￥【】等。
    请勿在变量中包含敏感词汇、IP地址等。
    模板变量实际内容：必须小于20字符以内，不支持传入链接
    
    ```
  
   
