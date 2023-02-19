package base_consts

import "github.com/kysion/base-library/base_model"

type global struct {
	OrmCacheConf []*base_model.TableCacheConf
}

var (
	Global = global{
		OrmCacheConf: []*base_model.TableCacheConf{},
	}
)
