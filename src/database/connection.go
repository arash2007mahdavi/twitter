package database

import (
	"context"
	"encoding/json"
	"fmt"
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