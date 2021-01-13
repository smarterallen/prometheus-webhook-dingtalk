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
```

---

#### Alarm example:
``` 
=====【告警】=====
【告警简述】：内存使用率超过95%  
【告警时间】：2020-07-29 21:53:30 
【告警级别】：warning 
【告警信息】：内存使用超过 95% (当前值95.49). 
【告警主机】：110.33.60.147  
【主机标签】：kehu-app  Linux 
--------------------------------
@所有人 

=====【恢复】=====
【告警简述】：内存使用率超过95%  
【告警时间】：2020-07-29 21:53:30 
【恢复时间】：2020-07-29 21:55:30
【告警级别】：warning 
【告警信息】：内存使用超过 95% (当前值95.49). 
【恢复主机】：110.33.60.147  
【主机标签】：kehu-app  Linux 
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
    # This setting specifies the port to use.
    Port = 18089
    # Delete the unwanted label in the alarm message
    dropLabel = "alertname,instance,job,severity,monitor,device,fstype,mountpoint"


    # setting dingtalk robot alarm interface 
    # important: Loop reading DingDingUrl$, if DingDingUrl2 is empty, it will not continue to fetch new DingDingUrl2+
    DingDingUrl0 = https://oapi.dingtalk.com/robot/send?access_token=6ee807cafb0b222a359604c77c555931658093fb5be2abffa5515292ad7
    secret1 = SECcb7ab8a6cced933c6cfeaede70cf7f7fdd2f7c847cc3251f0d8e9ae53e4bfxx
    # Whether the alarm message is @ everyone in the dingtalk group, isAtAll set "true" or "false"
    isAtAll0 = false
    # 开启阿里云电话告警会根据Mobile列表的电话号码进行拨号
    Mobile0 = ["132xxx78925"]
    # default http requests it receives. Cannot be modified!
    # Url0 = /alert0

    ############################################################################
    DingDingUrl1 = https://oapi.dingtalk.com/robot/send?access_token=1e767be4c7b770224008bd349fcf3b388e1f446b36ec4425b298aba1c180
    secret1 = SECcb7ab8a6cced933c6cfeaede70cf7f7fdd2f7c847cc3251f0d8e9ae53e4bxxx
    isAtAll1 = true
    Mobile1 = []
    # Url1 = /alert1

    DingDingUrl2 = https://oapi.dingtalk.com/robot/send?access_token=88e546e65fa5f557fad5ef2d9f208e792a17736c9d1eb942d036754d4769
    secret2 = SECcb7ab8a6cced933c6cfeaede70cf7f7fdd2f7c847cc3251f0d8e9ae53e4bxxx
    isAtAll2 = true
    Mobile2 = []
    # Url1 = /alert2



    # 是否开启阿里云电话告警,0为关闭,1为开启,恢复不产生电话告警
    openAlyDx=1
    # 地域
    ALY_RegionId=cn-hangzhou
    # 主账号AccessKey的ID
    ALY_AccessKeyId=
    # 主账号密钥
    ALY_AccessSecret=
    # 已购买的固定号码,为空则用公共池的号码!
    ALY_CalledShowNumber=
    # 文本转语音(TTS)模板ID
    ALY_TtsCode=

    # 告警时间段
    ALY_AlertTime=0~23
    # 告警级别, 在prometheus rule定义
    ALY_Level=severity:warning
    # 告警每个手机号码收到来电的间隔(分钟), 防止告警轰炸! 阿里云默认限制:1次/分,15次/时,30次/天
    ALY_AlertInterval=30

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
  
   
