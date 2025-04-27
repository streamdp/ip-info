## IP-INFO
[![Website ip-info.oncook.top](https://img.shields.io/website-up-down-green-red/https/ip-info.oncook.top.svg)](https://ip-info.oncook.top)
[![GitHub release](https://img.shields.io/github/release/streamdp/ip-info.svg)](https://github.com/streamdp/ip-info/releases/)
[![test](https://github.com/streamdp/ip-info/actions/workflows/test.yml/badge.svg)](https://github.com/streamdp/ip-info/actions/workflows/test.yml)
### Microservice for IP-based geolocation
This microservice is a small, independent software application designed to determine the geographic location of a device 
based on its IP address. It achieves this by using a free public database called [db-ip.com](https://db-ip.com)
(*free version provides about __77%__ accuracy*), which contains a vast amount
of information linking IP addresses to specific locations.
### Key features:
* **Automatic Database Updates:** The microservice regularly updates its local copy of the **db-ip.com** database to 
ensure that the location data is always accurate and up-to-date.
* **Fast Lookup:** It is optimized to perform quick searches within the database, allowing it to efficiently determine 
the location associated with a given IP address.
* **HTTP and gRPC Support:** The microservice can be accessed and interacted with using both protocols, providing 
flexibility in how it can be integrated into other systems or applications.
* **Rate limiting:** The microservice provides per-client rate limits and sends a **429** HTTP response when the client makes 
requests too frequently.
* **Caching:** The microservice implements caching to improve availability and reduce database load. 
### Usage example:
Start postgresql and ip-info containers:
```shell
$ docker-compose up -d
```
We are waiting for the database to be updated:
```shell
$ docker logs -f ip-info-container 
IP_INFO: 2024/09/22 16:41:44 HTTP server listening at :8080
IP_INFO: 2024/09/22 16:41:44 gRPC server listening at [::]:50051
IP_INFO: 2024/09/22 16:41:45 ip database update started
IP_INFO: 2024/09/22 16:41:45 truncate ip_to_city_one table before importing update
IP_INFO: 2024/09/22 16:41:45 droping ip_to_city_one_ip_start_gist_idx index
IP_INFO: 2024/09/22 16:41:45 import ip database updates
IP_INFO: 2024/09/22 16:43:41 creating  ip_to_city_one_ip_start_gist_idx index on ip_to_city_one table
IP_INFO: 2024/09/22 16:45:18 switching backup and working tables
IP_INFO: 2024/09/22 16:45:18 updating database config
IP_INFO: 2024/09/22 16:45:18 ip database update completed, next update through 223.2h
```
And we can make several test requests:
```shell
$ curl localhost:8080/ip-info?ip=8.8.8.8
{
  "error": "",
  "content": {
    "ip": "8.8.8.8",
    "continent": "NA",
    "country": "US",
    "state_prov": "California",
    "city": "Mountain View",
    "latitude": -122.085,
    "longitude": 37.4223
  }
}
```
```shell
$ grpcurl  -plaintext -d '{"ip": "211.27.38.98"}' 127.0.0.1:50051 IpInfo/GetIpInfo
{
  "ip": "211.27.38.98",
  "continent": "OC",
  "country": "AU",
  "stateProv": "New South Wales",
  "city": "Sydney",
  "latitude": 151.209,
  "longitude": -33.8688
}
```
### Benchmarking (i3-7100U CPU @ 2.40GHz)
IP randomization is not supported for security reasons, the difference in tests is about 10% for cases where 1 IP is 
requested and when the IP is requested randomly.
* **http** benchmarking with [hey - HTTP load generator tool](https://github.com/rakyll/hey) **without** cache:
```shell
$ hey -c 2 -n 10000 http://127.0.0.1:8080/ip-info?ip=8.8.8.8
  Total:        3.7542 secs
  Slowest:      0.0097 secs
  Fastest:      0.0005 secs
  Average:      0.0007 secs
  Requests/sec: 2663.6510
```
when **redis** cache used:
```shell
  Total:        2.1551 secs
  Slowest:      0.0355 secs
  Fastest:      0.0003 secs
  Average:      0.0004 secs
  Requests/sec: 4640.0855
```
when **memory** cache used:
```shell
  Total:        1.0665 secs
  Slowest:      0.0032 secs
  Fastest:      0.0001 secs
  Average:      0.0002 secs
  Requests/sec: 9376.7398
```
* **gRPC** benchmarking with [ghz - Simple gRPC load testing tool](https://github.com/bojand/ghz) **without** cache:
```shell
$ ghz -c 2 -n 10000 127.0.0.1:50051 --call IpInfo.GetIpInfo -d '{"ip":"8.8.8.8"}' --insecure 
  Total:        6.73 s
  Slowest:      8.95 ms
  Fastest:      0.61 ms
  Average:      1.09 ms
  Requests/sec: 1485.69
```
when **redis** cache used:
```shell
  Total:        4.40 s
  Slowest:      3.55 ms
  Fastest:      0.37 ms
  Average:      0.63 ms
  Requests/sec: 2270.86
```
when **memory** cache used:
```shell
  Total:        3.04 s
  Slowest:      2.41 ms
  Fastest:      0.19 ms
  Average:      0.37 ms
  Requests/sec: 3285.49
```
### Rate limiting
You could choose **limiter** between [redis_rate](https://github.com/go-redis/redis_rate) (_redis_ should be present)
and [golimiter](https://github.com/streamdp/golimiter), using **-limiter** flag or **IP_INFO_LIMITER**
environment variable. To enable rate limiting run _ip-info_ microservice with the **-enable-limiter** flag or 
**IP_INFO_ENABLE_LIMITER** environment variable. The default rate limit value is 10 requests per second per client, 
you can adjust it with the **-rate-limit** flag or **IP_INFO_RATE_LIMIT** environment variable. 
```shell
version: "3.4"
services:
   ip-info:
      image: streamdp/ip-info:v0.2.0
      container_name: ip-info
      environment:
         - IP_INFO_DATABASE_URL=postgresql://postgres:postgres@postgres:5432/dbip?sslmode=disable
         - IP_INFO_ENABLE_LIMITER=true
         - IP_INFO_RATE_LIMIT=15 # default 10 requests per second per client
         - IP_INFO_LIMITER=redis_rate # default "golimiter"
         - REDIS_URL=redis://:qwerty@redis:6379/0
      ports:
         - "8080:8080"
         - "50051:50051"
      restart: always
   
   redis:
      image: redis
      container_name: redis
      ports:
         - "6379:6379"
      command: redis-server --save "" --maxmemory 64mb --maxmemory-policy allkeys-lfu --requirepass qwerty
      restart: always
```
```shell
$ docker-compose up -d
```
### Caching
Caching in-memory with [microcache](https://github.com/streamdp/microcache) library is enabled by default, to disable you need 
to run _ip-info_ microservice with the **-disable-cache** flag or **IP_INFO_DISABLE_CACHE=true** environment variable. 
The default **TTL** value is **3600** seconds, you can adjust it with the **-cache-ttl** flag or **IP_INFO_CACHE_TTL** 
environment variable. You could choose cacher between **redis** and **microcache**, using **-cacher** 
flag or **IP_INFO_CACHER** environment variable. 
```shell
version: "3.4"
services:
   ip-info:      
      environment:
         - IP_INFO_CACHE_TTL=1800 # default 3600 seconds
         - IP_INFO_CACHER=redis # default "microcache"
         - REDIS_URL=redis://:qwerty@redis:6379/0
```
### Help 
You can see all available command flags when you run the application with the -h flag.
```shell
$ ./bin/app -h
ip-info is a microservice for IP location determination

Usage of ./bin/app:
  -cache-ttl int
        cache ttl in seconds (default 3600)
  -cacher string
        where to store cache entries: redis, microcache (default "microcache")
  -db-request-timeout int
        database request timeout in milliseconds (default 5000)
  -disable-cache
        disable cache
  -enable-limiter
        enable rate limiter
  -grpc-port int
        grpc server port (default 50051)
  -h    display help
  -http-port int
        http server port (default 8080)
  -http-read-timeout int
        http server read timeout (default 5000)
  -limiter string
        what use to limit queries: redis_rate, golimiter (default "golimiter")
  -rate-limit int
        rate limit, rps per client (default 10)
  -rate-limit-ttl int
        rate limit entries ttl in seconds (default 60)
  -read-header-timeout int
        http server read header timeout (default 5000)
  -redis-db int
        redis database
  -redis-host string
        redis host (default "127.0.0.1")
  -redis-port int
        redis port (default 6379)
  -v    display version
  -write-timeout int
        http server write timeout (default 5000)
```
### Contributing
Contributions are welcome! If you encounter any issues, have suggestions for new features, or want to improve **ip-info**, please feel free to open an issue or submit a pull request on the project's GitHub repository.
### License
**ip-info** is released under the _GPL 3.0_ License. See the [LICENSE](https://github.com/streamdp/ip-info/blob/master/LICENSE) file for complete license details.
