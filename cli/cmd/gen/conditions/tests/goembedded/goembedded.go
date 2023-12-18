package goembedded

import "github.com/ditrit/badaas/orm/model"

type ToBeEmbedded struct {
	Int int
}

type GoEmbedded struct {
	model.UIntModel

	Int int
	ToBeEmbedded
}
