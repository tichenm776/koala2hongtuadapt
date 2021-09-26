package http_api

import (
	"bytes"
	"encoding/json"
	"github.com/alecthomas/log4go"
	"github.com/gin-gonic/gin"
	"github.com/zhenorzz/snowflake"
	client "go-common/app/service/main/vip/dao/ele-api-client"
	"io/ioutil"
	"net/http"
	"strconv"
	"strings"
	"zhiyuan/koala2hongtuadapt/dao"
	hongtu2 "zhiyuan/koala2hongtuadapt/hongtu"
	"zhiyuan/koala2hongtuadapt/model"
	"zhiyuan/koala2hongtuadapt/server"
	"zhiyuan/koala2hongtuadapt/util"
	hongtu "zhiyuan/koala_api_go/hongtu_api"
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

func  CreateID()int{
	sf, err := snowflake.New(1)
	if err != nil {
		panic(err)
	}
	uuid,_ := sf.Generate()
	str_uuid := strconv.FormatUint(uuid, 10)
	//fmt.Println(str_uuid)
	value ,_ := strconv.Atoi(str_uuid)
	return value
}

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




	//log4go.Info(c.Request)
	////log4go.Info(c)
	//body, _ := ioutil.ReadAll(c.Request.Body)
	//m := make(map[string]interface{}, 0)
	//err := json.Unmarshal(body, &m)
	//if err != nil {
	//	//fmt.Println(err)
	//	log4go.Error("ShouldBind err", err)
	//	c.JSON(http.StatusOK, gin.H{
	//		"code":    -200,
	//		"err_msg": "无请求参数或请求参数有误!",
	//	})
	//	return
	//}
	//log4go.Info(m)
	loginparams := model.Login{}
	err := c.Bind(&loginparams)
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


func Compare(c *gin.Context){


	json_return := `{"face_info_1":{"rect":{"left":87,"top":109,"right":262,"bottom":283},"quality":0.9955986086279154,"brightness":116.19357429718876,"std_deviation":28.79753933868924},"face_info_2":{"rect":{"left":87,"top":109,"right":262,"bottom":283},"quality":0.9955986086279154,"brightness":116.19357429718876,"std_deviation":28.79753933868924},"same":true,"score":98.16325378417969,"thresholds":{"E3":41.35199737548828,"E4":48.90372848510742,"E5":55.6849365234375,"E6":62.1713752746582,"recognizing":68,"stranger":67,"verify":68,"gate":78}}`
	json_data := make(map[string]interface{},0)
	json.Unmarshal([]byte(json_return),&json_data)
	c.JSON(http.StatusOK, json_data)

}

func Subjectsgrouplist(c *gin.Context){
	resp4Device := model.Resp4Device{Code: 0, Err_msg: ""}
	data := make([]map[string]interface{},0)
	//json_return := `{"code":0,"data":[{"comment":"floor_15\u8bbf\u5ba2\u7ec4","id":99,"name":"floor_15\u8bbf\u5ba2\u7ec4","subject_count":1,"subject_type":1,"update_by":"admin@91zo.com","update_time":1632282494.0},{"comment":"2\u697c\u8bbf\u5ba2\u7ec4","id":97,"name":"2\u697c\u8bbf\u5ba2\u7ec4","subject_count":0,"subject_type":1,"update_by":"admin@91zo.com","update_time":1630893529.0},{"comment":"3F\u8bbf\u5ba2\u7ec4","id":95,"name":"3F\u8bbf\u5ba2\u7ec4","subject_count":0,"subject_type":1,"update_by":"admin@91zo.com","update_time":1630639023.0},{"comment":"1\u697c\u8bbf\u5ba2\u7ec4","id":89,"name":"1\u697c\u8bbf\u5ba2\u7ec4","subject_count":0,"subject_type":1,"update_by":"admin@91zo.com","update_time":1630544330.0},{"comment":"floor_20\u8bbf\u5ba2\u7ec4","id":17,"name":"floor_20\u8bbf\u5ba2\u7ec4","subject_count":0,"subject_type":1,"update_by":"admin@91zo.com","update_time":1629269644.0}],"page":{"count":5,"current":1,"size":10,"total":1},"timecost":55}`


	//json_data := make(map[string]interface{},0)
	//json.Unmarshal([]byte(json_return),&json_data)

	//log4go.Info(json_data["data"].([]interface{})[0])


	//temp_map := make(map[string]interface{},0)
	//temp_map["type"] = 1
	//temp_map["pageNum"] = 1
	//temp_map["pageSize"] = 10000
	Passgrouplist,err := dao.FindGroupShip("")
	//Passgrouplist,err := hongtu2.GetGroupsHongtuList(temp_map)
	if err != nil{
		resp4Device.Code = -100
		resp4Device.Err_msg = "查询失败"+err.Error()
		log4go.Error(resp4Device.Err_msg)
		//return model.Visitor{}, errors.New(resp4Device.Err_msg)
		c.JSON(200, util.Err(resp4Device.Code,resp4Device.Err_msg,err))
		return
	}

	json_return := `{"comment":"floor_15\u8bbf\u5ba2\u7ec4","id":99,"name":"floor_15\u8bbf\u5ba2\u7ec4","subject_count":1,"subject_type":1,"update_by":"admin@91zo.com","update_time":1632282494.0}`

	json_data := make(map[string]interface{},0)
	json.Unmarshal([]byte(json_return),&json_data)

	if len(Passgrouplist) > 0{
		for k,_ := range Passgrouplist{
			temp := json_data
			temp["comment"] = Passgrouplist[k].Name
			temp["id"] = Passgrouplist[k].ID
			temp["name"] = Passgrouplist[k].Name
			data = append(data,temp)
		}

	}
	go func() {
		GetGroupsHongtu()
	}()

	//c.JSON(http.StatusOK, json_data)
	c.JSON(http.StatusOK, gin.H{
		"code":    0,
		"err_msg": "",
		"data":    data,
		"page":map[string]interface{}{},
	})
}


func GetGroupsHongtu_value(HongtuGroup map[string]interface{}){

	uuid := ""
	Name := ""
	typevalue := 2
	personTotal := -1
	if value,ok := HongtuGroup["uuid"].(string);ok{
		uuid = value
	}
	if uuid == ""{
		log4go.Error("错误数据:",HongtuGroup)
		//continue
	}
	Group := util.G_map_GroupsLocal.Get(uuid)
	if Group !=nil{
		if groupvalue,ok := Group.(map[string]interface{});ok{
			if groupvalue["name"] == HongtuGroup["name"]{
				log4go.Info("equal")
			}else{
				log4go.Info("update name")
				dao.UpdateGroupShip(uuid, map[string]interface{}{
					"name":HongtuGroup["name"],
				})
			}
		}
	}else{
		if value,ok := HongtuGroup["type"].(int);ok{
			typevalue = value
		}
		if value,ok := HongtuGroup["name"].(string);ok{
			Name = value
		}
		if value,ok := HongtuGroup["personTotal"].(int);ok{
			personTotal = value
		}
		addgroupship := model.GroupShip{
			Uuid:uuid,
			Type:typevalue,
			Name:Name,
		}
		valuereturn,err := dao.CreateGroupShip(&addgroupship)
		if err != nil{
			log4go.Error("添加失败:",err)
			//continue
		}
		temp_map := make(map[string]interface{})
		temp_map["id"]=valuereturn.ID
		temp_map["comment"]=valuereturn.Name
		temp_map["Uuid"]=valuereturn.Uuid
		temp_map["name"]=valuereturn.Name
		temp_map["subject_count"]=personTotal
		temp_map["subject_type"]=1
		temp_map["update_by"]="admin@91zo.com"
		temp_map["update_time"]=1632282494.0
		temp_map["type"]=valuereturn.Type
		temp_map["datatype"]="local"
		util.G_map_GroupsLocal.Set(uuid,temp_map)
	}

}

func AddPhoto(c *gin.Context){
	//resp4Device := model.Resp4Device{Code: 0, Err_msg: ""}
	_, header, err1 := c.Request.FormFile("photo")
	if err1 != nil{
		log4go.Error("获取照片文件失败:",err1)
		c.JSON(http.StatusOK, gin.H{
			"code":    0,
			"err_msg": "获取照片文件失败",
			"page":map[string]interface{}{},
		})
		return
	}
	//resp4Device := model.Resp4Device{Code: 0, Err_msg: ""}
	//data := make([]map[string]interface{},0)
	//json_return := `{"code":0,"data":[{"comment":"floor_15\u8bbf\u5ba2\u7ec4","id":99,"name":"floor_15\u8bbf\u5ba2\u7ec4","subject_count":1,"subject_type":1,"update_by":"admin@91zo.com","update_time":1632282494.0},{"comment":"2\u697c\u8bbf\u5ba2\u7ec4","id":97,"name":"2\u697c\u8bbf\u5ba2\u7ec4","subject_count":0,"subject_type":1,"update_by":"admin@91zo.com","update_time":1630893529.0},{"comment":"3F\u8bbf\u5ba2\u7ec4","id":95,"name":"3F\u8bbf\u5ba2\u7ec4","subject_count":0,"subject_type":1,"update_by":"admin@91zo.com","update_time":1630639023.0},{"comment":"1\u697c\u8bbf\u5ba2\u7ec4","id":89,"name":"1\u697c\u8bbf\u5ba2\u7ec4","subject_count":0,"subject_type":1,"update_by":"admin@91zo.com","update_time":1630544330.0},{"comment":"floor_20\u8bbf\u5ba2\u7ec4","id":17,"name":"floor_20\u8bbf\u5ba2\u7ec4","subject_count":0,"subject_type":1,"update_by":"admin@91zo.com","update_time":1629269644.0}],"page":{"count":5,"current":1,"size":10,"total":1},"timecost":55}`
	json_return := `{"company_id":1,"id":123,"origin_url":"/static/upload/origin/2021-09-26/v2_4b4b6a8f681d68254427e80506426b6aaf68a86c.jpg","quality":null,"subject_id":null,"url":"/static/upload/photo/2021-09-26/v2_a94965c0110044ae16350644f047b517e8f13162.jpg","version":7}`

	json_data := make(map[string]interface{},0)
	json.Unmarshal([]byte(json_return),&json_data)
	photo,err :=header.Open()
	if err != nil{
		log4go.Error("获取照片文件失败:",err)
		c.JSON(http.StatusOK, gin.H{
			"code":    0,
			"err_msg": "获取照片文件失败",
			"page":map[string]interface{}{},
		})
		return
	}
	data_return,err := hongtu.AddPhoto(photo)
	if err != nil{
		log4go.Error("添加照片失败:",err)
		c.JSON(http.StatusOK, gin.H{
			"code":    -100,
			"err_msg": "添加照片失败",
			"page":map[string]interface{}{},
		})
		return
	}
	log4go.Info(data_return)
	uri,err := data_return.Get("data").Get("uri").String()
	if err != nil{
		log4go.Error("获取照片uri失败:",err)
		c.JSON(http.StatusOK, gin.H{
			"code":    -100,
			"err_msg": "获取照片uri失败",
			"page":map[string]interface{}{},
		})
		return
	}

	photoship ,err := server.AddPhotoShip(uri)
	if err != nil{
		log4go.Error("存储uri失败:",err)
		c.JSON(http.StatusOK, gin.H{
			"code":    -100,
			"err_msg": "存储uri失败",
			"page":map[string]interface{}{},
		})
		return
	}

	json_data["id"] = photoship.ID

	//c.JSON(http.StatusOK, json_data)
	c.JSON(http.StatusOK, gin.H{
		"code":    0,
		"err_msg": "",
		"data":    json_data,
		"page":map[string]interface{}{},
	})
}

