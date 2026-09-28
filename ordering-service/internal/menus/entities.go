package menus

import (
	"time"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

type MenuCategory struct {
	ID          uuid.UUID      `gorm:"type:uuid;default:gen_random_uuid();primaryKey" json:"id"`
	Name        string         `gorm:"type:varchar(255);unique;not null" json:"name"`
	Description string         `gorm:"type:text" json:"description"`
	SortOrder   int            `gorm:"type:int;default:0" json:"sort_order"`
	IsActive    bool           `gorm:"default:true" json:"is_active"`
	CreatedAt   time.Time      `json:"created_at"`
	UpdatedAt   time.Time      `json:"updated_at"`
	DeletedAt   gorm.DeletedAt `gorm:"index" json:"-"`
}

type MenuItem struct {
	ID          uuid.UUID      `gorm:"type:uuid;default:gen_random_uuid();primaryKey" json:"id"`
	CategoryID  uuid.UUID      `gorm:"type:uuid;not null" json:"category_id"`
	Category    MenuCategory   `gorm:"foreignKey:CategoryID" json:"-"`
	Name        string         `gorm:"type:varchar(255);unique;not null" json:"name"`
	Description string         `gorm:"type:text" json:"description"`
	Price       float64        `gorm:"type:decimal(10,2);not null" json:"price"`
	CostPrice   *float64       `gorm:"type:decimal(10,2)" json:"cost_price"`
	ImageURL    string         `gorm:"type:varchar(255)" json:"image_url"`
	IsAvailable bool           `gorm:"default:true" json:"is_available"`
	SortOrder   int            `gorm:"type:int;default:0" json:"sort_order"`
	CreatedAt   time.Time      `json:"created_at"`
	UpdatedAt   time.Time      `json:"updated_at"`
	DeletedAt   gorm.DeletedAt `gorm:"index" json:"-"`
}
