# mdocker

# What is mdocker
this is a tool to manage docker container. Detailed features:
1. Display container logs with websocket in browser.
2. Display container resource stats with websocket.
3. Exec container in browser with websocket.

# Dependency
1. docker-ce 20.10
2. go 1.19
3. gorrila websocket
4. docker go sdk



# feature
1. user and role management including password, role, access
2. multi user login and usage
3. security for websocket and docker
4. monitor websocket server status, provide API for metrics.
5. user, connections management. cache


# Is it secure to get docker.sock in the container?
Yes. It's very dangerous and we should avoid this.
1. Use Https instead of docker.sock file.
2. Alternative ways to meet your target.





# DEV Tools
websocket web
http://www.easyswoole.com/wstool.html


# Issues
1. why ubuntu cannot open html
https://askubuntu.com/questions/1184357/why-cant-chromium-suddenly-access-any-partition-except-for-home



### test
```
curl -H "Content-Type: application/json" -X POST -d '{"username": "hadoop", "password":"hadoop"}'  http://localhost:8081/login
```

# DONE
1. login is ok, but cookie seems not integrate in to frontend
2. container list cannot jump to exec dashboard
3. I think I should save real-time data in database.


# Install Influxdb
### Install
```bash
docker pull influxdb:2.7.1
mkdir -p  /data/dockerImages/influxdb/{data,config,log}

cat > /data/dockerImages/influxdb/docker-compose.yaml << "EOF"
version: '2'
services:
  influxdb:
    image: influxdb:2.7.1
    container_name: influxdb
    ports:
      - "8086:8086"
    environment:
      DOCKER_INFLUXDB_INIT_MODE: "setup"
      DOCKER_INFLUXDB_INIT_USERNAME: "hadoop"
      DOCKER_INFLUXDB_INIT_PASSWORD: "Hadoop.123" 
      DOCKER_INFLUXDB_INIT_ORG: "mdocker"
      DOCKER_INFLUXDB_INIT_BUCKET: "mdocker-bucket"
    volumes:
      - /data/dockerImages/influxdb/data:/var/lib/influxdb2
      - /data/dockerImages/influxdb/config:/etc/influxdb2
EOF

docker-compose  up -d
```

### Config
```
nXom-eG4c_SwYgoDsxFOrfz_hIg0_je6tSDTkC9MssBxl5F0eyVzSEIOazx_89513o4Y5Ld2NOwsllCK41L3xg==
```

# Development
1. add go package
go get github.com/influxdata/influxdb-client-go/v2

### Influxdb Usage
1. bucket查看
influx bucket list -o <org-name or org-id> -t <your-token>

2. 执行过后所有需要该token的指令就不需要指定token. File /etc/influxdb2/influx-configs exits.
influx config create --active -n mdocker-config  -t  nXom-eG4c_SwYgoDsxFOrfz_hIg0_je6tSDTkC9MssBxl5F0eyVzSEIOazx_89513o4Y5Ld2NOwsllCK41L3xg== -u http://localhost:8086 -o mdocker
```bash
root@aa9073f6d803:/etc/influxdb2# cat influx-configs 
[mdocker-config]
  url = "http://localhost:8086"
  token = "nXom-eG4c_SwYgoDsxFOrfz_hIg0_je6tSDTkC9MssBxl5F0eyVzSEIOazx_89513o4Y5Ld2NOwsllCK41L3xg=="
  org = "mdocker"
  active = true
```

3. list config info
influx config ls


4. 查看用户
root@aa9073f6d803:/etc/influxdb2# influx  user  list
ID			Name
0bef9742361cd000	hadoop

5. 查看已存在的 org
influx org list

6. view data in db
```bash
# influx query
from(bucket:"mdocker-bucket")
    |> range(start: -1d)
    |> filter(fn: (r) => r.tagname1 == "tagvalue1")
```

Ctrl + D will execute and display the data.

# Uusal CMD
```
sudo kill -TERM  $(ps -ef | grep mdocker | grep -v "grep" | awk '{print $2}')

```















