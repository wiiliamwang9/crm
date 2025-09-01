package services

import (
	"gorm.io/gorm"

	"crm/models"
)

// TagService 标签服务
type TagService struct {
	db *gorm.DB
}

// NewTagService 创建标签服务
func NewTagService(db *gorm.DB) *TagService {
	return &TagService{db: db}
}

// GetDimensions 获取所有维度及其标签
func (s *TagService) GetDimensions() ([]models.TagDimensionResponse, error) {
	var dimensions []models.TagDimension
	err := s.db.Preload("Tags", "is_deleted = ?", false).
		Where("is_deleted = ?", false).
		Order("sort_order ASC").Find(&dimensions).Error

	if err != nil {
		return nil, err
	}

	var responses []models.TagDimensionResponse
	for _, dimension := range dimensions {
		response := models.TagDimensionResponse{
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

// GetDimensionByID 根据ID获取维度
func (s *TagService) GetDimensionByID(id uint64) (*models.TagDimension, error) {
	var dimension models.TagDimension
	err := s.db.Preload("Tags", "is_deleted = ?", false).
		Where("id = ? AND is_deleted = ?", id, false).First(&dimension).Error

	return &dimension, err
}

// CreateDimension 创建维度
func (s *TagService) CreateDimension(req models.TagDimensionCreateRequest) (*models.TagDimension, error) {
	dimension := models.TagDimension{
		Name:        req.Name,
		Description: req.Description,
		SortOrder:   req.SortOrder,
	}

	err := s.db.Create(&dimension).Error
	return &dimension, err
}

// UpdateDimension 更新维度
func (s *TagService) UpdateDimension(id uint64, req models.TagDimensionUpdateRequest) error {
	return s.db.Model(&models.TagDimension{}).
		Where("id = ? AND is_deleted = ?", id, false).
		Updates(map[string]interface{}{
			"name":        req.Name,
			"description": req.Description,
			"sort_order":  req.SortOrder,
		}).Error
}

// DeleteDimension 删除维度（软删除）
func (s *TagService) DeleteDimension(id uint64) error {
	// 先删除该维度下的所有标签
	err := s.db.Model(&models.Tag{}).
		Where("dimension_id = ?", id).
		Update("is_deleted", true).Error
	if err != nil {
		return err
	}

	// 删除维度
	return s.db.Model(&models.TagDimension{}).
		Where("id = ?", id).
		Update("is_deleted", true).Error
}

// GetTags 获取标签列表
func (s *TagService) GetTags(req models.TagQueryRequest) ([]models.TagResponse, int64, error) {
	query := s.db.Table("tags t").
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
	var tags []models.TagResponse
	offset := (req.Page - 1) * req.PageSize
	err := query.Order("td.sort_order ASC, t.sort_order ASC").
		Offset(offset).Limit(req.PageSize).Find(&tags).Error

	return tags, total, err
}

// GetTagsByDimensionID 根据维度ID获取标签
func (s *TagService) GetTagsByDimensionID(dimensionID uint64) ([]models.Tag, error) {
	var tags []models.Tag
	err := s.db.Where("dimension_id = ? AND is_deleted = ?", dimensionID, false).
		Order("sort_order ASC").Find(&tags).Error

	return tags, err
}

// GetTagByID 根据ID获取标签
func (s *TagService) GetTagByID(id uint64) (*models.Tag, error) {
	var tag models.Tag
	err := s.db.Preload("Dimension").
		Where("id = ? AND is_deleted = ?", id, false).First(&tag).Error

	return &tag, err
}

// CreateTag 创建标签
func (s *TagService) CreateTag(req models.TagCreateRequest) (*models.Tag, error) {
	// 验证维度是否存在
	var dimension models.TagDimension
	err := s.db.Where("id = ? AND is_deleted = ?", req.DimensionID, false).
		First(&dimension).Error
	if err != nil {
		return nil, err
	}

	tag := models.Tag{
		DimensionID: req.DimensionID,
		Name:        req.Name,
		Color:       req.Color,
		Description: req.Description,
		SortOrder:   req.SortOrder,
	}

	err = s.db.Create(&tag).Error
	return &tag, err
}

// UpdateTag 更新标签
func (s *TagService) UpdateTag(id uint64, req models.TagUpdateRequest) error {
	return s.db.Model(&models.Tag{}).
		Where("id = ? AND is_deleted = ?", id, false).
		Updates(map[string]interface{}{
			"name":        req.Name,
			"color":       req.Color,
			"description": req.Description,
			"sort_order":  req.SortOrder,
		}).Error
}

// DeleteTag 删除标签（软删除）
func (s *TagService) DeleteTag(id uint64) error {
	return s.db.Model(&models.Tag{}).
		Where("id = ?", id).
		Update("is_deleted", true).Error
}

// GetAllActiveTags 获取所有活跃标签（用于下拉选择）
func (s *TagService) GetAllActiveTags() ([]models.TagResponse, error) {
	var tags []models.TagResponse
	err := s.db.Table("tags t").
		Select("t.id, t.dimension_id, td.name as dimension_name, t.name, t.color, t.description, t.sort_order").
		Joins("LEFT JOIN tag_dimensions td ON t.dimension_id = td.id").
		Where("t.is_deleted = ? AND td.is_deleted = ?", false, false).
		Order("td.sort_order ASC, t.sort_order ASC").Find(&tags).Error

	return tags, err
}
