# Project507-backend

## Mysql 配置

构建Mysql的Docker容器，配置主库从库，[详细步骤](docs/mysql.md#)。

## FFMpeg
```
docker build -t backend .
docker run -it  --privileged=true --name testff testff -e MINIO_ENDPOINT=http://192.168.198.132:32000 -e ACCESS_KEY_ID=htiNtOXcf1hPQ7ZI -e SECRET_ACCESS_KEY=PNuVJXz5L2qaYSk6xcJOFsFAaVELdTnq
``` 