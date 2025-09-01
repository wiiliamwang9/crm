package repository

import (
	"gorm.io/gorm"
)

// NewRepository 创建仓储聚合实例
func NewRepository(db *gorm.DB) *Repository {
	return &Repository{
		Customer:     NewCustomerRepository(db),
		Todo:         NewTodoRepository(db),
		Activity:     NewActivityRepository(db),
		User:         NewUserRepository(db),
		Tag:          NewTagRepository(db),
		TagDimension: NewTagDimensionRepository(db),
		Reminder:     NewReminderRepository(db),
	}
}
