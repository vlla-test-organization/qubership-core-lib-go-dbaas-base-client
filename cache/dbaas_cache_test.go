package cache

import (
	"encoding/json"
	"errors"
	"fmt"
	"strconv"
	"sync"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"github.com/vlla-test-organization/qubership-core-lib-go-dbaas-base-client/v3/model"
)

const dbType = "type"

type testDiscriminator struct {
	RoReplica string
}

func (d *testDiscriminator) GetValue() string {
	return d.RoReplica
}

func TestDbaaSCache_AddNewValue(t *testing.T) {
	classifier := map[string]interface{}{"key": "value"}
	classifierAsString, _ := json.Marshal(classifier)
	mapKey := Key{
		DbType:     dbType,
		Classifier: string(classifierAsString),
	}

	cache := DbaaSCache{LogicalDbCache: make(map[Key]interface{})}

	db, err := cache.Cache(mapKey, func() (interface{}, error) {
		return &model.LogicalDb{Id: "123"}, nil
	})
	assert.Nil(t, err)
	val, ok := cache.LogicalDbCache[mapKey]
	assert.True(t, ok)
	assert.Equal(t, db, val)
}

func TestDbaaSCache_GetValueFromCache(t *testing.T) {
	classifier := map[string]interface{}{"key": "value"}
	classifierAsString, _ := json.Marshal(classifier)
	mapKey := Key{
		DbType:     dbType,
		Classifier: string(classifierAsString),
	}
	initialLogicalDb := &model.LogicalDb{Id: "111"}
	cache := DbaaSCache{LogicalDbCache: make(map[Key]interface{})}
	key := NewKey(dbType, classifier)
	cache.LogicalDbCache[mapKey] = initialLogicalDb
	db, err := cache.Cache(key, func() (interface{}, error) {
		return &model.LogicalDb{Id: "123"}, nil
	})
	assert.Nil(t, err)
	assert.Equal(t, initialLogicalDb, db)
}

func TestDbaaSCache_GetValueFromCacheWithDiscriminator(t *testing.T) {
	discriminator := testDiscriminator{
		RoReplica: "roAccess:true",
	}
	classifier := map[string]interface{}{"key": "value"}
	classifierAsString, _ := json.Marshal(classifier)
	mapKey := Key{
		DbType:        dbType,
		Classifier:    string(classifierAsString),
		discriminator: discriminator.GetValue(),
	}
	initialLogicalDb := &model.LogicalDb{Id: "111"}
	cache := DbaaSCache{LogicalDbCache: make(map[Key]interface{})}
	key := NewKeyWithDiscriminator(dbType, classifier, &discriminator)
	cache.LogicalDbCache[mapKey] = initialLogicalDb
	db, err := cache.Cache(key, func() (interface{}, error) {
		return &model.LogicalDb{Id: "123"}, nil
	})
	assert.Nil(t, err)
	assert.Equal(t, initialLogicalDb, db)
}

func TestDbaaSCache_CacheReturnsError(t *testing.T) {
	classifier := map[string]interface{}{"key": "value"}

	cache := DbaaSCache{}
	key := NewKey(dbType, classifier)
	if _, err := cache.Cache(key, func() (interface{}, error) {
		return nil, errors.New("error during computing")
	}); assert.Error(t, err) {
		assert.Contains(t, err.Error(), "error during computing")
	}
}

func TestConcurrentReadWrite(t *testing.T) {
	key := NewKey(dbType, map[string]interface{}{"key": "value"})
	cache := DbaaSCache{LogicalDbCache: make(map[Key]interface{})}
	count := 10
	var wg sync.WaitGroup
	wg.Add(count)
	for i := 0; i < count; i++ {
		go func() {
			result, err := cache.Cache(key, func() (interface{}, error) {
				return "calculated", nil
			})
			assert.Nil(t, err)
			assert.Equal(t, "calculated", result.(string))
			wg.Done()
		}()
	}
	expired := waitWithTimeout(&wg, 5*time.Second)
	assert.False(t, expired, "timed out to wait for %d Cache() successful invocations to happen", count)

	result, err := cache.Cache(key, func() (interface{}, error) {
		return nil, fmt.Errorf("create function should not be executed at this moment")
	})
	assert.Nil(t, err)
	assert.Equal(t, "calculated", result.(string))
}

func TestConcurrentReadWriteAdnDelete(t *testing.T) {
	cache := DbaaSCache{LogicalDbCache: make(map[Key]interface{})}
	count := 10
	getKey := func(i int) Key {
		return NewKey(dbType, map[string]interface{}{"key": "value-" + strconv.Itoa(i)})
	}
	var wgReadWrite sync.WaitGroup
	var wgDelete sync.WaitGroup
	wgReadWrite.Add(count)
	wgDelete.Add(count)
	// readWrite
	func() {
		for i := 0; i < count; i++ {
			key := getKey(i)
			go func() {
				result, err := cache.Cache(key, func() (interface{}, error) {
					return "calculated", nil
				})
				assert.Nil(t, err)
				assert.Equal(t, "calculated", result.(string))
				wgReadWrite.Done()
			}()
		}
	}()
	expired := waitWithTimeout(&wgReadWrite, 5*time.Second)
	assert.False(t, expired, "timed out to wait for %d Cache() successful invocations to happen", count)
	// delete
	func() {
		for i := 0; i < count; i++ {
			key := getKey(i)
			go func() {
				cache.Delete(key)
				wgDelete.Done()
			}()
		}
	}()
	expired = waitWithTimeout(&wgDelete, 5*time.Second)
	assert.False(t, expired, "timed out to wait for %d Delete() successful invocations to happen", count)

	assert.Equal(t, 0, len(cache.LogicalDbCache))
}

func waitWithTimeout(wg *sync.WaitGroup, timeout time.Duration) bool {
	c := make(chan struct{})
	go func() {
		defer close(c)
		wg.Wait()
	}()
	select {
	case <-c:
		return false // completed normally
	case <-time.After(timeout):
		return true // timed out
	}
}
