#!/bin/bash
echo > /etc/apt/sources.list
# RUN echo "deb http://mirrors.aliyun.com/debian/ stretch main non-free contrib \ndeb-src http://mirrors.aliyun.com/debian/ stretch main non-free contrib \ndeb http://mirrors.aliyun.com/debian-security stretch/updates main \ndeb-src http://mirrors.aliyun.com/debian-security stretch/updates main \ndeb http://mirrors.aliyun.com/debian/ stretch-updates main non-free contrib \ndeb-src http://mirrors.aliyun.com/debian/ stretch-updates main non-free contrib \ndeb http://mirrors.aliyun.com/debian/ stretch-backports main non-free contrib \ndeb-src http://mirrors.aliyun.com/debian/ stretch-backports main non-free contrib" > /etc/apt/sources.list
echo "deb http://mirrors.aliyun.com/ubuntu/ bionic main restricted universe multiverse" >> /etc/apt/sources.list
echo "deb-src http://mirrors.aliyun.com/ubuntu/ bionic main restricted universe multiverse" >> /etc/apt/sources.list
apt-get update && \
    apt-get install -y \
    s3fs \
    && rm -rf /var/lib/apt/lists/*
echo $ACCESS_KEY_ID:$SECRET_ACCESS_KEY > /etc/passwd-s3fs
chmod 600 /etc/passwd-s3fs
mkdir data && mkdir /data/video && mkdir /data/img && mkdir /data/preview && mkdir /data/temp && mkdir /data/avatar
# s3fs video /data/video -o passwd_file=/etc/passwd-s3fs -o url=$MINIO_ENDPOINT -o use_path_request_style
# s3fs img /data/img -o passwd_file=/etc/passwd-s3fs -o url=$MINIO_ENDPOINT -o use_path_request_style
# s3fs preview /data/preview -o passwd_file=/etc/passwd-s3fs -o url=$MINIO_ENDPOINT -o use_path_request_style
# s3fs temp /data/temp -o passwd_file=/etc/passwd-s3fs -o url=$MINIO_ENDPOINT -o use_path_request_style
# s3fs avatar /data/avatar -o passwd_file=/etc/passwd-s3fs -o url=$MINIO_ENDPOINT -o use_path_request_style
s3fs data /data -o passwd_file=/etc/passwd-s3fs -o url=$MINIO_ENDPOINT -o use_path_request_style
/app