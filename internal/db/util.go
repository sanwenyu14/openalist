package db

import (
	"fmt"

	"github.com/AlliotTech/openalist/internal/conf"
	"gorm.io/gorm"
)

func columnName(name string) string {
	if conf.Conf.Database.Type == "postgres" {
		return fmt.Sprintf(`"%s"`, name)
	}
	return fmt.Sprintf("`%s`", name)
}

func addStorageOrder(db *gorm.DB) *gorm.DB {
	return db.Order(fmt.Sprintf("%s, %s", columnName("order"), columnName("id")))
}
