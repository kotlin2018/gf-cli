package initialize

import (
	"github.com/gogf/gf-cli/library/mlog"
	"os"
	"path"
	"strings"
)

func CreateApp(){
	appName := getProjectPath()
	os.MkdirAll(appName,0755) //根据项目名,创建项目
	os.Mkdir(path.Join(appName,"app"),0755)
	os.Mkdir(path.Join(appName,"boot"),0755)
	os.Mkdir(path.Join(appName,"config"),0755)
	os.Mkdir(path.Join(appName,"docker"),0755)
	os.Mkdir(path.Join(appName,"document"),0755)
	os.Mkdir(path.Join(appName,"i18n"),0755)
	os.Mkdir(path.Join(appName,"packed"),0755)
	os.Mkdir(path.Join(appName,"public"),0755)
	os.Mkdir(path.Join(appName,"router"),0755)
	os.Mkdir(path.Join(appName,"template"),0755)
	os.Mkdir(path.Join(appName,"library"),0755)

	//递归创建目录 app应用程序层
	os.MkdirAll(path.Join(appName, "/app/api/v1"), 0755)
	os.MkdirAll(path.Join(appName, "/app/api/request"), 0755)
	os.MkdirAll(path.Join(appName, "/app/dao"), 0755)
	os.MkdirAll(path.Join(appName, "/app/model"), 0755)
	os.MkdirAll(path.Join(appName, "/app/model/internal"), 0755)
	os.MkdirAll(path.Join(appName, "/app/dao/internal"), 0755)
	os.MkdirAll(path.Join(appName, "/app/service"), 0755)
	os.MkdirAll(path.Join(appName, "/public/html"), 0755)
	os.MkdirAll(path.Join(appName, "/public/plugin"), 0755)
	os.MkdirAll(path.Join(appName, "/public/resource/css"), 0755)
	os.MkdirAll(path.Join(appName, "/public/resource/img"), 0755)
	os.MkdirAll(path.Join(appName, "/public/resource/js"), 0755)

	// 生成.gitkeep文件
	createFile(path.Join(appName, "/app/dao", ".gitkeep"))
	createFile(path.Join(appName, "/app/model", ".gitkeep"))
	createFile(path.Join(appName, "/app/service", ".gitkeep"))
	createFile(path.Join(appName, "boot", ".gitkeep"))
	createFile(path.Join(appName, "config", ".gitkeep"))
	createFile(path.Join(appName, "docker", ".gitkeep"))
	createFile(path.Join(appName, "document", ".gitkeep"))
	createFile(path.Join(appName, "i18n", ".gitkeep"))
	createFile(path.Join(appName, "/public/html", ".gitkeep"))
	createFile(path.Join(appName, "/public/plugin", ".gitkeep"))
	createFile(path.Join(appName, "/public/resource/css", ".gitkeep"))
	createFile(path.Join(appName, "/public/resource/img", ".gitkeep"))
	createFile(path.Join(appName, "/public/resource/js", ".gitkeep"))
	createFile(path.Join(appName, "template", ".gitkeep"))

	// 生成模版代码
	writeToFile(path.Join(appName, "packed", "packed.go"),packed)
	writeToFile(path.Join(appName, ".gitattributes"), gitattributes)
	writeToFile(path.Join(appName, ".gitignore"), gitignore)
	writeToFile(path.Join(appName, "Dockerfile"), dockerfile)
	writeToFile(path.Join(appName, "README.MD"), readme)
	writeToFile(path.Join(appName, "config", "config.toml"), toml) // toml配置文件
	writeToFile(path.Join(appName, "/app/api/request", "params.go"), params)
	writeToFile(path.Join(appName, "/app/model/internal", "user.go"), user)

	// 使用插值表达式
	writeToFile(path.Join(appName, "/app/dao/", "user.go"), strings.Replace(daoUser, "{{.appName}}", appName, -1))
	writeToFile(path.Join(appName, "/app/model", "user.go"), strings.Replace(user2, "{{.appName}}", appName, -1))
	writeToFile(path.Join(appName, "/app/dao/internal", "user.go"), strings.Replace(dao, "{{.appName}}", appName, -1))
	writeToFile(path.Join(appName, "boot", "boot.go"), strings.Replace(boot, "{{.appName}}", appName, -1))
	writeToFile(path.Join(appName, "router", "router.go"),strings.Replace(router, "{{.appName}}", appName, -1))
	writeToFile(path.Join(appName, "/app/api/v1", "base.go"), strings.Replace(base, "{{.appName}}", appName, -1))
	writeToFile(path.Join(appName, "/app/service", "base.go"), strings.Replace(base2, "{{.appName}}", appName, -1))
	writeToFile(path.Join(appName, "main.go"), strings.Replace(main, "{{.appName}}", appName, -1))
	writeToFile(path.Join(appName, "go.mod"), strings.Replace(mod, "{{.appName}}", appName, -1))
	mlog.Print("初始化完成! ")
	mlog.Print("你现在可以在当前目录下运行 'gf run main.go' 来启动项目!")
}

// 将content写入fileName 文件中
// 将文本内容写入文件中
func writeToFile(fileName,content string){
	f, err := os.Create(fileName)
	if err !=nil {
		panic(err)
	}
	defer f.Close()
	_, err = f.WriteString(content)
	if err !=nil {
		panic(err)
	}
}

// 一个文件或者目录是否存在
func IsExist(path string)bool{
	_, err := os.Stat(path)
	return err == nil || os.IsExist(err)
}

func createFile(fileName string){
	f, err := os.Create(fileName)
	if err !=nil {
		panic(err)
	}
	defer f.Close()
}

