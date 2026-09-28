package menus

import "gorm.io/gorm"

type Repository interface {
	GetActiveCategories() ([]MenuCategory, error)
	GetAvailableItems() ([]MenuItem, error)
	GetMenuItemByID(id string) (*MenuItem, error)
}

type repository struct{ db *gorm.DB }

func NewRepository(db *gorm.DB) Repository { return &repository{db} }

func (r *repository) GetActiveCategories() ([]MenuCategory, error) {
	var categories []MenuCategory
	err := r.db.
		Where("is_active = ?", true).
		Order("sort_order asc").
		Find(&categories).Error
	return categories, err
}

func (r *repository) GetAvailableItems() ([]MenuItem, error) {
	var items []MenuItem
	err := r.db.
		Where("is_available = ?", true).
		Order("sort_order asc").
		Find(&items).Error
	return items, err
}

// 🔥 ເພີ່ມການເຮັດວຽກຂອງຟັງຊັນ
func (r *repository) GetMenuItemByID(id string) (*MenuItem, error) {
	var item MenuItem
	// ດຶງຂໍ້ມູນລາຍການດຽວ ທີ່ກົງກັບ id
	err := r.db.Where("id = ?", id).First(&item).Error
	return &item, err
}
