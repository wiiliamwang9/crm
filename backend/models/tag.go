package models

import (
	"time"

	"gorm.io/gorm"
)

// TagDimension 标签维度模型
type TagDimension struct {
	ID          uint64     `json:"id" gorm:"primaryKey;autoIncrement;comment:维度ID"`
	Name        string     `json:"name" gorm:"type:varchar(128);not null;uniqueIndex;comment:维度名称"`
	Description string     `json:"description" gorm:"type:text;comment:维度描述"`
	SortOrder   int        `json:"sort_order" gorm:"default:0;index;comment:排序顺序"`
	CreatedAt   time.Time  `json:"created_at" gorm:"index;comment:创建时间"`
	UpdatedAt   time.Time  `json:"updated_at" gorm:"comment:更新时间"`
	DeletedAt   *time.Time `json:"deleted_at" gorm:"comment:删除时间"`
	IsDeleted   bool       `json:"is_deleted" gorm:"default:false;index;comment:是否删除"`

	// 关联关系
	Tags []Tag `json:"tags,omitempty" gorm:"foreignKey:DimensionID"`
}

// Tag 标签模型
type Tag struct {
	ID          uint64     `json:"id" gorm:"primaryKey;autoIncrement;comment:标签ID"`
	DimensionID uint64     `json:"dimension_id" gorm:"not null;index;comment:维度ID"`
	Name        string     `json:"name" gorm:"type:varchar(128);not null;comment:标签名称"`
	Color       string     `json:"color" gorm:"type:varchar(32);comment:标签颜色"`
	Description string     `json:"description" gorm:"type:text;comment:标签描述"`
	SortOrder   int        `json:"sort_order" gorm:"default:0;index;comment:排序顺序"`
	CreatedAt   time.Time  `json:"created_at" gorm:"index;comment:创建时间"`
	UpdatedAt   time.Time  `json:"updated_at" gorm:"comment:更新时间"`
	DeletedAt   *time.Time `json:"deleted_at" gorm:"comment:删除时间"`
	IsDeleted   bool       `json:"is_deleted" gorm:"default:false;index;comment:是否删除"`

	// 关联关系
	Dimension TagDimension `json:"dimension,omitempty" gorm:"foreignKey:DimensionID"`
}

// TagDimensionResponse 标签维度响应
type TagDimensionResponse struct {
	ID          uint64 `json:"id"`
	Name        string `json:"name"`
	Description string `json:"description"`
	SortOrder   int    `json:"sort_order"`
	TagCount    int    `json:"tag_count"`
	Tags        []Tag  `json:"tags,omitempty"`
}

// TagResponse 标签响应
type TagResponse struct {
	ID            uint64 `json:"id"`
	DimensionID   uint64 `json:"dimension_id"`
	DimensionName string `json:"dimension_name"`
	Name          string `json:"name"`
	Color         string `json:"color"`
	Description   string `json:"description"`
	SortOrder     int    `json:"sort_order"`
}

// TagCreateRequest 创建标签请求
type TagCreateRequest struct {
	DimensionID uint64 `json:"dimension_id" binding:"required"`
	Name        string `json:"name" binding:"required,max=128"`
	Color       string `json:"color"`
	Description string `json:"description"`
	SortOrder   int    `json:"sort_order"`
}

// TagUpdateRequest 更新标签请求
type TagUpdateRequest struct {
	Name        string `json:"name" binding:"required,max=128"`
	Color       string `json:"color"`
	Description string `json:"description"`
	SortOrder   int    `json:"sort_order"`
}

// TagDimensionCreateRequest 创建标签维度请求
type TagDimensionCreateRequest struct {
	Name        string `json:"name" binding:"required,max=128"`
	Description string `json:"description"`
	SortOrder   int    `json:"sort_order"`
}

// TagDimensionUpdateRequest 更新标签维度请求
type TagDimensionUpdateRequest struct {
	Name        string `json:"name" binding:"required,max=128"`
	Description string `json:"description"`
	SortOrder   int    `json:"sort_order"`
}

// TagQueryRequest 查询标签请求
type TagQueryRequest struct {
	DimensionID uint64 `json:"dimension_id" form:"dimension_id"`
	Name        string `json:"name" form:"name"`
	Page        int    `json:"page" form:"page"`
	PageSize    int    `json:"page_size" form:"page_size"`
}

// TableName 指定表名
func (TagDimension) TableName() string {
	return "tag_dimensions"
}

// TableName 指定表名
func (Tag) TableName() string {
	return "tags"
}

// BeforeCreate GORM钩子：创建前
func (td *TagDimension) BeforeCreate(tx *gorm.DB) error {
	return nil
}

// BeforeCreate GORM钩子：创建前
func (t *Tag) BeforeCreate(tx *gorm.DB) error {
	if t.Color == "" {
		t.Color = "#2196F3" // 默认蓝色
	}
	return nil
}
