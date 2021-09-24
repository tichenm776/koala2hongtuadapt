package http_api

import (
	"bytes"
	"encoding/json"
	"github.com/alecthomas/log4go"
	"github.com/gin-gonic/gin"
	client "go-common/app/service/main/vip/dao/ele-api-client"
	"io/ioutil"
	"net/http"
	"strings"
	"zhiyuan/koala2hongtuadapt/model"
	"zhiyuan/koala2hongtuadapt/util"
	hongtu "zhiyuan/koala_api_go/hongtu_api"
	hongtu2 "zhiyuan/koala2hongtuadapt/hongtu"
	koala "zhiyuan/koala_api_go/koala_api"
	"zhiyuan/zyutil/config"
)



func LoginIn(c *gin.Context) {
	code := 0
	err_msg := ""

	// Read the Body content
	var bodyBytes []byte
	if c.Request.Body != nil {
		bodyBytes, _ = ioutil.ReadAll(c.Request.Body)
	}

	// Restore the io.ReadCloser to its original state
	c.Request.Body = ioutil.NopCloser(bytes.NewBuffer(bodyBytes))
	koala.LogRequest(c)

	username := c.PostForm("username")
	if username == "" {
		code = -100
		err_msg = "缺少用户名"
	}
	password := c.PostForm("password")
	if password == "" {
		code = -100
		err_msg = "缺少密码"
	}
	if code != 0 {
		c.JSON(http.StatusOK, gin.H{
			"code":    code,
			"err_msg": err_msg,
		})
		return
	}

	level := 0
	if strings.Contains(username, "@") {
		config.Init("./conf.yaml")
		hongtu.Init(config.Gconf.KoalaHost,80)
		err := hongtu.KoalaLogin(username, password)
		if err != nil {
			c.JSON(http.StatusOK, gin.H{
				"code":    -102,
				"err_msg": "账号密码错误!",
			})
			return
		}
		level = 1
		//configs.Init("./conf.yaml")
		_, err1 := config.Gconf.EditYaml(username, password)

		if err1 != nil {
			c.JSON(http.StatusOK, gin.H{
				"code":    -102,
				"err_msg": err1.Error(),
			})
			return
		}
	}

	if level == 1 {
		data := make(map[string]interface{})
		data["level"] = 1
		http.SetCookie(c.Writer, &http.Cookie{
			Name:     "username",
			Value:    username,
			MaxAge:   0,
			Path:     "/",
			Domain:   "",
			Secure:   false,
			HttpOnly: true,
		})
		c.JSON(http.StatusOK, gin.H{
			"code":    code,
			"err_msg": err_msg,
			"data":    data,
		})
		return
	} else {
		c.Request.Body = ioutil.NopCloser(bytes.NewBuffer(bodyBytes))
		koala.LoginZybox(level, c)
	}
}
var Purpose_map = map[int]string{0:"其他",1:"面试",2:"商务",3:"亲友",4:"快递送货"}

func AddPerson(c *gin.Context){
	resp4Device := model.Resp4Device{Code: 0, Err_msg: ""}
	subject :=	model.Subject{}
	temp_map := make(map[string]interface{},0)
	var (
		err  error
	)

	if err = c.Bind(&subject);err == nil{
		if err := util.Verify(subject,util.SubjectCreateDTOVerify);err != nil{
			c.JSON(200, util.Err(-1019,"参数错误",err))
			return
		}
	}
	if err != nil{
		c.JSON(200, util.Err(-1019,"账号或密码错误",err))
		return
	}
	uuid := client.UUID4()
	switch subject.Subject_type {
	case 0:
		//员工类型插入
		temp_map["type"] = 1
		temp_map["visitedUuid"] = 1
	case 1:
		//普通访客类型插入
		temp_map["type"] = 2
		temp_map["visitedUuid"] = 2
	case 2:
		//VIP访客插入
		temp_map["type"] = 2
	case 3:
		//黄名单插入
	}


	//temp_map["type"] = 1
	temp_map["name"] = subject.Name
	temp_map["uuid"] = uuid
	temp_map["phone"] = subject.Phone
	//添加照片与人员组
	//temp_map["imageUri"] = uri
	temp_map["identifyNum"] = subject.Remark
	temp_map["visitFirm"] = subject.Come_from
	temp_map["visitStartTimeStamp"] = subject.Start_time
	temp_map["visitEndTimeStamp"] = subject.End_time
	temp_map["visitReason"] = Purpose_map[subject.Purpose]
	temp_map["visitedUuid"] = uuid

	log4go.Debug("add koala temp", temp_map)

	list := []map[string]interface{}{temp_map}
	params := map[string]interface{}{
		"personList":list,
	}

	_, err = hongtu.AddPerson(params)
	if err != nil {
		resp4Device.Code = -100
		resp4Device.Err_msg = temp_map["name"].(string) + "申请失败"
		log4go.Error(resp4Device.Err_msg)
		//return model.Visitor{}, errors.New(resp4Device.Err_msg)
		c.JSON(200, util.Err(resp4Device.Code,resp4Device.Err_msg,err))
		return
	}
	//ship := model.SubjectShip{Uuid: uuid}
	//db.CreateSubjectShip()


}


func AuthLogin(c *gin.Context){

	loginparams := model.Login{}
	err := c.BindJSON(&loginparams)
	if err != nil{
		c.JSON(200, util.Err(-1019,"账号或密码错误",err))
		return
	}

	hongtu.Init(config.Gconf.KoalaHost, config.Gconf.KoalaPort)
	err = hongtu.KoalaLogin(loginparams.Username,loginparams.Password)
	if err != nil {

		c.JSON(200, util.Err(-1019,"账号或密码错误",err))
		return
	}

	json_return := `{"avatar":"","company":{"attendance_on":false,"attendance_weekdays":[1,2,3,4,5],"consigner":"雷","create_time":1619599999,"data_version":1,"deployment":1,"door_range":[[9,0],[21,0]],"door_weekdays":[1,2,3,4,5],"feature_version":7,"fmp_on":false,"full_day":false,"id":1,"is_certification":false,"lang":"中文简体","lang_code":"zh-Hans","logo":"/static/upload/logo/2021-07-07/v2_ab4254e1922a4eff462945d6a6ef421928bd5ac6.jpg","mask_detect":false,"max_temperature":37.3,"min_temperature":36.8,"name":"杭州之元科技有限公司","notdetermined_on":false,"remark":"","scenario":"正常使用","temperature_warn":false,"temperature_warn_last":2,"temperature_warn_open_door_limit":5,"upload":true,"yellowlist_warn":true},"company_id":1,"id":2,"lang":"中文简体","lang_code":"zh-Hans","organization_id":null,"password_reseted":true,"permission":[],"role_id":2,"username":"admin@91zo.com","verify":false}`
	json_data := make(map[string]interface{},0)
	json.Unmarshal([]byte(json_return),&json_data)
	go func() {
		_, err1 := config.Gconf.EditYaml(loginparams.Username, loginparams.Password)
		if err != nil{
			log4go.Error("修改配置文件失败:",err1)
		}
		log4go.Info("修改配置文件成功")
		return
	}()

	c.JSON(http.StatusOK, gin.H{
		"code":    0,
		"err_msg": "",
		"data":    json_data,
		"page":map[string]interface{}{},
	})

}





func GetEmployeeList(c *gin.Context){
	resp4Device := model.Resp4Device{Code: 0, Err_msg: ""}
	//subject :=	model.Subject{}
	temp_map := make(map[string]interface{},0)
	var (
		err  error
	)
	name,_ := c.Params.Get("name")

	personlist,err := hongtu2.GetEmployeeList(name)
	if err != nil {
		resp4Device.Code = -100
		resp4Device.Err_msg = temp_map["name"].(string) + "查询失败"
		log4go.Error(resp4Device.Err_msg)
		//return model.Visitor{}, errors.New(resp4Device.Err_msg)
		c.JSON(200, util.Err(resp4Device.Code,resp4Device.Err_msg,err))
		return
	}

	//类型转换 hongtu->koala

	//ship := model.SubjectShip{Uuid: uuid}
	//db.CreateSubjectShip()


	c.JSON(http.StatusOK, gin.H{
		"code":    0,
		"err_msg": "",
		"data":    personlist,
		"page":map[string]interface{}{},
	})



}



