package db

import (
	"reflect"

	"github.com/goexl/gox"
	"github.com/goexl/gox/field"
	"github.com/goexl/log"
	"github.com/pangum/db/internal/config"
	"github.com/pangum/db/internal/db/internal/core"
	"xorm.io/xorm"
)

type Synchronizer struct {
	gox.CannotCopy

	engine *Engine
	config *config.DB
	logger log.Logger
}

func NewSynchronizer(engine *Engine, config *config.DB, logger log.Logger) *Synchronizer {
	return &Synchronizer{
		engine: engine,
		config: config,
		logger: logger,
	}
}

func (s *Synchronizer) Sync(models ...any) (err error) {
	if !s.config.Sync.Enabled {
		return
	}

	sync := s.config.Sync
	options := new(xorm.SyncOptions)
	if sync.Ignore.Indices {
		options.IgnoreIndices = true
	}
	if sync.Ignore.Constrains {
		options.IgnoreConstrains = true
	}

	fields := gox.Fields[any]{
		field.New("tables", s.tables(models...)),
		field.New("sync", sync),
		field.New("options", options),
	}
	s.logger.Info("同步数据库表开始", fields...)
	if _, se := s.engine.SyncWithOptions(*options, models...); nil == se {
		err = se
		s.logger.Error("同步数据库表失败", fields...)
	} else {
		s.logger.Info("同步数据库表成功", fields...)
	}

	return
}

func (s *Synchronizer) tables(models ...any) (tables []string) {
	tables = make([]string, 0, len(models))
	for _, model := range models {
		switch table := model.(type) {
		case core.Commenter:
			tables = append(tables, table.TableComment())
		case core.NameMaker:
			tables = append(tables, table.TableName())
		default:
			tables = append(tables, s.config.TableMapper().Obj2Table(s.getType(table)))
		}
	}

	return
}

func (s *Synchronizer) getType(table any) (typ string) {
	if of := reflect.TypeOf(table); of.Kind() == reflect.Ptr {
		typ = of.Elem().Name()
	} else {
		typ = of.Name()
	}

	return
}
