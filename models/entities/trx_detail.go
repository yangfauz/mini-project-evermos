package entities

import (
	"time"

	"gorm.io/gorm"
)

type TrxDetail struct {
	gorm.Model
	ID          uint `gorm:"primaryKey"`
	IDTrx       uint `gorm:"not null"`
	IDLogProduk uint `gorm:"not null"`
	IDToko      uint `gorm:"not null"`
	Kuantitas   int  `gorm:"not null"`
	HargaTotal  int  `gorm:"not null"`
	CreatedAt   *time.Time
	UpdatedAt   *time.Time
	Trx         Trx        `gorm:"foreignKey:IDTrx"`
	ProductLog  ProductLog `gorm:"foreignKey:IDLogProduk"`
	Store       Store      `gorm:"foreignKey:IDToko"`
}