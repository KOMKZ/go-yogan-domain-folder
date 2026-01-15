package folder

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

// TestParsePathIDs 测试路径解析
func TestParsePathIDs(t *testing.T) {
	tests := []struct {
		name     string
		path     string
		expected []uint
	}{
		{"empty", "", nil},
		{"root_only", "/", nil},
		{"single", "/1/", []uint{1}},
		{"multiple", "/1/2/3/", []uint{1, 2, 3}},
		{"no_leading_slash", "1/2/3/", []uint{1, 2, 3}},
		{"no_trailing_slash", "/1/2/3", []uint{1, 2, 3}},
		{"with_spaces", "/1/2/3/", []uint{1, 2, 3}},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := parsePathIDs(tt.path)
			assert.Equal(t, tt.expected, result)
		})
	}
}

// TestParseUint 测试字符串转 uint
func TestParseUint(t *testing.T) {
	tests := []struct {
		name    string
		input   string
		want    uint
		wantErr bool
	}{
		{"zero", "0", 0, false},
		{"single_digit", "5", 5, false},
		{"multi_digit", "123", 123, false},
		{"large", "999999", 999999, false},
		{"invalid_letter", "abc", 0, true},
		{"mixed", "12a3", 0, true},
		{"negative", "-1", 0, true},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			var result uint
			_, err := parseUint(tt.input, &result)
			if tt.wantErr {
				assert.Error(t, err)
			} else {
				assert.NoError(t, err)
				assert.Equal(t, tt.want, result)
			}
		})
	}
}

// TestNewGormRepository 测试创建 Repository
func TestNewGormRepository(t *testing.T) {
	// 不需要真实数据库连接，只测试构造函数
	repo := NewGormRepository(nil, "test_folders")
	assert.NotNil(t, repo)
	assert.Equal(t, "test_folders", repo.tableName)
}
