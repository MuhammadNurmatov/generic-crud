package baserepo

import (
	"context"

	"gorm.io/gorm"
)

type Repository[T any] interface {
}

type repository[T any] struct {
	db *gorm.DB
}

func NewRepository[T any](db *gorm.DB) Repository[T] {
	return &repository[T]{db: db}
}

func (r *repository[T]) Create(ctx context.Context, entity *T) (*T, error) {
	err := r.db.WithContext(ctx).Create(&entity).Error
	if err != nil {
		return nil, err
	}

	return entity, err
}

func (r *repository[T]) Update(ctx context.Context, entity *T, params any) (*T, error) {
	err := r.db.WithContext(ctx).Model(&entity).Updates(params).Error
	if err != nil {
		return nil, err
	}

	return entity, err
}

func (r *repository[T]) UpdateBy(ctx context.Context, entity *T, condition any, params any) (*T, error) {
	err := r.db.WithContext(ctx).Model(&entity).Where(condition).Updates(params).Error
	if err != nil {
		return nil, err
	}

	return entity, err
}

func (r *repository[T]) Save(ctx context.Context, entity *T) (*T, error) {
	err := r.db.WithContext(ctx).Save(&entity).Error
	if err != nil {
		return nil, err
	}

	return entity, err
}

func (r *repository[T]) GetFirstByID(ctx context.Context, id int) (item T, err error) {
	err = r.db.WithContext(ctx).First(&item, id).Error
	if err != nil {
		return item, err
	}

	return item, err
}

func (r *repository[T]) GetFirstBy(cxt context.Context, condition interface{}, sort string) (entity *T, err error) {
	err = r.db.WithContext(cxt).Where(condition).Order(sort).First(&entity).Error
	if err != nil {
		return nil, err
	}

	return entity, nil
}

func (r *repository[T]) GetAll(ctx context.Context, limit, offset int, sort string) (items []*T, total int64, err error) {

	err = r.db.WithContext(ctx).Count(&total).Limit(limit).Offset(offset).Order(sort).Find(&items).Error
	if err != nil {
		return nil, 0, err
	}
	return items, total, err
}

func (r *repository[T]) GetAllBy(ctx context.Context, limit, offset int, condition any, sort string) (items []*T, total int64, err error) {
	err = r.db.WithContext(ctx).Where(condition).Count(&total).Limit(limit).Offset(offset).Order(sort).Find(&items).Error
	if err != nil {
		return nil, 0, err
	}
	return items, total, err
}

func (r *repository[T]) Delete(ctx context.Context, entity *T, ID any) error {
	err := r.db.WithContext(ctx).Delete(&entity, ID).Error
	if err != nil {
		return err
	}

	return nil
}

func (r *repository[T]) DeleteBy(ctx context.Context, entity *T, condition any) error {
	err := r.db.WithContext(ctx).Where(condition).Delete(&entity).Error
	if err != nil {
		return err
	}

	return nil
}
