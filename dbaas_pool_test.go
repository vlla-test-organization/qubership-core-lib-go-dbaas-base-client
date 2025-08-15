package dbaasbase

import (
	"context"
	"encoding/json"
	"errors"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"github.com/vlla-test-organization/qubership-core-lib-go-dbaas-base-client/v3/cache"
	"github.com/vlla-test-organization/qubership-core-lib-go-dbaas-base-client/v3/model"
	"github.com/vlla-test-organization/qubership-core-lib-go-dbaas-base-client/v3/model/rest"
)

const DB_TYPE = "dbType"

var testDbaasClient = new(mockDbaasClient)

type mockDbaasClient struct {
	mock.Mock
}

func (m *mockDbaasClient) GetOrCreateDb(ctx context.Context, dbType string, classifier map[string]interface{}, params rest.BaseDbParams) (*model.LogicalDb, error) {
	args := m.Called(dbType, classifier, params)
	return args.Get(0).(*model.LogicalDb), args.Error(1)
}

func (m *mockDbaasClient) GetConnection(ctx context.Context, dbType string, classifier map[string]interface{}, params rest.BaseDbParams) (map[string]interface{}, error) {
	args := m.Called(dbType, classifier, params)
	return args.Get(0).(map[string]interface{}), args.Error(1)
}

func TestNewDbaasPool_WithoutOptions(t *testing.T) {
	dbaasPool := NewDbaaSPool()
	assert.NotNil(t, dbaasPool)
}

func TestNewDbaaSPool_WithOptions(t *testing.T) {
	options := model.PoolOptions{LogicalDbProviders: []model.LogicalDbProvider{testCorrectLogicalDbProvider{}}}
	dbaasPool := NewDbaaSPool(options)
	assert.NotNil(t, dbaasPool)
	assert.NotNil(t, dbaasPool.Client)
}

func TestDbaaSPool_GetConnection(t *testing.T) {
	params := rest.BaseDbParams{}
	response := map[string]interface{}{"body": "success"}
	classifier := map[string]interface{}{"key": "value"}
	testDbaasClient.On("GetConnection", DB_TYPE, classifier, params).Return(response, nil).Once()

	dbaasPool := NewDbaaSPool()
	dbaasPool.Client = testDbaasClient
	actualResponse, _ := dbaasPool.GetConnection(context.Background(), DB_TYPE, classifier, params)
	assert.Equal(t, response, actualResponse)
}

func TestDbaaSPool_GetConnectionWithError(t *testing.T) {
	params := rest.BaseDbParams{}
	classifier := map[string]interface{}{"key": "value"}
	testDbaasClient.On("GetConnection", DB_TYPE, classifier, params).Return(make(map[string]interface{}), errors.New("error during acquiring connection")).Once()

	dbaasPool := NewDbaaSPool()
	dbaasPool.Client = testDbaasClient
	if _, err := dbaasPool.GetConnection(context.Background(), DB_TYPE, classifier, params); assert.Error(t, err) {
		assert.Contains(t, err.Error(), "error during acquiring connection")
	}
}

func TestDbaaSPool_CreateOrGetDatabase(t *testing.T) {
	classifier := map[string]interface{}{"key": "value"}
	params := rest.BaseDbParams{}
	testLogicalDb := &model.LogicalDb{Id: "1", Classifier: classifier, Type: DB_TYPE}
	testDbaasClient.On("GetOrCreateDb", DB_TYPE, classifier, params).Return(testLogicalDb, nil).Once()

	dbaasPool := NewDbaaSPool()
	dbaasPool.Client = testDbaasClient
	actualLogicalDb, _ := dbaasPool.GetOrCreateDb(context.Background(), DB_TYPE, classifier, params)
	assert.Equal(t, testLogicalDb, actualLogicalDb)
}

func TestDbaaSPool_CreateOrGetDatabaseFromCache(t *testing.T) {
	classifier := map[string]interface{}{"key": "value"}
	params := rest.BaseDbParams{}
	testLogicalDb := &model.LogicalDb{Id: "1", Classifier: classifier, Type: DB_TYPE}
	testDbaasClient.On("GetOrCreateDb", DB_TYPE, classifier, params).Return(testLogicalDb, nil).Once()

	dbaasPool := NewDbaaSPool()
	dbaasPool.Client = testDbaasClient
	assert.False(t, isCacheContainsDb(&dbaasPool.poolCache.LogicalDbCache, DB_TYPE, classifier))
	actualLogicalDb, _ := dbaasPool.GetOrCreateDb(context.Background(), DB_TYPE, classifier, params)
	assert.Equal(t, testLogicalDb, actualLogicalDb)

	assert.True(t, isCacheContainsDb(&dbaasPool.poolCache.LogicalDbCache, DB_TYPE, classifier))
	actualLogicalDbFromCache, _ := dbaasPool.GetOrCreateDb(context.Background(), DB_TYPE, classifier, params)
	assert.Equal(t, testLogicalDb, actualLogicalDbFromCache)
}

func TestDbaaSPool_CreateOrGetDatabaseReturnError(t *testing.T) {
	classifier := map[string]interface{}{"key": "value"}
	params := rest.BaseDbParams{}
	testLogicalDb := &model.LogicalDb{}
	testDbaasClient.On("GetOrCreateDb", DB_TYPE, classifier, params).Return(testLogicalDb, errors.New("error during database creation")).Once()

	dbaasPool := NewDbaaSPool()
	dbaasPool.Client = testDbaasClient
	if _, err := dbaasPool.GetOrCreateDb(context.Background(), DB_TYPE, classifier, params); assert.Error(t, err) {
		assert.Contains(t, err.Error(), "error during database creation")
	}
}

func isCacheContainsDb(localCache *map[cache.Key]interface{}, dbType string, classifier map[string]interface{}) bool {
	classifierString, _ := json.Marshal(classifier)
	key := cache.Key{
		DbType:     dbType,
		Classifier: string(classifierString),
	}
	_, ok := (*localCache)[key]
	return ok
}
