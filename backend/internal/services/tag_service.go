package services

import (
	"crm/internal/domain"
	"crm/internal/repository"
)

// TagService 标签服务
type TagService struct {
	tagRepo       repository.TagRepository
	dimensionRepo repository.TagDimensionRepository
}

// NewTagService 创建标签服务
func NewTagService(tagRepo repository.TagRepository, dimensionRepo repository.TagDimensionRepository) *TagService {
	return &TagService{
		tagRepo:       tagRepo,
		dimensionRepo: dimensionRepo,
	}
}

// GetDimensions 获取所有维度及其标签
func (s *TagService) GetDimensions() ([]domain.TagDimensionResponse, error) {
	return s.dimensionRepo.GetWithTags()
}

// GetDimensionByID 根据ID获取维度
func (s *TagService) GetDimensionByID(id uint64) (*domain.TagDimension, error) {
	return s.dimensionRepo.GetByID(id)
}

// CreateDimension 创建维度
func (s *TagService) CreateDimension(req domain.TagDimensionCreateRequest) (*domain.TagDimension, error) {
	dimension := domain.TagDimension{
		Name:        req.Name,
		Description: req.Description,
		SortOrder:   req.SortOrder,
	}

	err := s.dimensionRepo.Create(&dimension)
	return &dimension, err
}

// UpdateDimension 更新维度
func (s *TagService) UpdateDimension(id uint64, req domain.TagDimensionUpdateRequest) error {
	dimension, err := s.dimensionRepo.GetByID(id)
	if err != nil {
		return err
	}

	dimension.Name = req.Name
	dimension.Description = req.Description
	dimension.SortOrder = req.SortOrder

	return s.dimensionRepo.Update(dimension)
}

// DeleteDimension 删除维度（软删除）
func (s *TagService) DeleteDimension(id uint64) error {
	return s.dimensionRepo.Delete(id)
}

// GetTags 获取标签列表
func (s *TagService) GetTags(req domain.TagQueryRequest) ([]domain.TagResponse, int64, error) {
	return s.tagRepo.GetList(req)
}

// GetTagsByDimensionID 根据维度ID获取标签
func (s *TagService) GetTagsByDimensionID(dimensionID uint64) ([]domain.Tag, error) {
	return s.tagRepo.GetByDimensionID(dimensionID)
}

// GetTagByID 根据ID获取标签
func (s *TagService) GetTagByID(id uint64) (*domain.Tag, error) {
	return s.tagRepo.GetByID(id)
}

// CreateTag 创建标签
func (s *TagService) CreateTag(req domain.TagCreateRequest) (*domain.Tag, error) {
	// 验证维度是否存在
	_, err := s.dimensionRepo.GetByID(req.DimensionID)
	if err != nil {
		return nil, err
	}

	tag := domain.Tag{
		DimensionID: req.DimensionID,
		Name:        req.Name,
		Color:       req.Color,
		Description: req.Description,
		SortOrder:   req.SortOrder,
	}

	err = s.tagRepo.Create(&tag)
	return &tag, err
}

// UpdateTag 更新标签
func (s *TagService) UpdateTag(id uint64, req domain.TagUpdateRequest) error {
	tag, err := s.tagRepo.GetByID(id)
	if err != nil {
		return err
	}

	tag.Name = req.Name
	tag.Color = req.Color
	tag.Description = req.Description
	tag.SortOrder = req.SortOrder

	return s.tagRepo.Update(tag)
}

// DeleteTag 删除标签（软删除）
func (s *TagService) DeleteTag(id uint64) error {
	return s.tagRepo.Delete(id)
}

// GetAllActiveTags 获取所有活跃标签（用于下拉选择）
func (s *TagService) GetAllActiveTags() ([]domain.TagResponse, error) {
	return s.tagRepo.GetAllActiveTags()
}
