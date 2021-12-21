# Finance
  
This is an implementation of the backend of a finance application written in Go 1.17 created with a microservice architecture in mind.
  
This has been developed with containers in mind and is configured using environment variables. To use deploy this in a Go 1.17 container and pass the following environment variables;  
  
|Name|Value|
|-|-|
|DB_HOST|Hostname for MySQL Database (string)|
|DB_NAME|Database name for MySQL Database (string)|
|DB_PASSWORD|Password for MySQL Database (string)|
|DB_USER|Username for MySQL Database (string)|
|HTTP_HOST|Hostname to listen on (string)|
|HTTP_PORT|Port to listen on (int)|
|OPENSEARCH_HOST|Hostname for Opensearch (string)|
|OPENSEARCH_PORT|Port for Opensearch (int)|
|OPENSEARCH_USER|Username to authenticate Opensearch (string)|
|OPENSEARCH_PASSWORD|Password for Opensearch (string)|
|REDIS_DB|Redis database to use (int)|
|REDIS_HOST|Hostname for Redis cluster (string)|
|REDIS_PASSWORD|Password for Redis cluster (string)|

## To do

### Testing
- [ ] All code needs tests! *IMPORTANT*

### Authentication
- [ ] Implement authentication 

### Logging
- [ ] Loki integration

## Done
- [x] Implement Accounts model  
- [x] Have Redis caching  
- [x] Logging as JSON format for future Loki integration  
- [x] Prometheus metrics exporter  
- [x] Opensearch indexing  