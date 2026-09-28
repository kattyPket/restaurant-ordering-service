package table

import (
	"github.com/google/uuid"
	"gorm.io/gorm"
	"gorm.io/gorm/clause"
)

type Repository struct {
	db *gorm.DB
}

func NewRepository(db *gorm.DB) *Repository {
	return &Repository{db: db}
}

// ເພີ່ມຟັງຊັນ GetDB ເພື່ອໃຫ້ Service ໃຊ້ເລີ່ມ Transaction
func (r *Repository) GetDB() *gorm.DB {
	return r.db
}

// ປ່ຽນໃຫ້ຮັບ db *gorm.DB ເພື່ອຮອງຮັບ Transaction ດ້ວຍ FOR UPDATE
func (r *Repository) FindTableByQRTokenWithLock(db *gorm.DB, qrToken string) (*DiningTable, error) {
	var t DiningTable
	// 🔥 ໃຊ້ clause.Locking ຖ້າຕ້ອງການ, ແຕ່ໃນ GORM ສາມາດຂຽນດິບໄດ້:
	if err := db.Raw("SELECT * FROM dining_tables WHERE qr_token = ? FOR UPDATE", qrToken).Scan(&t).Error; err != nil {
		return nil, err
	}
	// ຖ້າຫາບໍ່ພົບ Scan() ຈະບໍ່ throw ErrRecordNotFound ແບບ First(), ຕ້ອງເຊັກ ID
	if t.ID.String() == "00000000-0000-0000-0000-000000000000" {
		return nil, gorm.ErrRecordNotFound
	}
	return &t, nil
}

func (r *Repository) FindTableByIDWithLock(db *gorm.DB, tableID uuid.UUID) (*DiningTable, error) {
	var t DiningTable
	if err := db.Clauses(clause.Locking{Strength: "UPDATE"}).First(&t, "id = ?", tableID).Error; err != nil {
		return nil, err
	}
	return &t, nil
}

func (r *Repository) UpdateTableStatus(db *gorm.DB, tableID string, status string) error {
	return db.Model(&DiningTable{}).
		Where("id = ?", tableID).
		Update("status", status).Error
}

func (r *Repository) FindActiveSessionByTableID(db *gorm.DB, tableID string) (*TableSession, error) {
	var s TableSession
	if err := db.Where("table_id = ? AND closed_at IS NULL", tableID).First(&s).Error; err != nil {
		return nil, err
	}
	return &s, nil
}

func (r *Repository) FindSessionByID(db *gorm.DB, sessionID uuid.UUID) (*TableSession, error) {
	var s TableSession
	if err := db.First(&s, "id = ?", sessionID).Error; err != nil {
		return nil, err
	}
	return &s, nil
}

func (r *Repository) CreateSession(db *gorm.DB, s *TableSession) error {
	return db.Create(s).Error
}

// (ຂອງເກົ່າສຳລັບ Validation ບໍ່ຕ້ອງໃຊ້ TX ກໍ່ໄດ້)
func (r *Repository) FindSessionByToken(token string) (*TableSession, error) {
	var s TableSession
	if err := r.db.Where("token = ?", token).First(&s).Error; err != nil {
		return nil, err
	}
	return &s, nil
}
