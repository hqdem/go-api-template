package main

import (
	"context"
	"fmt"
	"github.com/hqdem/go-api-template/cmd"
	"github.com/hqdem/go-api-template/pkg/xlog"
	_ "go.uber.org/automaxprocs"
	"go.uber.org/zap"
)

func main() {
	err := cmd.Execute()
	if err != nil {
		xlog.Fatal(context.Background(), fmt.Sprintf("can't run app: %s", err), zap.Error(err))
	}
}
