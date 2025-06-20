//go:build windows
// +build windows

package main

import (
	"fmt"
	"github.com/kardianos/service"
	log "github.com/sirupsen/logrus"
	"os"
	"path/filepath"
	"time"
)

func Daemonize(logfile string, proc func()) {
	if err := ChPwd(); err != nil {
		log.WithFields(log.Fields{"err": err}).Fatal("Unable to change working directory")
		return
	}
	svcConfig := &service.Config{
		Name:        "AIRIOT",
		DisplayName: "AirIot物联网平台进程管理",
	}

	prg := &windowsService{
		proc: proc,
	}
	s, err := service.New(prg, svcConfig)
	if err != nil {
		log.WithFields(log.Fields{"err": err}).Fatal("service init failed")
		return
	}
	if err := s.Run(); err != nil {
		log.WithFields(log.Fields{"err": err}).Fatal("service run failed")
		return
	}
}

type windowsService struct {
	proc func()
}

func (p *windowsService) Start(s service.Service) (err error) {
	go p.proc()
	return nil
}

func (p *windowsService) Stop(s service.Service) error {
	fmt.Println("服务准备停止...")
	if service.Interactive() {
		time.Sleep(10 * time.Second)
		os.Exit(0)
	}
	return nil
}

func ChPwd() error {
	pwd, _ := os.Getwd()
	fmt.Println("开始工作目录:", pwd)
	// 程序所在目录
	execDir, err := filepath.Abs(filepath.Dir(os.Args[0]))
	if err != nil {
		return err
	}
	fmt.Println("程序所在目录:", execDir)

	if pwd != execDir {
		fmt.Println("切换工作目录到", execDir)
		if err := os.Chdir(execDir); err != nil {
			return err
		}
		pwd, _ = os.Getwd()
		fmt.Println("切换后工作目录:", pwd)
	}
	return nil
}
