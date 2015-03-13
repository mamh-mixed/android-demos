# 常用命令

```shell

# 查看 jks 中证书列表
keytool -list -keypass cfca1234 -keystore temp/trust.jks

# 从 jks 中导出证书
keytool -exportcert -alias cfca_ev_oca -keypass cfca1234 -keystore temp/trust.jks -rfc -file cfca_ev_oca_crt.pem

```
