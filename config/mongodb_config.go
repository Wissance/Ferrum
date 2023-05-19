package config

const (
	OperationTimeout       DataSourceConnOption = "timeoutMS"
	ConnectionTimeout      DataSourceConnOption = "connectTimeoutMS"
	ConnectionsPool        DataSourceConnOption = "maxPoolSize"
	ReplicaSet             DataSourceConnOption = "replicaSet"
	MaxIdleTime            DataSourceConnOption = "maxIdleTimeMS"
	SocketTimeout          DataSourceConnOption = "socketTimeoutMS"
	ServerSelectionTimeout DataSourceConnOption = "serverSelectionTimeoutMS"
	HeartbeatFrequency     DataSourceConnOption = "heartbeatFrequencyMS"
	Tls                    DataSourceConnOption = "tls"
	WriteConcern           DataSourceConnOption = "w"
	DirectConnection       DataSourceConnOption = "directConnection"
)
