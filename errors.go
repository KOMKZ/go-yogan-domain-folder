package folder

import (
	"net/http"

	"github.com/KOMKZ/go-yogan-framework/errcode"
)

// 模块码：60 (folder)
const ModuleFolder = 60

// 领域错误定义
var (
	// ErrNotFound 文件夹不存在
	ErrNotFound = errcode.Register(errcode.New(
		ModuleFolder, 1001,
		"folder",
		"error.folder.not_found",
		"分类不存在",
		http.StatusNotFound,
	))

	// ErrParentNotFound 父文件夹不存在
	ErrParentNotFound = errcode.Register(errcode.New(
		ModuleFolder, 1002,
		"folder",
		"error.folder.parent_not_found",
		"父分类不存在",
		http.StatusBadRequest,
	))

	// ErrCircularReference 循环引用
	ErrCircularReference = errcode.Register(errcode.New(
		ModuleFolder, 1003,
		"folder",
		"error.folder.circular_reference",
		"不能移动到自己的子分类下",
		http.StatusBadRequest,
	))

	// ErrInvalidName 无效名称
	ErrInvalidName = errcode.Register(errcode.New(
		ModuleFolder, 1004,
		"folder",
		"error.folder.invalid_name",
		"无效的分类名称",
		http.StatusBadRequest,
	))

	// ErrDuplicateName 名称重复
	ErrDuplicateName = errcode.Register(errcode.New(
		ModuleFolder, 1005,
		"folder",
		"error.folder.duplicate_name",
		"同级目录下已存在同名分类",
		http.StatusConflict,
	))

	// ErrMaxDepthExceeded 超过深度限制
	ErrMaxDepthExceeded = errcode.Register(errcode.New(
		ModuleFolder, 1006,
		"folder",
		"error.folder.max_depth_exceeded",
		"超过最大层级限制",
		http.StatusBadRequest,
	))

	// ErrHasChildren 有子节点
	ErrHasChildren = errcode.Register(errcode.New(
		ModuleFolder, 1007,
		"folder",
		"error.folder.has_children",
		"该分类下有子分类，请先删除子分类",
		http.StatusBadRequest,
	))
)
