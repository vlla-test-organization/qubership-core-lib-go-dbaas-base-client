package model

import "github.com/netcracker/qubership-core-lib-go-dbaas-base-client/v3/model/rest"

// LogicalDbProvider allows use different sources as databases providers (for example zookeeper)
type LogicalDbProvider interface {
	GetOrCreateDb(dbType string, classifier map[string]interface{}, params rest.BaseDbParams) (*LogicalDb, error)
	GetConnection(dbType string, classifier map[string]interface{}, params rest.BaseDbParams) (map[string]interface{}, error)
}
