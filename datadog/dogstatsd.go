package main

import (
	"log"
	"time"

	"github.com/DataDog/datadog-go/v5/statsd"
)

func main() {
	c := NewDataDogStatsDClient()
	i := 0
	for i < 10 {
		c.BumpSum("test.sum", i)
		c.BumpAvg("test.avg", i)
		c.BumpHistogram("test.histogram", i)
		c.BumpTime("test.time")
		i += 1
	}
	c.client.Close()
}

type DatadogStatsClient struct {
	client      *statsd.Client
}

func NewDataDogStatsDClient() *DatadogStatsClient {
	c, err := statsd.New("127.0.0.1:8125", statsd.WithTags([]string{"env:prod", "service:testsvc"}))
	if err != nil {
		log.Fatal(err)
	}
	return &DatadogStatsClient{c}
}

func (c *DatadogStatsClient) BumpAvg(key string, val int) {
	err := c.client.Gauge(key, float64(val), []string{"tag2:value"}, 1)
	log.Printf("Avg error: %v\n", err)
}

func (c *DatadogStatsClient) BumpHistogram(key string, val int) {
	err := c.client.Histogram(key, float64(val), []string{"tag2:value"}, 1)
	log.Printf("Histogram error: %v\n", err)
}

func (c *DatadogStatsClient) BumpSum(key string, val int) {
	err := c.client.Count(key, int64(val), []string{"tag2:value"}, 1)
	log.Printf("Sum error: %v\n", err)
}

func (c *DatadogStatsClient) BumpTime(key string) interface {
	End()
} {
	return timeEnd{c, key, time.Now()}
}

type timeEnd struct {
	dataDogStatsClient *DatadogStatsClient
	key                string
	eventStartTime     time.Time
}

func (n timeEnd) End() {
	err := n.dataDogStatsClient.client.Gauge(n.key, float64(time.Since(n.eventStartTime).Nanoseconds()), []string{"tag2:value"}, 1)
	log.Printf("Time error: %v\n", err)
}