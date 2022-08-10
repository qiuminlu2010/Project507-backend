# 项目整体结构
|                    |                        |                                                                                                  |
| ------------------ | ---------------------- | ------------------------------------------------------------------------------------------------ |
| Kubernetes         | 分布式集群架构           | [https://kubernetes.io/](https://kubernetes.io/)                                             |
| Mysql              | 数据库管理系统  | [https://www.mysql.com/](https://www.mysql.com/)                                                 |
| Redis              | 缓存中间件     | [https://redis.io/](https://redis.io/)                                                 |
| Nsq                | 分布式消息系统     | [https://nsq.io/](https://nsq.io/)                                                 |
| Minio   | 分布式对象存储             | [https://min.io/](https://min.io/)                     |
| Backend            | 后端Http服务：基于Golang    | [https://go.dev/](https://go.dev/)                                             |
| Frontend           | 前端服务：基于Vue3+TypeScript      | [https://vuejs.org/](https://vuejs.org/)                                     |
| AutoTag            | 自动标签： 基于深度学习ASL模型   |                                          |

## ***Kubernetes***
Kubernetes集群部署步骤
[点击查看](docs/k8s.md)

## ***Mysql***
- kubernetes部署mysql主从服务
[查看配置文件](docs/mysql.md)
```bash
# namespace
kubectl apply -f mysql/1-ns.yaml
# configmap 
kubectl apply -f mysql/2-cm.yaml
# secret
kubectl apply -f mysql/3-secret.yaml
# service 
kubectl apply -f mysql/4-svc.yaml
# statefulSet
kubectl apply -f mysql/5-ss.yaml
```
- 扩容集群
```
kubectl -n mysql scale statefulset mysql -—replicas=3
```
## ***Backend***
### Dockerfile
[点击查看](backend/README.md)
### 功能介绍
1. 基础功能：用户注册/登录、关注/拉黑；发布动态、点赞、评论、回复点赞评论；热门搜索、关键词自动补全； 
2. 消息功能；基于WebSocket协议，实现用户私信、系统消息推送；包括关注内容推送、点赞评论推送、@用户推送；未读消息计数；离线消息缓存；
3. 数据缓存：基于Redis，缓存登录状态、用户基本信息、文章列表、搜索关键词、点赞评论数据；
4. 音视频功能：基于HLS协议，实现视频点播; 
5. 拓展功能：基于YOLOv5目标检测算法，自动生成图片标签；

## ***Frontend***
### 安装步骤
[点击查看](frontend/README.md)
### 技术依赖
|                    |                        |                                                                                                  |
| ------------------ | ---------------------- | ------------------------------------------------------------------------------------------------ |
| Vue3               | 渐进式 JavaScript 框架 | [https://v3.cn.vuejs.org/](https://v3.cn.vuejs.org/)                                             |
| TypeScript         | JavaScript 的一个超集  | [https://www.tslang.cn/](https://www.tslang.cn/)                                                 |
| Vite2              | 前端开发与构建工具     | [https://cn.vitejs.dev/](https://cn.vitejs.dev/)                                                 |
| Element Plus       | UI 组件库              | [https://element-plus.gitee.io/zh-CN/](https://element-plus.gitee.io/zh-CN/)                     |
| Pinia              | 新一代状态管理工具     | [https://pinia.vuejs.org/](https://pinia.vuejs.org/)                                             |
| Uno css            | 原子 CSS 引擎          | [https://github.com/unocss/unocss](https://github.com/unocss/unocss)                             |

## ***Minio***
[点击查看](docs/minio.md)
- longhorn部署
```bash
helm repo add longhorn https://charts.longhorn.io && helm repo update
helm install longhorn longhorn/longhorn --namespace longhorn-system --create-namespace
```
- minio部署
```bash
helm install minio \
  --namespace minio1 --create-namespace \
  --set accessKey=XXX \
  --set mode=distributed \
  --set replicas=3 \
  --set service.type=NodePort \
  --set persistence.size=200Gi \
  --set persistence.storageClass=longhorn \
  --set resources.requests.memory=4Gi \
  minio/minio
```

- minio挂载本地
```bash
apt install s3fs
echo $ACCESS_KEY_ID:$SECRET_ACCESS_KEY > /etc/passwd-s3fs
chmod 600 /etc/passwd-s3fs
s3fs data /data -o passwd_file=/etc/passwd-s3fs -o url=$MINIO_ENDPOINT -o use_path_request_style
```

## ***Redis*** 
[点击查看](docs/redis.md)

## 部分截图
![主页](docs/src/%E4%B8%BB%E9%A1%B5.PNG)
![消息](docs/src/%E8%81%8A%E5%A4%A92.PNG)
![播放](docs/src/%E6%92%AD%E6%94%BE.PNG)
