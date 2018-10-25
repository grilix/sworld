package main

import (
	"flag"
	"fmt"
	"log"
	"net/http"
	"os"
	"os/signal"
	"syscall"

	klog "github.com/go-kit/kit/log"
	"github.com/grilix/sworld/server"
	"github.com/grilix/sworld/sworld"
	"github.com/grilix/sworld/sworldservice"
)

func main() {
	var (
		httpAddr = flag.String("http.addr", ":8089", "HTTP listen address")
	)
	flag.Parse()

	var logger klog.Logger
	{
		logger = klog.NewLogfmtLogger(os.Stderr)
		logger = klog.With(logger, "ts", klog.DefaultTimestampUTC)
		logger = klog.With(logger, "caller", klog.DefaultCaller)
	}

	world := sworld.NewWorld()

	var service sworldservice.Service
	{
		service = sworldservice.NewService(world)
	}

	var h http.Handler
	{
		h = server.MakeHTTPServer(service, klog.With(logger, "component", "HTTP"))
	}

	errs := make(chan error)
	go func() {
		c := make(chan os.Signal)
		signal.Notify(c, syscall.SIGINT, syscall.SIGTERM)
		errs <- fmt.Errorf("%s", <-c)
	}()

	go func() {
		logger.Log("transport", "HTTP", "addr", *httpAddr)
		errs <- http.ListenAndServe(*httpAddr, h)
	}()

	logger.Log("exit", <-errs)
}
