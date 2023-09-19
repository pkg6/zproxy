## 安装

~~~
go install github.com/pkg6/zproxy@main
~~~

## 启动socks服务

~~~
zproxy start socks -H 0.0.0.0 -P 1080 -A admin@123456
~~~

## 启动http服务

~~~
zproxy start http -H 0.0.0.0 -P 1080 -A admin@123456
~~~

## 启动服务

~~~
zproxy run -C config.yml
~~~

## 加入启动项

~~~
root@zproxy-node:~# cat > /etc/systemd/system/zproxy.service << EOF
[Unit]
Description=Start zproxy
[Service]
ExecStart=/usr/local/bin/zproxy run
[Install]
WantedBy=multi-user.target
EOF

root@zproxy-node:~# cat > /etc/zproxy/config.yaml << EOF
Host: 0.0.0.0
HttpPort: 1081
SocksPort: 1080
Auth:
  admin: 123456
EOF

root@zproxy-node:~# mv $(go env GOPATH)/bin/zproxy /usr/local/bin
root@zproxy-node:~# systemctl daemon-reload
root@zproxy-node:~# systemctl start zproxy
~~~
