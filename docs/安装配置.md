```bash

cat ~/.ssh/id_rsa.pub | ssh USER@HOSTNAME 'umask 0077; mkdir -p .ssh; cat >> .ssh/authorized_keys'

ufw enable
ufw default deny
ufw allow ssh
ufw allow http
ufw allow https
# ufw allow 27017
# ufw allow 27018


cat > /etc/default/locale << EOF
LC_ALL="en_US.UTF-8"
LANGUAGE="en_US.UTF-8"
LC_CTYPE="en_US.UTF-8"
EOF

cat >> /etc/hosts << EOF
# 生产环境 hosts 配置
# 内网 IP           主机名   子域名               # 实例 ID        外网 IP          服务器
10.251.254.160      mongo1  mgo1.set.shou.money # i-2382ei0i4   121.43.226.227  MongoDB 副本集 1
10.252.136.28       mongo2  mgo2.set.shou.money # i-23uj89pb4   121.43.227.71   MongoDB 副本集 2
10.117.21.10        nsq     nsq1.set.shou.money # i-234ynman3   120.26.201.241  消息队列 NSQ
10.168.121.237      nginx   ngx1.set.shou.money # i-23d6jcevu   121.40.167.112  统一代理 Nginx
10.171.199.158      app1    app1.set.shou.money # i-23jva95w7   121.40.224.235  应用服务器 1
10.171.239.103      app2    app2.set.shou.money # i-234ynman3   121.40.225.122  应用服务器 2
EOF



wget http://nginx.org/download/nginx-1.8.0.tar.gz

apt-get update
apt-get install -y libpcre3 libpcre3-dev zlib1g zlib1g.dev libssl-dev openssl  \
    libgd2-xpm-dev libgoogle-perftools-dev


./configure --prefix=/opt/nginx \
    --with-http_image_filter_module \
    --with-google_perftools_module \
    --with-http_ssl_module \
    --with-http_spdy_module \
    --with-http_gunzip_module \
    --with-http_gzip_static_module \
    --with-http_stub_status_module \
    --with-http_realip_module \
    --with-http_addition_module

配置代理 quickpay/quick_master/quick_sett 和 Basic Auth


参照《MongoDB 安装配置.md》安装 Mongo

mkdir -p /opt/quickpay/logs
cd /opt/quickpay
nohup ./quickpay >> logs/quickpay.log 2>&1 &


快捷支付安装配置环境要求：

1. 服务器购买（Hyman）－－ 已完成
2. 域名配置（Hyman）－－ 已完成
3. MongoDB 安装配置（migo） －－ 已完成
4. Nginx 安装配置（migo） －－ 安装已完成，证书已配置，代理未配置
5. QuickPay 应用发布（migo）
6. 卡 Bin 数据（rui）
7. 应答码数据（wonsikin）
8. 中金渠道银行数据（wonsikin）
9. 中金生产环境证书和接口地址（Brian）
10. 线下联机接口地址和商户（Brian）
11. api.shou.money 证书（Brian）－－ 已完成


```
