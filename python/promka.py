"""

Small tool that reads metrics from a Kafka topic and exposes them into
the Prometheus metrics API so that they can be scraped.

In this case we only added a Gauge, but a Counter or a Histogram will work
just as well.

Assumptions:
* a Kafka metric is a JSON like {"name": "my_kafka_metric", "value": 1}

More: https://github.com/prometheus/client_python

TODO: daemonize and add logging

"""

from pykafka import KafkaClient
from prometheus_client import start_http_server, Gauge
import json


if __name__ == '__main__':
    # start up the metric API server
    start_http_server(8000)
    # define a Gauge to send to export into the metrics API
    g = Gauge('my_kafka_metric', 'Sample Kafka metric')

    # create a Kafka connectio to a particular topic
    client = KafkaClient("localhost:9092")
    topic = client.topics["test"]
    consumer = topic.get_simple_consumer(consumer_group=b"offset",
                                         auto_commit_enable=True)

    # keep looping through the Kafka messages and as long as we find a valid
    # JSON containing the fields 'name' and 'value', set the defined Gauge
    # to whatever value we get from Kafka
    while True:
        for msg in consumer:
            try:
                metric = json.loads(msg.value)
                if 'name' in metric and 'value' in metric:
                    g.set(metric['value'])
            except ValueError as e:
                # pass if we haven't found a valid JSON metric
                pass
