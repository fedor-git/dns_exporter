# dns_exporter
Golang DNS exporter

# About
This service is capable of querying specified DNS servers and domains for health checks. It is suitable for users with a specific DNS infrastructure who require active health checking of servers or domains. Additionally, the service has the capability to log information to files, which proves useful when a dedicated exporter is not needed, and metrics can be pulled using, for example, Node Exporter.

# Run App
```
./dns_exporter --config.path="./config.yaml"
```

# Exporter otput example
```prometheus
# HELP dns_availability DNS Server Status
# TYPE dns_availability gauge
dns_availability{dns_server="1.1.1.1"} 1
# HELP dns_lookup_time DNS lookup time measurement
# TYPE dns_lookup_time gauge
dns_lookup_time{dns_server="1.1.1.1",hostname="facebook.com"} 1.106679625
dns_lookup_time{dns_server="1.1.1.1",hostname="fwefwf.hh"} 0.102711166
dns_lookup_time{dns_server="1.1.1.1",hostname="google.com"} 0.127847708
dns_lookup_time{dns_server="1.1.1.1",hostname="microsoft.com"} 0.146392791
```
# Configuration file
```yaml
servers:
 - 1.1.1.1:
    foo: bar
 - 8.8.8.8

hosts:
 - google.com
 - facebook.com
 - microsoft.com
 - fwefwf.hh

configuration:
  writetofiles: False
  path: /tmp/exports
  timefile: lookuptime.prom
  upfile: availability.prom
  interval: 30
```
 * `servers` - set servers which you need to check 
 * `hosts` - hosts to check
 * `configuration`
    * `writetofiles` - if you need use own exporter for example NodeExporter
    * `path` - path to write metrics
    * `timefile`,`upfile` - files for write metrics
    * `interval` - time in seconds after which servers are polled 