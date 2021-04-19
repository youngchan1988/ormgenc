package interviewer

type DBInterviewer interface {
	Open() error
	Close() error
	GetDBTables() []*DBTable
}

type DBConnConfig struct {
	Host     string
	Database string
	User     string
	Password string
	Port     uint16
}

type DBTable struct {
	//表名称
	Name string
	//表字段
	Columns []*DBTableColumn
}

type DBTableColumn struct {
	//表字段名称
	Field string
	//表字段类型
	Type string
	//是否可空
	NotNull bool
	//表字段注释
	Comment string
	//是否主键
	Pk bool
}
