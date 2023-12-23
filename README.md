# dns_exporter
Golang DNS exporter

# About
This service is capable of querying specified DNS servers and domains for health checks. It is suitable for users with a specific DNS infrastructure who require active health checking of servers or domains. Additionally, the service has the capability to log information to files, which proves useful when a dedicated exporter is not needed, and metrics can be pulled using, for example, Node Exporter.

# Run App
```
./dns_exporter --config.path="./config.yaml"
```