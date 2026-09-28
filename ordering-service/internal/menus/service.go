package menus

type Service interface {
	GetMenu() ([]MenuCategoryResponse, error)
}

type service struct{ repo Repository }

func NewService(repo Repository) Service { return &service{repo} }

func (s *service) GetMenu() ([]MenuCategoryResponse, error) {
	categories, err := s.repo.GetActiveCategories()
	if err != nil {
		return nil, err
	}

	items, err := s.repo.GetAvailableItems()
	if err != nil {
		return nil, err
	}

	itemsByCategory := make(map[string][]MenuItemResponse, len(categories))
	for _, item := range items {
		categoryID := item.CategoryID.String()
		itemsByCategory[categoryID] = append(itemsByCategory[categoryID], MenuItemResponse{
			ID:          item.ID.String(),
			Name:        item.Name,
			Description: item.Description,
			Price:       item.Price,
			ImageURL:    item.ImageURL,
			IsAvailable: item.IsAvailable,
			SortOrder:   item.SortOrder,
		})
	}

	response := make([]MenuCategoryResponse, 0, len(categories))
	for _, category := range categories {
		response = append(response, MenuCategoryResponse{
			ID:          category.ID.String(),
			Name:        category.Name,
			Description: category.Description,
			SortOrder:   category.SortOrder,
			Items:       itemsByCategory[category.ID.String()],
		})
	}

	return response, nil
}
