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
   tar -zxvf dingdingalert-v2.1-linux-amd64.tar.gz  -C /usr/local/
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
        DingDingUrl0 = https://oapi.dingtalk.com/robot/send?access_token=6ee807cafb0b222a359604c77c555931658093fb5be2abffa5515292a
        secret0 = SECcb7ab8a6cced933c6cfeaede70cf7f7fdd2f7c847cc3251f0d8e9

        # Whether the alarm message is @ everyone in the dingtalk group, isAtAll set "true" or "false"
        isAtAll0 = true
        Mobile0 = ["132xx678925", "13035xx9308"]
        # default http requests it receives. Cannot be modified!
        # Url0 = /alert0


        ############################################################################
        DingDingUrl1 = https://oapi.dingtalk.com/robot/send?access_token=1e767be4c7b770224008bd349fcf3b388e1f446b36ec4425b298aba
        secret1 = SECcb7ab8a6cced933c6cfeaede70cf7f7fdd2f7c847cc3251f0d8e
        isAtAll1 = false
        Mobile1 = ["132xx678925", "189xxxx8325"]
        # Url1 = /alert1


        ############################################################################
        DingDingUrl2 = https://oapi.dingtalk.com/robot/send?access_token=88e546e65fa5f557fad5ef2d9f208e792a17736c9d1eb942d036754
        secret2 = SECfc2f185555526c78e0f95bef692a78fb445dcdd9e0d0624f93352a
        isAtAll2 = false
        Mobile2 = ["132xx678925", "189xxxx8325"]
        # Url2 = /alert2

        ############################################################################
        DingDingUrl3 = https://oapi.dingtalk.com/robot/send?access_token=1e767be4c7b770224008bd349fcf3b388e1f446b36ec4425b298ab
        secret3 = SECcb7ab8a6cced933c6cfeaede70cf7f7fdd2f7c847cc3251f0d8e
        isAtAll3 = false
        Mobile3 = ["132xx678925", "189xxxx8325"]
        # Url3 = /alert3


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
      group_wait: 10s
      group_interval: 10s
      repeat_interval: 10s
      receiver: 'criticalalert'
      routes:
      - receiver: 'criticalalert'
        match:
          severity: 'critical'
      - receiver: 'warningalert'
        match_re:
          severity: 'critical'

    receivers:
    - name: 'criticalalert'
      webhook_configs:
      - url: 'http://localhost:18089/alert1'
    - name: 'warningalert'
      webhook_configs:
      - url: 'http://localhost:18089/alert0'
      - url: 'http://localhost:18089/alert2'
    inhibit_rules:
      - source_match:
          severity: 'critical'
        target_match:
          severity: 'warning'
        equal: ['alertname', 'dev', 'instance']

    ```
