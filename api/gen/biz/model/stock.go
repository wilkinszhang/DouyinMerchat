package model

import (
	"context"
	"fmt"
	"gorm.io/gorm"
	"gorm.io/gorm/clause"
)

type Stock struct {
	Base
	ProductID   uint32 `gorm:"column:product_id;uniqueIndex"`
	Stock       int32  `gorm:"column:stock;default:0"`
	FrozenStock int32  `gorm:"column:frozen_stock;default:0"`
}

func (s Stock) TableName() string {
	return "stock"
}

// GetStockByProductID fetches a stock record by product ID
func GetStockByProductID(db *gorm.DB, ctx context.Context, productID uint32) (Stock, error) {
	var stock Stock
	err := db.WithContext(ctx).Where("product_id = ?", productID).First(&stock).Error
	return stock, err
}

// BatchGetStockByProductIDs fetches multiple stock records by product IDs
func BatchGetStockByProductIDs(db *gorm.DB, ctx context.Context, productIDs []uint32) ([]Stock, error) {
	var stocks []Stock
	err := db.WithContext(ctx).Where("product_id IN ?", productIDs).Find(&stocks).Error
	return stocks, err
}

// CreateOrUpdateStock creates or updates a stock record
func CreateOrUpdateStock(db *gorm.DB, ctx context.Context, stock *Stock) error {
	return db.WithContext(ctx).Clauses(clause.OnConflict{
		Columns:   []clause.Column{{Name: "product_id"}},
		DoUpdates: clause.AssignmentColumns([]string{"stock", "frozen_stock"}),
	}).Create(stock).Error
}

// BatchAddStock adds stock to multiple products
func BatchAddStock(db *gorm.DB, ctx context.Context, items map[uint32]int32) error {
	return db.WithContext(ctx).Transaction(func(tx *gorm.DB) error {
		for productID, quantity := range items {
			err := tx.WithContext(ctx).Exec(
				"INSERT INTO stock (product_id, stock) VALUES (?, ?) ON DUPLICATE KEY UPDATE stock = stock + ?",
				productID, quantity, quantity,
			).Error
			if err != nil {
				return err
			}
		}
		return nil
	})
}

// ReduceStock reduces the available stock and frozen stock
func ReduceStock(db *gorm.DB, ctx context.Context, productID uint32, quantity int32) error {
	result := db.WithContext(ctx).Exec(
		"UPDATE stock SET stock = stock - ?, frozen_stock = frozen_stock - ? WHERE product_id = ? AND frozen_stock >= ?",
		quantity, quantity, productID, quantity,
	)
	if result.Error != nil {
		return result.Error
	}
	if result.RowsAffected == 0 {
		return fmt.Errorf("insufficient frozen stock for product %d", productID)
	}
	return nil
}

// PreReduceStock reduces available stock and increases frozen stock
func PreReduceStock(db *gorm.DB, ctx context.Context, productID uint32, quantity int32) error {
	result := db.WithContext(ctx).Exec(
		"UPDATE stock SET stock = stock - ?, frozen_stock = frozen_stock + ? WHERE product_id = ? AND stock >= ?",
		quantity, quantity, productID, quantity,
	)
	if result.Error != nil {
		return result.Error
	}
	if result.RowsAffected == 0 {
		return fmt.Errorf("insufficient stock for product %d", productID)
	}
	return nil
}

// RollbackStock moves stock from frozen back to available
func RollbackStock(db *gorm.DB, ctx context.Context, productID uint32, quantity int32) error {
	result := db.WithContext(ctx).Exec(
		"UPDATE stock SET stock = stock + ?, frozen_stock = frozen_stock - ? WHERE product_id = ? AND frozen_stock >= ?",
		quantity, quantity, productID, quantity,
	)
	if result.Error != nil {
		return result.Error
	}
	if result.RowsAffected == 0 {
		return fmt.Errorf("insufficient frozen stock for rollback on product %d", productID)
	}
	return nil
}
