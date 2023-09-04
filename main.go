package main

import (
	"concise_douyin/config"
  "concise_douyin/router"
	"fmt"
)

func main() {
	r := router.Init()
	/// 用原生http的服务方式http.ListenAndServe(ip address,r)启动是一样的
	var err = r.Run(fmt.Sprintf(":%d", config.Global.Port))
	if err != nil {
		return
	}
}

