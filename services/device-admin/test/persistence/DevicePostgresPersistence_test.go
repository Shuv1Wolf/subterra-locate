package test_persistence

import (
	"context"
	"os"
	"testing"

	"github.com/Shuv1Wolf/subterra-locate/services/device-admin/persistence"
	cconf "github.com/pip-services4/pip-services4-go/pip-services4-components-go/config"
)

type DevicePostgresPersistenceTest struct {
	persistence *persistence.DevicePostgresPersistence
	fixture     *DevicePersistenceFixture
}

func newDevicePostgresPersistenceTest() *DevicePostgresPersistenceTest {
	postgresUri := os.Getenv("POSTGRES_URI")
	postgresHost := os.Getenv("POSTGRES_HOST")
	if postgresHost == "" {
		postgresHost = "localhost"
	}

	postgresPort := os.Getenv("POSTGRES_PORT")
	if postgresPort == "" {
		postgresPort = "5432"
	}

	postgresDatabase := os.Getenv("POSTGRES_DB")
	if postgresDatabase == "" {
		postgresDatabase = "postgres"
	}

	postgresUser := os.Getenv("POSTGRES_USER")
	if postgresUser == "" {
		postgresUser = "postgres"
	}
	postgresPassword := os.Getenv("POSTGRES_PASSWORD")
	if postgresPassword == "" {
		postgresPassword = "postgres"
	}

	if postgresUri == "" && postgresHost == "" {
		panic("Connection params not set")
	}

	dbConfig := cconf.NewConfigParamsFromTuples(
		"connection.uri", postgresUri,
		"connection.host", postgresHost,
		"connection.port", postgresPort,
		"connection.database", postgresDatabase,
		"credential.username", postgresUser,
		"credential.password", postgresPassword,
		"schema", "public",
	)

	persistence := persistence.NewDevicePostgresPersistence()
	persistence.Configure(context.Background(), dbConfig)

	fixture := NewDevicePersistenceFixture(persistence)

	return &DevicePostgresPersistenceTest{
		persistence: persistence,
		fixture:     fixture,
	}
}

func (c *DevicePostgresPersistenceTest) setup(t *testing.T) {
	err := c.persistence.Open(context.Background())
	if err != nil {
		t.Error("Failed to open persistence", err)
	}

	err = c.persistence.Clear(context.Background())
	if err != nil {
		t.Error("Failed to clear persistence", err)
	}
}

func (c *DevicePostgresPersistenceTest) teardown(t *testing.T) {
	err := c.persistence.Close(context.Background())
	if err != nil {
		t.Error("Failed to close persistence", err)
	}
}

func TestDevicePostgresPersistence(t *testing.T) {
	c := newDevicePostgresPersistenceTest()
	if c == nil {
		return
	}

	c.setup(t)
	t.Run("CRUD Operations", c.fixture.TestCrudOperations)
	c.teardown(t)

	c.setup(t)
	t.Run("Get With Filters", c.fixture.TestGetWithFilters)
	c.teardown(t)
}
