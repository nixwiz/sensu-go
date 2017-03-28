package etcd

import (
	"fmt"
	"log"
	"testing"

	"github.com/sensu/sensu-go/backend/store"
	"github.com/sensu/sensu-go/testing/util"
	"github.com/sensu/sensu-go/types"
	"github.com/stretchr/testify/assert"
)

func testWithEtcd(t *testing.T, f func(store.Store)) {
	util.WithTempDir(func(tmpDir string) {
		ports := make([]int, 2)
		err := util.RandomPorts(ports)
		if err != nil {
			log.Panic(err)
		}
		clURL := fmt.Sprintf("http://127.0.0.1:%d", ports[0])
		apURL := fmt.Sprintf("http://127.0.0.1:%d", ports[1])
		initCluster := fmt.Sprintf("default=%s", apURL)

		cfg := NewConfig()
		cfg.StateDir = tmpDir
		cfg.ClientListenURL = clURL
		cfg.PeerListenURL = apURL
		cfg.InitialCluster = initCluster

		e, err := NewEtcd(cfg)
		assert.NoError(t, err)
		if e != nil {
			defer e.Shutdown()
		}

		s, err := e.NewStore()
		assert.NoError(t, err)
		if err != nil {
			assert.FailNow(t, "failed to get store from etcd")
		}

		f(s)
	})
}

func TestEntityStorage(t *testing.T) {
	testWithEtcd(t, func(store store.Store) {
		entity := &types.Entity{
			ID: "0",
		}
		err := store.UpdateEntity(entity)
		assert.NoError(t, err)
		retrieved, err := store.GetEntityByID(entity.ID)
		assert.NoError(t, err)
		assert.Equal(t, entity.ID, retrieved.ID)
		entities, err := store.GetEntities()
		assert.NoError(t, err)
		assert.Equal(t, 1, len(entities))
		assert.Equal(t, entity.ID, entities[0].ID)
		err = store.DeleteEntity(entity)
		assert.NoError(t, err)
		retrieved, err = store.GetEntityByID(entity.ID)
		assert.Nil(t, retrieved)
		assert.NoError(t, err)
	})

}

func TestCheckStorage(t *testing.T) {
	testWithEtcd(t, func(store store.Store) {
		check := &types.Check{
			Name:        "check1",
			Interval:    60,
			Subscribers: []string{"subscription1"},
			Command:     "command1",
		}

		err := store.UpdateCheck(check)
		assert.NoError(t, err)
		retrieved, err := store.GetCheckByName("check1")
		assert.NoError(t, err)
		assert.NotNil(t, retrieved)

		assert.Equal(t, check.Name, retrieved.Name)
		assert.Equal(t, check.Interval, retrieved.Interval)
		assert.Equal(t, check.Subscribers, retrieved.Subscribers)
		assert.Equal(t, check.Command, retrieved.Command)

		checks, err := store.GetChecks()
		assert.NoError(t, err)
		assert.NotEmpty(t, checks)
		assert.Equal(t, 1, len(checks))
	})
}
