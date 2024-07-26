package db

import (
    "os"
    "testing"
    "github.com/stretchr/testify/assert"
)

func TestConnection(t *testing.T) {
    conn, err := NewConnection()
    assert.NoError(t, err)
    assert.NotNil(t, conn)
}

func TestGetDBConnection(t *testing.T) {
    conn, err := GetDBConnection()
    assert.NoError(t, err)
    assert.NotNil(t, conn)
}

func TestGetDatabaseConnectionStringNotTest(t *testing.T) {
    originalEnv := os.Getenv("ENV")

    os.Setenv("ENV", "dev")
    os.Setenv("DATABASE_URL", "foo")

    connStr := GetDatabaseConnectionString()

    os.Setenv("ENV", originalEnv)

    assert.Equal(t, "foo", connStr)
}

func TestGetDatabaseConnectionStringEmpty(t *testing.T) {
    originalEnv := os.Getenv("ENV")

    os.Setenv("ENV", "dev")
    os.Setenv("DATABASE_URL", "")

    connStr := GetDatabaseConnectionString()

    os.Setenv("ENV", originalEnv)

    assert.Equal(t, "", connStr)
}
