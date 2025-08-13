package repository

import (
	"context"
	"testing"

	"github.com/stretchr/testify/assert"
	"go.opentelemetry.io/otel"
	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
)

type UserTest struct {
	ID   uint64 `gorm:"primary_key"`
	Name string `gorm:"size:200;not null;check:name <> ''"`
}

func initDB(t *testing.T) (*gorm.DB, error) {
	db, err := gorm.Open(sqlite.Open(":memory:"), &gorm.Config{})
	assert.NoError(t, err)
	_ = db.AutoMigrate(&UserTest{})

	t.Cleanup(func() {
		sqlDB, err := db.DB()
		assert.NoError(t, err)
		sqlDB.Close()
	})

	return db, err
}

func TestRepository_Create(t *testing.T) {

	db, err := initDB(t)
	assert.NoError(t, err)

	tracer := otel.Tracer("Unit test repository")
	repo := repository[UserTest]{db, tracer}

	testTable := []struct {
		testName    string
		user        UserTest
		exceptName  string
		exceptError bool
	}{
		{testName: "error", user: UserTest{Name: ""}, exceptName: "", exceptError: true},
		{testName: "success", user: UserTest{Name: "Alice"}, exceptName: "Alice", exceptError: false},
	}

	for _, tt := range testTable {
		t.Run(tt.exceptName, func(t *testing.T) {
			result, err := repo.Create(context.Background(), &tt.user)

			if tt.exceptError {
				assert.Error(t, err)
				assert.Nil(t, result)
				return
			}

			assert.NoError(t, err)
			assert.NotNil(t, result)
			assert.Equal(t, tt.exceptName, result.Name)
		})
	}

}

func TestRepository_Update(t *testing.T) {
	db, err := initDB(t)
	assert.NoError(t, err)

	tracer := otel.Tracer("Unit test repository")
	repo := repository[UserTest]{db, tracer}

	initialUser := UserTest{Name: "Alice"}
	_, err = repo.Create(context.Background(), &initialUser)
	assert.NoError(t, err)

	testTable := []struct {
		testName    string
		user        UserTest
		params      UserTest
		exceptName  string
		exceptError bool
	}{
		{testName: "success", user: initialUser, params: UserTest{Name: "Update Alias"}, exceptName: "Update Alias", exceptError: false},
		{testName: "error", user: UserTest{Name: "Alias"}, params: UserTest{Name: ""}, exceptName: "", exceptError: true},
	}

	for _, tt := range testTable {
		t.Run(tt.testName, func(t *testing.T) {
			result, err := repo.Update(context.Background(), &tt.user, tt.params)

			if tt.exceptError {
				assert.Error(t, err)
				assert.Nil(t, result)
				return
			}

			assert.NoError(t, err)
			assert.NotNil(t, result)
			assert.Equal(t, tt.exceptName, result.Name)
		})
	}
}
