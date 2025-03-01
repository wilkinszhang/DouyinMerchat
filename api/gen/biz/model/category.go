package model

import (
	"context"
	"errors"
	"gorm.io/gorm"
	"time"
)

type Category struct {
	ID          uint `gorm:"primarykey"`
	CreatedAt   time.Time
	UpdatedAt   time.Time
	Name        string    `json:"name" gorm:"unique;type:varchar(255);index"`
	Description string    `json:"description"`
	Products    []Product `json:"product" gorm:"many2many:product_category"` //多对多连接表
}

// ProductWithCategories 包含商品信息和分类名称的结构体
type ProductWithCategories struct {
	gorm.Model
	Name        string
	Description string
	Picture     string
	Price       float32
	Categories  string // 存储GROUP_CONCAT的结果
}

func (c Category) TableName() string {
	return "category"
}

func GetProductsByCategoryName(db *gorm.DB, ctx context.Context, name string, page int32, pageSize int64) ([]ProductWithCategories, error) {
	offset := (int64)(page-1) * pageSize
	if offset < 0 {
		return nil, errors.New("invalid page param")
	}

	var products []ProductWithCategories

	err := db.WithContext(ctx).
		Table("product p").
		Select("p.*, GROUP_CONCAT(c.name) as category").
		Joins("JOIN product_category pc ON p.id = pc.product_id").
		Joins("JOIN category c ON pc.category_id = c.id").
		Where("c.name = ?", name).
		Group("p.id").
		Order("p.id").
		Offset(int(offset)).
		Limit(int(pageSize)).
		Find(&products).Error

	return products, err
}
