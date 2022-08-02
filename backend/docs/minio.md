
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
apt install s3fs
echo $ACCESS_KEY_ID:$SECRET_ACCESS_KEY > /etc/passwd-s3fs
chmod 600 /etc/passwd-s3fs
s3fs data /data -o passwd_file=/etc/passwd-s3fs -o url=$MINIO_ENDPOINT -o use_path_request_style
```