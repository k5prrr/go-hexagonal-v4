package domain

import (
	"time"
)

type Client struct {
	ID                  int64      `json:"id" gorm:"primaryKey;autoIncrement"`
	FamilyName          string     `json:"family_name" gorm:"size:128"`
	Name                string     `json:"name" gorm:"size:128"`
	MiddleName          string     `json:"middle_name" gorm:"size:128"`
	Phone               string     `json:"phone" gorm:"size:32"`
	Email               string     `json:"email" gorm:"size:64"`
	BirthDate           *time.Time `json:"birth_date" gorm:"type:date"`
	Age                 int        `json:"age"`
	ParentID            *int64     `json:"parent_id"` // nullable
	GenderID            *int       `json:"gender_id"` // nullable
	ShopID              string     `json:"shop_id" gorm:"size:128"`
	ProductCode         string     `json:"product_code" gorm:"size:64"`
	ProductName         string     `json:"product_name" gorm:"size:128"`
	OldProductID        *int       `json:"old_product_id"` // nullable
	OldProductName      string     `json:"old_product_name" gorm:"size:64"`
	OldProductModel     string     `json:"old_product_model" gorm:"size:64"`
	QuantityCartridge06 int        `json:"quantity_cartridge06"`
	QuantityCartridge08 int        `json:"quantity_cartridge08"`
	QuantityCartridge1  int        `json:"quantity_cartridge1"`
	CreatedAt           time.Time  `json:"created_at" gorm:"autoCreateTime"`
	UpdatedAt           time.Time  `json:"updated_at" gorm:"autoUpdateTime"`
	DeletedAt           *time.Time `json:"deleted_at" gorm:"index"` // soft delete
}
