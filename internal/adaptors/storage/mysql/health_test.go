package mysql

import (
	"context"
	"testing"

	"github.com/Calmantara/lis-backend/internal/core/configurations"
	"github.com/stretchr/testify/assert"
)

func TestHealthRepo_CheckMaster_Integration(t *testing.T) {
	// Load the configuration
	config := configurations.Load()
	// Create a new database client for the master
	masterClient := NewClient(config.DatabaseMaster)
	slaveClient := NewClient(config.DatabaseSlave)
	// Create a new health repository
	healthRepo := NewHealthRepo(NewConnection(masterClient, slaveClient))
	// Call the CheckMaster method
	err := healthRepo.CheckMaster(context.Background())
	// Assert that no error occurred
	assert.NoError(t, err)
}

func TestHealthRepo_CheckMaster_Error_Integration(t *testing.T) {
	// Load the configuration
	config := configurations.Load()
	// Create a new database client for the master
	masterClient := NewClient(config.DatabaseMaster)
	slaveClient := NewClient(config.DatabaseSlave)
	// disconnect
	db, _ := masterClient.connect()
	db.Close()
	// Create a new health repository
	healthRepo := NewHealthRepo(NewConnection(masterClient, slaveClient))
	// Call the CheckMaster method
	err := healthRepo.CheckMaster(context.Background())
	// Assert that no error occurred
	assert.Error(t, err)
}

func TestHealthRepo_CheckSlave_Integration(t *testing.T) {
	// Load the configuration
	config := configurations.Load()
	// Create a new database client for the slave
	masterClient := NewClient(config.DatabaseMaster)
	slaveClient := NewClient(config.DatabaseSlave)
	// Create a new health repository
	healthRepo := NewHealthRepo(NewConnection(masterClient, slaveClient))
	// Call the CheckSlave method
	err := healthRepo.CheckSlave(context.Background())
	// Assert that no error occurred
	assert.NoError(t, err)
}

func TestHealthRepo_CheckSlave_Error_Integration(t *testing.T) {
	// Load the configuration
	config := configurations.Load()
	// Create a new database client for the slave
	masterClient := NewClient(config.DatabaseMaster)
	slaveClient := NewClient(config.DatabaseSlave)
	// disconnect
	db, _ := slaveClient.connect()
	db.Close()
	// Create a new health repository
	healthRepo := NewHealthRepo(NewConnection(masterClient, slaveClient))
	// Call the CheckSlave method
	err := healthRepo.CheckSlave(context.Background())
	// Assert that no error occurred
	assert.Error(t, err)
}
