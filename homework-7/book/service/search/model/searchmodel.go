package model

import (
	"github.com/zeromicro/go-zero/core/stores/cache"
	"github.com/zeromicro/go-zero/core/stores/sqlx"
)

var _ SearchModel = (*customSearchModel)(nil)

type (
	// SearchModel is an interface to be customized, add more methods here,
	// and implement the added methods in customSearchModel.
	SearchModel interface {
		searchModel
	}

	customSearchModel struct {
		*defaultSearchModel
	}
)

// NewSearchModel returns a model for the database table.
func NewSearchModel(conn sqlx.SqlConn, c cache.CacheConf) SearchModel {
	return &customSearchModel{
		defaultSearchModel: newSearchModel(conn, c),
	}
}
