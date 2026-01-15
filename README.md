# go-yogan-domain-folder

通用文件夹/层级管理领域包，支持动态表名，实现数据隔离的同时复用业务逻辑。

## 特性

- **动态表名**：通过 Repository 构造参数指定表名，支持多业务复用
- **层级管理**：parent_id + path（物化路径）方案
- **树形操作**：创建、删除、移动、排序、获取树结构
- **深度限制**：可配置最大层级深度

## 安装

```bash
go get github.com/KOMKZ/go-yogan-domain-folder
```

## 使用示例

```go
package main

import (
    folder "github.com/KOMKZ/go-yogan-domain-folder"
    "gorm.io/gorm"
)

func main() {
    var db *gorm.DB // 初始化数据库连接

    // 创建 Repository（指定表名）
    repo := folder.NewGormRepository(db, "article_folders")

    // 创建 Service
    svc := folder.NewService(repo)

    // 创建文件夹
    input := &folder.CreateFolderInput{
        Name:     "技术文章",
        ParentID: nil, // 根节点
    }
    newFolder, err := svc.CreateFolder(ctx, input)
    if err != nil {
        panic(err)
    }

    // 获取树结构
    tree, err := svc.GetTree(ctx)
    if err != nil {
        panic(err)
    }
}
```

## 多业务复用

```go
// 文章分类
articleRepo := folder.NewGormRepository(db, "article_folders")
articleSvc := folder.NewService(articleRepo)

// 产品分类
productRepo := folder.NewGormRepository(db, "product_categories")
productSvc := folder.NewService(productRepo)

// 文档目录
docRepo := folder.NewGormRepository(db, "document_directories")
docSvc := folder.NewService(docRepo)
```

## 表结构

各应用需要自行创建对应的表，表结构如下：

```sql
CREATE TABLE your_table_name (
    id BIGINT UNSIGNED AUTO_INCREMENT PRIMARY KEY,
    name VARCHAR(255) NOT NULL,
    parent_id BIGINT UNSIGNED,
    sort_order INT DEFAULT 0,
    depth INT DEFAULT 0,
    path VARCHAR(1000),
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP,
    deleted_at TIMESTAMP NULL,

    INDEX idx_parent_id (parent_id),
    INDEX idx_path (path(255)),
    INDEX idx_deleted_at (deleted_at)
);
```

## License

MIT
