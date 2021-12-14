package contract

import (
	"gorm.io/gorm"
)

const ORMKey = "arms:orm"

type ORM interface {
	GetDB(option ...DBOption) (*gorm.DB, error)
}

type DBOption func(orm ORM)
