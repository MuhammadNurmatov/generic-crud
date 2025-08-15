package example

import (
	"baserepo/repository"
	"context"

	"go.opentelemetry.io/otel"
	"go.opentelemetry.io/otel/trace"
	"gorm.io/gorm"
)

type UseRepositoryExample interface {
	repository.Repository[UserExample]
	FindEmail(ctx context.Context, email string) (string, error)
}

type useRepositoryExample struct {
	repository.Repository[UserExample]
	db     *gorm.DB
	tracer trace.Tracer
}

func NewUseRepositoryExample(db *gorm.DB) UseRepositoryExample {
	tracer := otel.Tracer("use-repository-example")
	generic := repository.NewRepository[UserExample](db, tracer)
	return &useRepositoryExample{
		Repository: generic,
		db:         db,
		tracer:     tracer,
	}
}

func (u *useRepositoryExample) FindEmail(ctx context.Context, email string) (string, error) {
	ctx, span := u.tracer.Start(ctx, "example.FindEmail")
	defer span.End()

	//Todo ..

	total := "test@test"
	return total, nil
}
