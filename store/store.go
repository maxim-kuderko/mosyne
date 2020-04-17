package store

import "github.com/maxim-kuderko/mosyne/entities"

type Interface interface {
	ZGet(get entities.ZGetRequest) entities.ZGetResponse
	ZSet(get entities.ZSetRequest) entities.ZSetResponse
}

type Options struct {
}
