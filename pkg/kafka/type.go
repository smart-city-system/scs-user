// ProducerConfig specific configuration for producers.
package kafka_client

type Config struct {
	Brokers []string
	Topic   string
}

type ProducerConfig struct {
	BatchSize    int
	BatchTimeout int // In milliseconds
	Async        bool
	RequiredAcks int
}

// ConsumerConfig specific configuration for consumers.
type ConsumerConfig struct {
	GroupID        string
	Partition      int
	MinBytes       int
	MaxBytes       int
	CommitInterval int // In milliseconds
	StartOffset    int64
}
