package repository

import (
	"crm/internal/domain"

	"gorm.io/gorm"
)

type tagRepository struct {
	db *gorm.DB
}

// NewTagRepository 创建标签仓储实例
func NewTagRepository(db *gorm.DB) TagRepository {
	return &tagRepository{db: db}
}

// GetByID 根据ID获取标签
func (r *tagRepository) GetByID(id uint64) (*domain.Tag, error) {
	var tag domain.Tag
	err := r.db.Preload("Dimension").
		Where("id = ? AND is_deleted = ?", id, false).First(&tag).Error
	return &tag, err
}

// GetList 获取标签列表
func (r *tagRepository) GetList(req domain.TagQueryRequest) ([]domain.TagResponse, int64, error) {
	query := r.db.Table("tags t").
		Select("t.id, t.dimension_id, td.name as dimension_name, t.name, t.color, t.description, t.sort_order").
		Joins("LEFT JOIN tag_dimensions td ON t.dimension_id = td.id").
		Where("t.is_deleted = ? AND td.is_deleted = ?", false, false)

	// 添加筛选条件
	if req.DimensionID != 0 {
		query = query.Where("t.dimension_id = ?", req.DimensionID)
	}
	if req.Name != "" {
		query = query.Where("t.name LIKE ?", "%"+req.Name+"%")
	}

	// 获取总数
	var total int64
	query.Count(&total)

	// 分页查询
	var tags []domain.TagResponse
	offset := (req.Page - 1) * req.PageSize
	err := query.Order("td.sort_order ASC, t.sort_order ASC").
		Offset(offset).Limit(req.PageSize).Find(&tags).Error

	return tags, total, err
}

// Create 创建标签
func (r *tagRepository) Create(tag *domain.Tag) error {
	// 设置默认值
	if tag.Color == "" {
		tag.Color = "#2196F3" // 默认蓝色
	}
	return r.db.Create(tag).Error
}

// Update 更新标签
func (r *tagRepository) Update(tag *domain.Tag) error {
	return r.db.Save(tag).Error
}

// Delete 删除标签（软删除）
func (r *tagRepository) Delete(id uint64) error {
	return r.db.Model(&domain.Tag{}).
		Where("id = ?", id).
		Update("is_deleted", true).Error
}

// GetByDimensionID 根据维度ID获取标签
func (r *tagRepository) GetByDimensionID(dimensionID uint64) ([]domain.Tag, error) {
	var tags []domain.Tag
	err := r.db.Where("dimension_id = ? AND is_deleted = ?", dimensionID, false).
		Order("sort_order ASC").Find(&tags).Error
	return tags, err
}

// GetAllActiveTags 获取所有活跃标签
func (r *tagRepository) GetAllActiveTags() ([]domain.TagResponse, error) {
	var tags []domain.TagResponse
	err := r.db.Table("tags t").
		Select("t.id, t.dimension_id, td.name as dimension_name, t.name, t.color, t.description, t.sort_order").
		Joins("LEFT JOIN tag_dimensions td ON t.dimension_id = td.id").
		Where("t.is_deleted = ? AND td.is_deleted = ?", false, false).
		Order("td.sort_order ASC, t.sort_order ASC").Find(&tags).Error

	return tags, err
}

type tagDimensionRepository struct {
	db *gorm.DB
}

// NewTagDimensionRepository 创建标签维度仓储实例
func NewTagDimensionRepository(db *gorm.DB) TagDimensionRepository {
	return &tagDimensionRepository{db: db}
}

// GetByID 根据ID获取标签维度
func (r *tagDimensionRepository) GetByID(id uint64) (*domain.TagDimension, error) {
	var dimension domain.TagDimension
	err := r.db.Preload("Tags", "is_deleted = ?", false).
		Where("id = ? AND is_deleted = ?", id, false).First(&dimension).Error
	return &dimension, err
}

// GetList 获取标签维度列表
func (r *tagDimensionRepository) GetList(page, pageSize int) ([]domain.TagDimensionResponse, int64, error) {
	var dimensions []domain.TagDimension
	var total int64

	query := r.db.Model(&domain.TagDimension{}).Where("is_deleted = ?", false)

	// 获取总数
	query.Count(&total)

	// 分页查询
	offset := (page - 1) * pageSize
	err := query.Preload("Tags", "is_deleted = ?", false).
		Order("sort_order ASC").Offset(offset).Limit(pageSize).Find(&dimensions).Error

	if err != nil {
		return nil, 0, err
	}

	var responses []domain.TagDimensionResponse
	for _, dimension := range dimensions {
		response := domain.TagDimensionResponse{
			ID:          dimension.ID,
			Name:        dimension.Name,
			Description: dimension.Description,
			SortOrder:   dimension.SortOrder,
			TagCount:    len(dimension.Tags),
			Tags:        dimension.Tags,
		}
		responses = append(responses, response)
	}

	return responses, total, nil
}

// GetWithTags 获取所有维度及其标签
func (r *tagDimensionRepository) GetWithTags() ([]domain.TagDimensionResponse, error) {
	var dimensions []domain.TagDimension
	err := r.db.Preload("Tags", "is_deleted = ?", false).
		Where("is_deleted = ?", false).
		Order("sort_order ASC").Find(&dimensions).Error

	if err != nil {
		return nil, err
	}

	var responses []domain.TagDimensionResponse
	for _, dimension := range dimensions {
		response := domain.TagDimensionResponse{
			ID:          dimension.ID,
			Name:        dimension.Name,
			Description: dimension.Description,
			SortOrder:   dimension.SortOrder,
			TagCount:    len(dimension.Tags),
			Tags:        dimension.Tags,
		}
		responses = append(responses, response)
	}

	return responses, nil
}

// Create 创建标签维度
func (r *tagDimensionRepository) Create(dimension *domain.TagDimension) error {
	return r.db.Create(dimension).Error
}

// Update 更新标签维度
func (r *tagDimensionRepository) Update(dimension *domain.TagDimension) error {
	return r.db.Save(dimension).Error
}

// Delete 删除标签维度（软删除）
func (r *tagDimensionRepository) Delete(id uint64) error {
	tx := r.db.Begin()
	defer func() {
		if r := recover(); r != nil {
			tx.Rollback()
		}
	}()

	// 先删除该维度下的所有标签
	err := tx.Model(&domain.Tag{}).
		Where("dimension_id = ?", id).
		Update("is_deleted", true).Error
	if err != nil {
		tx.Rollback()
		return err
	}

	// 删除维度
	err = tx.Model(&domain.TagDimension{}).
		Where("id = ?", id).
		Update("is_deleted", true).Error
	if err != nil {
		tx.Rollback()
		return err
	}

	return tx.Commit().Error
}
