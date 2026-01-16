package folder

import (
	"context"
	"errors"
	"strings"

	"github.com/KOMKZ/go-yogan-domain-folder/model"
	"gorm.io/gorm"
)

// GormRepository GORM 实现的 Repository
type GormRepository struct {
	db        *gorm.DB
	tableName string
}

// NewGormRepository 创建 GORM Repository
// tableName 参数允许不同业务使用不同的表
func NewGormRepository(db *gorm.DB, tableName string) *GormRepository {
	return &GormRepository{
		db:        db,
		tableName: tableName,
	}
}

// table 返回指定表名的 DB 实例
func (r *GormRepository) table(ctx context.Context) *gorm.DB {
	return r.db.WithContext(ctx).Table(r.tableName)
}

// Create 创建文件夹
func (r *GormRepository) Create(ctx context.Context, folder *model.Folder) error {
	return r.table(ctx).Create(folder).Error
}

// Update 更新文件夹
func (r *GormRepository) Update(ctx context.Context, folder *model.Folder) error {
	return r.table(ctx).Save(folder).Error
}

// Delete 删除文件夹（软删除）
func (r *GormRepository) Delete(ctx context.Context, id uint) error {
	return r.table(ctx).Delete(&model.Folder{}, id).Error
}

// FindByID 根据 ID 查询
func (r *GormRepository) FindByID(ctx context.Context, id uint) (*model.Folder, error) {
	var folder model.Folder
	err := r.table(ctx).First(&folder, id).Error
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, ErrFolderNotFound
		}
		return nil, err
	}
	return &folder, nil
}

// FindByParentID 根据父 ID 查询子节点
func (r *GormRepository) FindByParentID(ctx context.Context, parentID *uint) ([]*model.Folder, error) {
	var folders []*model.Folder
	query := r.table(ctx)
	if parentID == nil {
		query = query.Where("parent_id IS NULL")
	} else {
		query = query.Where("parent_id = ?", *parentID)
	}
	err := query.Order("sort_order ASC, id ASC").Find(&folders).Error
	return folders, err
}

// FindChildren 查询直接子节点
func (r *GormRepository) FindChildren(ctx context.Context, parentID uint) ([]*model.Folder, error) {
	var folders []*model.Folder
	err := r.table(ctx).
		Where("parent_id = ?", parentID).
		Order("sort_order ASC, id ASC").
		Find(&folders).Error
	return folders, err
}

// FindRoots 查询所有根节点
func (r *GormRepository) FindRoots(ctx context.Context) ([]*model.Folder, error) {
	var folders []*model.Folder
	err := r.table(ctx).
		Where("parent_id IS NULL").
		Order("sort_order ASC, id ASC").
		Find(&folders).Error
	return folders, err
}

// FindByPath 根据路径前缀查询所有子孙节点
func (r *GormRepository) FindByPath(ctx context.Context, pathPrefix string) ([]*model.Folder, error) {
	var folders []*model.Folder
	err := r.table(ctx).
		Where("path LIKE ?", pathPrefix+"%").
		Order("depth ASC, sort_order ASC, id ASC").
		Find(&folders).Error
	return folders, err
}

// FindAncestors 根据路径查询所有祖先节点
func (r *GormRepository) FindAncestors(ctx context.Context, path string) ([]*model.Folder, error) {
	// 解析路径获取祖先 ID 列表
	ids := parsePathIDs(path)
	if len(ids) == 0 {
		return []*model.Folder{}, nil
	}

	var folders []*model.Folder
	err := r.table(ctx).
		Where("id IN ?", ids).
		Order("depth ASC").
		Find(&folders).Error
	return folders, err
}

// FindAll 查询所有文件夹
func (r *GormRepository) FindAll(ctx context.Context) ([]*model.Folder, error) {
	var folders []*model.Folder
	err := r.table(ctx).
		Order("depth ASC, sort_order ASC, id ASC").
		Find(&folders).Error
	return folders, err
}

// UpdateSortOrder 更新排序
func (r *GormRepository) UpdateSortOrder(ctx context.Context, id uint, sortOrder int) error {
	return r.table(ctx).
		Where("id = ?", id).
		Update("sort_order", sortOrder).Error
}

// FindMaxSortOrder 查询同级下最大排序号
func (r *GormRepository) FindMaxSortOrder(ctx context.Context, parentID *uint) (int, error) {
	var maxOrder int
	query := r.table(ctx).Select("COALESCE(MAX(sort_order), 0)")
	if parentID == nil {
		query = query.Where("parent_id IS NULL")
	} else {
		query = query.Where("parent_id = ?", *parentID)
	}
	err := query.Scan(&maxOrder).Error
	return maxOrder, err
}

// UpdatePathAndDepth 更新单个节点的路径和深度
func (r *GormRepository) UpdatePathAndDepth(ctx context.Context, id uint, path string, depth int) error {
	return r.table(ctx).
		Where("id = ?", id).
		Updates(map[string]interface{}{
			"path":  path,
			"depth": depth,
		}).Error
}

// UpdateChildrenPathAndDepth 批量更新子孙节点的路径和深度
func (r *GormRepository) UpdateChildrenPathAndDepth(ctx context.Context, oldPathPrefix, newPathPrefix string, depthDiff int) error {
	// 使用 SQL 表达式更新
	return r.table(ctx).
		Where("path LIKE ?", oldPathPrefix+"%").
		Where("path != ?", oldPathPrefix). // 排除自身
		Updates(map[string]interface{}{
			"path":  gorm.Expr("REPLACE(path, ?, ?)", oldPathPrefix, newPathPrefix),
			"depth": gorm.Expr("depth + ?", depthDiff),
		}).Error
}

// ExistsByNameAndParent 检查同级下是否存在相同名称
func (r *GormRepository) ExistsByNameAndParent(ctx context.Context, name string, parentID *uint, excludeID *uint) (bool, error) {
	var count int64
	query := r.table(ctx).Where("name = ?", name)
	if parentID == nil {
		query = query.Where("parent_id IS NULL")
	} else {
		query = query.Where("parent_id = ?", *parentID)
	}
	if excludeID != nil {
		query = query.Where("id != ?", *excludeID)
	}
	err := query.Count(&count).Error
	return count > 0, err
}

// HasChildren 检查是否有子节点
func (r *GormRepository) HasChildren(ctx context.Context, id uint) (bool, error) {
	var count int64
	err := r.table(ctx).Where("parent_id = ?", id).Count(&count).Error
	return count > 0, err
}

// IncrementItemCount 更新指定文件夹的直接子项数量
// delta 可以为正数（增加）或负数（减少）
func (r *GormRepository) IncrementItemCount(ctx context.Context, id uint, delta int) error {
	return r.table(ctx).
		Where("id = ?", id).
		Update("item_count", gorm.Expr("GREATEST(item_count + ?, 0)", delta)).Error
}

// IncrementTotalItemCount 更新祖先链上所有节点的 total_item_count
// path 为当前节点的路径，会更新 path 上所有祖先节点
func (r *GormRepository) IncrementTotalItemCount(ctx context.Context, path string, delta int) error {
	// 解析路径获取所有祖先 ID
	ids := parsePathIDs(path)
	if len(ids) == 0 {
		return nil
	}

	return r.table(ctx).
		Where("id IN ?", ids).
		Update("total_item_count", gorm.Expr("GREATEST(total_item_count + ?, 0)", delta)).Error
}

// parsePathIDs 解析路径中的 ID 列表
// 路径格式："/1/3/5/" -> [1, 3, 5]
func parsePathIDs(path string) []uint {
	path = strings.Trim(path, "/")
	if path == "" {
		return nil
	}

	parts := strings.Split(path, "/")
	ids := make([]uint, 0, len(parts))
	for _, p := range parts {
		var id uint
		if _, err := parseUint(p, &id); err == nil {
			ids = append(ids, id)
		}
	}
	return ids
}

// parseUint 解析字符串为 uint
func parseUint(s string, v *uint) (bool, error) {
	var n uint
	for _, c := range s {
		if c < '0' || c > '9' {
			return false, errors.New("invalid number")
		}
		n = n*10 + uint(c-'0')
	}
	*v = n
	return true, nil
}
