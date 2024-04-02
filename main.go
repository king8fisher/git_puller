package main

import (
	"github.com/king8fisher/git_puller/config"
	"github.com/king8fisher/git_puller/git"
	"github.com/king8fisher/git_puller/output"
	"os"
	"os/signal"
	"syscall"
	"time"
)

func main() {
	destinationPath := config.CloneDestinationPath
	git.Clone(destinationPath)
	gitTicker := time.NewTicker(time.Second * config.CapSeconds)
	done := make(chan bool)
	go func() {
		for {
			select {
			case <-done:
				return
			case <-gitTicker.C:
				_ = git.PullCapped(destinationPath)
			}
		}
	}()
	sigChan := make(chan os.Signal, 1)
	signal.Notify(sigChan, os.Interrupt, syscall.SIGTERM)
	<-sigChan
	done <- true
	output.Info("info", "interrupted gracefully")
}
