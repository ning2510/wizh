package main

import (
	"fmt"
	"log"
	"wizh/cmd/api/router"
	"wizh/pkg/viper"
	"wizh/pkg/zap"

	"github.com/gin-gonic/gin"
)

var (
	config     = viper.InitConf("api")
	serverAddr = fmt.Sprintf("%s:%d", config.Viper.GetString("server.host"), config.Viper.GetInt("server.port"))
)

func main() {
	logger := zap.InitLogger()

	r := gin.Default()
	router.InitRouter(r)
	if err := r.Run(serverAddr); err != nil {
		logger.Fatalln(err)
	}
	log.Printf("Listen to %v", serverAddr)
}
