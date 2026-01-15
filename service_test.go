package folder

import (
	"context"
	"testing"

	"github.com/KOMKZ/go-yogan-domain-folder/model"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

// MockRepository mock 实现
type MockRepository struct {
	mock.Mock
}

func (m *MockRepository) Create(ctx context.Context, folder *model.Folder) error {
	args := m.Called(ctx, folder)
	// 模拟创建后自动设置 ID
	if folder.ID == 0 {
		folder.ID = 1
	}
	return args.Error(0)
}

func (m *MockRepository) Update(ctx context.Context, folder *model.Folder) error {
	args := m.Called(ctx, folder)
	return args.Error(0)
}

func (m *MockRepository) Delete(ctx context.Context, id uint) error {
	args := m.Called(ctx, id)
	return args.Error(0)
}

func (m *MockRepository) FindByID(ctx context.Context, id uint) (*model.Folder, error) {
	args := m.Called(ctx, id)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*model.Folder), args.Error(1)
}

func (m *MockRepository) FindByParentID(ctx context.Context, parentID *uint) ([]*model.Folder, error) {
	args := m.Called(ctx, parentID)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).([]*model.Folder), args.Error(1)
}

func (m *MockRepository) FindChildren(ctx context.Context, parentID uint) ([]*model.Folder, error) {
	args := m.Called(ctx, parentID)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).([]*model.Folder), args.Error(1)
}

func (m *MockRepository) FindRoots(ctx context.Context) ([]*model.Folder, error) {
	args := m.Called(ctx)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).([]*model.Folder), args.Error(1)
}

func (m *MockRepository) FindByPath(ctx context.Context, pathPrefix string) ([]*model.Folder, error) {
	args := m.Called(ctx, pathPrefix)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).([]*model.Folder), args.Error(1)
}

func (m *MockRepository) FindAncestors(ctx context.Context, path string) ([]*model.Folder, error) {
	args := m.Called(ctx, path)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).([]*model.Folder), args.Error(1)
}

func (m *MockRepository) FindAll(ctx context.Context) ([]*model.Folder, error) {
	args := m.Called(ctx)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).([]*model.Folder), args.Error(1)
}

func (m *MockRepository) UpdateSortOrder(ctx context.Context, id uint, sortOrder int) error {
	args := m.Called(ctx, id, sortOrder)
	return args.Error(0)
}

func (m *MockRepository) FindMaxSortOrder(ctx context.Context, parentID *uint) (int, error) {
	args := m.Called(ctx, parentID)
	return args.Int(0), args.Error(1)
}

func (m *MockRepository) UpdatePathAndDepth(ctx context.Context, id uint, path string, depth int) error {
	args := m.Called(ctx, id, path, depth)
	return args.Error(0)
}

func (m *MockRepository) UpdateChildrenPathAndDepth(ctx context.Context, oldPathPrefix, newPathPrefix string, depthDiff int) error {
	args := m.Called(ctx, oldPathPrefix, newPathPrefix, depthDiff)
	return args.Error(0)
}

func (m *MockRepository) ExistsByNameAndParent(ctx context.Context, name string, parentID *uint, excludeID *uint) (bool, error) {
	args := m.Called(ctx, name, parentID, excludeID)
	return args.Bool(0), args.Error(1)
}

func (m *MockRepository) HasChildren(ctx context.Context, id uint) (bool, error) {
	args := m.Called(ctx, id)
	return args.Bool(0), args.Error(1)
}

// TestCreateFolder_Success 测试创建根文件夹成功
func TestCreateFolder_Success(t *testing.T) {
	mockRepo := new(MockRepository)
	svc := NewService(mockRepo)
	ctx := context.Background()

	// 设置期望
	mockRepo.On("ExistsByNameAndParent", ctx, "技术文章", (*uint)(nil), (*uint)(nil)).Return(false, nil)
	mockRepo.On("FindMaxSortOrder", ctx, (*uint)(nil)).Return(0, nil)
	mockRepo.On("Create", ctx, mock.AnythingOfType("*model.Folder")).Return(nil)
	mockRepo.On("Update", ctx, mock.AnythingOfType("*model.Folder")).Return(nil)

	// 执行
	input := &CreateFolderInput{
		Name:     "技术文章",
		ParentID: nil,
	}
	folder, err := svc.CreateFolder(ctx, input)

	// 断言
	assert.NoError(t, err)
	assert.NotNil(t, folder)
	assert.Equal(t, "技术文章", folder.Name)
	assert.Nil(t, folder.ParentID)
	assert.Equal(t, 0, folder.Depth)
	assert.Equal(t, 1, folder.SortOrder)
	mockRepo.AssertExpectations(t)
}

// TestCreateFolder_WithParent 测试创建子文件夹成功
func TestCreateFolder_WithParent(t *testing.T) {
	mockRepo := new(MockRepository)
	svc := NewService(mockRepo)
	ctx := context.Background()

	parentID := uint(1)
	parent := &model.Folder{
		ID:    1,
		Name:  "技术文章",
		Depth: 0,
		Path:  "/1/",
	}

	// 设置期望
	mockRepo.On("ExistsByNameAndParent", ctx, "Go语言", &parentID, (*uint)(nil)).Return(false, nil)
	mockRepo.On("FindByID", ctx, parentID).Return(parent, nil)
	mockRepo.On("FindMaxSortOrder", ctx, &parentID).Return(0, nil)
	mockRepo.On("Create", ctx, mock.AnythingOfType("*model.Folder")).Return(nil)
	mockRepo.On("Update", ctx, mock.AnythingOfType("*model.Folder")).Return(nil)

	// 执行
	input := &CreateFolderInput{
		Name:     "Go语言",
		ParentID: &parentID,
	}
	folder, err := svc.CreateFolder(ctx, input)

	// 断言
	assert.NoError(t, err)
	assert.NotNil(t, folder)
	assert.Equal(t, "Go语言", folder.Name)
	assert.Equal(t, &parentID, folder.ParentID)
	assert.Equal(t, 1, folder.Depth)
	mockRepo.AssertExpectations(t)
}

// TestCreateFolder_DuplicateName 测试名称重复
func TestCreateFolder_DuplicateName(t *testing.T) {
	mockRepo := new(MockRepository)
	svc := NewService(mockRepo)
	ctx := context.Background()

	// 设置期望
	mockRepo.On("ExistsByNameAndParent", ctx, "技术文章", (*uint)(nil), (*uint)(nil)).Return(true, nil)

	// 执行
	input := &CreateFolderInput{
		Name:     "技术文章",
		ParentID: nil,
	}
	folder, err := svc.CreateFolder(ctx, input)

	// 断言
	assert.ErrorIs(t, err, ErrDuplicateName)
	assert.Nil(t, folder)
	mockRepo.AssertExpectations(t)
}

// TestCreateFolder_InvalidName 测试无效名称
func TestCreateFolder_InvalidName(t *testing.T) {
	tests := []struct {
		name     string
		input    string
		wantErr  bool
		errorMsg string
	}{
		{"empty", "", true, "empty name"},
		{"spaces_only", "   ", true, "spaces only"},
		{"valid", "技术文章", false, "valid name"},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			mockRepo := new(MockRepository)
			svc := NewService(mockRepo)
			ctx := context.Background()

			if !tt.wantErr {
				mockRepo.On("ExistsByNameAndParent", ctx, tt.input, (*uint)(nil), (*uint)(nil)).Return(false, nil)
				mockRepo.On("FindMaxSortOrder", ctx, (*uint)(nil)).Return(0, nil)
				mockRepo.On("Create", ctx, mock.AnythingOfType("*model.Folder")).Return(nil)
				mockRepo.On("Update", ctx, mock.AnythingOfType("*model.Folder")).Return(nil)
			}

			input := &CreateFolderInput{
				Name:     tt.input,
				ParentID: nil,
			}
			_, err := svc.CreateFolder(ctx, input)

			if tt.wantErr {
				assert.ErrorIs(t, err, ErrInvalidName)
			} else {
				assert.NoError(t, err)
			}
		})
	}
}

// TestUpdateFolder_Success 测试更新文件夹成功
func TestUpdateFolder_Success(t *testing.T) {
	mockRepo := new(MockRepository)
	svc := NewService(mockRepo)
	ctx := context.Background()

	existingFolder := &model.Folder{
		ID:       1,
		Name:     "旧名称",
		ParentID: nil,
		Depth:    0,
		Path:     "/1/",
	}

	id := uint(1)
	mockRepo.On("FindByID", ctx, id).Return(existingFolder, nil)
	mockRepo.On("ExistsByNameAndParent", ctx, "新名称", (*uint)(nil), &id).Return(false, nil)
	mockRepo.On("Update", ctx, mock.AnythingOfType("*model.Folder")).Return(nil)

	input := &UpdateFolderInput{
		ID:   1,
		Name: "新名称",
	}
	folder, err := svc.UpdateFolder(ctx, input)

	assert.NoError(t, err)
	assert.Equal(t, "新名称", folder.Name)
	mockRepo.AssertExpectations(t)
}

// TestDeleteFolder_Success 测试删除文件夹成功
func TestDeleteFolder_Success(t *testing.T) {
	mockRepo := new(MockRepository)
	svc := NewService(mockRepo)
	ctx := context.Background()

	existingFolder := &model.Folder{
		ID:   1,
		Name: "技术文章",
	}

	mockRepo.On("FindByID", ctx, uint(1)).Return(existingFolder, nil)
	mockRepo.On("HasChildren", ctx, uint(1)).Return(false, nil)
	mockRepo.On("Delete", ctx, uint(1)).Return(nil)

	err := svc.DeleteFolder(ctx, 1)

	assert.NoError(t, err)
	mockRepo.AssertExpectations(t)
}

// TestDeleteFolder_HasChildren 测试删除有子节点的文件夹
func TestDeleteFolder_HasChildren(t *testing.T) {
	mockRepo := new(MockRepository)
	svc := NewService(mockRepo)
	ctx := context.Background()

	existingFolder := &model.Folder{
		ID:   1,
		Name: "技术文章",
	}

	mockRepo.On("FindByID", ctx, uint(1)).Return(existingFolder, nil)
	mockRepo.On("HasChildren", ctx, uint(1)).Return(true, nil)

	err := svc.DeleteFolder(ctx, 1)

	assert.Error(t, err)
	assert.Contains(t, err.Error(), "children")
	mockRepo.AssertExpectations(t)
}

// TestGetTree_Success 测试获取树结构
func TestGetTree_Success(t *testing.T) {
	mockRepo := new(MockRepository)
	svc := NewService(mockRepo)
	ctx := context.Background()

	parentID := uint(1)
	folders := []*model.Folder{
		{ID: 1, Name: "技术文章", ParentID: nil, Depth: 0, Path: "/1/"},
		{ID: 2, Name: "Go语言", ParentID: &parentID, Depth: 1, Path: "/1/2/"},
		{ID: 3, Name: "Python", ParentID: &parentID, Depth: 1, Path: "/1/3/"},
	}

	mockRepo.On("FindAll", ctx).Return(folders, nil)

	tree, err := svc.GetTree(ctx)

	assert.NoError(t, err)
	assert.Len(t, tree, 1) // 只有一个根节点
	assert.Equal(t, "技术文章", tree[0].Name)
	assert.Len(t, tree[0].Children, 2) // 两个子节点
	mockRepo.AssertExpectations(t)
}

// TestMoveFolder_Success 测试移动文件夹成功
func TestMoveFolder_Success(t *testing.T) {
	mockRepo := new(MockRepository)
	svc := NewService(mockRepo)
	ctx := context.Background()

	folder := &model.Folder{
		ID:       2,
		Name:     "Go语言",
		ParentID: nil,
		Depth:    0,
		Path:     "/2/",
	}

	newParentID := uint(1)
	newParent := &model.Folder{
		ID:    1,
		Name:  "技术文章",
		Depth: 0,
		Path:  "/1/",
	}

	mockRepo.On("FindByID", ctx, uint(2)).Return(folder, nil)
	mockRepo.On("FindByID", ctx, newParentID).Return(newParent, nil)
	mockRepo.On("FindByPath", ctx, "/2/").Return([]*model.Folder{folder}, nil)
	mockRepo.On("FindMaxSortOrder", ctx, &newParentID).Return(0, nil)
	mockRepo.On("Update", ctx, mock.AnythingOfType("*model.Folder")).Return(nil)
	mockRepo.On("UpdateChildrenPathAndDepth", ctx, "/2/", "/1/2/", 1).Return(nil)

	err := svc.MoveFolder(ctx, 2, &newParentID)

	assert.NoError(t, err)
	mockRepo.AssertExpectations(t)
}

// TestMoveFolder_CircularReference 测试循环引用
func TestMoveFolder_CircularReference(t *testing.T) {
	mockRepo := new(MockRepository)
	svc := NewService(mockRepo)
	ctx := context.Background()

	folder := &model.Folder{
		ID:    1,
		Name:  "技术文章",
		Depth: 0,
		Path:  "/1/",
	}

	childID := uint(2)
	child := &model.Folder{
		ID:    2,
		Name:  "Go语言",
		Depth: 1,
		Path:  "/1/2/",
	}

	mockRepo.On("FindByID", ctx, uint(1)).Return(folder, nil)
	mockRepo.On("FindByID", ctx, childID).Return(child, nil)

	// 尝试将父节点移动到子节点下
	err := svc.MoveFolder(ctx, 1, &childID)

	assert.ErrorIs(t, err, ErrCircularReference)
	mockRepo.AssertExpectations(t)
}

// TestReorderFolder_Success 测试排序成功
func TestReorderFolder_Success(t *testing.T) {
	mockRepo := new(MockRepository)
	svc := NewService(mockRepo)
	ctx := context.Background()

	folder := &model.Folder{
		ID:        1,
		Name:      "技术文章",
		SortOrder: 1,
	}

	mockRepo.On("FindByID", ctx, uint(1)).Return(folder, nil)
	mockRepo.On("UpdateSortOrder", ctx, uint(1), 5).Return(nil)

	err := svc.ReorderFolder(ctx, 1, 5)

	assert.NoError(t, err)
	mockRepo.AssertExpectations(t)
}

// TestGetAncestors_Success 测试获取祖先节点
func TestGetAncestors_Success(t *testing.T) {
	mockRepo := new(MockRepository)
	svc := NewService(mockRepo)
	ctx := context.Background()

	folder := &model.Folder{
		ID:    3,
		Name:  "并发编程",
		Depth: 2,
		Path:  "/1/2/3/",
	}

	parentID := uint(1)
	ancestors := []*model.Folder{
		{ID: 1, Name: "技术文章", Depth: 0, Path: "/1/"},
		{ID: 2, Name: "Go语言", ParentID: &parentID, Depth: 1, Path: "/1/2/"},
	}

	mockRepo.On("FindByID", ctx, uint(3)).Return(folder, nil)
	mockRepo.On("FindAncestors", ctx, "/1/2/3/").Return(ancestors, nil)

	result, err := svc.GetAncestors(ctx, 3)

	assert.NoError(t, err)
	assert.Len(t, result, 2)
	mockRepo.AssertExpectations(t)
}

// TestGetFolder_Success 测试获取单个文件夹
func TestGetFolder_Success(t *testing.T) {
	mockRepo := new(MockRepository)
	svc := NewService(mockRepo)
	ctx := context.Background()

	folder := &model.Folder{
		ID:   1,
		Name: "技术文章",
	}

	mockRepo.On("FindByID", ctx, uint(1)).Return(folder, nil)

	result, err := svc.GetFolder(ctx, 1)

	assert.NoError(t, err)
	assert.Equal(t, "技术文章", result.Name)
	mockRepo.AssertExpectations(t)
}

// TestGetFolder_NotFound 测试获取不存在的文件夹
func TestGetFolder_NotFound(t *testing.T) {
	mockRepo := new(MockRepository)
	svc := NewService(mockRepo)
	ctx := context.Background()

	mockRepo.On("FindByID", ctx, uint(999)).Return(nil, ErrFolderNotFound)

	result, err := svc.GetFolder(ctx, 999)

	assert.ErrorIs(t, err, ErrFolderNotFound)
	assert.Nil(t, result)
	mockRepo.AssertExpectations(t)
}

// TestGetChildren_Success 测试获取子节点
func TestGetChildren_Success(t *testing.T) {
	mockRepo := new(MockRepository)
	svc := NewService(mockRepo)
	ctx := context.Background()

	parentID := uint(1)
	children := []*model.Folder{
		{ID: 2, Name: "Go语言", ParentID: &parentID},
		{ID: 3, Name: "Python", ParentID: &parentID},
	}

	mockRepo.On("FindByParentID", ctx, &parentID).Return(children, nil)

	result, err := svc.GetChildren(ctx, &parentID)

	assert.NoError(t, err)
	assert.Len(t, result, 2)
	mockRepo.AssertExpectations(t)
}

// TestGetSubTree_Success 测试获取子树
func TestGetSubTree_Success(t *testing.T) {
	mockRepo := new(MockRepository)
	svc := NewService(mockRepo)
	ctx := context.Background()

	root := &model.Folder{
		ID:    1,
		Name:  "技术文章",
		Depth: 0,
		Path:  "/1/",
	}

	parentID := uint(1)
	descendants := []*model.Folder{
		{ID: 1, Name: "技术文章", Depth: 0, Path: "/1/"},
		{ID: 2, Name: "Go语言", ParentID: &parentID, Depth: 1, Path: "/1/2/"},
	}

	mockRepo.On("FindByID", ctx, uint(1)).Return(root, nil)
	mockRepo.On("FindByPath", ctx, "/1/").Return(descendants, nil)

	result, err := svc.GetSubTree(ctx, 1)

	assert.NoError(t, err)
	assert.Len(t, result, 1) // Go语言是子节点
	mockRepo.AssertExpectations(t)
}

// TestUpdateFolder_NotFound 测试更新不存在的文件夹
func TestUpdateFolder_NotFound(t *testing.T) {
	mockRepo := new(MockRepository)
	svc := NewService(mockRepo)
	ctx := context.Background()

	mockRepo.On("FindByID", ctx, uint(999)).Return(nil, ErrFolderNotFound)

	input := &UpdateFolderInput{
		ID:   999,
		Name: "新名称",
	}
	folder, err := svc.UpdateFolder(ctx, input)

	assert.ErrorIs(t, err, ErrFolderNotFound)
	assert.Nil(t, folder)
	mockRepo.AssertExpectations(t)
}

// TestDeleteFolder_NotFound 测试删除不存在的文件夹
func TestDeleteFolder_NotFound(t *testing.T) {
	mockRepo := new(MockRepository)
	svc := NewService(mockRepo)
	ctx := context.Background()

	mockRepo.On("FindByID", ctx, uint(999)).Return(nil, ErrFolderNotFound)

	err := svc.DeleteFolder(ctx, 999)

	assert.ErrorIs(t, err, ErrFolderNotFound)
	mockRepo.AssertExpectations(t)
}

// TestMoveFolder_ToSelf 测试移动到自身
func TestMoveFolder_ToSelf(t *testing.T) {
	mockRepo := new(MockRepository)
	svc := NewService(mockRepo)
	ctx := context.Background()

	folder := &model.Folder{
		ID:    1,
		Name:  "技术文章",
		Depth: 0,
		Path:  "/1/",
	}

	selfID := uint(1)
	mockRepo.On("FindByID", ctx, uint(1)).Return(folder, nil)

	err := svc.MoveFolder(ctx, 1, &selfID)

	assert.ErrorIs(t, err, ErrCircularReference)
	mockRepo.AssertExpectations(t)
}

// TestMoveFolder_ToRoot 测试移动到根节点
func TestMoveFolder_ToRoot(t *testing.T) {
	mockRepo := new(MockRepository)
	svc := NewService(mockRepo)
	ctx := context.Background()

	parentID := uint(1)
	folder := &model.Folder{
		ID:       2,
		Name:     "Go语言",
		ParentID: &parentID,
		Depth:    1,
		Path:     "/1/2/",
	}

	mockRepo.On("FindByID", ctx, uint(2)).Return(folder, nil)
	mockRepo.On("FindByPath", ctx, "/1/2/").Return([]*model.Folder{folder}, nil)
	mockRepo.On("FindMaxSortOrder", ctx, (*uint)(nil)).Return(0, nil)
	mockRepo.On("Update", ctx, mock.AnythingOfType("*model.Folder")).Return(nil)
	mockRepo.On("UpdateChildrenPathAndDepth", ctx, "/1/2/", "/2/", -1).Return(nil)

	err := svc.MoveFolder(ctx, 2, nil)

	assert.NoError(t, err)
	mockRepo.AssertExpectations(t)
}

// TestCreateFolder_ParentNotFound 测试父节点不存在
func TestCreateFolder_ParentNotFound(t *testing.T) {
	mockRepo := new(MockRepository)
	svc := NewService(mockRepo)
	ctx := context.Background()

	parentID := uint(999)
	mockRepo.On("ExistsByNameAndParent", ctx, "子文件夹", &parentID, (*uint)(nil)).Return(false, nil)
	mockRepo.On("FindByID", ctx, parentID).Return(nil, ErrFolderNotFound)

	input := &CreateFolderInput{
		Name:     "子文件夹",
		ParentID: &parentID,
	}
	folder, err := svc.CreateFolder(ctx, input)

	assert.ErrorIs(t, err, ErrParentNotFound)
	assert.Nil(t, folder)
	mockRepo.AssertExpectations(t)
}

// TestNewServiceWithConfig 测试带配置的服务创建
func TestNewServiceWithConfig(t *testing.T) {
	mockRepo := new(MockRepository)
	config := ServiceConfig{MaxDepth: 5}
	svc := NewServiceWithConfig(mockRepo, config)

	assert.NotNil(t, svc)
	assert.Equal(t, 5, svc.config.MaxDepth)
}

// TestCreateFolder_MaxDepthExceeded 测试超过最大深度
func TestCreateFolder_MaxDepthExceeded(t *testing.T) {
	mockRepo := new(MockRepository)
	config := ServiceConfig{MaxDepth: 2}
	svc := NewServiceWithConfig(mockRepo, config)
	ctx := context.Background()

	parentID := uint(1)
	parent := &model.Folder{
		ID:    1,
		Name:  "技术文章",
		Depth: 1, // 已经是深度1，再创建子节点会达到深度2，等于MaxDepth
		Path:  "/1/",
	}

	mockRepo.On("ExistsByNameAndParent", ctx, "子文件夹", &parentID, (*uint)(nil)).Return(false, nil)
	mockRepo.On("FindByID", ctx, parentID).Return(parent, nil)

	input := &CreateFolderInput{
		Name:     "子文件夹",
		ParentID: &parentID,
	}
	folder, err := svc.CreateFolder(ctx, input)

	assert.ErrorIs(t, err, ErrMaxDepthExceeded)
	assert.Nil(t, folder)
	mockRepo.AssertExpectations(t)
}

// TestBuildTree_Empty 测试空列表构建树
func TestBuildTree_Empty(t *testing.T) {
	result := buildTree([]*model.Folder{}, nil)
	assert.Len(t, result, 0)
}

// TestBuildTree_MultiLevel 测试多级树构建
func TestBuildTree_MultiLevel(t *testing.T) {
	parentID1 := uint(1)
	parentID2 := uint(2)
	folders := []*model.Folder{
		{ID: 1, Name: "根1", ParentID: nil, Depth: 0},
		{ID: 2, Name: "子1-1", ParentID: &parentID1, Depth: 1},
		{ID: 3, Name: "子1-2", ParentID: &parentID1, Depth: 1},
		{ID: 4, Name: "孙1-1-1", ParentID: &parentID2, Depth: 2},
		{ID: 5, Name: "根2", ParentID: nil, Depth: 0},
	}

	result := buildTree(folders, nil)

	assert.Len(t, result, 2) // 两个根节点
	assert.Equal(t, "根1", result[0].Name)
	assert.Len(t, result[0].Children, 2) // 根1有2个子节点
	assert.Len(t, result[0].Children[0].Children, 1) // 子1-1有1个孙节点
	assert.Equal(t, "根2", result[1].Name)
}
