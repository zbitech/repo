package mongodb

import (
	"context"
	"testing"

	"github.com/zbitech/common/pkg/utils"
	"github.com/zbitech/common/pkg/vars"
)

func Test_LoadDatabase(t *testing.T) {

	ctx := context.Background()
	conn := NewMongoDBConnection(vars.MONGODB_URL)
	conn.OpenConnection(ctx)
	defer conn.CloseConnection(ctx)

	LoadDatabase(ctx, conn, true)
}

func Test_GetUsers(t *testing.T) {

	// test with 0 user in db - get array of 0 (no error)
	// test with 1 user in db - get array of 1
	// test with 2+ users in db - get array of 2+
	// test with db unavailable - get db error
	ctx := context.Background()
	conn := NewMongoDBConnection(vars.MONGODB_URL)
	conn.OpenConnection(ctx)
	defer conn.CloseConnection(ctx)

	repo := NewAdminMongoRepository(conn)
	users := repo.GetUsers(ctx)
	if users == nil || len(users) == 0 {
		t.Fatalf("Expected users but got nil")
	}

	t.Logf("Users - %s", utils.MarshalObject(users))
}

func Test_GetUser(t *testing.T) {

	// test with invalid user
	// test with valid user
	// test with db unavailable (wrong url+port)
	ctx := context.Background()
	conn := NewMongoDBConnection(vars.MONGODB_URL)
	conn.OpenConnection(ctx)
	defer conn.CloseConnection(ctx)

	repo := NewAdminMongoRepository(conn)
	user, err := repo.GetUser(ctx, "admin")
	if err != nil {
		t.Fatalf("Expected admin user but got err - %s", err)
	}

	t.Logf("Users - %s", utils.MarshalObject(user))
}

func Test_GetAPIKeys(t *testing.T) {

	ctx := context.Background()
	conn := NewMongoDBConnection(vars.MONGODB_URL)
	conn.OpenConnection(ctx)
	defer conn.CloseConnection(ctx)

	repo := NewAdminMongoRepository(conn)
	keys, err := repo.GetAPIKeys(ctx, "jakinyele")
	if err != nil {
		t.Fatalf("Expected apikeys but got err - %s", err)
	}

	t.Logf("API Keys - %s", keys)
}
func Test_GetAPIKey(t *testing.T) {

	ctx := context.Background()
	conn := NewMongoDBConnection(vars.MONGODB_URL)
	conn.OpenConnection(ctx)
	defer conn.CloseConnection(ctx)

	repo := NewAdminMongoRepository(conn)
	key, err := repo.GetAPIKey(ctx, "9664a36c-58bd-4968-82d1-a8fb4c3502bf")
	if err != nil {
		t.Fatalf("Expected apikeys but got err - %s", err)
	}

	t.Logf("API Keys - %s", utils.MarshalObject(key))
}

func Test_GetUserPolicy(t *testing.T) {

	ctx := context.Background()
	conn := NewMongoDBConnection(vars.MONGODB_URL)
	conn.OpenConnection(ctx)
	defer conn.CloseConnection(ctx)

	repo := NewAdminMongoRepository(conn)
	u_policy, err := repo.GetUserPolicy(ctx, "jakinyele")
	if err != nil {
		t.Fatalf("Expected user policy but got error - %s", err)
	}

	t.Logf("User Policy - %s", utils.MarshalObject(u_policy))
}
