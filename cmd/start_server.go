package cmd

import (
	"context"
	"fmt"
	"github.com/tmc/langchaingo/llms/openai"
	"go.uber.org/zap"
	"os"
	"os/signal"
	"syscall"
	"time"
	config "walrus_llm_project/common/conf"
	"walrus_llm_project/common/llm"
	"walrus_llm_project/common/utils"
	"walrus_llm_project/log"
)

func Execute() {
	defer func() {
		if err := recover(); err != nil {
			panicMsg := fmt.Sprintf("[panic] err: %v\nstack: %s\n", err, utils.GetCurrentGoroutineStack())
			log.Logger.Error("panic msg ", zap.String("panic msg", panicMsg))
		}
	}()
	err := initLLMClient()
	if err != nil {
		panic(err)
	}
	StartHttpServer(config.GetConfig().Server)
	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGHUP, syscall.SIGQUIT, syscall.SIGTERM, syscall.SIGINT)
	select {
	case <-quit:
		log.Logger.Warn("Shutdown Server ...")

		ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
		defer cancel()

		if err := ShutdownHttpServer(ctx); err != nil {
			log.Logger.Sugar().Error("Web Http Server Shutdown:", err)
		}

		log.Logger.Info("Server exited")
	}
}
func initLLMClient() error {
	openaiClient, err := openai.New(
		openai.WithModel("deepseek-chat"),
		openai.WithBaseURL("https://api.deepseek.com"),
		openai.WithToken("sk-"),
	)
	if err != nil {
		return err
	}
	llm.OpenaiClient = openaiClient
	return nil
}
