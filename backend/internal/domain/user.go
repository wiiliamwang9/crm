package domain

import (
	"time"
)

// User 用户模型
type User struct {
	ID                 uint64     `json:"id" gorm:"primaryKey;autoIncrement;comment:用户ID"`
	Name               string     `json:"name" gorm:"type:varchar(256);not null;comment:用户名"`
	ManagerID          *uint64    `json:"manager_id" gorm:"comment:主管ID"`
	Email              string     `json:"email" gorm:"type:varchar(256);comment:邮箱"`
	Phone              string     `json:"phone" gorm:"type:varchar(32);comment:手机号"`
	Department         string     `json:"department" gorm:"type:varchar(128);comment:部门"`
	DepartmentLeaderID *uint64    `json:"department_leader_id" gorm:"comment:部门领导ID"`
	Position           string     `json:"position" gorm:"type:varchar(128);comment:职位"`
	WechatWorkID       string     `json:"wechat_work_id" gorm:"type:varchar(128);comment:企业微信ID"`
	WechatID           string     `json:"wechat_id" gorm:"type:varchar(128);comment:微信ID"`
	Status             string     `json:"status" gorm:"type:varchar(32);default:active;comment:状态"`
	AvatarURL          string     `json:"avatar_url" gorm:"type:varchar(512);comment:头像URL"`
	LastLoginAt        *time.Time `json:"last_login_at" gorm:"comment:最后登录时间"`
	CreatedAt          time.Time  `json:"created_at" gorm:"index;comment:创建时间"`
	UpdatedAt          time.Time  `json:"updated_at" gorm:"comment:更新时间"`
	DeletedAt          *time.Time `json:"deleted_at" gorm:"comment:删除时间"`
	IsDeleted          bool       `json:"is_deleted" gorm:"default:false;index;comment:是否删除"`

	// 关联关系
	Manager          *User `json:"manager,omitempty" gorm:"foreignKey:ManagerID"`
	DepartmentLeader *User `json:"department_leader,omitempty" gorm:"foreignKey:DepartmentLeaderID"`
}

// UserResponse 用户响应
type UserResponse struct {
	ID          uint64 `json:"id"`
	Name        string `json:"name"`
	Department  string `json:"department"`
	Position    string `json:"position"`
	Email       string `json:"email"`
	Phone       string `json:"phone"`
	Status      string `json:"status"`
	ManagerName string `json:"manager_name,omitempty"`
	AvatarURL   string `json:"avatar_url"`
}

// HomepageUserResponse 首页用户信息响应
type HomepageUserResponse struct {
	ID             uint64 `json:"id"`
	Name           string `json:"name"`
	ShopName       string `json:"shop_name"`
	TodayRevenue   int    `json:"today_revenue"`
	TodayFollowUps int    `json:"today_follow_ups"`
}

// UserQueryRequest 查询用户请求
type UserQueryRequest struct {
	Department string `json:"department" form:"department"`
	Status     string `json:"status" form:"status"`
	Name       string `json:"name" form:"name"`
	Page       int    `json:"page" form:"page"`
	PageSize   int    `json:"page_size" form:"page_size"`
}

// TableName 指定表名
func (User) TableName() string {
	return "users"
}
