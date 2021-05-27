package gen

import (
	"fmt"
	"github.com/valyala/fasttemplate"
	"github.com/youngchan1988/gocommon"
	"github.com/youngchan1988/gocommon/fileutils"
	"github.com/youngchan1988/gocommon/log"
	"github.com/youngchan1988/gocommon/stringutils"
	"ormgenc/interviewer"
	"ormgenc/template"
	"strings"
)

const logTag = "GormGen"

type GormGen struct {
	DBInterviewer interviewer.DBInterviewer
}

func (g *GormGen) GenerateModel(outPath string) {
	dbTables := g.DBInterviewer.GetDBTables()
	err := fileutils.WriteFile(outPath+"/dbmodel_gen.go", template.PostgresqlGormTpl, true)
	if err != nil {
		log.Error(logTag, err, 1, "write [dbmodel_gen.go] error")
		return
	}
	err = fileutils.WriteFile(outPath+"/dbutils_gen.go", template.PostgresqlGormUtilsTpl, true)
	if err != nil {
		log.Error(logTag, err, 1, "write [dbutils_gen.go] error")
		return
	}

	if !gocommon.IsEmpty(dbTables) {
		for _, t := range dbTables {
			tplString := g.generateModelForTable(t)
			if !gocommon.IsEmpty(tplString) {
				err = fileutils.WriteFile(fmt.Sprintf("%s/%smodel_gen.go", outPath, strings.ToLower(stringutils.TocamelCase(t.Name))), tplString, true)
			}
		}
	}
}

//生成数据model
func (g *GormGen) generateModelForTable(table *interviewer.DBTable) string {
	tpl := fasttemplate.New(template.PostgresqlGormModelTpl, "{{", "}}")
	tplParams := make(map[string]interface{})
	tplParams["model_name"] = stringutils.ToCamelCase(table.Name)
	tplParams["model_underscore_name"] = table.Name
	dbFields := strings.Builder{}
	for _, c := range table.Columns {
		dbFields.WriteString(template.PostgresqlGormField(c) + "\n")
	}
	tplParams["model_fields"] = dbFields.String()
	return tpl.ExecuteString(tplParams)
}
