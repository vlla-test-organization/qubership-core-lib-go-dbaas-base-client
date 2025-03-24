package rest

// BaseDbParams provides some parameters to customize database
type BaseDbParams struct {
	// This is a prefix of the database name. Prefix depends on the type of the database and
	// it should be less than 27 characters if dbName is not specified.
	NamePrefix string `json:"namePrefix,omitempty"`

	// Additional settings for creating database. There is a possibility to update settings after database creation.
	Settings map[string]interface{} `json:"settings,omitempty"`

	// Specifies the identificator of physical database where a logical database will be created.
	// If it is not specified then logical database will be created in default physical database.
	PhysicalDatabaseId string `json:"physicalDatabaseId,omitempty"`

	//Requested role for database creation. Required for v3. Default is admin.
	Role string `json:"userRole,omitempty"`
}

// CreateDbRequest is a request model for adding database to DBaaS
type CreateDbRequest struct {
	// Struct with database parameters
	BaseDbParams

	// Classifier describes the purpose of database and distinguishes this database from other databases in the same namespace.
	// It contains such keys as dbClassifier, scope, microserviceName, namespace. Setting keys depends on the database type.
	// If database with such classifier exists, then this database will be given away.
	Classifier map[string]interface{} `json:"classifier"`

	// Describes the type of database in which you want to create a database. For example MongoDB or PostgreSQL
	Type string `json:"type"`
}
