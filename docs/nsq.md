### nsqlookup
```bash
docker run --name --rm lookupd -p 4160:4160 -p 4161:4161 nsqio/nsq /nsqlookupd
```
### nsqd
```bash
docker run --name nsqd --rm -p 4150:4150 -p 4151:4151 \
    nsqio/nsq /nsqd \
    --broadcast-address=192.168.198.132 \
    --lookupd-tcp-address=192.168.198.132:4160
    --data-path=/data/nsq
```
```
docker run --name nsqadmin --rm nsqio/nsq /nsqadmin --lookupd-http-address=192.168.198.132:4161
```
```
docker run --name nsq_to_file --rm nsqio/nsq /nsq_to_file --topic=test --output-dir=/tmp --lookupd-http-address=192.168.198.132:4161
```