# Project507-backend

## Mysql配置
https://blog.csdn.net/galoiszhou/article/details/118359174
https://wenku.baidu.com/view/7ce144e32fc58bd63186bceb19e8b8f67c1ceffe?aggId=a7f9e2fed25abe23482fb4daa58da0116c171fce

### 拉取mysql镜像
``` bash
docker pull mysql
```

### 启动容器
``` bash
docker run -d \
-p 3307:3306 \
--name mysql \
-v /home/qiu/mysql/conf:/etc/mysql/conf.d \
-v /home/qiu/mysql/logs:/logs \
-v /home/qiu/mysql/data:/var/lib/mysql \
-e MYSQL_ROOT_PASSWORD=06278611 \
mysql \
--character-set-server=utf8mb4 \
--collation-server=utf8mb4_general_ci 
```