## IP-INFO
[![Website ip-info.run.place](https://img.shields.io/website-up-down-green-red/https/ip-info.run.place.svg)](https://ip-info.run.place)
[![GitHub release](https://img.shields.io/github/release/streamdp/ip-info.svg)](https://github.com/streamdp/ip-info/releases/)
[![test](https://github.com/streamdp/ip-info/actions/workflows/test.yml/badge.svg)](https://github.com/streamdp/ip-info/actions/workflows/test.yml)
[![GitHub license](https://img.shields.io/github/license/streamdp/ip-info.svg)](https://github.com/streamdp/ip-info/blob/master/LICENSE)
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
requested and when the IP is requested randomly (currently internal query caches are not implemented).
* **http** benchmarking with [hey - HTTP load generator tool](https://github.com/rakyll/hey):
```shell
$ hey -c 2 -n 10000 http://127.0.0.1:8080/ip-info?ip=8.8.8.8
Summary:
  Total:        3.7542 secs
  Slowest:      0.0097 secs
  Fastest:      0.0005 secs
  Average:      0.0007 secs
  Requests/sec: 2663.6510
```
* **gRPC** benchmarking with [ghz - Simple gRPC load testing tool](https://github.com/bojand/ghz):
```shell
$ ghz -c 2 -n 10000 127.0.0.1:50051 --call IpInfo.GetIpInfo -d '{"ip":"8.8.8.8"}' --insecure 
Summary:
  Count:        10000
  Total:        6.73 s
  Slowest:      8.95 ms
  Fastest:      0.61 ms
  Average:      1.09 ms
  Requests/sec: 1485.69
```
