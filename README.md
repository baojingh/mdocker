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

# Ubuntu install Graphviz
for pprof visual
```
apt-get install -y graphviz
```
# prof analyzer
### Solution1
```go
import "runtime/pprof" 

func monitorMem() {
	f, err := os.Create("cpu.pprof")
	if err != nil {
		log.Fatal("could not create memory profile: ", err)
	}
	defer f.Close() // error handling omitted for example
	runtime.GC()    // get up-to-date statistics
	if err := pprof.WriteHeapProfile(f); err != nil {
		log.Fatal("could not write memory profile: ", err)
	}
}
```
go tool pprof -http 127.0.0.1:8848 ./cpu.pprof  
go tool pprof cpu.pprof

### Solution2
```go
	"net/http"
	_ "net/http/pprof"

  go func() {
		http.HandleFunc("/test", Test)
		http.ListenAndServe(":9988", nil)
	}()
```
go tool pprof -http=":8080" http://localhost:9988/debug/pprof/profile?seconds=10

go tool pprof http://localhost:9988/debug/pprof/profile?seconds=10

visit 
http://localhost:9988/debug/pprof/


```
allocs: A sampling of all past memory allocations
block: Stack traces that led to blocking on synchronization primitives
cmdline: The command line invocation of the current program
goroutine: Stack traces of all current goroutines
heap: A sampling of memory allocations of live objects. You can specify the gc GET parameter to run GC before taking the heap sample.
mutex: Stack traces of holders of contended mutexes
profile: CPU profile. You can specify the duration in the seconds GET parameter. After you get the profile file, use the go tool pprof command to investigate the profile.
threadcreate: Stack traces that led to the creation of new OS threads
trace: A trace of execution of the current program. You can specify the duration in the seconds GET parameter. After you get the trace file, use the go tool trace command to investigate the trace.
```


```
# web方式查看 cpu 30秒内的占用情况
go tool pprof -http=":8080" https://xxx.com/debug/pprof/profile?seconds=30

# web方式查看 memory 30秒内占用情况
go tool pprof -http=":8080" https://xxx.com/debug/pprof/heap?seconds=30

# web方式查看 goroutine 30秒内的占用情况
go tool pprof -http=":8080" https://xxx.com/debug/pprof/goroutine?seconds=30

# 查看火焰图
http://localhost:8080/ui/flamegraph
```


# pprof Usage
### CPU
flat 	当前函数占用 cpu 耗时
flat % 	当前函数占用 cpu 耗时百分比
sum% 	函数占用 cpu 时间累积占比，从小到大一直累积到 100%
cum 	当前函数加上调用当前函数的函数占用 cpu 的总耗时
%cum 	当前函数加上调用当前函数的函数占用 cpu 的总耗时占比




# 引用
https://www.cnblogs.com/jiujuan/p/14588185.html




