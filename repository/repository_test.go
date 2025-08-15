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
		{testName: "success", user: UserTest{Name: "Tester"}, exceptName: "Tester", exceptError: false},
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

	initialUser := UserTest{Name: "Tester"}
	_, err = repo.Create(context.Background(), &initialUser)
	assert.NoError(t, err)

	testTable := []struct {
		testName    string
		user        *UserTest
		params      any
		exceptName  string
		exceptError bool
	}{
		{testName: "success", user: &initialUser, params: UserTest{Name: "Update Alias"}, exceptName: "Update Alias", exceptError: false},
		{testName: "error", user: nil, params: "error tesst", exceptName: "", exceptError: true},
	}

	for _, tt := range testTable {
		t.Run(tt.testName, func(t *testing.T) {
			result, err := repo.Update(context.Background(), tt.user, tt.params)

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

func TestRepository_Save(t *testing.T) {
	db, err := initDB(t)
	assert.NoError(t, err)

	tracer := otel.Tracer("Unit test repository")
	repo := repository[UserTest]{db, tracer}

	assert.NoError(t, err)

	testTable := []struct {
		testName    string
		user        *UserTest
		exceptName  string
		exceptError bool
	}{
		{testName: "success", user: &UserTest{Name: "Tester"}, exceptName: "Tester", exceptError: false},
		{testName: "error", user: nil, exceptName: "", exceptError: true},
	}

	for _, tt := range testTable {
		t.Run(tt.testName, func(t *testing.T) {
			result, err := repo.Save(context.Background(), tt.user)

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

func TestRepository_GetFirstByID(t *testing.T) {
	db, err := initDB(t)
	assert.NoError(t, err)

	tracer := otel.Tracer("Unit test repository")
	repo := repository[UserTest]{db, tracer}

	initialUser := UserTest{Name: "Tester"}
	_, err = repo.Create(context.Background(), &initialUser)
	assert.NoError(t, err)

	testTable := []struct {
		testName    string
		ID          any
		exceptName  string
		exceptError bool
	}{
		{testName: "success", ID: initialUser.ID, exceptName: "Tester", exceptError: false},
		{testName: "error", ID: -1, exceptName: "", exceptError: true},
	}

	for _, tt := range testTable {
		t.Run(tt.testName, func(t *testing.T) {
			result, err := repo.GetFirstByID(context.Background(), tt.ID)

			if tt.exceptError {
				assert.Error(t, err)
				assert.Empty(t, result)
				return
			}

			assert.NoError(t, err)
			assert.Equal(t, tt.exceptName, result.Name)
		})
	}
}

func TestRepository_GetFirstBy(t *testing.T) {
	db, err := initDB(t)
	assert.NoError(t, err)

	tracer := otel.Tracer("Unit test repository")
	repo := repository[UserTest]{db, tracer}

	initialUser := UserTest{Name: "Tester"}
	_, err = repo.Create(context.Background(), &initialUser)
	assert.NoError(t, err)

	testTable := []struct {
		testName    string
		entity      *UserTest
		condition   any
		sort        string
		exceptName  string
		exceptError bool
	}{
		{testName: "success", entity: &UserTest{}, condition: initialUser, sort: "ID desc", exceptName: "Tester", exceptError: false},
		{testName: "error", entity: nil, condition: nil, sort: "test desc", exceptName: "test", exceptError: true},
	}

	for _, tt := range testTable {
		t.Run(tt.testName, func(t *testing.T) {
			result, err := repo.GetFirstBy(context.Background(), tt.condition, tt.sort)

			if tt.exceptError {
				assert.Error(t, err)
				assert.Empty(t, result)
				return
			}

			assert.NoError(t, err)
			assert.NotNil(t, result)
			assert.Equal(t, tt.exceptName, result.Name)
		})
	}
}

func TestRepository_GetAll(t *testing.T) {
	db, err := initDB(t)
	assert.NoError(t, err)

	tracer := otel.Tracer("Unit test repository")
	repo := repository[UserTest]{db, tracer}

	initialUser := UserTest{Name: "Tester"}
	_, err = repo.Create(context.Background(), &initialUser)
	assert.NoError(t, err)

	testTable := []struct {
		testName    string
		sort        string
		exceptTotal int64
		exceptError bool
	}{
		{testName: "success", sort: "ID desc", exceptTotal: 1, exceptError: false},
		{testName: "error", sort: "Code desc", exceptTotal: 0, exceptError: true},
	}

	for _, tt := range testTable {
		t.Run(tt.testName, func(t *testing.T) {
			result, total, err := repo.GetAll(context.Background(), 10, 10, tt.sort)

			if tt.exceptError {
				assert.Error(t, err)
				assert.Nil(t, result)
				return
			}

			assert.NoError(t, err)
			assert.NotNil(t, result)
			assert.Equal(t, tt.exceptTotal, total)

		})
	}

}

func TestRepository_GetAllBy(t *testing.T) {
	db, err := initDB(t)
	assert.NoError(t, err)

	tracer := otel.Tracer("Unit test repository")
	repo := repository[UserTest]{db, tracer}

	initialUser := UserTest{Name: "Tester"}
	_, err = repo.Create(context.Background(), &initialUser)
	assert.NoError(t, err)

	testTable := []struct {
		testName    string
		condition   any
		sort        string
		exceptTotal int64
		exceptError bool
	}{
		{testName: "success", condition: UserTest{ID: initialUser.ID}, sort: "ID desc", exceptTotal: 1, exceptError: false},
		{testName: "error", condition: UserTest{ID: initialUser.ID}, sort: "Code desc", exceptTotal: 0, exceptError: true},
	}

	for _, tt := range testTable {
		t.Run(tt.testName, func(t *testing.T) {
			result, total, err := repo.GetAll(context.Background(), 10, 10, tt.sort)

			if tt.exceptError {
				assert.Error(t, err)
				assert.Nil(t, result)
				return
			}

			assert.NoError(t, err)
			assert.NotNil(t, result)
			assert.Equal(t, tt.exceptTotal, total)

		})
	}
}

func TestRepository_Delete(t *testing.T) {
	db, err := initDB(t)
	assert.NoError(t, err)

	tracer := otel.Tracer("Unit test repository")
	repo := repository[UserTest]{db, tracer}

	initialUser := UserTest{Name: "Tester"}
	_, err = repo.Create(context.Background(), &initialUser)
	assert.NoError(t, err)

	testTable := []struct {
		testName string
		ID       uint64
	}{
		{testName: "success", ID: initialUser.ID},
		{testName: "error", ID: 2},
	}

	for _, tt := range testTable {
		t.Run(tt.testName, func(t *testing.T) {
			err = repo.Delete(context.Background(), tt.ID)
			assert.Nil(t, err)
		})
	}
}
