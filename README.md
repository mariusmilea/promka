# Export metrics sent via Kafka into Prometheus

Prometheus is a pull based monitoring system, meaning that hosts expose the metrics on o port and the Prometheus server regularly fetches those metrics and updates its database.
In order to expose the metrics on a server, one needs to open a port in their firewall to allow this. In some cases, for sensitive applications, this is not allowed.

To circumvent this way of working, the following can be done:
* send the metrics from the server to a kafka topic
* use promka.py to read the metrics and expose them on the metrics API
* configure prometheus to scrape the newly created metrics API target

# Installation

To test this, we need to have a fully functional kafka setup. For this, we're going to use a docker instance running Kafka.

All from below should be done on the same machine.


For Python, install pykafka and prometheus_client:

`sudo pip install pykafka`

`sudo pip install prometheus_client`

and then run the promka.py script in a SSH session and leave it running:

`python promka.py`


For Golang, install the prometheus and kafka clients: 

`go get -u "github.com/prometheus/client_golang/prometheus"`

`go get -u "github.com/optiopay/kafka"`

and then run the promka.go script in a SSH session and leave it running:

`go run promka.go`


Now, open another SSH session and do the following:

Clone the docker kafka setup from Spotify:

`git clone https://github.com/spotify/docker-kafka.git`

Go into the cloned folder:

`cd docker-kafka`

Run the kafka container:

`docker run -p 2181:2181 -p 9092:9092 --env ADVERTISED_HOST=$(hostname -i) --env ADVERTISED_PORT=9092 spotify/kafka &`

Run ps to get the container ID:

`docker ps`

Connect to the container:

`docker exec -ti b0a163b866ab bash`

Create a topic called test:

`/opt/kafka_2.11-0.8.2.1/bin/kafka-topics.sh --zookeeper localhost:2181 --create --topic test --partitions 1 --replication-factor 1`

**NOTE:** adding multiple partitions per topic will causee un unordered sequence of metrics, because Kafka only preserves the order for a single partition 

List the current topics to see you topic has been created

`/opt/kafka_2.11-0.8.2.1/bin/kafka-topics.sh --zookeeper localhost:2181 --list`

Run the console producer:

`/opt/kafka_2.11-0.8.2.1/bin/kafka-console-producer.sh --topic test --broker-list localhost:9092`

and start typing messages like:

`{"name": "test_kafka_metric", "value": 1}`


Go to http://localhost:8000/metrics and see that your metric "my_kafka_metric" is in there. Type more messages into the console producer and you'll see the values change in the API.
