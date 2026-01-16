package folder

import (
	"context"

	"github.com/KOMKZ/go-yogan-domain-folder/model"
)

// Repository 文件夹仓储接口
type Repository interface {
	// 基础 CRUD
	Create(ctx context.Context, folder *model.Folder) error
	Update(ctx context.Context, folder *model.Folder) error
	Delete(ctx context.Context, id uint) error
	FindByID(ctx context.Context, id uint) (*model.Folder, error)

	// 层级查询
	FindByParentID(ctx context.Context, parentID *uint) ([]*model.Folder, error)
	FindChildren(ctx context.Context, parentID uint) ([]*model.Folder, error)
	FindRoots(ctx context.Context) ([]*model.Folder, error)

	// 路径查询
	FindByPath(ctx context.Context, pathPrefix string) ([]*model.Folder, error)
	FindAncestors(ctx context.Context, path string) ([]*model.Folder, error)

	// 全量查询
	FindAll(ctx context.Context) ([]*model.Folder, error)

	// 排序
	UpdateSortOrder(ctx context.Context, id uint, sortOrder int) error
	FindMaxSortOrder(ctx context.Context, parentID *uint) (int, error)

	// 批量更新
	UpdatePathAndDepth(ctx context.Context, id uint, path string, depth int) error
	UpdateChildrenPathAndDepth(ctx context.Context, oldPathPrefix, newPathPrefix string, depthDiff int) error

	// 检查
	ExistsByNameAndParent(ctx context.Context, name string, parentID *uint, excludeID *uint) (bool, error)
	HasChildren(ctx context.Context, id uint) (bool, error)

	// 计数操作
	IncrementItemCount(ctx context.Context, id uint, delta int) error      // 更新直接子项数量
	IncrementTotalItemCount(ctx context.Context, path string, delta int) error // 更新祖先链的总数量
}
