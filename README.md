> 通过反转数据库表信息自动生成ORM代码
# 1. 安装

` go get github.com/youngchan1988/ormgenc`

# 2. 使用

Clone : `git@github.com:youngchan1988/ormgenc.git`

执行 : 

```shell
$ cd ormgenc
$ go install
```

```shell
$ ormgenc --help
Usage of ormgenc:
  -config string
        set config file path (default "./config.yaml")
  -db string
        set target database which like 'postgresql', 'mysql', 'sqlite' and so on (default "postgresql")
  -debug
        sets log level to debug
  -help
        ormgenc usage
  -orm string
        set target orm which like 'gorm' and 'xorm' (default "gorm")
  -out string
        set output path (default ".")
  -version
        show current version
```

# 3. 配置

配置文件中包含连接数据库的信息:

```yaml
db_host: "127.0.0.1"
db_name: "apptoygodb"
db_user: "newcore"
db_password: "xinheyun2014@0711"
db_port: 5433
```

# 4. Entity <-> Model

生成的ORM代码中包含两个方法：`EntityToModel` 和 `ModelToEntity`  ，方便在业务Entity和数据库Model 相互转换。

比如：

```go
//Define User struct
type UserEntity struct {
	ID int `json:"id" ormgen:"Id"`
	Name     string `json:"name" ormgen:"Name"`
	MobilePhone string `json:"mobilePhone" ormgen:"MobilePhone"`
	Email string `json:"email" ormgen:"Email"`
	State int `json:"state" ormgen:"State"`
	Info *UserInfoEntity `json:"info" ormgen:"Info"`
	CreatedDate time.Time `json:"createdDate" ormgen:"CreatedDate"`
	UpdatedDate time.Time `json:"updatedDate" ormgen:"UpdatedDate"`
	DeletedDate time.Time `json:"deleteDate" ormgen:"DeletedDate"`
}

//UserInfoEntity 
type UserInfoEntity struct {
	Avatar string `json:"avatar"`
}

```

```go
type UserDBModel struct {

	Id int `gorm:"column:id;primaryKey"`

	Name string `gorm:"column:name"`

	Password string `gorm:"column:password"`

	MobilePhone string `gorm:"column:mobile_phone"`

	Email string `gorm:"column:email"`

	State int `gorm:"column:state"`

	Info *pgtype.JSONB `gorm:"column:info"`

	CreatedDate time.Time `gorm:"column:created_date;autoCreateTime"`

	UpdatedDate time.Time `gorm:"column:updated_date;autoUpdateTime"`

	DeletedDate gorm.DeletedAt `gorm:"column:deleted_date;index"`
}

```

在Entity中使用`ormgen` struct tag 来标识这个字段在数据库Model中对应的名称，针对`Postgresql`数据支持`jsonb`数据类型的转换



# 5.Todo

### orm

- [x] gorm

- [ ] xorm

### database

- [x] postgresql

- [ ] mysql