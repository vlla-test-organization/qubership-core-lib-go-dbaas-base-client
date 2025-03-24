package model

// LogicalDb is a way to store information about databases locally
type LogicalDb struct {
	// A unique identifier of the document in the database.
	Id string `json:"id"`

	// Classifier describes the purpose of database and distinguishes this database from other databases in the same namespace.
	// It contains such keys as dbClassifier, scope, microserviceName, namespace. Setting keys depends on the database type.
	// If database with such classifier exists, then this database will be given away.
	Classifier map[string]interface{} `json:"classifier"`

	// This is an information about connection to database. It contains such keys as url, authDbName, username, password, port, host.
	//Setting keys depends on the database type.
	ConnectionProperties map[string]interface{} `json:"connectionProperties"`

	// Namespace where database is placed.
	Namespace string `json:"namespace"`

	// Name of the database
	Name string `json:"name"`

	// Type of database, for example PostgreSQL or MongoDB
	Type string `json:"type"`

	// Additional settings for creating a database.
	Settings map[string]interface{} `json:"settings"`
}

// СlientOptions are options for dbaas Client creation
type СlientOptions struct {
	// LogicalDbProviders stores list of available logical db providers
	LogicalDbProviders []LogicalDbProvider
}

// PoolOptions are options for connection pool configuring
type PoolOptions struct {
	// LogicalDbProviders stores list of available logical db providers
	LogicalDbProviders []LogicalDbProvider
}
