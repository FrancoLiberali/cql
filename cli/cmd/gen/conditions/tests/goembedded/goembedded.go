package goembedded

import "github.com/ditrit/badaas/orm"

type ToBeEmbedded struct {
	EmbeddedInt int
}

type GoEmbedded struct {
	orm.UIntModel

	ToBeEmbedded
}
