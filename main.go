package main

import (
	"BackendTemplate/pkg/api"
	"BackendTemplate/pkg/database"
	"BackendTemplate/pkg/utils"
	"embed"
	"flag"
	"fmt"
	"os"
	"strconv"

	"github.com/gin-gonic/gin"
)

//go:embed dist
var embedFS embed.FS

func main() {
	utils.InitFunction()
	gin.SetMode(gin.DebugMode)
	var bindPort = flag.Int("p", 8089, "Specify alternate port")
	flag.Parse()
	if *bindPort > 65535 || *bindPort < 0 {
		flag.Usage()
		os.Exit(0)
	}
	database.ConnectDateBase()
	defer database.Engine.Close()

	database.Engine.Update(&database.Listener{Status: 2})
	database.Engine.Update(&database.WebDelivery{Status: 2})

	r := gin.New()

	// 配置所有路由和中间件
	api.SetupRoutes(r, embedFS)

	fmt.Println("Listening on port ", *bindPort)
	r.Run("0.0.0.0:" + strconv.Itoa(*bindPort)) // 启动服务
}
