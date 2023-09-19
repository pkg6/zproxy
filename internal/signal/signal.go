package signal

import (
	"log"
	"os"
	"os/signal"
	"syscall"
)

func Clean() {
	signalChan := make(chan os.Signal, 1)
	cleanupDone := make(chan bool)
	//os.Interrupt: 相当于SIGINT信号,通常表示用户中断,比如Ctrl + C触发。
	//syscall.SIGHUP: 表示终端断开连接,通常用于管道命令、后台运行进程的关闭。
	//syscall.SIGINT: 中断信号,和os.Interrupt等价,表示用户中断。
	//syscall.SIGTERM: 程序终止信号,表示正常终止程序。
	//syscall.SIGQUIT: 和SIGTERM类似,表示终止程序运行。
	//signal.Notify(signalChan, os.Interrupt, syscall.SIGHUP, syscall.SIGINT, syscall.SIGTERM, syscall.SIGQUIT)
	signal.Notify(signalChan, os.Interrupt, syscall.SIGINT, syscall.SIGTERM, syscall.SIGQUIT)
	go func() {
		for range signalChan {
			log.Println("stopping services...")
			cleanupDone <- true
		}
	}()
	<-cleanupDone
}
