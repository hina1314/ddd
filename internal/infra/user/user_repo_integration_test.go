package user_test

import (
	"context"
	"database/sql"
	"os"
	"testing"

	"github.com/google/uuid"
	"github.com/hina1314/ddd/db/model"
	"github.com/hina1314/ddd/internal/domain/user/entity"
	userrepo "github.com/hina1314/ddd/internal/infra/user"
	_ "github.com/jackc/pgx/v5/stdlib"
	"github.com/stretchr/testify/require"
)

func TestRepositorySavesEmailOnlyUser(t *testing.T) {
	dsn := os.Getenv("TEST_DATABASE_URL")
	if dsn == "" {
		t.Skip("TEST_DATABASE_URL is not set")
	}

	db, err := sql.Open("pgx", dsn)
	require.NoError(t, err)
	t.Cleanup(func() { require.NoError(t, db.Close()) })
	require.NoError(t, db.PingContext(context.Background()))

	repository := userrepo.NewUserRepository(model.NewStore(db))
	suffix := uuid.NewString()
	email := "sqlc-" + suffix + "@example.com"
	record, err := entity.NewUser("", email, "user_"+suffix[:12], "test-hash")
	require.NoError(t, err)

	err = repository.Save(context.Background(), record)
	require.NoError(t, err)
	t.Cleanup(func() {
		_, cleanupErr := db.ExecContext(context.Background(), `DELETE FROM "user" WHERE id = $1`, record.ID)
		require.NoError(t, cleanupErr)
	})

	saved, err := repository.GetByEmail(context.Background(), email)
	require.NoError(t, err)
	require.Empty(t, saved.Phone)
	require.Equal(t, email, saved.Email.String())
}
