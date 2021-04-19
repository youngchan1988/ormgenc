package gen

type Gen interface {
	//GenerateModel 生成表结构Model代码
	GenerateModel(outPath string)
}
