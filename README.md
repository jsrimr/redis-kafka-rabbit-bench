# Benchmark for various messaging services

# Docker setup

## install docker
```
sudo snap install docker
```
 - sudo 명령없이 docker 명령어 실행하기 
 ```
 sudo usermod -aG docker $USER
 sudo addgroup --system docker
 sudo adduser $USER docker
 newgrp docker
 sudo snap disable docker
 sudo snap enable docker
 ```

## redis
```
docker run --name benchredis -d redis:5.0.5-alpine
```

## rabbitmq
```
docker run -d --hostname benchrabbitmq --name benchrabbitmq -e RABBITMQ_DEFAULT_USER=guest -e RABBITMQ_DEFAULT_PASS=guest -p 15672:15672 -p 5671:5671 -p 5672:5672 -p 15671:15671 -p 25672:25672 rabbitmq:management
```

## kafka
```
docker-compose up -d
```
docker-compose.yml
```
version: '2'
services:
  zookeeper:
    image: confluentinc/cp-zookeeper:5.3.1
    container_name: zookeper
    environment:
      ZOOKEEPER_CLIENT_PORT: 2181

  kafka:
    image: confluentinc/cp-kafka:5.3.1
    container_name: kafka
    depends_on:
      - zookeeper
    ports:
      - 9092:9092
    environment:
      KAFKA_BROKER_ID: 1
      KAFKA_ZOOKEEPER_CONNECT: zookeeper:2181
      KAFKA_ADVERTISED_LISTENERS: PLAINTEXT://kafka:29092,PLAINTEXT_HOST://localhost:9092
      KAFKA_LISTENER_SECURITY_PROTOCOL_MAP: PLAINTEXT:PLAINTEXT,PLAINTEXT_HOST:PLAINTEXT
      KAFKA_INTER_BROKER_LISTENER_NAME: PLAINTEXT
      KAFKA_OFFSETS_TOPIC_REPLICATION_FACTOR: 1
```

``
