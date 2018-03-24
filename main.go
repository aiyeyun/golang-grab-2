package main

import (
	"runtime"
	"time"
	"xmn_2/core/algorithm/ssc"
	_ "xmn_2/core/model"
	"log"
	"os"
	"xmn_2/core/algorithm/shishicai/CustomPackage"
	//"xmn_2/core/algorithm/shishicai/play1"
	"xmn_2/core/algorithm/shishicai/play2"
)

func main(){
	log.Println("服务启动中．．．　进程ID:", os.Getpid())
	runtime.GOMAXPROCS(runtime.NumCPU())
	for {
		select {
		case <-time.After(1 * time.Minute):
			//todo
			// 时时彩　包含数据包　算法　邮件报警
			go ssc.Contain()
			// 时时彩 连号 算法 邮件报警
			go ssc.Consecutive()
			// 时时彩 连续AB表报警
			go ssc.ContailMultiple()
			// 时时彩 AB包 自定义A包周期 报警
			go CustomPackage.Calculation()
			// 时时彩 a出现几期的b
			//go play1.Calculation()
			// 时时彩 间隔几连号
			go play2.Consecutive()
		}
	}
}
