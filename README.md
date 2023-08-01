# GO-KAFAK-EXAMPLE
## Library 
- kafka
- zookeeper
- kowl
- postgesql
- prometheus
- node_exporter
- grafana
- jaeger
- gin
- go-kakfa
- go-migrate
- sqlx
- viper
- kakfa-python
- poetry

## API
Build up auth/user system that can register and login.
If user sign up, the system will send a verification code's mail to your mail.

Listen in http://localhost:8000/api

```bash
make run-producer
make run-consumer
make py-consumer
make py-producer
make docker-up
make docker-down
```
## Jaeger
Trace server

http://localhost:16686


## Prometheus
Metrics for monitoring server
http://localhost:9090

## Grafana
Dashboard for visualizing metrics

http://localhost:3001

## Kafka
Using wurstmeister/zookeeper, wurstmeister/kafka

zookeeper: http://localhost:2181
kafka: http://localhost:9092


## Kowl
UI dashboard for kakfa

http://localhost:8080

