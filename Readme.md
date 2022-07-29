# 项目整体结构
|                    |                        |                                                                                                  |
| ------------------ | ---------------------- | ------------------------------------------------------------------------------------------------ |
| Kubernetes         | 分布式集群架构           | [https://kubernetes.io/](https://kubernetes.io/)                                             |
| Mysql              | 数据库管理系统  | [https://www.mysql.com/](https://www.mysql.com/)                                                 |
| Redis              | 缓存中间件     | [https://redis.io/](https://redis.io/)                                                 |
| Minio              | 兼容AWS-S3的分布式对象存储             | [https://min.io/](https://min.io/)                     |
| Backend            | 后端服务：基于Golang    | [https://go.dev/](https://go.dev/)                                             |
| Frontend           | 前端服务：基于Vue3+TypeScript      | [https://vuejs.org/](https://vuejs.org/)                                     |

## Kubernetes
[点击查看](docs/k8s.md)

## Backend
### 部署
[点击查看](docs/backend.md)
### 功能介绍
1. 基础功能：用户注册/登录、关注/拉黑；发布动态、点赞、评论、回复点赞评论；热门搜索、关键词自动补全； 
2. 消息功能；基于WebSocket协议，实现用户私信、系统消息推送；包括关注内容推送、点赞评论推送、@用户推送；延时推送；未读消息计数；离线消息缓存；
3. 数据缓存：基于Redis，缓存登录状态、用户基本信息、文章列表、搜索关键词、点赞评论数据；
4. 音视频功能：基于WebRTC，推流视频
5. 拓展功能：基于YOLOv7目标检测算法，自动生成图片标签；

## Frontend
### 部署
[点击查看](docs/backend.md)
### 技术依赖
|                    |                        |                                                                                                  |
| ------------------ | ---------------------- | ------------------------------------------------------------------------------------------------ |
| Vue3               | 渐进式 JavaScript 框架 | [https://v3.cn.vuejs.org/](https://v3.cn.vuejs.org/)                                             |
| TypeScript         | JavaScript 的一个超集  | [https://www.tslang.cn/](https://www.tslang.cn/)                                                 |
| Vite2              | 前端开发与构建工具     | [https://cn.vitejs.dev/](https://cn.vitejs.dev/)                                                 |
| Element Plus       | UI 组件库              | [https://element-plus.gitee.io/zh-CN/](https://element-plus.gitee.io/zh-CN/)                     |
| Pinia              | 新一代状态管理工具     | [https://pinia.vuejs.org/](https://pinia.vuejs.org/)                                             |
| Uno css            | 原子 CSS 引擎          | [https://github.com/unocss/unocss](https://github.com/unocss/unocss)                             |

## Minio 
[点击查看](docs/minio.md)

## Mysql
[点击查看](docs/mysql.md)

## Redis 
[点击查看](docs/redis.md)

