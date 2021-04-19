package main

import (
	"flag"
	"fmt"
	"github.com/spf13/viper"
	"github.com/youngchan1988/gocommon/log"
	"ormgenc/gen"
	"ormgenc/interviewer"
	"os/exec"
)

const version = "v0.0.1"

func main() {
	debug := flag.Bool("debug", false, "sets log level to debug")
	targetOrmArg := flag.String("orm", "gorm", "set target orm which like gorm and xorm")
	targetDbArg := flag.String("db", "postgresql", "set target database which like postgresql and mysql ans sqlite and so on")
	configFile := flag.String("config", "./config.yaml", "set config file path")
	outPath := flag.String("out", ".", "set output path")
	versionArg := flag.Bool("version", false, "get current version")
	help := flag.Bool("help", false, "ormgenc usage")
	flag.Parse()
	if *help {
		flag.Usage()
		return
	} else if *versionArg {
		fmt.Println(version)
		return
	}

	log.Init(*debug, "", "")
	viper.SetConfigFile(*configFile)
	//viper.SetConfigName("config")
	viper.SetConfigType("yaml")
	//viper.AddConfigPath(".")
	err := viper.ReadInConfig()
	if err != nil {
		panic(fmt.Errorf("Fatal error config file: %s \n", err))
	}
	targetOrm := *targetOrmArg
	targetDb := *targetDbArg
	if targetOrm != "gorm" && targetOrm != "xorm" {
		panic(fmt.Errorf("unresolved orm %s", targetOrm))
	}
	if targetDb != "postgresql" && targetDb != "mysql" && targetDb != "sqlite" {
		panic(fmt.Errorf("unresolved database %s", targetDb))
	}
	var ormGen gen.Gen
	var dbInterviewer interviewer.DBInterviewer
	if targetDb == "postgresql" {
		dbInterviewer = &interviewer.PostgresqlInterviewer{
			ConnConfig: &interviewer.DBConnConfig{
				Host:     viper.GetString("db_host"),
				Database: viper.GetString("db_name"),
				User:     viper.GetString("db_user"),
				Password: viper.GetString("db_password"),
				Port:     uint16(viper.GetUint("db_port")),
			},
		}

	}
	genDir := *outPath + "/dbmodel"
	if targetOrm == "gorm" {
		ormGen = &gen.GormGen{
			DBInterviewer: dbInterviewer,
		}
		ormGen.GenerateModel(genDir)
	}
	execcmd := exec.Command("gofmt", "-w", genDir)
	fmt.Println(execcmd.String())
	out, err1 := execcmd.CombinedOutput()
	fmt.Println(string(out))
	if err1 != nil {
		panic(fmt.Sprintf("gofmt failed: %v", err1))
	}
}
