package cmd

import (
	"context"
	"errors"
	"fmt"
	"github.com/gin-gonic/gin"
	"go.uber.org/zap"
	"net/http"
	"walrus_llm_project/api/route"
	config "walrus_llm_project/common/conf"
	"walrus_llm_project/log"
)

var (
	httpServer *http.Server
)

func StartHttpServer(conf config.Server) {

	gin.SetMode(conf.RunMode)
	engine := gin.New()
	route.Setup(engine)

	httpServer = &http.Server{
		Addr:    fmt.Sprintf(":%d", conf.HttpPort),
		Handler: engine,
	}
	go func() {
		log.Logger.Info("http server listen:", zap.Int("port", conf.HttpPort))
		fmt.Println("httpsvr server start")
		if err := httpServer.ListenAndServe(); err != nil && !errors.Is(err, http.ErrServerClosed) {
			log.Logger.Error("httpsvr server listen failed ", zap.Error(err))
			panic(err)
		}

	}()
}

func ShutdownHttpServer(ctx context.Context) []error {
	var errList []error

	if httpServer != nil {
		err := httpServer.Shutdown(ctx)
		if err != nil {
			errList = append(errList, err)
		}
	}

	return errList
}
