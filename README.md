## IP-INFO
### Microservice for IP-based geolocation
This microservice is a small, independent software application designed to determine the geographic location of a device based on its IP address. It achieves this by using a free public database called db-ip.com, which contains a vast amount of information linking IP addresses to specific locations.
### Key features:
* **Automatic Database Updates:** The microservice regularly updates its local copy of the db-ip.com database to ensure that the location data is always accurate and up-to-date.
* **Fast Lookup:** It is optimized to perform quick searches within the database, allowing it to efficiently determine the location associated with a given IP address.
* **HTTP and gRPC Support:** The microservice can be accessed and interacted with using both protocols, providing flexibility in how it can be integrated into other systems or applications.
### Usage example:
```shell
docker-compose up -d .
```
```shell
curl localhost:8080/ip-info?ip=8.8.8.8
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