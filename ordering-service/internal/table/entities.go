package table

import (
	"time"

	"github.com/google/uuid"
)

type DiningTable struct {
	ID       uuid.UUID `gorm:"type:uuid;primaryKey" json:"id"`
	TableNo  string    `gorm:"column:table_no" json:"table_no"`
	Capacity int       `json:"capacity"`
	QRToken  uuid.UUID `gorm:"column:qr_token;unique" json:"qr_token"`
	Status   string    `json:"status"`
	CreateAt time.Time `json:"created_at"`
	UpdateAt time.Time `json:"updated_at"`
}

func (DiningTable) TableName() string { return "dining_tables" }

type TableSession struct {
	ID        uuid.UUID  `gorm:"type:uuid;primaryKey;default:gen_random_uuid()" json:"id"`
	TableID   uuid.UUID  `gorm:"type:uuid;index;not null" json:"table_id"`
	Token     string     `gorm:"uniqueIndex;not null" json:"token"`
	CreateAt  time.Time  `json:"created_at"`
	ExpiresAt time.Time  `gorm:"not null" json:"expires_at"`
	ClosedAt  *time.Time `json:"closed_at,omitempty"`
}

func (TableSession) TableName() string { return "table_sessions" }
