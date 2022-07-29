## [安装kbs](https://blog.csdn.net/SHELLCODE_8BIT/article/details/122192034)
### 配置免密登录
```
sudo vi ~/.ssh/config
# 指定别名为 master1
Host master1
# 指定目标 ip
hostname 192.168.56.11
# 指定登录用户名
user root
# 复制到 master1 主机
scp root@worker1:~/.ssh/id_rsa.pub /home
# 写入到 authorized_keys 中
cat /home/id_rsa.pub >> ~/.ssh/authorized_keys
```
### 安装docker
```
sudo apt install docker.io
sudo vi /etc/docker/daemon.json
{
"registry-mirrors": [
"https://docker.mirrors.ustc.edu.cn"
],
"exec-opts": [ "native.cgroupdriver=systemd" ]
}
sudo systemctl daemon-reload
sudo systemctl restart docker
```
### 添加Kubernetes签名密钥
```bash
curl -s https://packages.cloud.google.com/apt/doc/apt-key.gpg | sudo apt-key add
```

### 添加Kubernetes软件包存储库
```bash
sudo apt-get update
sudo apt-get install -y apt-transport-https ca-certificates curl
# 需要翻墙 或者下载本地拷贝
sudo curl -fsSLo /usr/share/keyrings/kubernetes-archive-keyring.gpg https://packages.cloud.google.com/apt/doc/apt-key.gpg
echo "deb [signed-by=/usr/share/keyrings/kubernetes-archive-keyring.gpg] https://apt.kubernetes.io/ kubernetes-xenial main" | sudo tee /etc/apt/sources.list.d/kubernetes.list
# or 国内源
echo "deb https://mirrors.aliyun.com/kubernetes/apt/ kubernetes-xenial main" | sudo tee /etc/apt/sources.list.d/kubernetes.list
sudo apt-get update
sudo apt-get install -y kubectl
sudo apt-get install -y kubelet=1.23.4-00 kubeadm=1.23.4-00 kubectl=1.23.4-00
```
### 自动补全
```bash
apt-get install bash-completion
kubectl completion bash | sudo tee /etc/bash_completion.d/kubectl > /dev/null
```
### 安装 kubectl convert 插件
```bash
curl -LO https://dl.k8s.io/release/$(curl -L -s https://dl.k8s.io/release/stable.txt)/bin/linux/amd64/kubectl-convert
curl -LO "https://dl.k8s.io/$(curl -L -s https://dl.k8s.io/release/stable.txt)/bin/linux/amd64/kubectl-convert.sha256"
echo "$(cat kubectl-convert.sha256) kubectl-convert" | sha256sum --check
sudo install -o root -g root -m 0755 kubectl-convert /usr/local/bin/kubectl-convert
kubectl convert --help
```
```
apiVersion: v1
kind: Config

proxy-url: https://proxy.host:3188

clusters:
- cluster:
  name: 507

users:
- name: qiu

contexts:
- context:
  name: development
```
```
kubeadm init \
--apiserver-advertise-address=192.168.198.132 \
--image-repository registry.aliyuncs.com/google_containers \
--pod-network-cidr=10.244.0.0/16

kubeadm init --image-repository registry.cn-hangzhou.aliyuncs.com/google_containers \
  --apiserver-advertise-address=192.168.198.132

kubeadm init \
  --pod-network-cidr=10.244.0.0/16 \
  --service-cidr=10.96.0.0/16  \
  --kubernetes-version v1.24.3 \ 


kubeadm init \
--image-repository registry.aliyuncs.com/google_containers \
--service-cidr=10.96.0.0/16  \
--pod-network-cidr=10.244.0.0/16

kubeadm init \
--apiserver-advertise-address=192.168.198.132 \
--image-repository registry.aliyuncs.com/google_containers \
--service-cidr=10.96.0.0/16  \
--pod-network-cidr=10.244.0.0/16
kubeadm token create --ttl 0 --print-join-command
openssl x509 -pubkey -in /etc/kubernetes/pki/ca.crt | openssl rsa -pubin -outform der 2>/dev/null | openssl dgst -sha256 -hex | sed 's/^.* //'

kubeadm join 192.168.192.132:6443 --token ulrvkf.be5sovl7ddu4qmzt --discovery-token-ca-cert-hash sha256:c901e5dd420020ac7729591e6cd420acfe817cef9a4f436c06ca38aafab36710
#!/bin/bash


netstat -ntlup|grep 10250
systemctl status kubelet
journalctl -xefu kubelet
```
### 加入节点
```
kubeadm join 192.168.56.11:6443 --token wbryr0.am1n476fgjsno6wa --discovery-token-ca-cert-hash sha256:7640582747efefe7c2d537655e428faa6275dbaff631de37822eb8fd4c054807
```
### 重启服务
```
systemctl stop kubelet.service && \
systemctl daemon-reload && \
systemctl start kubelet.service
```

```
kubectl get apiservices
```
### 通过 YAML 文件创建一个 Deployment
```
kubectl apply -f https://k8s.io/examples/application/deployment.yaml
kubectl describe deployment nginx-deployment
kubectl get pods -l app=nginx
kubectl delete deployment nginx-deployment
```
```
kubectl get nodes
kubectl get cs
kubectl get deploy
kubectl get rs
kubectl -n kube-system get pod
kubectl get po -o wide

kubectl rollout status deployment/nginx-deployment
```
### 移除污点
```
 kubectl taint nodes --all node-role.kubernetes.io/master-
kubectl taint nodes --all node.kubernetes.io/disk-pressure- 
kubectl describe nodes k8s-master |grep Taints
kubectl taint node k8s-master node.kubernetes.io/disk-pressure:NoSchedule-
```

### 
```bash 
#mysql客户端
kubectl run -it --rm --image=mysql:5.6 --restart=Never mysql-client -- mysql -h mysql -ppassword
```

### 
```
kubectl expose deployment hello-world --type=NodePort --name=example-service
```




kubectl describe secret dashboard-admin-token-ht2vx -n kube-system|grep '^token'|awk '{print $2}'

### minio
```
kubectl minio tenant create tenant1 --servers 2 --volumes 2 --capacity 2G
kubectl minio tenant create tenant-1 \
	--servers 1                            \
	--volumes 2                           \
	--capacity 2G                         \
	--namespace minio-tenant-1              \
	--storage-class local-storage        

kubectl minio tenant create minio-tenant-1   \
  --servers                 1                \
  --volumes                 4               \
  --capacity                4Gi             \
  --storage-class           local-storage    \
  --namespace               minio-tenant-1
kubectl config set-context --current --namespace=minio-tenant-1

  --namespace               minio-tenant-1
docker run \
  -p 9000:9000 \
  -p 9001:9001 \
  -e "MINIO_ROOT_USER=qiu" \
  -e "MINIO_ROOT_PASSWORD=06278611" \
  minio/minio server /data --console-address ":9001"
kubectl get pod kube-controller-manager-k8s-master-control-plane -n kube-system -o yaml

helm install minio \
  --namespace minio1 --create-namespace \
  --set accessKey=htiNtOXcf1hPQ7ZI,secretKey=PNuVJXz5L2qaYSk6xcJOFsFAaVELdTnq \
  --set mode=standalone \
  --set service.type=NodePort \
  --set persistence.enabled=true \
  --set persistence.size=10Gi \
  --set persistence.storageClass=longhorn  \
  --set resources.requests.memory=2Gi \
  minio/minio

helm install minio \
  --namespace minio1 --create-namespace \
  --set accessKey=htiNtOXcf1hPQ7ZI,secretKey=PNuVJXz5L2qaYSk6xcJOFsFAaVELdTnq \
  --set mode=distributed \
  --set replicas=2 \
  --set service.type=NodePort \
  --set persistence.size=20Gi \
  --set persistence.storageClass=longhorn \
  --set resources.requests.memory=1Gi \
  minio/minio

helm -n minio uninstall minio
kubectl -n minio delete pvc --all 

kubectl -n tenant1-ns patch svc minio -p '{"spec": {"type": "NodePort"}}'
kubectl -n longhorn-system patch svc longhorn-frontend -p '{"spec": {"type": "NodePort"}}'

```
# [部署minio](https://blog.csdn.net/networken/article/details/111469223)

```
docker container run -d -p 13306:3306 -v /data/mysql_wordpress:/var/lib/mysql --rm --name mysql_wordpress --env MYSQL_ROOT_PASSWORD=123456  --env MYSQL_DATABASE=wordpress mysql:5.7

```
```
docker container run \
  -d \
  -p 8081:80 \
  --rm \
  --name wordpress \
  --env WORDPRESS_DB_PASSWORD=123456 \
  --link mysql_wordpress:mysql \
  --volume "/data/wordpress":/var/www/html \
  wordpress

$ docker run -d --privileged=true --rm --name mysql_wordpress -p 3310:3306 -v /data/mysql:/var/lib/mysql -e MYSQL_DATABASE=wordpress -e MYSQL_ROOT_PASSWORD=123456  mysql:5.7
$ docker run -d --name wordpress --rm -e WORDPRESS_DB_HOST=172.20.127.220:3310 -e WORDPRESS_DB_USER=root -e WORDPRESS_DB_PASSWORD=123456 -p 8081:80 --link mysql_wordpress:mysql wordpress

```