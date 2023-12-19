package gormembedded

import "github.com/FrancoLiberali/cql/model"

type ToBeGormEmbedded struct {
	Int int
}

type GormEmbedded struct {
	model.UIntModel

	Int                  int
	GormEmbedded         ToBeGormEmbedded `gorm:"embedded;embeddedPrefix:gorm_embedded_"`
	GormEmbeddedNoPrefix ToBeGormEmbedded `gorm:"embedded"`
}
