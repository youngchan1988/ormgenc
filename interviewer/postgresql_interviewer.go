package interviewer

import (
	"context"
	"errors"
	"fmt"
	"github.com/jackc/pgx/v4"
	"github.com/youngchan1988/gocommon"
	"github.com/youngchan1988/gocommon/cast"
	"github.com/youngchan1988/gocommon/log"
)

const logTag = "PostgresqlInterviewer"

type PostgresqlInterviewer struct {
	conn *pgx.Conn
	//DB连接配置信息
	ConnConfig *DBConnConfig
}

func (pgc *PostgresqlInterviewer) Open() error {
	connConfig, _ := pgx.ParseConfig(fmt.Sprintf("postgresql://%s:%s@%s:%d/%s", pgc.ConnConfig.User, pgc.ConnConfig.Password, pgc.ConnConfig.Host, pgc.ConnConfig.Port, pgc.ConnConfig.Database))
	var err error
	pgc.conn, err = pgx.ConnectConfig(context.Background(), connConfig)
	if err != nil {
		return err
	}
	return nil
}

func (pgc *PostgresqlInterviewer) Close() error {
	if pgc.conn != nil && !pgc.conn.IsClosed() {
		return pgc.conn.Close(context.Background())
	}
	return errors.New("postgresql connection is unavailable or have already closed")
}

func (pgc *PostgresqlInterviewer) GetDBTables() []*DBTable {
	//获取每个表字段信息
	dbTables := make([]*DBTable, 0)

	tableNames := pgc.getDbTableNames()
	if !gocommon.IsEmpty(tableNames) {
		for _, name := range tableNames {
			dbColumns := pgc.getDbTableColumns(name)
			if !gocommon.IsEmpty(dbColumns) {
				dbTable := &DBTable{
					Name:    name,
					Columns: dbColumns,
				}
				dbTables = append(dbTables, dbTable)
			}
		}
	}

	return dbTables
}

func (pgc *PostgresqlInterviewer) dbEngine() (*pgx.Conn, error) {
	if pgc.conn == nil {
		err := pgc.Open()
		if err != nil {
			log.Error(logTag, err, 1, "Open db error")
			return nil, err
		}
	}
	return pgc.conn, nil
}

//获取数据库所有表名称
func (pgc *PostgresqlInterviewer) getDbTableNames() []string {
	db, err := pgc.dbEngine()
	if err != nil {
		return nil
	}
	var tableNames []string

	rows, err1 := db.Query(context.Background(), "SELECT tablename FROM pg_tables WHERE schemaname='public'")
	if err1 != nil {
		log.Error(logTag, err1, 1, "query tablenames error")
		return nil
	}
	tableNames = make([]string, 0)
	for rows.Next() {
		var name string
		err1 = rows.Scan(&name)
		if err1 != nil {
			log.Error(logTag, err1, 1, "scan row to field error")
			return nil
		}
		tableNames = append(tableNames, name)
	}

	return tableNames
}

//获取表字段信息
func (pgc *PostgresqlInterviewer) getDbTableColumns(tableName string) []*DBTableColumn {
	db, err := pgc.dbEngine()
	if err != nil {
		return nil
	}

	pk, err2 := db.Query(context.Background(), fmt.Sprintf("SELECT pg_attribute.attname AS colname FROM\n    pg_constraint  INNER JOIN pg_class\n                              ON pg_constraint.conrelid = pg_class.oid\n                   INNER JOIN pg_attribute ON pg_attribute.attrelid = pg_class.oid\n        AND  pg_attribute.attnum = pg_constraint.conkey[1]\n                   INNER JOIN pg_type ON pg_type.oid = pg_attribute.atttypid\nWHERE pg_class.relname = '%s'\n  AND pg_constraint.contype='p'", tableName))
	if err2 != nil {
		log.Error(logTag, err2, 1, "query table pk error")
		return nil
	}
	pk.Next()
	var pkFieldName string
	err2 = pk.Scan(&pkFieldName)
	if err2 != nil {
		log.Error(logTag, err2, 1, "scan row to pk field error")
		return nil
	}
	pk.Close()
	rows, err1 := db.Query(context.Background(), fmt.Sprintf("SELECT a.attname AS name,\n       pg_type.typname AS typename,\n       a.attnotnull AS notnull,\n       col_description(a.attrelid,a.attnum) AS comment\nFROM pg_class AS c,pg_attribute AS a INNER JOIN pg_type ON pg_type.oid = a.atttypid\nWHERE c.relname = '%s' AND a.attrelid = c.oid AND a.attnum>0", tableName))
	if err1 != nil {
		log.Error(logTag, err1, 1, "query table fields error")
		return nil
	}

	dbColumns := make([]*DBTableColumn, 0)
	for rows.Next() {
		var fieldName string
		var fieldType string
		var notNull bool
		var comment interface{}
		err1 = rows.Scan(&fieldName, &fieldType, &notNull, &comment)
		if err1 != nil {
			log.Error(logTag, err1, 1, "scan row to field error")
			return nil
		}
		column := &DBTableColumn{
			Field:   fieldName,
			Type:    fieldType,
			NotNull: notNull,
			Comment: cast.InterfaceToStringWithDefault(comment),
			Pk:      fieldName == pkFieldName,
		}
		dbColumns = append(dbColumns, column)
	}
	return dbColumns
}
