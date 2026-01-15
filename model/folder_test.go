package model

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

// TestFolder_ToNode 测试转换为节点
func TestFolder_ToNode(t *testing.T) {
	parentID := uint(1)
	folder := &Folder{
		ID:        2,
		Name:      "Go语言",
		ParentID:  &parentID,
		SortOrder: 1,
		Depth:     1,
	}

	node := folder.ToNode()

	assert.Equal(t, uint(2), node.ID)
	assert.Equal(t, "Go语言", node.Name)
	assert.Equal(t, &parentID, node.ParentID)
	assert.Equal(t, 1, node.SortOrder)
	assert.Equal(t, 1, node.Depth)
	assert.Nil(t, node.Children)
}

// TestFolder_ToNode_Root 测试根节点转换
func TestFolder_ToNode_Root(t *testing.T) {
	folder := &Folder{
		ID:        1,
		Name:      "技术文章",
		ParentID:  nil,
		SortOrder: 0,
		Depth:     0,
	}

	node := folder.ToNode()

	assert.Equal(t, uint(1), node.ID)
	assert.Equal(t, "技术文章", node.Name)
	assert.Nil(t, node.ParentID)
	assert.Equal(t, 0, node.SortOrder)
	assert.Equal(t, 0, node.Depth)
}
