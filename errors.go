package folder

import "errors"

var (
	// ErrFolderNotFound 文件夹不存在
	ErrFolderNotFound = errors.New("folder not found")

	// ErrParentNotFound 父文件夹不存在
	ErrParentNotFound = errors.New("parent folder not found")

	// ErrCircularReference 循环引用（移动到自己的子节点下）
	ErrCircularReference = errors.New("cannot move folder to its own descendant")

	// ErrInvalidName 无效的文件夹名称
	ErrInvalidName = errors.New("invalid folder name")

	// ErrDuplicateName 同级目录下名称重复
	ErrDuplicateName = errors.New("folder name already exists in the same parent")

	// ErrMaxDepthExceeded 超过最大深度限制
	ErrMaxDepthExceeded = errors.New("max folder depth exceeded")

	// ErrHasChildren 文件夹下有子节点，无法删除
	ErrHasChildren = errors.New("cannot delete folder with children")
)
