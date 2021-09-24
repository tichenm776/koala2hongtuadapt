package http_api

import (
	"encoding/json"
	"fmt"
	"github.com/DeanThompson/ginpprof"
	"github.com/alecthomas/log4go"
	"github.com/bitly/go-simplejson"
	"github.com/gin-gonic/gin"
	"github.com/gorilla/websocket"
	"github.com/robfig/cron"
	"github.com/unrolled/secure"
	"net/http"
	"net/url"
	"strconv"
	"strings"
	"time"
	"zhiyuan/koala2hongtuadapt/dao"
	"zhiyuan/koala2hongtuadapt/hongtu"
	"zhiyuan/koala2hongtuadapt/server"
	"zhiyuan/koala2hongtuadapt/util"
	"zhiyuan/koala_api_go/b3r_api"
	koala "zhiyuan/koala_api_go/koala_api"
	"zhiyuan/zyutil"
	"zhiyuan/zyutil/config"
)

func New(port string) {
	go func() {
		koalainit()
		//TemplateIdInit()
	}()

	engine := gin.Default()
	initRouter(engine)
	//ModifySubject()
	gin.SetMode(gin.ReleaseMode)
	err := engine.Run(":" + port)
	if err != nil {
		log4go.Error("启动服务失败，原因：" + err.Error())
	} else {
		log4go.Info("服务启动成功...")
	}
}


func koalainit ()bool{
	config.Init("conf.yaml")
	koalalogin := KoalaLogin() //2021.5.28 临时注释
	if koalalogin {
		//go	Identify_Init(config.Gconf.KoalaHost,config.Gconf.KoalaUsername,config.Gconf.KoalaPassword)
		return koalalogin
	}else {
		time.Sleep(30*time.Second)
		return koalainit()
	}

}



func Identify_Init(koala_host,koala_username,koala_password string)(){

	client := &http.Client{
		Jar: nil,
	}
	params := map[string]interface{}{
		"koala_host":koala_host,
		"koala_username":koala_username,
		"koala_password":koala_password,
		"callbackurl":"http://127.0.0.1/device/v1/visit/signin",
	}

	jsonBytes, err := json.Marshal(params)
	log4go.Info(string(jsonBytes))
	req, err := http.NewRequest("POST","http://127.0.0.1:9052/v1/voice/screens" , strings.NewReader(string(jsonBytes)))
	if err != nil {
		log4go.Error(err)
		return
	}
	log4go.Debug(req.URL)
	log4go.Debug(req.Method)
	req.Header.Set("Content-Type", "application/json")
	resp, err := client.Do(req)
	if err != nil {
		log4go.Error(err.Error())
		return
	}
	log4go.Error(resp.StatusCode)
	return
}



func GetSubjectsLocal(){
	dao.FindSubjectShip("")
}

func GetSubjectsHongtu(){
	hongtu.GetEmployeeList("")
}

func GetGroupsLocal(){
	dao.FindGroupShip("")
}

func GetGroupsHongtu(){
	temp_map := make(map[string]interface{},0)
	temp_map["type"] = 1
	temp_map["pageNum"] = 1
	temp_map["pageSize"] = 10000
	hongtu.GetGroupsHongtuList(temp_map)
}





func initRouter(e *gin.Engine) {
	ginpprof.Wrap(e)
	//e.Use(TlsHandler())
	//e.GET("/cilent",Websocketcilent)
	system := e.Group("")
	{
		system.GET("/start", howToStart)
		system.POST("/auth/login", AuthLogin)
		system.POST("/subject", AddPerson)
		system.DELETE("/subject", AuthLogin)
		system.POST("/subject/photo", AuthLogin)
		system.GET("/mobile-admin/subjects/list", GetEmployeeList)
		system.GET("/subjects/group/list", AuthLogin)


		//system.GET("/client", Websocketclient)
	}
	koalamatepython := e.Group("/v1")
	{
		//koalamatepython.GET("/start", howToStart)
		koalamatepython.POST("/login", LoginIn)
		koalamatepython.POST("/logout", koala.ReverseProxy)
		koalamatepython.GET("/status", koala.ReverseProxy)
		koalamatepython.PUT("/tatus/sys_time", koala.ReverseProxy)
		koalamatepython.GET("/status/check_network", koala.ReverseProxy)
		koalamatepython.PUT("/status/reboot", koala.ReverseProxy)
		koalamatepython.GET("/config/ip",koala.ReverseProxy)
		koalamatepython.PUT("/config/ip", koala.ReverseProxy)
		koalamatepython.GET("/config/face_server", koala.ReverseProxy)
		koalamatepython.PUT("/config/face_server", koala.ReverseProxy)
		koalamatepython.GET("/config/cameras", koala.ReverseProxy)
		koalamatepython.PUT("/config/cameras", koala.ReverseProxy)
		koalamatepython.GET("/config/vocies ", koala.ReverseProxy)
		koalamatepython.PUT("/config/voices", koala.ReverseProxy)
		koalamatepython.POST("/config/voice_upload", koala.ReverseProxy)
		koalamatepython.GET("/config/work_attendance", koala.ReverseProxy)
		koalamatepython.PUT("/config/work_attendance", koala.ReverseProxy)
		koalamatepython.GET("/logs", koala.ReverseProxy)
		koalamatepython.GET("/log/", koala.ReverseProxy)
		koalamatepython.POST("/update", koala.ReverseProxy)
		koalamatepython.GET("/work_attendance/export", koala.ReverseProxy)
	}
	apis := e.Group("/v1")
	{
		apis.GET("/start", howToStart)
		// apis.GET("/subject/purpose", GetPurpose)
		// apis.GET("/subject/group/list", GetPersonGroupList)
		// apis.POST("/subject/list", GetStaffs)
		// apis.POST("/subject",VisitorApplication)
		//apis.POST("/subject/binding", StaffsBinding)
		// apis.POST("/subject/photo/check",PhotoCheck)
		// apis.POST("/local/subject/delete",DeleteSubject)
		// apis.POST("/local/subjects",GetSubject)
		// apis.POST("/subject/update",DeleteNextTime)

		//apis.GET("/subject/binding", QueryStaffWechatBinding) //1.员工查询微信是否已绑定
		//apis.POST("/subject/binding", StaffWechatBinding)     //1.1员工微信绑定
		//
		//apis.GET("/visitor/binding", QueryVisitorWechatBinding) //2.访客查询微信是否已绑定
		//apis.POST("/visitor/binding", VisitorWechatBinding)     //2.1访客微信绑定
		//
		//apis.POST("/subject", VisitorApplication) //3.访客申请（员工或访客）
		//apis.DELETE("/subject", VisitorApplicationDelete) //3.访客申请删除（员工或访客）
		//
		//apis.GET("/subject/purpose", GetPurpose) //4.选择来访目的
		//
		//apis.GET("/visit/records", GetVisitors) //5.查询访客申请记录（员工或访客用户，都用本接口）
		//apis.GET("/visit/:id", GetVisitorById) //5.查询访客申请记录（员工或访客用户，都用本接口）
		//
		//apis.GET("/auth/records", GetVisitors) //6.员工查询授权记录
		//apis.GET("/qrcode/record", GetVisitor) //6.员工查询授权记录
		//
		//apis.PUT("/auth/:id", VisitorAuth) //7.员工授权
		//
		//apis.DELETE("/subject/binding/:weixinUserId", WechatUnBinding) //8.微信解绑
		//apis.GET("/wechat/is_binding", QueryWechatIsBinding) //9 查询微信是否已绑定
		//
		//apis.GET("/config/servers", GetConfig) //1.1 读取服务器参数 2021.6.4
		//apis.PUT("/config/servers", SetConfig) //1.2 设置服务器参数 2021.6.4
		//
		//apis.GET("/config/templates", GetTemp) //2.1 读取微信短信模板设置参数
		//apis.PUT("/config/templates", SetTemp) //2.2 设置微信短信模板参数
		//apis.DELETE("/visit/records/:id", DeleteVisitor) //3. 访客申请记录删除
		//apis.POST("/visit/signin", VisitorSign) //4.1 访客签到确认，生成访客二维码
		//apis.POST("/visit/qrcode", GetVisitor) //4.2 访客二维码验证



	}
	// test := e.Group("/v1")
	// {
	// 	//test.GET("/subject/test",Createtimes)
	// }
	//go Init() //2021.5.28 临时注释

	// hub := server.NewWebSocketServer()
	// go hub.Run()
	// e.GET("/ws", func(c *gin.Context) {
	// 	server.ServeWs(hub, c.Writer, c.Request)
	// })

}
func Connect() {
}



func CreateWebSocket()(){

















}
func startEvent(data []map[string]interface{}) {
	var(
		koala_ip string
		rtsp string
		camera_position string

	)
	var complete chan int = make(chan int)
	for _,v :=range data{
		koala_ip = v["box_address"].(string)
		rtsp = v["camera_address"].(string)
		camera_position = v["camera_position"].(string)
		log4go.Info("koala_ip is ",koala_ip)
		log4go.Info("rtsp is ",rtsp)
		enable := v["enable"].(float64)
		log4go.Info("enable is ",enable)
		//enable_string := strconv.FormatFloat(enable,'E',-1,64)
		if enable == 0{
			log4go.Info("enable is 0",enable)
			continue
		}
		var addr = "ws://" +koala_ip + ":9000/video?url="
		log4go.Info("addr is 0",addr)
		var camera_address1 = url.QueryEscape(rtsp)
		log4go.Info("camera_address1 is 0",camera_address1)
		camera_address1 = strings.TrimSpace(camera_address1)
		if len(camera_address1) > 0 {
			is_b3r := b3r_api.IsB3r(koala_ip)
			if is_b3r == false { // koala
				go startWebsocket(addr + camera_address1)
			} else { // B3R
				addr = "ws://" + koala_ip + "/StartFrame"
				go startB3rWebsocket(addr, camera_address1, camera_position)
			}
		}
	}

	// 等待
	<-complete
}

func startWebsocket(url_str string) {
	for {
		log4go.Info("connecting to %s", url_str)
		fmt.Println("connecting to %s", url_str)
		c, _, err := websocket.DefaultDialer.Dial(url_str, nil)
		if err != nil {
			log4go.Error("error dial:", err)
			time.Sleep(time.Second * 10)
			continue
		}
		defer c.Close()

		//c.SetReadLimit(512)
		//c.SetReadDeadline(time.Now().Add(pongWait))
		//c.SetPongHandler(func(string) error { c.SetReadDeadline(time.Now().Add(pongWait)); return nil })

		for {
			time.Sleep(time.Millisecond)
			_, message, err := c.ReadMessage()
			if err != nil {
				log4go.Info("read:", err)
				c.Close()
				time.Sleep(time.Second * 10)
				break
			}
			go logMessage(message)

			go process(message)
		}
	}
}
func startB3rWebsocket(url_str, camera_ip, camera_position string) {
	for {
		log4go.Info("connecting to %s", url_str)
		c, _, err := websocket.DefaultDialer.Dial(url_str, nil)
		if err != nil {
			log4go.Error("error dial:", err)
			time.Sleep(time.Second * 10)
			continue
		}
		defer c.Close()

		//c.SetReadLimit(512)
		//c.SetReadDeadline(time.Now().Add(pongWait))
		//c.SetPongHandler(func(string) error { c.SetReadDeadline(time.Now().Add(pongWait)); return nil })

		// 连接成功后，必须发type=0 的消息
		msg := []byte("{\"ws_type\":0, \"chn\":1}")
		err = c.WriteMessage(websocket.TextMessage,msg)
		if err != nil {
			c.Close()
			log4go.Error("error send message:", err)
			time.Sleep(time.Second * 10)
			continue
		}

		// websocket keepalive chec
		b3_count := 0
		ticker := time.NewTicker(time.Second * 8)
		defer ticker.Stop()
		go func() {
			for {
				<-ticker.C
				if b3_count > 60 {
					log4go.Debug("B3R websocket keepalive timeout")
					break
				} else {
					b3_count++
					err = c.WriteMessage(websocket.TextMessage,[]byte("{\"ws_type\":3}"))
					if err != nil {
						log4go.Error("B3R websocket error send message:", err)
						break
					}
				}
			}
			ticker.Stop()
			c.Close()
		}()

		for {
			time.Sleep(time.Millisecond)
			_, message, err := c.ReadMessage()
			if err != nil {
				log4go.Info("read:", err)
				c.Close()
				time.Sleep(time.Second * 10)
				break
			}

			go logMessage(message)

			go processB3r(message, camera_position)
		}
	}
}
func process(message []byte) {
	jdata, err := simplejson.NewJson(message)
	if err != nil {
		log4go.Error(err.Error())
		return //nil, errors.New("Face++返回报文错误")
	}

	data := jdata.Get("data")

	var recognize_status string
	status := data.Get("status").MustString()
	if status != "" {
		recognize_status = status
		fmt.Printf(string(status))
	} else {
		recognize_status = data.Get("status").Get("recognize_status").MustString()
	}

	if recognize_status == "recognized" || recognize_status == "unrecognized" {
		util.Process_recognized(string(message))
	}

	/*
		code, _ := jdata.Get("code").Int()
		if code != 0 {
			desc, _ := jdata.Get("desc").String()
			log4go.Error(desc)
			return nil, errors.New(desc)
		}
	*/

}

func processB3r(message []byte, camera_position string) {
	jdata, err := simplejson.NewJson(message)
	if err != nil {
		log4go.Error(err.Error())
		return
	}

	dataType := jdata.Get("dataType").MustString()
	if dataType != "Alert" {
		log4go.Error("报文里dataType != 'Alert'")
		return
	}

	fmpErr := jdata.Get("fmpErr").MustInt(1)
	if fmpErr == 1 {
		log4go.Error("报文里fmpErr = 1")
		return
	}

	timestamp := jdata.Get("timestamp").MustInt(1)

	var msg  gin.H
	aboveHigh := jdata.Get("aboveHigh").MustInt(0)
	if aboveHigh == 1 {  // recognized
		description := jdata.Get("description").MustString()
		desc := strings.Split(description, ",")
		subject_id, err:= strconv.Atoi(desc[0])
		if err != nil {
			subject_id = -1
		}
		subject_type, err:= strconv.Atoi(desc[1])
		if err != nil {
			subject_type = -1
		}
		msg = gin.H {
			"error_code": 0,
			"error": "",
			"type": "recognized",
			"screen": gin.H{
				"camera_position": camera_position,
			},
			"person": gin.H{
				"subject_type": subject_type,
				"id": subject_id,
				"entry_date": 0,
				"birthday": 0,
				"name": desc[4],
				"description": desc[6],
			},
			"data": gin.H{
				"timestamp": timestamp,
				"status": "recognized",
				"person": gin.H{
					"subject_id": subject_id,
				},
			},
		}
	}else{// unrecognized
		msg = gin.H {
			"error_code": 0,
			"error": "",
			"type": "unrecognized",
			"screen": gin.H{
				"camera_position": camera_position,
			},
			"data": gin.H{
				"status": "unrecognized",
				"timestamp": timestamp,
			},
		}
	}


	jsonByte, err := json.Marshal(msg)
	if err != nil {
		fmt.Printf("Marshal with error: %+v\n", err)
		return
	}

	util.Process_recognized(string(jsonByte))

}
func logMessage(message []byte) {
	jdata, err := simplejson.NewJson(message)
	if err != nil {
		log4go.Error(err.Error())
		return //nil, errors.New("Face++返回报文错误")
	}

	msg_type := jdata.Get("type").MustString()
	if msg_type == "attr" {
		jdata.Get("data").Set("image", "")
	} else {

		jdata.Get("data").Get("face").Set("image", "")
		jdata.Get("person").Set("src", "")
	}
	resp, err := jdata.MarshalJSON()
	if err != nil {
		log4go.Error(err.Error())
		return //nil, errors.New("Face++返回报文错误")
	}
	log4go.Info("%s", resp)
}



func KoalaLogin() bool {
	zyutil.Recover()
	hongtu.Init(config.Gconf.KoalaHost, config.Gconf.KoalaPort)
	err := hongtu.KoalaLogin(config.Gconf.KoalaUsername, config.Gconf.KoalaPassword)
	if err != nil {
		return false
	}
	go Keepalive(config.Gconf.KoalaUsername, config.Gconf.KoalaPassword)
	return true
}

func Keepalive(username, password string) {
	for {
		time.Sleep(30 * time.Minute)
		hongtu.KoalaLogin(username, password)
	}
}

func ModifySubject() {
	cronTarget := cron.New()
	spec := "*/30 * * * * ?"
	cronTarget.AddFunc(spec, func() {
		server.DeleteNextTime()
	})
	cronTarget.Start()

}

// example for http request handler.
func howToStart(c *gin.Context) {
	c.String(0, "Golang 大法好 !!!")
}






func TlsHandler() gin.HandlerFunc {
	return func(c *gin.Context) {
		secureMiddleware := secure.New(secure.Options{
			SSLRedirect: true,
			SSLHost:     "localhost:8080",
		})
		err := secureMiddleware.Process(c.Writer, c.Request)

		// If there was an error, do not continue.
		if err != nil {
			return
		}

		c.Next()
	}
}
