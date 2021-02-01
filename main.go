package main

import (
	"fmt"
	_ "github.com/gogf/gf-cli/boot"
	"github.com/gogf/gf-cli/commands/build"
	"github.com/gogf/gf-cli/commands/docker"
	"github.com/gogf/gf-cli/commands/env"
	"github.com/gogf/gf-cli/commands/fix"
	"github.com/gogf/gf-cli/commands/gen"
	"github.com/gogf/gf-cli/commands/get"
	"github.com/gogf/gf-cli/commands/initialize"
	"github.com/gogf/gf-cli/commands/install"
	"github.com/gogf/gf-cli/commands/mod"
	"github.com/gogf/gf-cli/commands/pack"
	"github.com/gogf/gf-cli/commands/run"
	"github.com/gogf/gf-cli/commands/swagger"
	"github.com/gogf/gf-cli/commands/update"
	"github.com/gogf/gf-cli/library/allyes"
	"github.com/gogf/gf-cli/library/mlog"
	"github.com/gogf/gf-cli/library/proxy"
	"github.com/gogf/gf/errors/gerror"
	"github.com/gogf/gf/os/gbuild"
	"github.com/gogf/gf/os/gcmd"
	"github.com/gogf/gf/os/genv"
	"github.com/gogf/gf/os/gfile"
	"github.com/gogf/gf/text/gstr"
	"strings"
)

const (
	VERSION = "v1.15.0"
)

func init() {
	// Automatically sets the golang proxy for all commands.
	// 自动为所有命令设置golang代理
	proxy.AutoSet()
}

var (
	// 从字符串的开头去除空白（或其他字符）
	helpContent = gstr.TrimLeft(`
usage
    gf command [argument] [option] 

command
    env        	显示当前Golang环境变量
    get        	默认情况下安装或更新GF到系统
    gen        	为ORM模型自动生成go文件
    mod        	go模块的额外功能
    run        	运行具有热编译特性的go代码
    init  	创建并初始化一个空的GF项目
    new		创建一个具备基础功能的GF项目
    help       	显示有关指定命令的详细信息
    pack       	将任何文件/目录打包到资源文件或go文件...
    build      	跨平台交叉编译...
    docker     	为当前GF项目创建docker映像
    swagger    	当前项目的swagger功能...
    update     	更新当前gf二进制文件（可能需要root/admin权限）
    install    	将gf二进制文件安装到系统（可能需要root/admin权限）
    version    	显示当前二进制版本信息

option
    -y         all yes for all命令，无需提示询问 
    -?,-h      显示指定命令的帮助或详细信息 
    -v,-i      显示版本信息

additional(额外的)
    使用“gf help command”或“gf command-h”了解有关命令的详细信息，这些命令的注释末尾有“…”
`)
)

func main() {
	// Force using configuration file in current working directory.
	// 强制使用当前工作目录中的配置文件
	genv.Set("GF_GCFG_PATH", gfile.Pwd()) //当前工作目录的绝对路径

	defer func() {
		if exception := recover(); exception != nil {
			if err, ok := exception.(error); ok {
				mlog.Print(gerror.Current(err).Error())
			} else {
				panic(exception)
			}
		}
	}()

	allyes.Init()

	command := gcmd.GetArg(1)
	// Help information
	if gcmd.ContainsOpt("h") && command != "" {
		help(command)
		return
	}
	switch command {
	case "help":
		help(gcmd.GetArg(2))
	case "version":
		version()
	case "env":
		env.Run() // 执行env命令的逻辑
	case "get":
		get.Run() // 执行get命令的逻辑
	case "gen":
		gen.Run()
	case "fix":
		fix.Run()
	case "mod":
		mod.Run()
	case "init":
		initialize.Run()
	case "new":
		initialize.CreateApp()
	case "pack":
		pack.Run()
	case "docker":
		docker.Run()
	case "swagger":
		swagger.Run()
	case "update":
		update.Run()
	case "install":
		install.Run()
	case "build":
		build.Run()
	case "run":
		run.Run()
	default:
		for k := range gcmd.GetOptAll() {
			switch k {
			case "?", "h":
				mlog.Print(helpContent)
				return
			case "i", "v":
				version()
				return
			}
		}
		// No argument or option, do installation checks.
		// 没有参数或选项，执行安装检查
		if !install.IsInstalled() {
			mlog.Print("你好，这是你第一次安装gf cli")
			s := gcmd.Scanf("是否要在系统中安装gf二进制文件? [y/n]: ")
			if strings.EqualFold(s, "y") {
				install.Run()
				gcmd.Scan("按<Enter>退出...")
				return
			}
		}
		mlog.Print(helpContent)
	}
}

// help shows more information for specified command.
// 显示指定命令的详细信息(这个函数封装得很奈斯)
func help(command string) {
	switch command {
	case "get":
		get.Help()
	case "gen":
		gen.Help()
	case "new","init":
		initialize.Help()
	case "docker":
		docker.Help()
	case "swagger":
		swagger.Help()
	case "build":
		build.Help()
	case "pack":
		pack.Help()
	case "run":
		run.Help()
	case "mod":
		mod.Help()
	default:
		mlog.Print(helpContent)
	}
}

// version prints the version information of the cli tool.
func version() {
	info := gbuild.Info()
	if info["git"] == "" {
		info["git"] = "none"
	}
	mlog.Printf(`GoFrame CLI Tool %s, https://goframe.org`, VERSION)
	mlog.Printf(`Install Path: %s`, gfile.SelfPath())
	if info["gf"] == "" {
		mlog.Print(`当前版本为自定义安装版本，没有安装信息.`)
		return
	}

	mlog.Print(gstr.Trim(fmt.Sprintf(`
Build Detail:
  Go Version:  %s
  GF Version:  %s
  Git Commit:  %s
  Build Time:  %s
`, info["go"], info["gf"], info["git"], info["time"])))
}
