
## longhorn部署
```
helm repo add longhorn https://charts.longhorn.io
helm repo update
helm install longhorn longhorn/longhorn --namespace longhorn-system --create-namespace
```
## minio部署
```
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
## minio挂载
```
mkdir /data/video && mkdir /data/img && mkdir /data/preview && mkdir /data/temp && mkdir /data/avatar
apt install s3fs
echo ACCESS_KEY_ID:SECRET_ACCESS_KEY > /etc/passwd-s3fs
chmod 600 /etc/passwd-s3fs
s3fs video /data/video -o passwd_file=/etc/passwd-s3fs -o url=http://192.168.198.132:32000/ -o use_path_request_style
s3fs img /mnt/data/mountminio/img -o passwd_file=/etc/passwd-s3fs -o url=http://192.168.198.132:32000/ -o use_path_request_style
s3fs preview /mnt/data/mountminio/preview -o passwd_file=/etc/passwd-s3fs -o url=http://192.168.198.132:32000/ -o use_path_request_style
```