package postgres

import "gorm.io/gorm"

func preloadHierarchy(d *gorm.DB) *gorm.DB {
	return d.Preload("Childrens", preloadHierarchy)
}
