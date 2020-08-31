package conn

//"github.com/apache/pulsar/pulsar-client-go/pulsar"

type Producer struct {
	//pub   pulsar.Producer
	//pri   pulsar.Producer
	//depth pulsar.Producer
}

type PulsarConfig struct {
	Url          string
	PublicTopic  string
	PrivateTopic string
	DepthTopic   string
}
