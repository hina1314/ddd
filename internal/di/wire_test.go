package di

import (
	"context"
	"database/sql"
	"database/sql/driver"
	"errors"
	"sync/atomic"
	"testing"

	"github.com/hina1314/ddd/config"
	"github.com/stretchr/testify/require"
)

type cleanupDriver struct {
	closed   atomic.Bool
	verified atomic.Bool
}
type cleanupConn struct{ owner *cleanupDriver }

func (d *cleanupDriver) Open(string) (driver.Conn, error)  { return &cleanupConn{owner: d}, nil }
func (c *cleanupConn) Prepare(string) (driver.Stmt, error) { return nil, errors.New("unused") }
func (c *cleanupConn) Begin() (driver.Tx, error)           { return nil, errors.New("unused") }
func (c *cleanupConn) Close() error                        { c.owner.closed.Store(true); return nil }
func (c *cleanupConn) Ping(context.Context) error          { c.owner.verified.Store(true); return nil }

func TestDependencyAssemblyFailureClosesVerifiedPool(t *testing.T) {
	d := &cleanupDriver{}
	sql.Register(t.Name(), d)
	// The test working directory has no config/i18n: translation assembly fails
	// after the database provider has successfully opened and verified its pool.
	deps, err := NewDependencies(config.Config{
		DBDriver: t.Name(), DBSource: "fixture", DBMaxOpenConns: 4, DBMaxIdleConns: 1,
	})
	require.Error(t, err)
	require.Nil(t, deps)
	require.True(t, d.verified.Load())
	require.True(t, d.closed.Load(), "Wire must clean up earlier providers on a later assembly failure")
}
