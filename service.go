package folder

import (
	"context"
	"fmt"
	"strings"

	"github.com/KOMKZ/go-yogan-domain-folder/model"
)

// ServiceConfig 服务配置
type ServiceConfig struct {
	MaxDepth int // 最大深度限制，0 表示无限制
}

// DefaultServiceConfig 默认配置
var DefaultServiceConfig = ServiceConfig{
	MaxDepth: 10,
}

// Service 文件夹服务
type Service struct {
	repo   Repository
	config ServiceConfig
}

// NewService 创建服务
func NewService(repo Repository) *Service {
	return &Service{
		repo:   repo,
		config: DefaultServiceConfig,
	}
}

// NewServiceWithConfig 创建带配置的服务
func NewServiceWithConfig(repo Repository, config ServiceConfig) *Service {
	return &Service{
		repo:   repo,
		config: config,
	}
}

// CreateFolderInput 创建文件夹输入
type CreateFolderInput struct {
	Name     string
	ParentID *uint
}

// CreateFolder 创建文件夹
func (s *Service) CreateFolder(ctx context.Context, input *CreateFolderInput) (*model.Folder, error) {
	// 验证名称
	if err := s.validateName(input.Name); err != nil {
		return nil, err
	}

	// 检查名称唯一性
	exists, err := s.repo.ExistsByNameAndParent(ctx, input.Name, input.ParentID, nil)
	if err != nil {
		return nil, err
	}
	if exists {
		return nil, ErrDuplicateName
	}

	// 计算 depth 和 path
	var depth int
	var path string

	if input.ParentID != nil {
		parent, err := s.repo.FindByID(ctx, *input.ParentID)
		if err != nil {
			return nil, ErrParentNotFound
		}
		depth = parent.Depth + 1
		path = parent.Path // 先使用父节点的路径，后面再追加自己的 ID
	} else {
		depth = 0
		path = "/"
	}

	// 检查深度限制
	if s.config.MaxDepth > 0 && depth >= s.config.MaxDepth {
		return nil, ErrMaxDepthExceeded
	}

	// 获取排序号
	maxOrder, err := s.repo.FindMaxSortOrder(ctx, input.ParentID)
	if err != nil {
		return nil, err
	}

	folder := &model.Folder{
		Name:      input.Name,
		ParentID:  input.ParentID,
		Depth:     depth,
		Path:      path, // 临时路径
		SortOrder: maxOrder + 1,
	}

	// 创建文件夹
	if err := s.repo.Create(ctx, folder); err != nil {
		return nil, err
	}

	// 更新 path 包含自身 ID
	folder.Path = path + fmt.Sprintf("%d/", folder.ID)
	if err := s.repo.Update(ctx, folder); err != nil {
		return nil, err
	}

	return folder, nil
}

// UpdateFolderInput 更新文件夹输入
type UpdateFolderInput struct {
	ID   uint
	Name string
}

// UpdateFolder 更新文件夹
func (s *Service) UpdateFolder(ctx context.Context, input *UpdateFolderInput) (*model.Folder, error) {
	// 查找文件夹
	folder, err := s.repo.FindByID(ctx, input.ID)
	if err != nil {
		return nil, err
	}

	// 验证名称
	if err := s.validateName(input.Name); err != nil {
		return nil, err
	}

	// 检查名称唯一性（排除自身）
	exists, err := s.repo.ExistsByNameAndParent(ctx, input.Name, folder.ParentID, &input.ID)
	if err != nil {
		return nil, err
	}
	if exists {
		return nil, ErrDuplicateName
	}

	folder.Name = input.Name
	if err := s.repo.Update(ctx, folder); err != nil {
		return nil, err
	}

	return folder, nil
}

// DeleteFolder 删除文件夹
func (s *Service) DeleteFolder(ctx context.Context, id uint) error {
	// 检查是否存在
	_, err := s.repo.FindByID(ctx, id)
	if err != nil {
		return err
	}

	// 检查是否有子节点
	hasChildren, err := s.repo.HasChildren(ctx, id)
	if err != nil {
		return err
	}
	if hasChildren {
		return ErrHasChildren
	}

	return s.repo.Delete(ctx, id)
}

// GetFolder 获取单个文件夹
func (s *Service) GetFolder(ctx context.Context, id uint) (*model.Folder, error) {
	return s.repo.FindByID(ctx, id)
}

// GetChildren 获取子节点
func (s *Service) GetChildren(ctx context.Context, parentID *uint) ([]*model.Folder, error) {
	return s.repo.FindByParentID(ctx, parentID)
}

// GetTree 获取完整树结构
func (s *Service) GetTree(ctx context.Context) ([]*model.FolderNode, error) {
	folders, err := s.repo.FindAll(ctx)
	if err != nil {
		return nil, err
	}
	return buildTree(folders, nil), nil
}

// GetSubTree 获取子树
func (s *Service) GetSubTree(ctx context.Context, rootID uint) ([]*model.FolderNode, error) {
	root, err := s.repo.FindByID(ctx, rootID)
	if err != nil {
		return nil, err
	}

	descendants, err := s.repo.FindByPath(ctx, root.Path)
	if err != nil {
		return nil, err
	}

	return buildTree(descendants, &rootID), nil
}

// GetAncestors 获取祖先节点（面包屑）
func (s *Service) GetAncestors(ctx context.Context, id uint) ([]*model.Folder, error) {
	folder, err := s.repo.FindByID(ctx, id)
	if err != nil {
		return nil, err
	}

	return s.repo.FindAncestors(ctx, folder.Path)
}

// MoveFolder 移动文件夹
func (s *Service) MoveFolder(ctx context.Context, id uint, newParentID *uint) error {
	folder, err := s.repo.FindByID(ctx, id)
	if err != nil {
		return err
	}

	// 检查是否移动到自己或子节点下
	if newParentID != nil {
		if *newParentID == id {
			return ErrCircularReference
		}

		newParent, err := s.repo.FindByID(ctx, *newParentID)
		if err != nil {
			return ErrParentNotFound
		}

		// 检查新父节点是否是当前节点的子孙
		if strings.HasPrefix(newParent.Path, folder.Path) {
			return ErrCircularReference
		}
	}

	// 计算新的 depth 和 path
	var newDepth int
	var newPath string

	if newParentID != nil {
		newParent, _ := s.repo.FindByID(ctx, *newParentID)
		newDepth = newParent.Depth + 1
		newPath = newParent.Path + fmt.Sprintf("%d/", folder.ID)
	} else {
		newDepth = 0
		newPath = fmt.Sprintf("/%d/", folder.ID)
	}

	// 检查深度限制
	if s.config.MaxDepth > 0 {
		// 计算子树的最大深度
		descendants, err := s.repo.FindByPath(ctx, folder.Path)
		if err != nil {
			return err
		}
		maxChildDepth := 0
		for _, d := range descendants {
			relativeDepth := d.Depth - folder.Depth
			if relativeDepth > maxChildDepth {
				maxChildDepth = relativeDepth
			}
		}
		if newDepth+maxChildDepth >= s.config.MaxDepth {
			return ErrMaxDepthExceeded
		}
	}

	oldPath := folder.Path
	depthDiff := newDepth - folder.Depth

	// 更新当前节点
	folder.ParentID = newParentID
	folder.Depth = newDepth
	folder.Path = newPath

	// 获取新的排序号
	maxOrder, err := s.repo.FindMaxSortOrder(ctx, newParentID)
	if err != nil {
		return err
	}
	folder.SortOrder = maxOrder + 1

	if err := s.repo.Update(ctx, folder); err != nil {
		return err
	}

	// 更新所有子孙节点的 path 和 depth
	if err := s.repo.UpdateChildrenPathAndDepth(ctx, oldPath, newPath, depthDiff); err != nil {
		return err
	}

	return nil
}

// ReorderFolder 调整排序
func (s *Service) ReorderFolder(ctx context.Context, id uint, newOrder int) error {
	_, err := s.repo.FindByID(ctx, id)
	if err != nil {
		return err
	}
	return s.repo.UpdateSortOrder(ctx, id, newOrder)
}

// GetDescendantIDs 获取指定文件夹及其所有子孙的 ID 列表
// 用于实现"选择父分类时搜索所有子分类下的内容"功能
func (s *Service) GetDescendantIDs(ctx context.Context, id uint) ([]uint, error) {
	folder, err := s.repo.FindByID(ctx, id)
	if err != nil {
		return nil, err
	}

	// 使用 path 前缀查询所有子孙节点
	descendants, err := s.repo.FindByPath(ctx, folder.Path)
	if err != nil {
		return nil, err
	}

	ids := make([]uint, len(descendants))
	for i, d := range descendants {
		ids[i] = d.ID
	}

	return ids, nil
}

// validateName 验证名称
func (s *Service) validateName(name string) error {
	name = strings.TrimSpace(name)
	if name == "" {
		return ErrInvalidName
	}
	if len(name) > 255 {
		return ErrInvalidName
	}
	return nil
}

// IncrementItemCount 更新文件夹的子项数量
// 同时更新 item_count（直接子项）和祖先链的 total_item_count
func (s *Service) IncrementItemCount(ctx context.Context, folderID uint, delta int) error {
	folder, err := s.repo.FindByID(ctx, folderID)
	if err != nil {
		return err
	}

	// 1. 更新当前文件夹的 item_count
	if err := s.repo.IncrementItemCount(ctx, folderID, delta); err != nil {
		return err
	}

	// 2. 更新祖先链的 total_item_count
	return s.repo.IncrementTotalItemCount(ctx, folder.Path, delta)
}

// buildTree 构建树结构
func buildTree(folders []*model.Folder, rootParentID *uint) []*model.FolderNode {
	nodeMap := make(map[uint]*model.FolderNode)
	var roots []*model.FolderNode

	// 创建所有节点
	for _, f := range folders {
		nodeMap[f.ID] = f.ToNode()
	}

	// 建立父子关系
	for _, f := range folders {
		node := nodeMap[f.ID]
		if f.ParentID == nil {
			if rootParentID == nil {
				roots = append(roots, node)
			}
		} else if rootParentID != nil && *f.ParentID == *rootParentID {
			roots = append(roots, node)
		} else if parent, ok := nodeMap[*f.ParentID]; ok {
			parent.Children = append(parent.Children, node)
		}
	}

	return roots
}
