package mysql

import (
	"testing"
	"time"

	"github.com/Calmantara/lis-backend/internal/core/configurations"
	"github.com/stretchr/testify/assert"
	"gorm.io/gorm"
)

const (
	pwd    = "mysecretpassword"
	user   = "postgres"
	dbName = "postgres"
	driver = "postgres"
)

func TestNewClient(t *testing.T) {
	t.Run("success connect", func(t *testing.T) {
		env := configurations.Load()
		assert.NotPanics(t, func() {
			NewClient(env.DatabaseMaster)
		})
	})
	t.Run("not connect", func(t *testing.T) {
		assert.Panics(t, func() {
			NewClient(configurations.Database{})
		})
	})

	t.Run("invalid driver", func(t *testing.T) {
		assert.Panics(t, func() {
			NewClient(configurations.Database{
				Host:                    "invalidaddress",
				Password:                pwd,
				Port:                    uint64(1234),
				Username:                user,
				SSL:                     "disable",
				Driver:                  "invalid",
				Name:                    dbName,
				MaxConnection:           1,
				MaxIDleConnection:       1,
				MaxIDleTime:             time.Minute,
				MaxLifeTime:             time.Minute,
				FallbackApplicationName: "hubble-aura-device_test",
			})
		})
	})
	t.Run("panic ping", func(t *testing.T) {
		env := configurations.Load()
		assert.Panics(t, func() {
			cln := NewClient(env.DatabaseMaster)
			db, _ := cln.db.DB()
			db.Close()
			cln.loadGorm(db)
		})
	})
}

func TestMasterSlaveConnection(t *testing.T) {
	t.Run("new connection", func(t *testing.T) {

		var cln *Client
		env := configurations.Load()
		assert.NotPanics(t, func() {
			cln = NewClient(env.DatabaseSlave)
		})

		conn := NewConnection(cln, cln)
		assert.EqualValues(t, cln.db, conn.Master())
		assert.EqualValues(t, cln.db, conn.Slave())
	})
}

func TestBeforeFind(t *testing.T) {
	t.Run("failed zero value", func(t *testing.T) {
		db := &gorm.DB{}
		test := struct {
			Slice  int    `filter:"slice"`
			Sort   int    `sort:"true"`
			Filter string `filter:"filter;>"`
		}{}
		res := beforeFind(db, test)
		assert.EqualValues(t, db, res)
	})
	t.Run("ok", func(t *testing.T) {
		env := configurations.Load()
		cln := NewClient(env.DatabaseMaster)
		test := struct {
			Val int `filter:"slice"`
		}{Val: 10}
		res := beforeFind(cln.db, test)
		assert.NotEqual(t, cln.db, res)
	})
}
func TestGormTransaction(t *testing.T) {
	t.Run("begin transaction", func(t *testing.T) {
		env := configurations.Load()
		cln := NewClient(env.DatabaseMaster)

		ctx := t.Context()
		ctx, txn := begin(ctx, cln.db)

		assert.NotNil(t, txn)
		assert.True(t, txn.Error == nil)
		assert.True(t, txn.Statement.ConnPool != nil)

		// second process
		ctx, txn = begin(ctx, cln.db)

		assert.NotNil(t, txn)
		assert.True(t, txn.Error == nil)
		assert.True(t, txn.Statement.ConnPool != nil)
	})

	t.Run("Get existing transaction", func(t *testing.T) {
		env := configurations.Load()
		cln := NewClient(env.DatabaseMaster)

		ctx := t.Context()
		ctx, txn := begin(ctx, cln.db)

		existingTxn := getTransaction(ctx)
		assert.Equal(t, txn, existingTxn)
	})

	t.Run("commit transaction", func(t *testing.T) {
		env := configurations.Load()
		cln := NewClient(env.DatabaseMaster)

		ctx := t.Context()
		ctx, txn := begin(ctx, cln.db.Begin())

		err := commit(ctx)
		assert.Nil(t, err)
		assert.True(t, txn.Error == nil)
	})

	t.Run("commit without transaction", func(t *testing.T) {
		ctx := t.Context()

		err := commit(ctx)
		assert.Nil(t, err)
	})

	t.Run("rollback transaction", func(t *testing.T) {
		env := configurations.Load()
		cln := NewClient(env.DatabaseMaster)

		ctx := t.Context()
		ctx, txn := begin(ctx, cln.db.Begin())

		err := rollback(ctx)
		assert.Nil(t, err)
		assert.True(t, txn.Error == nil)
	})

	t.Run("rollback without transaction", func(t *testing.T) {
		ctx := t.Context()

		err := rollback(ctx)
		assert.Nil(t, err)
	})
}
