package goembedded

import "github.com/ditrit/badaas/orm"

type ToBeEmbedded struct {
	Int int
}

type GoEmbedded struct {
	orm.UIntModel

	Int int
	ToBeEmbedded
}
