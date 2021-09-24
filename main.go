package main

import (
	"fmt"
	"github.com/alecthomas/log4go"
	"runtime"
	"strconv"
	db "zhiyuan/koala2hongtuadapt/dao"
	"zhiyuan/koala2hongtuadapt/http_api"
	"zhiyuan/zyutil/config"
)

func main() {
	if runtime.GOOS == "linux" {
		log4go.LoadConfiguration("./koala2hongtuadapt.xml")
	}
	//} else {
	//	log4go.LoadConfiguration("./log4wechatservice.xml")
	//}

	err := db.New()
	if err != nil {
		log4go.Error(err.Error())
		fmt.Println(err.Error())
		return
	}
	config.Init("conf.yaml")
	//go util.G_map.CronDelete(600)
	//http.ListenAndServe(":7300", r)
	http_api.New(strconv.Itoa(config.Gconf.ServerPort_theme))
	//http_api.New("9010")

}
