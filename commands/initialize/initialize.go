package initialize

import (
	"github.com/gogf/gf-cli/library/allyes"
	"github.com/gogf/gf-cli/library/mlog"
	"github.com/gogf/gf/encoding/gcompress"
	"github.com/gogf/gf/frame/g"
	"github.com/gogf/gf/net/ghttp"
	"github.com/gogf/gf/os/gcmd"
	"github.com/gogf/gf/os/gfile"
	"github.com/gogf/gf/text/gstr"
	"strings"
)

const (
	emptyProject     = "github.com/gogf/gf-empty"
	emptyProjectName = "gf-empty"
)

var (
	cdnUrl  = g.Config("url").GetString("cdn.url")
	homeUrl = g.Config("url").GetString("home.url")
)

func init() {
	if cdnUrl == "" {
		mlog.Fatal("cdn配置不能为空")
	}
	if homeUrl == "" {
		mlog.Fatal("home配置不能为空")
	}
}

func Help() {
	mlog.Print(gstr.TrimLeft(`
usage    
    gf new projectName

argument 
    projectName  项目的名称。它将在当前目录中创建一个 name 的文件夹.
				 projectName也将是项目的模块名称.

examples 
    gf new my-app
    gf new my-project-name
`))
}

func Run() {

	projectName := getProjectPath()
	dirPath := projectName
	mlog.Print("初始化...")
	// --------以下这段逻辑就是从远程服务器地址 "homeUrl + "/cli/project" 获取go-empty项目，并用projectName替换go-empty----------
	// MD5 检索.
	respMd5, err := ghttp.Get(homeUrl + "/cli/project/md5")
	if err != nil {
		mlog.Fatalf("获取项目zip md5失败: %s", err.Error())
	}
	defer respMd5.Close()
	md5DataStr := respMd5.ReadAllString()
	if md5DataStr == "" {
		mlog.Fatal("获取项目zip md5失败：md5值为空。可能是网络问题，请重试?")
	}

	// 压缩数据检索.
	respData, err := ghttp.Get(cdnUrl + "/cli/project/zip?" + md5DataStr)
	if err != nil {
		mlog.Fatal("获取项目zip数据失败: %s", err.Error())
	}
	defer respData.Close()
	zipData := respData.ReadAll()
	if len(zipData) == 0 {
		mlog.Fatal("获取项目数据失败：数据值为空。可能是网络问题，请重试?")
	}

	// 解压压缩数据.
	if err = gcompress.UnZipContent(zipData, dirPath, emptyProjectName+"-master"); err != nil {
		mlog.Fatal("解压缩项目数据失败,", err.Error())
	}
	// 替换项目名称.
	if err = gfile.ReplaceDir(emptyProject, projectName, dirPath, "Dockerfile,*.go,*.MD,*.mod", true); err != nil {
		mlog.Fatal("内容替换失败,", err.Error())
	}
	if err = gfile.ReplaceDir(emptyProjectName, projectName, dirPath, "Dockerfile,*.go,*.MD,*.mod", true); err != nil {
		mlog.Fatal("内容替换失败,", err.Error())
	}
	mlog.Print("初始化完成! ")
	mlog.Print("你现在可以在当前目录下运行 'gf run main.go' 来启动项目!")
}

func getProjectPath()string {
	parser, err := gcmd.Parse(nil) // 创建并返回一个新的解析器
	if err != nil {
		mlog.Fatal(err)
	}
	projectName := parser.GetArg(2) //解析器获取第二个参数
	if projectName == "" {
		mlog.Fatal("项目名称不应为空")
	}
	dirPath := projectName
	if !gfile.IsEmpty(dirPath) && !allyes.Check() {
		s := gcmd.Scanf(`文件夹“%s”不为空，文件可能已重写，是否继续? [y/n]: `, projectName)
		if strings.EqualFold(s, "n") {
			return ""
		}
	}
	return dirPath
}
