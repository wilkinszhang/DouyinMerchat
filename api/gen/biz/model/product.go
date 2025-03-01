package model

import (
	"context"
	"fmt"
	"github.com/cloudwego/hertz/pkg/common/json"
	"github.com/redis/go-redis/v9"
	"gorm.io/gorm"
	"time"
)

type Product struct {
	ID          uint `gorm:"primarykey"`
	CreatedAt   time.Time
	UpdatedAt   time.Time
	Name        string     `json:"name"`
	Description string     `json:"description"`
	Picture     string     `json:"picture"`
	Price       float32    `json:"price"`
	Categories  []Category `json:"categories" gorm:"many2many:product_category"`
}

type ProductQuery struct {
	ctx context.Context
	db  *gorm.DB
}

func NewProductQuery(ctx context.Context, db *gorm.DB) ProductQuery {
	return ProductQuery{ctx: ctx, db: db}
}

type CachedProductQuery struct {
	productQuery ProductQuery
	cacheClient  *redis.Client
	prefix       string
}

func (p Product) TableName() string {
	return "product"
}

func NewCachedProductQuery(pq ProductQuery, cacheClient *redis.Client) CachedProductQuery {
	return CachedProductQuery{productQuery: pq, cacheClient: cacheClient, prefix: "douyin_merchant"}
}

func (p ProductQuery) GetById(productId int) (product Product, err error) {
	err = p.db.WithContext(p.ctx).Model(&Product{}).Where(&Product{ID: uint(productId)}).First(&product).Error
	return
}

func (c CachedProductQuery) GetById(productId int) (product Product, err error) {
	//douyin_merchant_product_by_id_商品id
	cacheKey := fmt.Sprintf("%s_%s_%d", c.prefix, "product_by_id", productId)
	//Get Redis `GET key` command 返回key对应的value
	cacheResult := c.cacheClient.Get(c.productQuery.ctx, cacheKey)

	//立即执行函数 避免变量污染全局作用域
	err = func() error {
		err1 := cacheResult.Err()
		if err1 != nil {
			return err1
		}
		cacheResultByte, err2 := cacheResult.Bytes()
		if err2 != nil {
			return err2
		}
		//将redis返回的字节数组数据写到product结构体
		err3 := json.Unmarshal(cacheResultByte, &product)
		if err3 != nil {
			return err3
		}
		return nil
	}()
	if err != nil {
		//redis没找到
		product, err = c.productQuery.GetById(productId)
		if err != nil {
			return Product{}, err
		}
		encoded, err := json.Marshal(product)
		if err != nil {
			return product, nil
		}
		//重新放到redis，设置过期一个小时
		_ = c.cacheClient.Set(c.productQuery.ctx, cacheKey, encoded, time.Hour)
	}
	return
}

type Result struct {
	Product
	Category string
}

// model/product.go
func (p ProductQuery) GetProductsByQuery(query string) ([]ProductWithCategories, error) {
	var results []ProductWithCategories

	err := p.db.WithContext(p.ctx).
		Table("product p").
		Select("p.*, GROUP_CONCAT(c.name) as category").
		Joins("LEFT JOIN product_category pc ON p.id = pc.product_id").
		Joins("LEFT JOIN category c ON pc.category_id = c.id").
		Where("MATCH(p.name, p.description) AGAINST (? IN BOOLEAN MODE)", query+"*").
		Group("p.id").
		Order(p.db.Raw("MATCH(p.name) AGAINST (?) DESC", query)).
		Order(p.db.Raw("MATCH(p.description) AGAINST (?) DESC", query)).
		Limit(100).
		Find(&results).Error

	if err != nil {
		return nil, err
	}
	return results, nil
}
