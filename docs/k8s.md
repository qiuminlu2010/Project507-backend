## 事前准备
* 至少三台物理机，一台做master，其余做node。 内网IP例如：
  - master: 192.168.198.132
  - node1: 192.168.198.133
  - node2: 192.168.198.134
* 准备基于Debian的Linux系统，如Ubuntu。（其他类型的系统的指令有所区别）
* 配置SSH密钥免密登录

## 修改host
以master为例
```
vi ~/.ssh/config
```
填入
```bash
# 指定别名为 master
Host master
# 指定目标 ip
hostname 192.168.198.132
# 指定登录用户名
user root
```
## 安装docker
```
apt install docker.io
vi /etc/docker/daemon.json
```
```bash
# 修改/etc/docker/daemon.json
{
"registry-mirrors": [
"https://docker.mirrors.ustc.edu.cn"
],
"exec-opts": [ "native.cgroupdriver=systemd" ]
}
```
```
systemctl daemon-reload
systemctl restart docker
```
## 安装Kubernetes
```bash
apt-get install -y apt-transport-https ca-certificates curl
#添加Kubernetes软件包存储库国内源
echo "deb https://mirrors.aliyun.com/kubernetes/apt/ kubernetes-xenial main" | sudo tee /etc/apt/sources.list.d/kubernetes.list
#指定版本
apt-get install -y kubelet=1.23.4-00 kubeadm=1.23.4-00 kubectl=1.23.4-00
```
## master初始化集群
```
kubeadm init \
--apiserver-advertise-address=192.168.198.132 \
--image-repository registry.aliyuncs.com/google_containers \
--service-cidr=10.96.0.0/16  \
--pod-network-cidr=10.244.0.0/16
```

## node加入集群
```bash 
# token和hash由master给出
kubeadm join 192.168.192.132:6443 --token TOKENXXX --discovery-token-ca-cert-hash SHA256XXX
```

## 重新生成token和hash
```bash
kubeadm token create --ttl 0 --print-join-command
openssl x509 -pubkey -in /etc/kubernetes/pki/ca.crt | openssl rsa -pubin -outform der 2>/dev/null | openssl dgst -sha256 -hex | sed 's/^.* //'
```