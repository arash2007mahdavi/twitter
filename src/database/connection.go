package database

import (
	"context"
	"database/sql"
	"encoding/json"
	"fmt"
	"strconv"
	"time"
	"twitter/src/logger"

	"gorm.io/gorm"
)

func TypeComverter[T any](req any) (*T, error) {
	var result T
	json_sample, err := json.Marshal(&req)
	if err != nil {
		return nil, err
	}
	err = json.Unmarshal(json_sample, &result)
	if err != nil {
		return nil ,err
	}
	return &result, nil
}

type BaseService[T any, Tc any, Tu any, Tr any] struct {
	DB *gorm.DB
	Logger logger.Logger
}

func NewBaseService[T any, Tc any, Tu any, Tr any]() *BaseService[T,Tc,Tu,Tr] {
	return &BaseService[T, Tc, Tu, Tr]{
		DB: GetDB(),
		Logger: logger.NewLogger(),
	}
}

func (service *BaseService[T, Tc, Tu, Tr]) Create(ctx context.Context, req *Tc) (*Tr, error) {
	data, err := TypeComverter[T](req)
	if err != nil {
		return nil, fmt.Errorf("error in comverting model")
	}
	tx := service.DB.WithContext(ctx).Begin()
	err = tx.Create(data).Error
	if err != nil {
		tx.Rollback()
		return nil, fmt.Errorf("error in creating new model")
	}
	tx.Commit()
	return TypeComverter[Tr](data)
}

func (service *BaseService[T, Tc, Tu, Tr]) GetAll(ctx context.Context) (*[]T, error) {
	slice := make([]T, 0)
	service.DB.Where("enabled = ?", true).Find(&slice)
	return &slice, nil
}

func (service *BaseService[T, Tc, Tu, Tr]) Update(ctx context.Context, req *Tu) (*Tr, error) {
	data, err := TypeComverter[map[string]interface{}](req)
	int_m, _:= strconv.Atoi(ctx.Value("modified_by").(string))
	int_id, _:= strconv.Atoi(ctx.Value("id").(string))
	(*data)["modified_by"] = sql.NullInt64{Int64: int64(int_m), Valid: true}
	(*data)["modified_at"] = sql.NullTime{Time: time.Now().UTC(), Valid: true}
	if err != nil {
		return nil, fmt.Errorf("error in comverting model")
	}
	tx := service.DB.WithContext(ctx).Begin()
	model := new(T)
	err = tx.Model(model).
		Where("id = ? AND deleted_by is null", int_id).
		Updates(*data).Error
	if err != nil {
		tx.Rollback()
		return nil, fmt.Errorf("error in updating model")
	}
	tx.Commit()
	return service.GetById(ctx, int_id)
}

func (service *BaseService[T, Tc, Tu, Tr]) Delete(ctx context.Context) error {
	data := map[string]interface{}{}
	int_d, _:= strconv.Atoi(ctx.Value("deleted_by").(string))
	int_id, _:= strconv.Atoi(ctx.Value("id").(string))
	(data)["deleted_by"] = sql.NullInt64{Int64: int64(int_d), Valid: true}
	(data)["deleted_at"] = sql.NullTime{Time: time.Now().UTC(), Valid: true}
	(data)["enabled"] = false
	tx := service.DB.WithContext(ctx).Begin()
	model := new(T)
	err := tx.Model(&model).
		Where("id = ? AND deleted_by is null", int_id).
		Updates(data).Error
	if err != nil {
		tx.Rollback()
		return fmt.Errorf("error in deleting model")
	}
	tx.Commit()
	return nil
}

func (service *BaseService[T, Tc, Tu, Tr]) GetById(ctx context.Context, id int) (*Tr, error) {
	model := new(T)
	service.DB.Model(model).
		Where("id = ? AND deleted_by is null", id).
		First(model)
	return TypeComverter[Tr](model)
}