package baserepo

import (
	"context"
	"fmt"

	"go.opentelemetry.io/otel/attribute"
	"go.opentelemetry.io/otel/trace"
	"gorm.io/gorm"
)

type Repository[T any] interface {
	Create(ctx context.Context, entity *T) (*T, error)
	Update(ctx context.Context, entity *T, params any) (*T, error)
	UpdateBy(ctx context.Context, entity *T, condition any, params any) (*T, error)
	Save(ctx context.Context, entity *T) (*T, error)
	GetFirstByID(ctx context.Context, id int) (T, error)
	GetFirstBy(ctx context.Context, condition any, sort string) (*T, error)
	GetAll(ctx context.Context, limit, offset int, sort string) ([]*T, int64, error)
	GetAllBy(ctx context.Context, limit, offset int, condition any, sort string) ([]*T, int64, error)
	Delete(ctx context.Context, entity *T, ID any) error
}

type repository[T any] struct {
	db     *gorm.DB
	tracer trace.Tracer
}

func NewRepository[T any](db *gorm.DB, tracer trace.Tracer) Repository[T] {
	return &repository[T]{db: db, tracer: tracer}
}

func addEvent(span trace.Span, name string, attrs ...attribute.KeyValue) {
	span.AddEvent(name, trace.WithAttributes(attrs...))
}

func (r *repository[T]) Create(ctx context.Context, entity *T) (*T, error) {
	ctx, span := r.tracer.Start(ctx, "repository.Create")
	defer span.End()

	addEvent(span, "Creating entity")

	if err := r.db.WithContext(ctx).Create(entity).Error; err != nil {
		span.RecordError(err)
		return nil, err
	}

	addEvent(span, "Entity created successfully")
	return entity, nil
}

func (r *repository[T]) Update(ctx context.Context, entity *T, params any) (*T, error) {
	ctx, span := r.tracer.Start(ctx, "repository.Update")
	defer span.End()

	addEvent(span, "Updating entity")

	if err := r.db.WithContext(ctx).Model(entity).Updates(params).Error; err != nil {
		span.RecordError(err)
		return nil, err
	}

	addEvent(span, "Entity updated successfully")
	return entity, nil
}

func (r *repository[T]) UpdateBy(ctx context.Context, entity *T, condition any, params any) (*T, error) {
	ctx, span := r.tracer.Start(ctx, "repository.UpdateBy")
	defer span.End()

	addEvent(span, "Updating entity by condition")

	if err := r.db.WithContext(ctx).Model(entity).Where(condition).Updates(params).Error; err != nil {
		span.RecordError(err)
		return nil, err
	}

	addEvent(span, "Entity updated successfully")
	return entity, nil
}

func (r *repository[T]) Save(ctx context.Context, entity *T) (*T, error) {
	ctx, span := r.tracer.Start(ctx, "repository.Save")
	defer span.End()

	addEvent(span, "Saving entity")

	if err := r.db.WithContext(ctx).Save(entity).Error; err != nil {
		span.RecordError(err)
		return nil, err
	}

	addEvent(span, "Entity saved successfully")
	return entity, nil
}

func (r *repository[T]) GetFirstByID(ctx context.Context, id int) (T, error) {
	ctx, span := r.tracer.Start(ctx, "repository.GetFirstByID")
	defer span.End()

	addEvent(span, "Getting entity by ID", attribute.Int("entity.id", id))

	var item T
	if err := r.db.WithContext(ctx).First(&item, id).Error; err != nil {
		span.RecordError(err)
		return item, err
	}

	addEvent(span, "Entity retrieved successfully")
	return item, nil
}

func (r *repository[T]) GetFirstBy(ctx context.Context, condition any, sort string) (*T, error) {
	ctx, span := r.tracer.Start(ctx, "repository.GetFirstBy")
	defer span.End()

	addEvent(span, "Getting first entity by condition")

	var entity T
	if err := r.db.WithContext(ctx).Where(condition).Order(sort).First(&entity).Error; err != nil {
		span.RecordError(err)
		return nil, err
	}

	addEvent(span, "Entity retrieved successfully")
	return &entity, nil
}

func (r *repository[T]) GetAll(ctx context.Context, limit, offset int, sort string) ([]*T, int64, error) {
	ctx, span := r.tracer.Start(ctx, "repository.GetAll")
	defer span.End()

	addEvent(span, "Getting all entities")

	var total int64
	if err := r.db.WithContext(ctx).Model(new(T)).Count(&total).Error; err != nil {
		span.RecordError(err)
		return nil, 0, err
	}

	var items []*T
	if err := r.db.WithContext(ctx).Limit(limit).Offset(offset).Order(sort).Find(&items).Error; err != nil {
		span.RecordError(err)
		return nil, 0, err
	}

	addEvent(span, "Entities retrieved successfully", attribute.Int64("total", total))
	return items, total, nil
}

func (r *repository[T]) GetAllBy(ctx context.Context, limit, offset int, condition any, sort string) ([]*T, int64, error) {
	ctx, span := r.tracer.Start(ctx, "repository.GetAllBy")
	defer span.End()

	addEvent(span, "Getting all entities by condition")

	var total int64
	if err := r.db.WithContext(ctx).Model(new(T)).Where(condition).Count(&total).Error; err != nil {
		span.RecordError(err)
		return nil, 0, err
	}

	var items []*T
	if err := r.db.WithContext(ctx).Where(condition).Limit(limit).Offset(offset).Order(sort).Find(&items).Error; err != nil {
		span.RecordError(err)
		return nil, 0, err
	}

	addEvent(span, "Entities retrieved successfully", attribute.Int64("total", total))
	return items, total, nil
}

func (r *repository[T]) Delete(ctx context.Context, entity *T, ID any) error {
	ctx, span := r.tracer.Start(ctx, "repository.Delete")
	defer span.End()

	addEvent(span, "Deleting entity", attribute.String("entity.id", fmt.Sprintf("%v", ID)))

	if err := r.db.WithContext(ctx).Delete(entity, ID).Error; err != nil {
		span.RecordError(err)
		return err
	}

	addEvent(span, "Entity deleted successfully")
	return nil
}
