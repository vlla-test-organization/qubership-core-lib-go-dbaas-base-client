package dbaasbase

import (
	"context"

	"github.com/vlla-test-organization/qubership-core-lib-go-dbaas-base-client/v3/cache"
	"github.com/vlla-test-organization/qubership-core-lib-go-dbaas-base-client/v3/model"
	"github.com/vlla-test-organization/qubership-core-lib-go-dbaas-base-client/v3/model/rest"
)

type DbaaSPool struct {
	poolCache *cache.DbaaSCache
	Client    DbaaSClient
}

func NewDbaaSPool(options ...model.PoolOptions) *DbaaSPool {
	poolCache := &cache.DbaaSCache{LogicalDbCache: make(map[cache.Key]interface{})}
	var clientOpts model.ClientOptions
	if options != nil {
		clientOpts = model.ClientOptions{LogicalDbProviders: options[0].LogicalDbProviders}
	} else {
		clientOpts = model.ClientOptions{}
	}
	client := NewDbaasClient(clientOpts)
	return &DbaaSPool{poolCache: poolCache, Client: client}
}

func (p *DbaaSPool) GetOrCreateDb(ctx context.Context, dbType string, classifier map[string]interface{}, params rest.BaseDbParams) (*model.LogicalDb, error) {
	key := cache.NewKey(dbType, classifier)
	db, err := p.poolCache.Cache(key, func() (interface{}, error) {
		return p.Client.GetOrCreateDb(ctx, dbType, classifier, params)
	})
	if err != nil {
		logger.Errorf("Can neither create db with classifier %+v nor get it from poolCache", classifier)
		return nil, err
	}
	return db.(*model.LogicalDb), nil
}

func (p *DbaaSPool) GetConnection(ctx context.Context, dbType string, classifier map[string]interface{}, params rest.BaseDbParams) (map[string]interface{}, error) {
	connection, err := p.Client.GetConnection(ctx, dbType, classifier, params)
	if err != nil {
		logger.Errorf("Can not get connection to db with classifier %+v", classifier)
		return nil, err
	}
	return connection, nil
}
