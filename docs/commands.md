# 常用命令

```shell

# 查看 jks 中证书列表
keytool -list -keypass cfca1234 -keystore temp/trust.jks

# 从 jks 中导出证书
keytool -exportcert -alias cfca_ev_oca -keypass cfca1234 -keystore temp/trust.jks -rfc -file cfca_ev_oca_crt.pem

# 生成 docker 镜像，并上传
docker build -t quickpay .
docker save -o quickpay.tar quickpay
rsync -rcv --progress quickpay.tar root@121.40.86.222:/opt/quickpay.tar

# 启动 MongoDB 容器
docker run -d -p 27017:27017 -p 27018:27018 -v /opt/data/mongo:/data/db --name mongo mongo
# 启动 Quickpay 容器
docker run -d -p 80:3009 --link mongo:mongo --name quickpay quickpay

```
