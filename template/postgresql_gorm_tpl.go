package template

import (
	"fmt"
	"github.com/youngchan1988/gocommon"
	"github.com/youngchan1988/gocommon/stringutils"
	"ormgenc/interviewer"
)

const PostgresqlGormTpl = `
// Package dbmodel Generated code. DO NOT modify by hand!
package dbmodel

import (
	"errors"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
	"sync"
	"time"
)

var once sync.Once

type gormDB struct {
	db *gorm.DB
}

var instance *gormDB

func Gorm() *gormDB {
	once.Do(func() {
		instance = &gormDB{}
	})
	return instance
}

func (g *gormDB) Open(dsn string) error {
	var err error
	g.db, err = gorm.Open(postgres.New(postgres.Config{
		DSN:                  dsn,
		PreferSimpleProtocol: true,
	}), &gorm.Config{})
	if err == nil {
		sqlDB, _ := g.db.DB()
		sqlDB.SetMaxIdleConns(10)
		sqlDB.SetMaxOpenConns(100)
		sqlDB.SetConnMaxLifetime(time.Hour)
	}
	return err
}

func (g *gormDB) Close() error {
	if g.db != nil {
		sqlDB, _ := g.db.DB()
		return sqlDB.Close()
	}
	return errors.New("database instance is nil")
}
`

const PostgresqlGormModelTpl = `
package dbmodel

import (
	"errors"
	"fmt"
	"gorm.io/gorm"
)

type {{model_name}}DBModel struct {
	db *gorm.DB
	{{model_fields}}
}

func (g *gormDB){{model_name}}() *{{model_name}}DBModel {
	if g.db == nil {
		panic(errors.New("unresolved gormDB.db is nil"))
	}
	g.db = g.db.Table("{{model_underscore_name}}")
	m := &{{model_name}}DBModel{db: g.db}
	return m
}

func (m *{{model_name}}DBModel) InsertOne(model *{{model_name}}DBModel) error {
	if m.db == nil {
		panic(errors.New("unresolved gormDB.db is nil"))
	}
	return m.db.Create(model).Error
}

func (m *{{model_name}}DBModel) InsertMany(models []*{{model_name}}DBModel) error {
	if m.db == nil {
		panic(errors.New("unresolved gormDB.db is nil"))
	}
	return m.db.Create(models).Error
}

func (m *{{model_name}}DBModel) Update(model *{{model_name}}DBModel) error {
	if m.db == nil {
		panic(errors.New("unresolved gormDB.db is nil"))
	}
	return m.db.Updates(model).Error
}

func (m *{{model_name}}DBModel) UpdateColumn(name string, value interface{}) error {
	if m.db == nil {
		panic(errors.New("unresolved gormDB.db is nil"))
	}
	return m.db.Update(name, value).Error
}

func (m *{{model_name}}DBModel) Delete() error {
	if m.db == nil {
		panic(errors.New("unresolved gormDB.db is nil"))
	}
	return m.db.Delete(&{{model_name}}DBModel{}).Error
}

func (m *{{model_name}}DBModel) DeleteOne(id uint) error {
	if m.db == nil {
		panic(errors.New("unresolved gormDB.db is nil"))
	}
	return m.db.Delete(&{{model_name}}DBModel{}, id).Error
}

func (m *{{model_name}}DBModel) DeleteMany(ids []uint) error {
	if m.db == nil {
		panic(errors.New("unresolved gormDB.db is nil"))
	}
	return m.db.Delete(&{{model_name}}DBModel{}, ids).Error
}

func (m *{{model_name}}DBModel) FindOne() (*{{model_name}}DBModel, error) {
	if m.db == nil {
		panic(errors.New("unresolved gormDB.db is nil"))
	}
	model := &{{model_name}}DBModel{}
	err := m.db.First(model).Error
	return model, err
}

func (m *{{model_name}}DBModel) FindMany() ([]*{{model_name}}DBModel, error) {
	if m.db == nil {
		panic(errors.New("unresolved gormDB.db is nil"))
	}
	models := make([]*{{model_name}}DBModel, 0)
	err := m.db.Find(&models).Error
	return models, err
}

func (m *{{model_name}}DBModel) Select(columns ...interface{}) *{{model_name}}DBModel {
	if m.db == nil {
		panic(errors.New("unresolved gormDB.db is nil"))
	}
	m.db = m.db.Select(columns)
	return m
}

func (m *{{model_name}}DBModel) ByID(id uint) *{{model_name}}DBModel {
	if m.db == nil {
		panic(errors.New("unresolved gormDB.db is nil"))
	}
	m.db = m.db.Where("id = ?", id)
	return m
}

func (m *{{model_name}}DBModel) Where(query interface{}, args ...interface{}) *{{model_name}}DBModel {
	if m.db == nil {
		panic(errors.New("unresolved gormDB.db is nil"))
	}
	m.db = m.db.Where(query, args)
	return m
}

func (m *{{model_name}}DBModel) Limit(limit int) *{{model_name}}DBModel {
	if m.db == nil {
		panic(errors.New("unresolved gormDB.db is nil"))
	}
	m.db = m.db.Limit(limit)
	return m
}

func (m *{{model_name}}DBModel) Offset(offset int) *{{model_name}}DBModel {
	if m.db == nil {
		panic(errors.New("unresolved gormDB.db is nil"))
	}
	m.db = m.db.Offset(offset)
	return m
}

func (m *{{model_name}}DBModel) OrderBy(column string, sort string) *{{model_name}}DBModel {
	if m.db == nil {
		panic(errors.New("unresolved gormDB.db is nil"))
	}
	m.db = m.db.Order(fmt.Sprintf("%s %s", column, sort))
	return m
}

func (m *{{model_name}}DBModel) GroupBy(column string) *{{model_name}}DBModel {
	if m.db == nil {
		panic(errors.New("unresolved gormDB.db is nil"))
	}
	m.db = m.db.Group(column)
	return m
}

func (m *{{model_name}}DBModel) Having(query interface{}, args ...interface{}) *{{model_name}}DBModel {
	if m.db == nil {
		panic(errors.New("unresolved gormDB.db is nil"))
	}
	m.db = m.db.Having(query, args)
	return m
}

func (m *{{model_name}}DBModel) Count() (int64, error) {
	if m.db == nil {
		panic(errors.New("unresolved gormDB.db is nil"))
	}
	var count int64
	r := m.db.Count(&count)
	return count, r.Error
}
`

func PostgresqlGormField(tableColumn *interviewer.DBTableColumn) string {
	name := stringutils.ToCamelCase(tableColumn.Field)

	switch tableColumn.Type {
	case "int2":
	case "int8":
		if tableColumn.Pk {
			return fmt.Sprintf("//%s\n %s  int `gorm:\"column:%s\";primaryKey`", tableColumn.Comment, name, tableColumn.Field)
		} else {
			return fmt.Sprintf("//%s\n %s  int `gorm:\"column:%s\"`", tableColumn.Comment, name, tableColumn.Field)
		}

	case "varchar":
		if tableColumn.Pk {
			return fmt.Sprintf("//%s\n %s string `gorm:\"column:%s\";primaryKey`", tableColumn.Comment, name, tableColumn.Field)
		} else {
			return fmt.Sprintf("//%s\n %s string `gorm:\"column:%s\"`", tableColumn.Comment, name, tableColumn.Field)
		}

	case "timestamptz":
	case "timestamp":
		if gocommon.IsContains(tableColumn.Field, "created") {
			return fmt.Sprintf("//%s\n %s time.Time  `gorm:\"autoCreateTime\";column:%s`", tableColumn.Comment, name, tableColumn.Field)
		} else if gocommon.IsContains(tableColumn.Field, "updated") {
			return fmt.Sprintf("//%s\n %s time.Time  `gorm:\"autoUpdateTime\";column:%s`", tableColumn.Comment, name, tableColumn.Field)
		} else if gocommon.IsContains(tableColumn.Field, "deleted") {
			return fmt.Sprintf("//%s\n %s gorm.DeletedAt  `gorm:\"index\";column:%s`", tableColumn.Comment, name, tableColumn.Field)
		}
	}

	return ""
}
