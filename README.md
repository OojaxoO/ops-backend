# 运维平台后端 

**提供cmdb,容器中心,部署监控,运维文档等功能**
## 设计
https://note.youdao.com/web/#/file/recent/markdown/WEBd7168340fa1170e0e4f8851e037331ea/ 
## 配置
vim /etc/profile  
export GO111MODULE=on  
GOPROXY=https://goproxy.io  
export GOPROXY  

source /etc/profile  

## 运行
make && make install  
cd /opt/ops-backend/  
./ops-backend  

## 同步主机  
go run cron/cron.go  

## 数据库迁移  
go run migrate/migrate.go  



