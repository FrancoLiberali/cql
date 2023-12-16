package gormembedded

import "github.com/ditrit/badaas/orm"

type ToBeGormEmbedded struct {
	Int int
}

type GormEmbedded struct {
	orm.UIntModel

	Int                  int
	GormEmbedded         ToBeGormEmbedded `gorm:"embedded;embeddedPrefix:gorm_embedded_"`
	GormEmbeddedNoPrefix ToBeGormEmbedded `gorm:"embedded"`
}
