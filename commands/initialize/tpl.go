package initialize

var gitattributes = `* linguist-language=GO`

var gitignore = `.buildpath
.hgignore.swp
.project
.orig
.swp
.idea/
.settings/
.vscode/
vender/
log/
composer.lock
gitpush.sh
pkg/
bin/
cbuild
*/.DS_Store
main
.vscode
go.sum`

var dockerfile = `FROM loads/alpine:3.8

LABEL maintainer="john@goframe.org"

###############################################################################
#                                INSTALLATION
###############################################################################

# 设置固定的项目路径
ENV WORKDIR /var/www/gf-empty

# 添加应用可执行文件，并设置执行权限
ADD ./bin/linux_amd64/main   $WORKDIR/main
RUN chmod +x $WORKDIR/main

# 添加I18N多语言文件、静态文件、配置文件、模板文件
ADD i18n     $WORKDIR/i18n
ADD public   $WORKDIR/public
ADD config   $WORKDIR/config
ADD template $WORKDIR/template

###############################################################################
#                                   START
###############################################################################
WORKDIR $WORKDIR
CMD ./main`

var toml = `# web服务器配置
[server]
    # 基础配置(常规配置)
    Address         = ":8000"
    ServerRoot      = "public"  #静态文件地址
    DumpRouterMap   = true
    # 系统访问日志
    AccessLogEnable = true
    # 系统异常日志panic
    ErrorLogEnable  = true
    PProfEnable     = true
    # 系统日志目录，启动，访问，异常
    LogPath         = "./log/server_log"
    # 其他可选配置
    NameToUriType   = 3

# Logger.
[logger]
    Path        = "./log/run_log"
    Level       = "all"
    Stdout      = true

# Database.
[database]
    link  = "mysql:root:root@tcp(127.0.0.1:3306)/copyele"
    debug = true
    # Database logger.
    [database.logger]
        Path   = "./log/sql_log"
        Level  = "all"
        Stdout = true

# Template.
[viewer]
    Path        = "template"
    DefaultFile = "index.html"
    Delimiters  =  ["${", "}"]

[gfcli.gen.dao]
    link   = "mysql:root:root@tcp(127.0.0.1:3306)/copyele"
    group  = ""      # 数据库分组名
    prefix = ""      # 数据库对象及文件的前缀
    tables = ""		  # 当前数据库中需要执行代码生成的数据表。(如果为空，表示数据库的所有表都会生成)

[jwt]
    ExpiresAt   = 1             # 过期时间1天
    RefreshAt   = 7             # 刷新时间7天 (24 * 7 = 168)
    SigningKey  = "SliverHorn"  # 签名(密钥)

[captcha]
    KeyLong 	= 6
    ImgWidth 	= 240
    ImgHeight 	= 80
`

var boot = `package boot

import (
	_ "{{.appName}}/packed"
)

func init() {
	
}
func init(){

}
`

var router = `package router

import (
	"github.com/gogf/gf/frame/g"
	"github.com/gogf/gf/net/ghttp"
	v1 "{{.appName}}/app/api/v1"
	"time"
)

func Error(r *ghttp.Request) {
	r.Middleware.Next()
	if err := r.GetError(); err != nil {
		// 记录到自定义错误日志文件
		g.Log("exception").Error(err)
		//返回固定的友好信息
		r.Response.ClearBuffer()
		r.Response.WriteJson(g.Map{"msg":"服务器居然开小差了，请稍后再试吧！"})
	}
}

func CORS(r *ghttp.Request) {
	r.Response.CORSDefault()
	r.Middleware.Next()
}

func JwtAuth(r *ghttp.Request) {
	v1.GfJWTMiddleware.MiddlewareFunc()(r)
	r.Middleware.Next()
}

func init() {
	g.I18n("zh-CN")
	s := g.Server()
	s.SetReadTimeout(10 * time.Second)
	s.SetWriteTimeout(10 * time.Second)
	s.SetMaxHeaderBytes(1 << 20)
	s.SetIndexFolder(true)
	s.AddStaticPath("", "public")

	s.Use(Error,CORS)
	s.BindHandler("POST:/login",v1.GfJWTMiddleware.LoginHandler)
	// base路由组下的不加JwtAuth权限认证
	s.Group("/base", func(g *ghttp.RouterGroup) {
		g.POST("/register",v1.Register)
		g.ALL("/captcha",v1.Captcha)
		g.ALL("/refresh_token",v1.GfJWTMiddleware.RefreshHandler)
		//g.Middleware(middleware.JwtAuth)
	})
}`

var main = `package main

import (
	"github.com/gogf/gf/frame/g"
	_ "{{.appName}}/boot"
	_ "{{.appName}}/router"
)

func main() {
	g.Server().Run()
}`

var base = `package v1

import (
	"errors"
	jwt "github.com/gogf/gf-jwt"
	"github.com/gogf/gf-jwt/example/api"
	"github.com/gogf/gf/frame/g"
	"github.com/gogf/gf/net/ghttp"
	"github.com/gogf/gf/util/gconv"
	"{{.appName}}/app/api/request"
	"{{.appName}}/app/model"
	"{{.appName}}/app/service"
	"time"
)

var GfJWTMiddleware *jwt.GfJWTMiddleware // 底层的JWT中间件

// 读取配置文件中jwt配置
var (
	signingKey 	= g.Cfg().GetString("jwt.SigningKey")
	timeout  	= gconv.Duration(g.Cfg().GetInt("jwt.ExpiresAt")) * time.Hour *24
	maxRefresh  = gconv.Duration(g.Cfg().GetInt("jwt.RefreshAt")) * time.Hour *24
)

// 重写此函数以自定义您自己的JWT设置。
func init(){
	authMiddleWare, err := jwt.New(&jwt.GfJWTMiddleware{
		Realm:         signingKey,                                          // 用于展示中间件的名称
		Key:           []byte(signingKey),                                 	// 密钥
		Timeout:       timeout,                              				// token过期时间
		MaxRefresh:    maxRefresh,                              			// token刷新时间
		IdentityKey:   "id",                                               	// 身份验证的key值
		TokenHeadName: "Bearer",                                           	// token在请求头时的名称，默认值为Bearer
		TokenLookup:   "header: Authorization, query: token, cookie: jwt", 	// token检索模式，用于提取token-> Authorization
		TimeFunc:       time.Now,                                          	// 测试或服务器在其他时区可设置该属性

		Authenticator:   Authenticator,                                    		// 根据登录信息对用户进行身份验证的回调函数
		LoginResponse:   LoginResponse,                                    		// 完成登录后返回的信息，用户可自定义返回数据，默认返回

		RefreshResponse: api.RefreshResponse,                                  	// 刷新token后返回的信息，用户可自定义返回数据，默认返回
		Unauthorized:    api.Unauthorized,                                     	// 处理不进行授权的逻辑
		IdentityHandler: api.IdentityHandler,                                  	// 解析并设置用户身份信息
		PayloadFunc:     api.PayloadFunc,                                      	// 登录期间的回调的函数
	})

	if err != nil {
		g.Log().Fatal("JWT Error:" + err.Error())
	}
	GfJWTMiddleware = authMiddleWare
}

// 检测身份信息是否正常
func Authenticator(r *ghttp.Request)(interface{},error) {
	// 解析请求参数到结构体
	l := (*request.Login)(nil)
	if err := r.Parse(&l); err !=nil {
		r.Response.WriteJson(g.Map{"msg":err.Error()})
		r.Exit()
	}
	// 验证码校对
	if !service.Store.Verify(l.CaptchaId,l.Captcha,true) {
		return nil,errors.New("验证码错误")
	}

	user,err := service.Login(l)
	if err != nil {
		r.Response.WriteJson(g.Map{"msg":err.Error()})
		r.ExitAll()
	}

	// 设置参数保存到请求中
	r.SetParam("user",user)
	return g.Map{"data":user},nil
}

// LoginResponse 自定义的登录成功回调函数
func LoginResponse(r *ghttp.Request, code int, token string, expire time.Time) {
	u := (*model.User)(nil)
	if err := gconv.Struct(r.GetParam("user"), &u); err != nil {
		r.Response.WriteJson(g.Map{"msg":"登录失败"})
		r.Exit()
	}
	//r.Response.WriteJson(g.Map{"msg":"登录成功!","data":g.Map{"user":u,"token":token,"expiresAt":expire.Unix() * 1000}})
	r.Response.WriteJson(g.Map{"msg":"登录成功!","data":g.Map{"user":u,"token":token,"expiresAt":expire.Format(time.RFC3339)}})
}

// 生成验证码
func Captcha(r *ghttp.Request) {
	id,picPath,err := service.Captcha()
	if err !=nil {
		r.Response.WriteJson(g.Map{"msg":"获取验证码失败","data":err.Error()})
		r.Exit()
	}
	r.Response.WriteJson(g.Map{"msg":"验证码获取成功","captchaId":id,"picPath":picPath})
}

func Login(r *ghttp.Request){

}

func Register(r *ghttp.Request) {
	p := (*request.Register)(nil)
	if err := r.Parse(&p); err != nil {
		r.Response.WriteJson(g.Map{"msg":err.Error()})
		r.Exit()
	}
	if err := service.Register(p); err != nil {
		r.Response.WriteJson(g.Map{"msg":err.Error()})
		r.ExitAll()
	}
	r.Response.WriteJson(g.Map{"msg":"注册成功!"})
}
`

var params = `package request

type Login struct {
	Username  string  //p:"username" v:"required|passport#请输入用户名称|您输入用户名称非法"
	Password  string  //p:"password" v:"required|password#请输入密码|您输入用户密码非法"
	Email	  string  //p:"email" v:"email#邮箱格式错误"
	Phone	  string  //p:"phone" v:"phone#手机号格式错误"
	Captcha   string  //json:"captcha" valid:"required#请输入正确的验证码"
	CaptchaId string  //json:"captchaId" valid:"required|length:20,20#请输入captchaId|您输入captchaId长度非法"
}

type Register struct {
	Username  		string //p:"username" v:"required|passport#请输入用户名称|用户名格式错误"
	Password  		string //p:"password" v:"required|password#请输入密码|密码格式错误"
	ConfirmPwd   	string //p:"confirmPwd" v:"required|same:password#请输入密码|两次输入密码不一致"
	Email	  		string //p:"email" v:"email#邮箱格式错误"
	Phone	  		string //p:"phone" v:"phone#手机号格式错误"
}`

var user = `package internal

import (
    "github.com/gogf/gf/os/gtime"
)

// User is the golang structure for table user.
type User struct {
    Id            string                                                           
    UserName      string                                                    
    NickName      string                                           
    Password      string                                                  
    Email         string                                         
    Mobile        string      
    Avatar        string                                                 
    Birthday      string                                              
    Sex           int8                                        
    Status        int8                                  
    CreateAt      *gtime.Time                                                 
    LastLoginTime *gtime.Time                                           
    LastLoginIp   string                                                   
    DeptId        string                                                   
    Remark        string                                               
    IsAdmin       bool                                     
}`

var user2 = `package model

import (
	"{{.appName}}/app/model/internal"
)

// User is the golang structure for table user.
type User internal.User

// Fill with you ideas below.`

var base2 = `package service

import (
	"errors"
	"github.com/gogf/gf/crypto/gmd5"
	"github.com/gogf/gf/frame/g"
	"github.com/mojocn/base64Captcha"
	"{{.appName}}/app/api/request"
	"{{.appName}}/app/dao"
	"{{.appName}}/app/model"
)

var Store = base64Captcha.DefaultMemStore

// 生成验证码
func Captcha()(id,picPath string,err error) {
	imgHeight := g.Cfg("captcha").GetInt("captcha.ImgHeight")
	imgWidth := g.Cfg("captcha").GetInt("captcha.ImgWidth")
	keyLong := g.Cfg("captcha").GetInt("captcha.KeyLong")
	driver := base64Captcha.DriverDigit{Height: imgHeight, Width: imgWidth, Length: keyLong, MaxSkew: 0.7, DotCount: 80}
	captcha := base64Captcha.NewCaptcha(&driver, Store)
	id, picPath, err = captcha.Generate()
	return
}

//
func Login(l *request.Login) (user *model.User,err error) {
	u := (*model.User)(nil)
	err = dao.User.M.Where("username=? OR phone=? OR email=?", l.Username, l.Phone, l.Email).Scan(&u)
	if err != nil {
		return nil,errors.New("用户不存在")
	}
	pass,_:= gmd5.Encrypt(l.Password)
	// 检查密码是否正确
	record, err := dao.User.M.One("password=?", pass)
	if record.IsEmpty() || err != nil {
		return nil,errors.New("密码错误")
	}
	return u,err
}

func Register(p *request.Register) (err error) {
	one, err := dao.User.M.FindOne(g.Map{"username": p.Username})
	if !one.IsEmpty() && err == nil {
		return errors.New("用户已存在,注册失败")
	}
	pass,_ := gmd5.Encrypt(p.Password)
	u := model.User{
		UserName: p.Username,
		Password: pass,
	}
	if _, err := dao.User.M.Insert(&u); err != nil{
		return errors.New("注册失败")
	}
	return
}
`

var mod = `module {{.appName}}

go 1.15

require (
	github.com/gogf/gf v1.15.1
	github.com/gogf/gf-jwt v1.1.1
	github.com/mojocn/base64Captcha v1.3.1
)
`

var readme = `# GoFrame Project

https://goframe.org

├── app          #业务逻辑层	 所有的业务逻辑存放目录。
│   ├── api		 #业务接口	 接收/解析用户输入参数的入口/接口层。
│   ├── dao		 #数据访问    数据库的访问操作，仅包含最基础的数据库CURD方法
│   ├── model    #结构模型    数据结构管理模块，管理数据实体对象，以及输入与输出数据结构定义
│   └── service  #逻辑封装    业务逻辑封装管理，实现特定的业务逻辑实现和封装
├── boot		 #初始化包    用于项目初始化参数设置，往往作为main.go中第一个被import的包
├── config		 #配置管理    所有的配置文件存放目录。
├── docker	     #镜像文件    Docker镜像相关依赖文件，脚本文件等等
├── document	 #项目文档    Documentation项目文档，如: 设计文档、帮助文档等等
├── i18n		 #I18N国际化  I18N国际化配置文件目录
├── library		 #公共库包    公共的功能封装包，往往不包含业务需求实现
├── packed	     #打包目录    将资源文件打包的Go文件存放在这里，boot包初始化时会自动调用
├── public		 #静态目录    仅有该目录下的文件才能对外提供静态服务访问
├── router		 #路由注册    用于路由统一的注册管理
├── template     #模板文件    MVC模板文件存放的目录
├── Dockerfile   #镜像描述    云原生时代用于编译生成Docker镜像的描述文件
├── go.mod		 #依赖管理    使用Go Module包管理的依赖描述文件
└── main.go		 #入口文件    程序入口文件
`

var dao = `package internal

import (
	"context"
	"{{.appName}}/app/model"
	"database/sql"
	"github.com/gogf/gf/database/gdb"
	"github.com/gogf/gf/frame/g"
	"github.com/gogf/gf/frame/gmvc"
	"time"
)

// UserDao 是用于逻辑模型数据访问和自定义数据操作功能管理的管理器
type UserDao struct {
	gmvc.M
	DB      gdb.DB
	Table   string
	Columns userColumns
}

// User列定义并存储表user的列名
type userColumns struct {
    Id             string //                                                          
    UserName       string // 用户名                                                   
    NickName       string // 用户昵称                                                 
    Password       string // 登陆密码                                                 
    Email          string // 用户登录邮箱                                             
    Mobile         string // 中国手机不带国家代码，国际手机号格式为：国家代码-手机号  
    Avatar         string // 用户头像                                                 
    Birthday       string // 生日                                                     
    Sex            string // 性别;0:保密,1:男,2:女                                    
    Status         string // 用户状态;0:禁用,1:正常,2:未验证                          
    CreateAt       string // 注册时间                                                 
    LastLoginTime  string // 最后登录时间                                             
    LastLoginIp    string // 最后登录ip                                               
    DeptId         string // 部门id                                                   
    Remark         string // 备注                                                     
    IsAdmin        string // 是否后台管理员 1 是  0   否
}

var (
	// User 是表user操作的全局公共可访问对象
	User = &UserDao{
		M:     g.DB("default").Model("user").Safe(),
		DB:    g.DB("default"),
		Table: "user",
		Columns: userColumns{
			Id:            "id",               
            UserName:      "user_name",        
            NickName:      "nick_name",        
            Password:      "password",         
            Email:         "email",            
            Mobile:        "mobile",           
            Avatar:        "avatar",           
            Birthday:      "birthday",         
            Sex:           "sex",              
            Status:        "status",           
            CreateAt:      "create_at",        
            LastLoginTime: "last_login_time",  
            LastLoginIp:   "last_login_ip",    
            DeptId:        "dept_id",          
            Remark:        "remark",           
            IsAdmin:       "is_admin",
		},
	}
)

// Ctx 是一个链接函数，它创建并返回一个新的DB，该DB是当前DB对象的浅层副本，其中包含给定的上下文。
// 请注意，返回的DB对象只能使用一次，因此不要将其分配给全局或包变量以供长期使用。
func (d *UserDao) Ctx(ctx context.Context) *UserDao {
	return &UserDao{M: d.M.Ctx(ctx)}
}

// As 设置当前表的别名。
func (d *UserDao) As(as string) *UserDao {
	return &UserDao{M: d.M.As(as)}
}

// TX 设置当前操作的事务。
func (d *UserDao) TX(tx *gdb.TX) *UserDao {
	return &UserDao{M: d.M.TX(tx)}
}

// Master 指定操作是在主节点上进行。
func (d *UserDao) Master() *UserDao {
	return &UserDao{M: d.M.Master()}
}

// Slave 指定操作是在从节点上执行。(请注意，只有在配置了任何从属节点时才有意义)
func (d *UserDao) Slave() *UserDao {
	return &UserDao{M: d.M.Slave()}
}

// Args 为模型操作设置自定义参数。
func (d *UserDao) Args(args ...interface{}) *UserDao {
	return &UserDao{M: d.M.Args(args ...)}
}

// LeftJoin 对模型执行“left join ... on ...”语句。
// 参数<table>可以是联接表及其联接条件，也可以是其别名，例如:
// Table("user").LeftJoin("user_detail", "user_detail.uid=user.uid")
// Table("user", "u").LeftJoin("user_detail", "ud", "ud.uid=u.uid")
func (d *UserDao) LeftJoin(table ...string) *UserDao {
	return &UserDao{M: d.M.LeftJoin(table...)}
}

// RightJoin 对模型执行“right join ... on ...”语句。
// 参数<table>可以是联接表及其联接条件，也可以是其别名，例如:
// Table("user").RightJoin("user_detail", "user_detail.uid=user.uid")
// Table("user", "u").RightJoin("user_detail", "ud", "ud.uid=u.uid")
func (d *UserDao) RightJoin(table ...string) *UserDao {
	return &UserDao{M: d.M.RightJoin(table...)}
}

// InnerJoin 对模型执行“inner join ... on ...”语句。
// 参数<table>可以是联接表及其联接条件，也可以是其别名，例如:
// Table("user").InnerJoin("user_detail", "user_detail.uid=user.uid")
// Table("user", "u").InnerJoin("user_detail", "ud", "ud.uid=u.uid")
func (d *UserDao) InnerJoin(table ...string) *UserDao {
	return &UserDao{M: d.M.InnerJoin(table...)}
}

// Fields 指定需要操作的表字段，多个字段使用字符','连接。
// 参数<fieldNamesOrMapStruct>的类型可以是string/map/*map/struct/*struct
func (d *UserDao) Fields(fieldNamesOrMapStruct ...interface{}) *UserDao {
	return &UserDao{M: d.M.Fields(fieldNamesOrMapStruct...)}
}

// FieldsEx 指定例外的字段，(不被操作的字段)，多个字段使用字符','连接。
// 参数<fieldNamesOrMapStruct>的类型可以是string/map/*map/struct/*struct。
func (d *UserDao) FieldsEx(fieldNamesOrMapStruct ...interface{}) *UserDao {
	return &UserDao{M: d.M.FieldsEx(fieldNamesOrMapStruct...)}
}

// Option 设置模型的“额外操作”选项。
func (d *UserDao) Option(option int) *UserDao {
	return &UserDao{M: d.M.Option(option)}
}

// OmitEmpty 空值过滤，(过滤输入参数中的空值: nil,"",0)。
func (d *UserDao) OmitEmpty() *UserDao {
	return &UserDao{M: d.M.OmitEmpty()}
}

// Filter 过滤提交参数中不符合表结构的数据项。
func (d *UserDao) Filter() *UserDao {
	return &UserDao{M: d.M.Filter()}
}

// Where 设置模型的条件语句。参数<where>可以是string/map/gmap/slice/struct/*struct等类型。
// 请注意，如果多次调用它，将使用“AND”将多个条件连接到where语句中。
// 例如:
// Where("uid=10000")
// Where("uid", 10000)
// Where("money>? AND name like ?", 99999, "vip_%")
// Where("uid", 1).Where("name", "john")
// Where("status IN (?)", g.Slice{1,2,3})
// Where("age IN(?,?)", 18, 50)
// Where(User{ Id : 1, UserName : "john"})
func (d *UserDao) Where(where interface{}, args ...interface{}) *UserDao {
	return &UserDao{M: d.M.Where(where, args...)}
}

// WherePri方法的功能同Where，但提供了对表主键的智能识别。
// 如果主键是“id”，并且给定<where>参数为“123”，则WherePri函数将条件视为“id=123”，而M.where将条件视为字符串“123”。
func (d *UserDao) WherePri(where interface{}, args ...interface{}) *UserDao {
	return &UserDao{M: d.M.WherePri(where, args...)}
}

// And 在where语句中添加“AND”条件。
func (d *UserDao) And(where interface{}, args ...interface{}) *UserDao {
	return &UserDao{M: d.M.And(where, args...)}
}

// Or 在where语句中添加“OR”条件。
func (d *UserDao) Or(where interface{}, args ...interface{}) *UserDao {
	return &UserDao{M: d.M.Or(where, args...)}
}

// Group 分组 (设置模型的“group by”语句)。
func (d *UserDao) Group(groupBy string) *UserDao {
	return &UserDao{M: d.M.Group(groupBy)}
}

// Order 排序 (设置模型的“order by”语句)。
func (d *UserDao) Order(orderBy ...string) *UserDao {
	return &UserDao{M: d.M.Order(orderBy...)}
}

// Limit 设置模型的“limit”语句。
// 参数<limit>可以是一个或两个数字，如果传递了两个数字，则为模型设置“limit limit[0]、limit[1]”语句，否则设置“limit limit[0]”语句。
func (d *UserDao) Limit(limit ...int) *UserDao {
	return &UserDao{M: d.M.Limit(limit...)}
}

// Offset 设置模型的“offset”语句。
// 它只适用于某些数据库，如SQLServer、PostgreSQL等。
func (d *UserDao) Offset(offset int) *UserDao {
	return &UserDao{M: d.M.Offset(offset)}
}

// Page 设置模型的页码。参数<page>从1开始分页。
// 注意，对于“Limit”语句，Limit函数从0开始是不同的。
func (d *UserDao) Page(page, limit int) *UserDao {
	return &UserDao{M: d.M.Page(page, limit)}
}

// Batch 指定批量操作中分批操作的条数数量 (默认是10)
func (d *UserDao) Batch(batch int) *UserDao {
	return &UserDao{M: d.M.Batch(batch)}
}

// Cache 设置模型的缓存功能。  
// 它缓存sql的结果，这意味着如果有另一个相同的sql请求，它只是从缓存中读取并返回结果，而不是提交并执行到数据库中。
// 如果参数<duration><0，这意味着它用给定的<name>清除缓存。
// 如果参数<duration>=0，则表示它永不过期。
// 如果参数<duration>>0，则表示它在<duration>之后过期。
// 可选参数<name>用于将名称绑定到缓存，这意味着您可以稍后控制缓存，如更改<duration>或使用指定的<name>清除缓存。
// 请注意，如果模型正在对事务进行操作，则缓存功能将被禁用。
func (d *UserDao) Cache(duration time.Duration, name ...string) *UserDao {
	return &UserDao{M: d.M.Cache(duration, name...)}
}

// Data 设置模型的操作数据。
// 参数<data>可以是string/map/gmap/slice/struct/*struct等类型。
// 例如:
// Data("uid=10000")
// Data("uid", 10000)
// Data(g.Map{"uid": 10000, "name":"john"})
// Data(g.Slice{g.Map{"uid": 10000, "name":"john"}, g.Map{"uid": 20000, "name":"smith"})
func (d *UserDao) Data(data ...interface{}) *UserDao {
	return &UserDao{M: d.M.Data(data...)}
}

// Delete 用于数据的永久删除，被删除的数据不可恢复，请慎重使用。
// 往往需要结合Where、Order、Limit等方法共同使用，也可以直接给Delete方法传递where参数。
func (d *UserDao) Delete(data ...interface{}) (sql.Result, error) {
	res, err := d.M.Unscoped().Delete(data...)
    if err !=nil {
        return nil, err
    }
    return res, nil
}

// All 对模型执行“SELECT FROM…”语句。
// 它从表中检索记录并将结果返回为[]*model.User。如果没有使用表中给定的条件检索到记录，则返回nil。
// 可选参数<where>与M.where函数的参数相同，请参见M.where。
func (d *UserDao) All(where ...interface{}) ([]*model.User, error) {
	all, err := d.M.All(where...)
	if err != nil {
		return nil, err
	}
	var entities []*model.User
	if err = all.Structs(&entities); err != nil && err != sql.ErrNoRows {
		return nil, err
	}
	return entities, nil
}

// One 从表中检索一条记录，并将结果返回为*model.User。
// 如果没有使用表中给定的条件检索到记录，则返回nil。
// 可选参数<where>与M.where函数的参数相同，请参见M.where。
func (d *UserDao) One(where ...interface{}) (*model.User, error) {
	one, err := d.M.One(where...)
	if err != nil {
		return nil, err
	}
	var entity *model.User
	if err = one.Struct(&entity); err != nil && err != sql.ErrNoRows {
		return nil, err
	}
	return entity, nil
}

// FindOne 通过M.WherePri和M.One检索并返回单个记录。另见M.WherePri和M.One。
func (d *UserDao) FindOne(where ...interface{}) (*model.User, error) {
	one, err := d.M.FindOne(where...)
	if err != nil {
		return nil, err
	}
	var entity *model.User
	if err = one.Struct(&entity); err != nil && err != sql.ErrNoRows {
		return nil, err
	}
	return entity, nil
}

// FindAll 通过M.WherePri和M.All检索并返回结果集。另见M.WherePri和M.All。
func (d *UserDao) FindAll(where ...interface{}) ([]*model.User, error) {
	all, err := d.M.FindAll(where...)
	if err != nil {
		return nil, err
	}
	var entities []*model.User
	if err = all.Structs(&entities); err != nil && err != sql.ErrNoRows {
		return nil, err
	}
	return entities, nil
}

// Scan 根据参数<pointer>的类型自动调用Struct或Structs函数
// 例如:
// user  := new(User)
// err   := dao.User.Where("id", 1).Scan(user)
//
// user  := (*User)(nil)
// err   := dao.User.Where("id", 1).Scan(&user)
//
// users := ([]User)(nil)
// err   := dao.User.Scan(&users)
//
// users := ([]*User)(nil)
// err   := dao.User.Scan(&users)
func (d *UserDao) Scan(pointer interface{}, where ...interface{}) error {
	return d.M.Scan(pointer, where...)
}

// Chunk 使用给定的大小和回调函数迭代表。
func (d *UserDao) Chunk(limit int, callback func(entities []*model.User, err error) bool) {
	d.M.Chunk(limit, func(result gdb.Result, err error) bool {
		var entities []*model.User
		err = result.Structs(&entities)
		if err == sql.ErrNoRows {
			return false
		}
		return callback(entities, err)
	})
}

// LockUpdate 用于创建FOR UPDATE锁，避免选择行被其它共享锁修改或删除，FOR UPDATE会阻塞其他锁定性读对锁定行的读取
// 例如: db.Table("users").Where("votes>?", 100).LockUpdate().All() 等价于:
// SELECT * FROM 'users'' WHERE 'votes' > 100 FOR UPDATE
func (d *UserDao) LockUpdate() *UserDao {
	return &UserDao{M: d.M.LockUpdate()}
}

// LockShared 使用LockShared方法，在运行sql语句时带一把”共享锁“，共享锁可以避免被选择的行被修改，直到事务提交。
// 例如: db.Table("users").Where("votes>?", 100).LockShared().All() 等价于:
// SELECT * FROM 'users'' WHERE 'votes' > 100 LOCK IN SHARE MODE
func (d *UserDao) LockShared() *UserDao {
	return &UserDao{M: d.M.LockShared()}
}

// Unscoped 启用/禁用软删除功能。
func (d *UserDao) Unscoped() *UserDao {
	return &UserDao{M: d.M.Unscoped()}
}`

var daoUser = `package dao

import (
	"{{.appName}}/app/dao/internal"
)

// User 是表user操作的全局公共可访问对象。
var User = &userDao{internal.User}

// userDao 是用于逻辑模型数据访问和自定义数据操作功能管理的管理器
// 您可以在其上定义方法来扩展其功能
type userDao struct {
	*internal.UserDao
}`

var packed = `package packed`
