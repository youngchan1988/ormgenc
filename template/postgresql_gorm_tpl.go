package template

import (
	"fmt"
	"github.com/youngchan1988/gocommon"
	"github.com/youngchan1988/gocommon/stringutils"
	"ormgenc/interviewer"
)

const PostgresqlGormTpl = `
// Package dbmodel 
// Generated code. DO NOT modify by hand!
package dbmodel

import (
	"database/sql"
	"errors"
	"github.com/youngchan1988/gocommon/log"
	"reflect"
	"sync"
	"time"

	"github.com/jackc/pgtype"
	"github.com/youngchan1988/gocommon"
	"github.com/youngchan1988/gocommon/cast"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

var once sync.Once

type GormDB struct {
	db *gorm.DB
}

var instance *GormDB
var mutex sync.Mutex

//Gorm gromDB实例
func Gorm() *GormDB {
	once.Do(func() {
		instance = &GormDB{}
	})
	return instance
}

//Open 连接数据库，程序启动时调用
func (g *GormDB) Open(dsn string, idleConn int, openConn int) error {
	var err error
	g.db, err = gorm.Open(postgres.New(postgres.Config{
		DSN:                  dsn,
		PreferSimpleProtocol: true,
	}), &gorm.Config{})
	if err == nil {
		sqlDB, _ := g.db.DB()
		sqlDB.SetMaxIdleConns(idleConn)
		sqlDB.SetMaxOpenConns(openConn)
		sqlDB.SetConnMaxLifetime(time.Hour)
	}
	return err
}

//Close 关闭数据库
func (g *GormDB) Close() error {
	if g.db != nil {
		sqlDB, _ := g.db.DB()
		return sqlDB.Close()
	}
	return errors.New("database instance is nil")
}

//Transaction 使用事务, 加了同步信号锁
func (g *GormDB) Transaction(body func(tx *GormDB) error, opts ...*sql.TxOptions) error {
	defer mutex.Unlock()
	mutex.Lock()

	err := g.db.Transaction(func(tx *gorm.DB) error {
		gormTx := &GormDB{
			db: tx,
		}
		return body(gormTx)
	}, opts...)
	if err != nil {
		log.Error("dbmodel", err, 1, "transaction error")
		return err
	}

	return nil
}

//EntityToModel 根据tag：'ormgen'，将entity struct 结构转换db model struct 结构
func EntityToModel(entity interface{}, model interface{}) error {
	et := reflect.TypeOf(entity)
	mt := reflect.TypeOf(model)

	if et.Kind() == reflect.Ptr && mt.Kind() == reflect.Ptr {
		et = et.Elem()
		mt = mt.Elem()
		ev := reflect.ValueOf(entity).Elem()
		mv := reflect.ValueOf(model).Elem()
		if et.Kind() == reflect.Struct && mt.Kind() == reflect.Struct {
			for i := 0; i < et.NumField(); i++ {
				etField := et.Field(i)
				etFieldTagName := etField.Tag.Get("ormgen")
				if gocommon.IsEmpty(etFieldTagName) {
					etFieldTagName = etField.Name
				}
				mtField, have := mt.FieldByName(etFieldTagName)
				mfValue := mv.FieldByName(etFieldTagName)
				//判断model字段是否存在，数据类型是否一致
				if have {
					efValue := ev.Field(i)
					var mtFieldType reflect.Type
					var etFieldType reflect.Type
					if mtField.Type.Kind() == reflect.Ptr {
						mtFieldType = mtField.Type.Elem()
					} else {
						mtFieldType = mtField.Type
					}
					if etField.Type.Kind() == reflect.Ptr {
						etFieldType = etField.Type.Elem()
					} else {
						etFieldType = etField.Type
					}
					//对model字段赋值
					if mtFieldType.Name() == etFieldType.Name() {
						mfValue.Set(ev.Field(i))
					} else if etField.Type.Kind() == reflect.Bool {
						//处理bool类型
						if mtField.Type.Kind() == reflect.Uint || mtField.Type.Kind() == reflect.Uint64 {
							mfValue.SetUint(cast.InterfaceToUInt64WithDefault(efValue.Interface(), 0))
						} else if mtField.Type.Kind() == reflect.Int || mtField.Type.Kind() == reflect.Int64 {
							mfValue.SetInt(cast.InterfaceToInt64WithDefault(efValue.Interface(), 0))
						}
					} else if  mtField.Type.Kind() == reflect.Ptr &&
						mtField.Type.Elem().Name() == "JSONB" &&
						mfValue.IsValid() &&
						efValue.IsValid() &&
						!efValue.IsZero() &&
						!efValue.IsNil() {
						//jsonb 类型的处理
						jsonb := pgtype.JSONB{}
						v := efValue.Interface()
						if etField.Type.Kind() == reflect.Ptr {
							v = efValue.Elem().Interface()
						}
						err := jsonb.Set(v)
						if err != nil {
							return err
						}
						mfValue.Set(reflect.ValueOf(&jsonb))
					}
				}
			}
			return nil
		}
		return errors.New("entity and model must be struct ptr")
	}
	return errors.New("entity and model must be struct ptr")
}

//ModelToEntity 根据tag：'ormgen'，将db model struct 结构转换entity struct 结构
func ModelToEntity(model interface{}, entity interface{}) error {
	et := reflect.TypeOf(entity)
	mt := reflect.TypeOf(model)

	if et.Kind() == reflect.Ptr && mt.Kind() == reflect.Ptr {
		et = et.Elem()
		mt = mt.Elem()
		ev := reflect.ValueOf(entity).Elem()
		mv := reflect.ValueOf(model).Elem()
		if et.Kind() == reflect.Struct && mt.Kind() == reflect.Struct {
			for i := 0; i < et.NumField(); i++ {
				etField := et.Field(i)
				etFieldTagName := etField.Tag.Get("ormgen")
				if gocommon.IsEmpty(etFieldTagName) {
					etFieldTagName = etField.Name
				}
				mtField, have := mt.FieldByName(etFieldTagName)
				mfValue := mv.FieldByName(etFieldTagName)
				//判断model字段是否存在，数据类型是否一致

				if have {
					efValue := ev.Field(i)
					var mtFieldType reflect.Type
					var etFieldType reflect.Type
					if mtField.Type.Kind() == reflect.Ptr {
						mtFieldType = mtField.Type.Elem()
					} else {
						mtFieldType = mtField.Type
					}
					if etField.Type.Kind() == reflect.Ptr {
						etFieldType = etField.Type.Elem()
					} else {
						etFieldType = etField.Type
					}
					//对entity字段赋值
					if mtFieldType.Name() == etFieldType.Name() {
						efValue.Set(mfValue)
					} else if etField.Type.Kind() == reflect.Bool {
						//处理bool类型
						boolValue := cast.InterfaceToBoolWithDefault(mfValue.Interface(), false)
						efValue.SetBool(boolValue)
					} else if mtField.Type.Kind() == reflect.Ptr &&
						mtField.Type.Elem().Name() == "JSONB" &&
						mfValue.IsValid() &&
						!mfValue.IsZero() &&
						!mfValue.IsNil() &&
						efValue.IsValid() {
						//jsonb 类型的处理
						jsonb := mfValue.Elem().Interface().(pgtype.JSONB)
						t := etField.Type
						if etField.Type.Kind() == reflect.Ptr {
							t = etField.Type.Elem()
						}
						entityValue := reflect.New(t)
						err := jsonb.AssignTo(entityValue.Interface())
						if err != nil {
							return err
						}
						if etField.Type.Kind() == reflect.Ptr {
							efValue.Set(entityValue)
						}else{
							efValue.Set(entityValue.Elem())
						}
					}

				}
			}
			return nil
		}
		return errors.New("entity and model must be struct ptr")
	}
	return errors.New("entity and model must be struct ptr")
}

func ToJsonb(obj interface{}) (*pgtype.JSONB, error) {
	jsonb := &pgtype.JSONB{}
	err := jsonb.Set(obj)
	if err != nil {
		return nil, err
	}
	return jsonb, nil
}

func FromJsonb(jsonb *pgtype.JSONB, obj interface{}) error {
	return jsonb.AssignTo(obj)
}

`

const PostgresqlGormModelTpl = `
// Package dbmodel 
// Generated code. DO NOT modify by hand!
package dbmodel

import (
	"errors"
	"fmt"
	"github.com/jackc/pgtype"
	"gorm.io/gorm"
	"time"
)

type {{model_name}}DBSelector struct {
	session *gorm.DB
}

type {{model_name}}DBModel struct {
	{{model_fields}}
}

func (g *GormDB){{model_name}}() *{{model_name}}DBSelector {
	if g.db == nil {
		panic(errors.New("unresolved GormDB.db is nil"))
	}
	//新建数据库操作会话, SQL响应时间限制5s以内
	timeoutCtx, _ := context.WithTimeout(context.Background(), 5*time.Second)
	session := g.db.WithContext(timeoutCtx).Table("{{model_underscore_name}}")
	s := &{{model_name}}DBSelector{session: session}
	return s
}

func (s *{{model_name}}DBSelector) InsertOne(model *{{model_name}}DBModel) error {
	if s.session == nil {
		panic(errors.New("unresolved GormDB.db is nil"))
	}
	result := s.session.Create(model)
	if result.Error == gorm.ErrRecordNotFound {
		return nil
	}
	return result.Error
}

func (s *{{model_name}}DBSelector) InsertMany(models []*{{model_name}}DBModel) error {
	if s.session == nil {
		panic(errors.New("unresolved GormDB.db is nil"))
	}
	result := s.session.Create(models)
	if result.Error == gorm.ErrRecordNotFound {
		return nil
	}
	return result.Error
}

func (s *{{model_name}}DBSelector) Update(model *{{model_name}}DBModel) error {
	if s.session == nil {
		panic(errors.New("unresolved GormDB.db is nil"))
	}
	result := s.session.Updates(model)
	if result.Error == gorm.ErrRecordNotFound {
		return nil
	}
	return result.Error
}

func (s *{{model_name}}DBSelector) UpdateColumn(name string, value interface{}) error {
	if s.session == nil {
		panic(errors.New("unresolved GormDB.db is nil"))
	}
	result := s.session.Update(name, value)
	if result.Error == gorm.ErrRecordNotFound {
		return nil
	}
	return result.Error
}

func (s *{{model_name}}DBSelector) Delete() error {
	if s.session == nil {
		panic(errors.New("unresolved GormDB.db is nil"))
	}
	result := s.session.Delete(&{{model_name}}DBModel{})
	if result.Error == gorm.ErrRecordNotFound {
		 return nil
	}
	return result.Error
}

func (s *{{model_name}}DBSelector) DeleteOne(id int) error {
	if s.session == nil {
		panic(errors.New("unresolved GormDB.db is nil"))
	}
	result := s.session.Delete(&{{model_name}}DBModel{}, id)
	if result.Error == gorm.ErrRecordNotFound {
		return nil
	}
	return result.Error
}

func (s *{{model_name}}DBSelector) DeleteMany(ids []int) error {
	if s.session == nil {
		panic(errors.New("unresolved GormDB.db is nil"))
	}
	result := s.session.Delete(&{{model_name}}DBModel{}, ids)
	if result.Error == gorm.ErrRecordNotFound {
		return nil
	}
	return result.Error
}

func (s *{{model_name}}DBSelector) FindOne() (*{{model_name}}DBModel, error) {
	if s.session == nil {
		panic(errors.New("unresolved GormDB.db is nil"))
	}
	model := &{{model_name}}DBModel{}
	result := s.session.Limit(1).Find(model)
	if result.Error == gorm.ErrRecordNotFound {
		return nil, nil
	}
	if result.Error != nil {
		return nil, result.Error
	}
	if result.RowsAffected == 0 {
		return nil, nil
	}
	return model, nil
}

func (s *{{model_name}}DBSelector) FindMany() ([]*{{model_name}}DBModel, error) {
	if s.session == nil {
		panic(errors.New("unresolved GormDB.db is nil"))
	}
	models := make([]*{{model_name}}DBModel, 0)
	result := s.session.Find(&models)
	if result.Error == gorm.ErrRecordNotFound {
		return nil, nil
	}
	if result.Error != nil {
		return nil, result.Error
	}
	if result.RowsAffected == 0 {
		return nil, nil
	}
	return models, result.Error
}

func (s *{{model_name}}DBSelector) Select(columns ...interface{}) *{{model_name}}DBSelector {
	if s.session == nil {
		panic(errors.New("unresolved GormDB.db is nil"))
	}
	s.session = s.session.Select(columns)
	return s
}

func (s *{{model_name}}DBSelector) ByID(id int) *{{model_name}}DBSelector {
	if s.session == nil {
		panic(errors.New("unresolved GormDB.db is nil"))
	}
	s.session = s.session.Where("id = ?", id)
	return s
}

func (s *{{model_name}}DBSelector) Exist() (bool, error) {
	if s.session == nil {
		panic(errors.New("unresolved GormDB.db is nil"))
	}
	result := s.session.Limit(1).Find(&{{model_name}}DBModel{})
	if result.Error == gorm.ErrRecordNotFound {
		return false, nil
	}
	return result.RowsAffected > 0, result.Error
}

func (s *{{model_name}}DBSelector) In(column string, data interface{}) *{{model_name}}DBSelector {
	if s.session == nil {
		panic(errors.New("unresolved GormDB.db is nil"))
	}
	s.session = s.session.Where(column + " IN ?", data)
	return s
}

func (s *{{model_name}}DBSelector) Where(query interface{}, args ...interface{}) *{{model_name}}DBSelector {
	if s.session == nil {
		panic(errors.New("unresolved GormDB.db is nil"))
	}
	s.session = s.session.Where(query, args...)
	return s
}

func (s *{{model_name}}DBSelector) Search(column string, key string) *{{model_name}}DBSelector {
	if s.session == nil {
		panic(errors.New("unresolved GormDB.db is nil"))
	}
	if len(key) == 0 {
		return s
	}
	s.session = s.session.Where(column + " ~ ?", key)
	return s
}

func (s *{{model_name}}DBSelector) Limit(limit int) *{{model_name}}DBSelector {
	if s.session == nil {
		panic(errors.New("unresolved GormDB.db is nil"))
	}
	s.session = s.session.Limit(limit)
	return s
}

func (s *{{model_name}}DBSelector) Offset(offset int) *{{model_name}}DBSelector {
	if s.session == nil {
		panic(errors.New("unresolved GormDB.db is nil"))
	}
	s.session = s.session.Offset(offset)
	return s
}

func (s *{{model_name}}DBSelector) OrderBy(column string, sort string) *{{model_name}}DBSelector {
	if s.session == nil {
		panic(errors.New("unresolved GormDB.db is nil"))
	}
	s.session = s.session.Order(fmt.Sprintf("%s %s", column, sort))
	return s
}

func (s *{{model_name}}DBSelector) GroupBy(column string) *{{model_name}}DBSelector {
	if s.session == nil {
		panic(errors.New("unresolved GormDB.db is nil"))
	}
	s.session = s.session.Group(column)
	return s
}

func (s *{{model_name}}DBSelector) Having(query interface{}, args ...interface{}) *{{model_name}}DBSelector {
	if s.session == nil {
		panic(errors.New("unresolved GormDB.db is nil"))
	}
	s.session = s.session.Having(query, args...)
	return s
}

func (s *{{model_name}}DBSelector) Count() (int64, error) {
	if s.session == nil {
		panic(errors.New("unresolved GormDB.db is nil"))
	}
	var count int64
	result := s.session.Count(&count)
	if result.Error == gorm.ErrRecordNotFound {
		return 0, nil
	}
	return count, result.Error
}

`

func PostgresqlGormField(tableColumn *interviewer.DBTableColumn) string {
	name := stringutils.ToCamelCase(tableColumn.Field)

	switch tableColumn.Type {
	case "int2", "int4", "int8", "int", "bigint":
		if tableColumn.Pk {
			return fmt.Sprintf("/*%s\n*/\n %s  int `gorm:\"column:%s;primaryKey\"`", tableColumn.Comment, name, tableColumn.Field)
		} else {
			return fmt.Sprintf("/*%s\n*/\n %s  int `gorm:\"column:%s\"`", tableColumn.Comment, name, tableColumn.Field)
		}

	case "varchar":
		if tableColumn.Pk {
			return fmt.Sprintf("/*%s\n*/\n %s string `gorm:\"column:%s;primaryKey\"`", tableColumn.Comment, name, tableColumn.Field)
		} else {
			return fmt.Sprintf("/*%s\n*/\n %s string `gorm:\"column:%s\"`", tableColumn.Comment, name, tableColumn.Field)
		}

	case "timestamptz", "timestamp":
		if gocommon.IsContains(tableColumn.Field, "created") {
			return fmt.Sprintf("/*%s\n*/\n %s time.Time  `gorm:\"column:%s;autoCreateTime\"`", tableColumn.Comment, name, tableColumn.Field)
		} else if gocommon.IsContains(tableColumn.Field, "updated") {
			return fmt.Sprintf("/*%s\n*/\n %s time.Time  `gorm:\"column:%s;autoUpdateTime\"`", tableColumn.Comment, name, tableColumn.Field)
		} else if gocommon.IsContains(tableColumn.Field, "deleted") {
			return fmt.Sprintf("/*%s\n*/\n %s gorm.DeletedAt  `gorm:\"column:%s;index\"`", tableColumn.Comment, name, tableColumn.Field)
		}

	case "jsonb":
		return fmt.Sprintf("/*%s\n*/\n %s *pgtype.JSONB `gorm:\"column:%s\"`", tableColumn.Comment, name, tableColumn.Field)
	}

	return ""
}
