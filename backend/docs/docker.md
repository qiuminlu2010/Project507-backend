### 安装存储库
```yum install -y yum-utils
```
### 配置镜像源
```
yum-config-manager --add-repo http://mirrors.aliyun.com/docker-ce/linux/centos/docker-ce.repo
```
### 安装docker
```
yum install docker-ce docker-ce-cli containerd.io
```
### 开启docker服务
```
systemctl start docker
```
### 验证
```
docker run hello-world
```
### 搭建博客网站
```
docker run -d --privileged=true --rm --name mysql_wordpress -p 3310:3306 -v /data/mysql:/var/lib/mysql -e MYSQL_DATABASE=wordpress -e MYSQL_ROOT_PASSWORD=123456  mysql:5.7
docker run -d --name wordpress --rm -e WORDPRESS_DB_HOST=172.20.127.220:3310 -e WORDPRESS_DB_USER=root -e WORDPRESS_DB_PASSWORD=123456 -p 8081:80 --link mysql_wordpress:mysql wordpress
```