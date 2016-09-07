package main

import (
        "log"
        "flag"
        "net/http"
        "encoding/json"
        "github.com/prometheus/client_golang/prometheus"
        "github.com/optiopay/kafka"
)

// define the JSON structure for the metrics received via Kafka
type JsonMetric struct {
    Metric   string     `json:"metric"`
    Value    float64    `json:"value"`
}

// define the listening address for the HTTP server
var addr = flag.String("listen-address", ":8001", "The address to listen on for HTTP requests.")

// define a kafka topic and partition
const (
    topic     = "test"
    partition = 0
)

var kafkaAddrs = []string{"localhost:9092"}

// define a simple prometheus gauge
var (
    my_kafka_metric = prometheus.NewGauge(prometheus.GaugeOpts{
        Name: "my_kafka_metric",
        Help: "Sample Kafka metric.",
    })
)

func init() {
    // a metric needs to be registered before being exposed
    prometheus.MustRegister(my_kafka_metric)
}

// read messages from kafka and set the prometheus gauge accordingly
func get_metrics(broker kafka.Client) {
    conf := kafka.NewConsumerConf(topic, partition)
    conf.StartOffset = kafka.StartOffsetNewest
    res := JsonMetric{}
    consumer, err := broker.Consumer(conf)
    if err != nil {
        log.Fatalf("cannot create kafka consumer for %s:%d: %s", topic, partition, err)
    }

    for {
        msg, err := consumer.Consume()
        if err != nil {
            if err != kafka.ErrNoData {
                log.Printf("cannot consume %q topic message: %s", topic, err)
            }
            break
        }
        json.Unmarshal([]byte(msg.Value), &res)
        my_kafka_metric.Set(res.Value)
    }
    log.Print("consumer quit")
}


func main() {
    broker, err := kafka.Dial(kafkaAddrs, kafka.NewBrokerConf("test-client"))
    if err != nil {
        log.Fatalf("cannot connect to kafka cluster: %s", err)
    }
    defer broker.Close()

    go get_metrics(broker)
    
    // expose metrics
    flag.Parse()
    http.Handle("/metrics", prometheus.Handler())
    http.ListenAndServe(*addr, nil)


}
