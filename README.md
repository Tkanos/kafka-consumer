# kafka-consumer

With this project I wanted to architecture a kafka consumer in the go-kit philosophy.

### Pre-requesite :
- (mandatory) start a Kafka cluster (I use docker for that https://sookocheff.com/post/kafka/kafka-quick-start/, if you have your own cluster, please change the brokers on the config.toml
- (not mandatory) start a Zipkin server (for the zipkin metrics) ( docker run -d -p 9411:9411 openzipkin/zipkin )

### Run
On main folder type :
``` bash
ENVIRONMENT=DEV go run main.go
```

then in another terminal go to the folder producer-test and type :

```bash
go run main.go
```