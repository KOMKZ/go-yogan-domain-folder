package model

import (
	"time"

	"gorm.io/gorm"
)

// Folder 通用文件夹/层级节点模型
// 注意：不实现 TableName() 方法，表名由 Repository 动态指定
type Folder struct {
	ID        uint           `gorm:"primaryKey" json:"id"`
	Name      string         `gorm:"size:255;not null" json:"name"`
	ParentID  *uint          `gorm:"index" json:"parentId"`
	SortOrder int            `gorm:"default:0" json:"sortOrder"`
	Depth     int            `gorm:"default:0" json:"depth"`
	Path      string         `gorm:"size:1000" json:"path"` // 物化路径，如 "/1/3/5/"
	CreatedAt time.Time      `json:"createdAt"`
	UpdatedAt time.Time      `json:"updatedAt"`
	DeletedAt gorm.DeletedAt `gorm:"index" json:"deletedAt,omitempty"`

	// 非数据库字段
	Children []Folder `gorm:"-" json:"children,omitempty"`
}

// FolderNode 用于树形结构展示
type FolderNode struct {
	ID        uint          `json:"id"`
	Name      string        `json:"name"`
	ParentID  *uint         `json:"parentId"`
	SortOrder int           `json:"sortOrder"`
	Depth     int           `json:"depth"`
	Children  []*FolderNode `json:"children,omitempty"`
}

// ToNode 将 Folder 转换为 FolderNode
func (f *Folder) ToNode() *FolderNode {
	return &FolderNode{
		ID:        f.ID,
		Name:      f.Name,
		ParentID:  f.ParentID,
		SortOrder: f.SortOrder,
		Depth:     f.Depth,
		Children:  nil,
	}
}
