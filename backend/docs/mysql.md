## 拉取mysql镜像
```shell
docker pull mysql
```

## 主库的配置文件
```shell 
vim /usr/local/mysql/master/conf/my.cnf
```
```bash
[mysqld]

#[必须]服务器唯一ID，默认是1
server-id=1

#打开Mysql日志，日志格式为二进制
log-bin=/var/lib/mysql/binlog

#每次执行写入就与硬盘同步
sync-binlog=1

#关闭名称解析
skip-name-resolve

#只保留7天的二进制日志，以防磁盘被日志占满
expire-logs-days = 7

#需要同步的二进制数据库名
binlog-do-db=blog

#不备份的数据库
binlog-ignore-db=information_schema
binlog-ignore-db=performation_schema
binlog-ignore-db=sys
binlog-ignore-db=mysql
```

## 从库的配置文件
```shell 
vim /usr/local/mysql/master/conf/my.cnf
```
```bash
[mysqld]
#配置server-id，让从服务器有唯一ID号
server-id=2

#开启从服务器二进制日志
log-bin=/var/lib/mysql/binlog

#打开mysql中继日志，日志格式为二进制
relay_log=/var/lib/mysql/mysql-relay-bin

#设置只读权限
read_only=1

#使得更新的数据写进二进制日志中
log_slave_updates=1

#如果salve库名称与master库名相同，使用本配置
replicate-do-db=blog
```

## 启动主库master的容器
```shell
docker run -d \
-p 3307:3306 \
--name mysql_master \
-v /usr/local/mysql/master/conf:/etc/mysql/conf.d \
-v /usr/local/mysql/master/logs:/logs \
-v /usr/local/mysql/master/data:/var/lib/mysql \
-e MYSQL_ROOT_PASSWORD=password \
mysql \
--character-set-server=utf8mb4 \
--collation-server=utf8mb4_general_ci 
```

## 启动从库slave1的容器
```shell
docker run -d \
-p 3308:3306 \
--name mysql_slave1 \
-v /usr/local/mysql/slave/conf:/etc/mysql/conf.d \
-v /usr/local/mysql/slave/logs:/logs \
-v /usr/local/mysql/slave/data:/var/lib/mysql \
-e MYSQL_ROOT_PASSWORD=password \
mysql \
--character-set-server=utf8mb4 \
--collation-server=utf8mb4_general_ci  
```

## 查看主库的binlog信息
```bash 
#记录下File字段和Pos字段
show master status;
```
## 开启主从复制
```bash 
docker exec -it mysql_slave1 bash 
mysql -u root -p 
# source_log_file对应上面File字段，source_log_pos对应上面的Pos字段
change replication source to source_host='192.168.198.132', source_user='root', source_password='password', source_port=3307, source_log_file='xxxxx', source_log_pos=xxx, source_connect_retry=30;
start replica;
```