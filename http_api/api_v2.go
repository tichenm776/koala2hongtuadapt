package http_api

import (
	"bytes"
	"encoding/base64"
	"encoding/json"
	"github.com/alecthomas/log4go"
	"github.com/gin-gonic/gin"
	"github.com/zhenorzz/snowflake"
	client "go-common/app/service/main/vip/dao/ele-api-client"
	"io"
	"io/ioutil"
	"mime/multipart"
	"net/http"
	"os"
	"os/exec"
	"reflect"
	"strconv"
	"strings"
	"time"
	"zhiyuan/koala2hongtuadapt/dao"
	hongtu2 "zhiyuan/koala2hongtuadapt/hongtu"
	"zhiyuan/koala2hongtuadapt/model"
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
			"desc": err_msg,
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
				"desc": "账号密码错误!",
			})
			return
		}
		level = 1
		//configs.Init("./conf.yaml")
		_, err1 := config.Gconf.EditYaml(username, password)

		if err1 != nil {
			c.JSON(http.StatusOK, gin.H{
				"code":    -102,
				"desc": err1.Error(),
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
			"desc": err_msg,
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
func typeof(v interface{}) string {
	return reflect.TypeOf(v).String()
}
func AddPerson(c *gin.Context){

	log4go.Info(c.Request.Header)
	clientip := c.Request.Header.Get("X-Real-IP")
	//clientip := c.ClientIP()
	log4go.Info(c)
	body, _ := ioutil.ReadAll(c.Request.Body)
	m := make(map[string]interface{}, 0)
	json.Unmarshal(body, &m)
	log4go.Info("m", m)
	resp4Device := model.Resp4Device{Code: 0, Err_msg: ""}
	//subject :=	model.Subject{}
	temp_map := make(map[string]interface{},0)
	var (
		//err  error
	)
	come_from := ""
	department := ""
	name := ""
	phone := ""
	remark := ""
	interviewee_id := ""
	end_time := int64(0)
	start_time := int64(0)
	subject_type := 0
	Purpose := 0
	group_ids := []int{}

	//log4go.Info(typeof(m["group_ids"].([]interface{})[0]))
	//log4go.Info(m["group_ids"])


	if value,ok := m["come_from"].(string);ok{
		come_from = value
	}
	if value,ok := m["department"].(string);ok{
		department = value
	}
	//if value,ok := m["end_time"].(int);ok{
	//	end_time = value
	//}
	if value,ok := m["end_time"].(float64);ok{
		end_time = int64(value)
	}
	if value,ok := m["group_ids"].([]interface{});ok{
		if len(value) > 0{
			for k,_ := range value{
				if value2,ok := value[k].(float64);ok{
					group_ids = append(group_ids, int(value2))
				}
			}
		}
	}
	if value,ok := m["name"].(string);ok{
		name = value
	}
	if value,ok := m["phone"].(string);ok{
		phone = value
	}
	if value,ok := m["remark"].(string);ok{
		remark = value
		remarks_arr := strings.Split(value,"-")
		log4go.Info("remarks_arr is",remarks_arr)
		if len(remarks_arr) >= 2{
			remark = remarks_arr[0]
			interviewee_id = remarks_arr[1]
		}
	}
	//if value,ok := m["start_time"].(int);ok{
	//	start_time = value
	//}
	if value,ok := m["start_time"].(float64);ok{
		start_time = int64(value)
	}
	if value,ok := m["subject_type"].(float64);ok{
		subject_type = int(value)
	}
	if value,ok := m["purpose"].(float64);ok{
		Purpose = int(value)
	}
	//if err = c.Bind(&subject);err != nil{
	//	log4go.Error("获取参数错误"+err.Error())
	//		c.JSON(200, util.Err(-1019,"获取参数错误",err))
	//		return
	//}
	//log4go.Info(come_from)
	log4go.Info(department)
	//log4go.Info(name)
	//log4go.Info(phone)
	//log4go.Info(remark)
	//log4go.Info(interviewee_id)
	//log4go.Info(start_time)
	//log4go.Info(end_time)


	//log4go.Info(subject)
	//log4go.Info(subject.Group_ids)
	//log4go.Info(subject.Photo_ids)

	//if err != nil{
	//	c.JSON(200, util.Err(-1019,"账号或密码错误",err))
	//	return
	//}
	uuid := client.UUID4()
	switch subject_type {
	case 0:
		//员工类型插入
		temp_map["type"] = 1
		id ,_:= strconv.Atoi(interviewee_id)
		Subject,err := dao.FindSubjectShip2("","",id)
		if err != nil{
			resp4Device.Code = -100
			resp4Device.Err_msg = temp_map["name"].(string) + "申请失败"
			log4go.Error(resp4Device.Err_msg)
			//return model.Visitor{}, errors.New(resp4Device.Err_msg)
			c.JSON(200, util.Err(resp4Device.Code,resp4Device.Err_msg,err))
			return
		}
		temp_map["visitedUuid"] = Subject.Uuid
		if Subject.Uuid == ""{

		}
	case 1:
		//普通访客类型插入
		temp_map["type"] = 2
		id ,_:= strconv.Atoi(interviewee_id)
		Subject,err := dao.FindSubjectShip2("","",id)
		if err != nil{
				resp4Device.Code = -100
				resp4Device.Err_msg = temp_map["name"].(string) + "申请失败"
				log4go.Error(resp4Device.Err_msg)
				//return model.Visitor{}, errors.New(resp4Device.Err_msg)
				c.JSON(200, util.Err(resp4Device.Code,resp4Device.Err_msg,err))
				return
		}
		temp_map["visitedUuid"] = Subject.Uuid
		log4go.Info(Subject)

	case 2:
		//VIP访客插入
		temp_map["type"] = 2
		id ,_:= strconv.Atoi(interviewee_id)
		Subject,err := dao.FindSubjectShip2("","",id)
		if err != nil{
			resp4Device.Code = -100
			resp4Device.Err_msg = temp_map["name"].(string) + "申请失败"
			log4go.Error(resp4Device.Err_msg)
			//return model.Visitor{}, errors.New(resp4Device.Err_msg)
			c.JSON(200, util.Err(resp4Device.Code,resp4Device.Err_msg,err))
			return
		}
		temp_map["visitedUuid"] = Subject.Uuid
	case 3:
		//黄名单插入
	}
	groupList := []string{}

	if len(group_ids) > 0{
		for k,_ := range group_ids{
			result,err := dao.FindGroupShip2(group_ids[k])
			if err != nil{
				resp4Device.Code = -100
				resp4Device.Err_msg = temp_map["name"].(string) + "申请失败:查询访客组失败"
				log4go.Error(resp4Device.Err_msg)
				//return model.Visitor{}, errors.New(resp4Device.Err_msg)
				c.JSON(200, util.Err(resp4Device.Code,resp4Device.Err_msg,err))
				return
			}
			groupList = append(groupList, result.Uuid)
		}
	}

	//log4go.Info("temp_map[\"visitedUuid\"] is",temp_map["visitedUuid"])
	temp_map["name"] = name
	temp_map["uuid"] = uuid
	temp_map["phone"] = phone
	//添加照片与人员组
	//temp_map["imageUri"] = uri
	temp_map["identifyNum"] = remark
	temp_map["visitFirm"] = come_from
	temp_map["visitStartTimeStamp"] = start_time*1000
	//temp_map["visitStartTimeStamp"] = start_time*1000
	temp_map["visitEndTimeStamp"] = end_time*1000
	//temp_map["visitEndTimeStamp"] = end_time*1000
	temp_map["visitReason"] = Purpose_map[Purpose]
	temp_map["groupList"] = groupList
	log4go.Info("visitStartTimeStamp is",temp_map["visitStartTimeStamp"])
	log4go.Info("visitEndTimeStamp is",temp_map["visitEndTimeStamp"])
	//log4go.Debug("add koala temp", temp_map)

	list := []map[string]interface{}{temp_map}
	params := map[string]interface{}{
		"personList":list,
	}
	//log4go.Info("add params",params)
	_, err := hongtu.AddPerson(params)
	if err != nil {
		resp4Device.Code = -100
		resp4Device.Err_msg = temp_map["name"].(string) + "申请失败"
		log4go.Error(resp4Device.Err_msg)
		//return model.Visitor{}, errors.New(resp4Device.Err_msg)
		c.JSON(200, util.Err(resp4Device.Code,resp4Device.Err_msg,err))

		return
	}
	ship := &model.SubjectShip{Uuid: uuid,Name: name,
		ClientIp: clientip,IdentifyNum:remark,
		VisitType:temp_map["type"].(int),Phone:phone,Uri:"",}
	visitor,err := dao.CreateSubjectShip(ship)

	if err != nil {
		resp4Device.Code = -100
		resp4Device.Err_msg = temp_map["name"].(string) + "申请失败"
		log4go.Error(resp4Device.Err_msg)
		//return model.Visitor{}, errors.New(resp4Device.Err_msg)
		c.JSON(200, util.Err(resp4Device.Code,resp4Device.Err_msg,err))
		return
	}
	log4go.Info("visitor is",visitor)
	c.JSON(http.StatusOK, gin.H{
		"code":    0,
		"desc": "",
		"data":    visitor,
		"page":map[string]interface{}{},
	})
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
	//		"desc": "无请求参数或请求参数有误!",
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
		"desc": "",
		"data":    json_data,
		"page":map[string]interface{}{},
	})

}





func GetEmployeeList(c *gin.Context){
	resp4Device := model.Resp4Device{Code: 0, Err_msg: ""}
	////subject :=	model.Subject{}
	temp_map := make([]map[string]interface{},0)
	//var (
	//	err  error
	//)
	clientip := c.ClientIP()

	json_return := `{"avatar":"","birthday":null,"building":null,"come_from":"","company_id":1,"create_time":1632472126,"credential_no":"","credential_type":null,"department":"","description":"","domicile_address":"","domicile_city_code":null,"domicile_district_code":null,"domicile_province_code":null,"domicile_street_code":null,"education_code":null,"email":"","end_time":0,"entrance_people_type":null,"entry_date":null,"extra_id":null,"gender":1,"groups":[{"id":98,"name":"floor_15\u5458\u5de5\u7ec4"}],"house":null,"house_rel_code":1,"id":158,"interviewee":"","interviewee_pinyin":"","is_use":true,"job_number":"","marital_status_code":null,"name":"\u738b\u529b\u5b8f2","nation":null,"nationality_code":"","origin":"","password_reseted":false,"people_type":null,"phone":"18357036166","photos":[],"pinyin":"wanglihong2","purpose":0,"remark":"","residence_address":"","residence_city_code":null,"residence_district_code":null,"residence_province_code":null,"residence_street_code":null,"source":1,"start_time":0,"subject_type":0,"title":"","village":null,"visit_notify":false,"wg_number":""}`
	json_data := make(map[string]interface{},0)
	json.Unmarshal([]byte(json_return),&json_data)
	tempmap := json_data
	name := c.Query("name")
	size := c.Query("size")
	page := c.Query("page")
	log4go.Info("get size",size)
	log4go.Info("get page",page)
	size_int,err := strconv.Atoi(size)
	page_int,err := strconv.Atoi(page)
	log4go.Info("get size_int",size_int)
	log4go.Info("get page_int",page_int)
	//dao.FindSubjectShip("",name)
	if size_int == 0{
		size_int = 5
	}
	if page_int == 0{
		page_int = 1
	}
	personlist,pages,err := hongtu2.GetEmployeeList(name,page_int,size_int)
	if err != nil {
		resp4Device.Code = -100
		resp4Device.Err_msg = "查询失败"
		log4go.Error(resp4Device.Err_msg)
		//return model.Visitor{}, errors.New(resp4Device.Err_msg)
		c.JSON(200, util.Err(resp4Device.Code,resp4Device.Err_msg,err))
		return
	}

	//dao.DeleteShip(clientip)
	for k,_ := range personlist{
		//tmp_map := make(map[string]interface{},0)
		uuid := ""
		uri := ""
		name := ""
		identifyNum := ""
		phone := ""
		visitType := 1

		if value,ok := personlist[k]["uuid"].(string);ok{
			uuid = value
		}
		if value,ok := personlist[k]["imageUri"].(string);ok{
			if value != ""{
				index1 := strings.Index(value,"?")
				index2 := strings.Index(value,"pub/")
				uri = value[index2+4:index1]
			}
		}
		if value,ok := personlist[k]["name"].(string);ok{
			name = value
		}
		if value,ok := personlist[k]["identifyNum"].(string);ok{
			identifyNum = value
		}
		if value,ok := personlist[k]["visitType"].(int);ok{
			visitType = value
		}
		if value,ok := personlist[k]["phone"].(string);ok{
			phone = value
		}
		subject := model.SubjectShip{}
		temp_sub := model.SubjectShip{
			Uuid: uuid,
			Uri:uri,
			Name: name,
			Phone: phone,
			VisitType: visitType,
			IdentifyNum: identifyNum,
			ClientIp: clientip,
		}
		subject,err = dao.FindSubjectShip3(clientip,uuid)
		if err != nil{
			if strings.Index(err.Error(),"not found") != -1{
				subject,err = dao.CreateSubjectShip(&temp_sub)
				if err != nil{
					resp4Device.Code = -100
					resp4Device.Err_msg = "添加失败"+err.Error()
					log4go.Error(resp4Device.Err_msg)
					continue
				}
			}else{
				log4go.Error("查询失败",err)
				continue
			}
		}
		//log4go.Info("get subject",subject)
		tmp_map := Copymap(tempmap)
		id := subject.ID
		subjectname := subject.Name
		tmp_map.(map[string]interface{})["id"] = id
		tmp_map.(map[string]interface{})["name"] = subjectname

		temp_map = append(temp_map, tmp_map.(map[string]interface{}))
	}

	//log4go.Info("temp_map is",temp_map)

	//类型转换 hongtu->koala

	//ship := model.SubjectShip{Uuid: uuid}
	//db.CreateSubjectShip()
	//json_return := `{"code":0,"data":[{"avatar":"","birthday":null,"building":null,"come_from":"","company_id":1,"create_time":1632472126,"credential_no":"","credential_type":null,"department":"","description":"","domicile_address":"","domicile_city_code":null,"domicile_district_code":null,"domicile_province_code":null,"domicile_street_code":null,"education_code":null,"email":"","end_time":0,"entrance_people_type":null,"entry_date":null,"extra_id":null,"gender":1,"groups":[{"id":98,"name":"floor_15\u5458\u5de5\u7ec4"}],"house":null,"house_rel_code":1,"id":158,"interviewee":"","interviewee_pinyin":"","is_use":true,"job_number":"","marital_status_code":null,"name":"\u738b\u529b\u5b8f2","nation":null,"nationality_code":"","origin":"","password_reseted":false,"people_type":null,"phone":"18357036166","photos":[],"pinyin":"wanglihong2","purpose":0,"remark":"","residence_address":"","residence_city_code":null,"residence_district_code":null,"residence_province_code":null,"residence_street_code":null,"source":1,"start_time":0,"subject_type":0,"title":"","village":null,"visit_notify":false,"wg_number":""},{"avatar":"","birthday":null,"building":null,"come_from":"","company_id":1,"create_time":1632450256,"credential_no":"","credential_type":null,"department":"","description":"","domicile_address":"","domicile_city_code":null,"domicile_district_code":null,"domicile_province_code":null,"domicile_street_code":null,"education_code":null,"email":"","end_time":0,"entrance_people_type":null,"entry_date":null,"extra_id":null,"gender":1,"groups":[{"id":98,"name":"floor_15\u5458\u5de5\u7ec4"}],"house":null,"house_rel_code":1,"id":157,"interviewee":"","interviewee_pinyin":"","is_use":true,"job_number":"","marital_status_code":null,"name":"\u738b\u529b\u5b8f2","nation":null,"nationality_code":"","origin":"","password_reseted":false,"people_type":null,"phone":"18357036165","photos":[],"pinyin":"wanglihong2","purpose":0,"remark":"","residence_address":"","residence_city_code":null,"residence_district_code":null,"residence_province_code":null,"residence_street_code":null,"source":1,"start_time":0,"subject_type":0,"title":"","village":null,"visit_notify":false,"wg_number":""},{"avatar":"","birthday":null,"building":null,"come_from":"","company_id":1,"create_time":1632282560,"credential_no":"","credential_type":null,"department":"","description":"","domicile_address":"","domicile_city_code":null,"domicile_district_code":null,"domicile_province_code":null,"domicile_street_code":null,"education_code":null,"email":"","end_time":0,"entrance_people_type":null,"entry_date":null,"extra_id":null,"gender":1,"groups":[{"id":98,"name":"floor_15\u5458\u5de5\u7ec4"}],"house":null,"house_rel_code":1,"id":155,"interviewee":"","interviewee_pinyin":"","is_use":true,"job_number":"","marital_status_code":null,"name":"\u738b\u529b\u5b8f3","nation":null,"nationality_code":"","origin":"","password_reseted":false,"people_type":null,"phone":"18357036164","photos":[],"pinyin":"wanglihong3","purpose":0,"remark":"","residence_address":"","residence_city_code":null,"residence_district_code":null,"residence_province_code":null,"residence_street_code":null,"source":1,"start_time":0,"subject_type":0,"title":"","village":null,"visit_notify":false,"wg_number":""},{"avatar":"","birthday":null,"building":null,"come_from":"","company_id":1,"create_time":1630897914,"credential_no":"","credential_type":null,"department":"","description":"","domicile_address":"","domicile_city_code":null,"domicile_district_code":null,"domicile_province_code":null,"domicile_street_code":null,"education_code":null,"email":"","end_time":0,"entrance_people_type":null,"entry_date":null,"extra_id":null,"gender":1,"groups":[],"house":null,"house_rel_code":1,"id":151,"interviewee":"","interviewee_pinyin":"","is_use":true,"job_number":"","marital_status_code":null,"name":"\u738b\u529b\u5b8f2","nation":null,"nationality_code":"","origin":"","password_reseted":false,"people_type":null,"phone":"18357036165","photos":[{"company_id":1,"id":116,"origin_url":"/static/upload/origin/2021-09-06/v2_32dded04f9e166d2a07fce8738cad0e4bb8c895f.jpg","quality":0.99557,"subject_id":151,"url":"/static/upload/photo/2021-09-06/v2_40bd8bcaf0b94ac8947c7ad729435baa42426275.jpg","version":7}],"pinyin":"wanglihong2","purpose":0,"remark":"","residence_address":"","residence_city_code":null,"residence_district_code":null,"residence_province_code":null,"residence_street_code":null,"source":1,"start_time":0,"subject_type":0,"title":"","village":null,"visit_notify":false,"wg_number":""},{"avatar":"","birthday":null,"building":null,"come_from":"","company_id":1,"create_time":1630639100,"credential_no":"","credential_type":null,"department":"","description":"","domicile_address":"","domicile_city_code":null,"domicile_district_code":null,"domicile_province_code":null,"domicile_street_code":null,"education_code":null,"email":"","end_time":0,"entrance_people_type":null,"entry_date":null,"extra_id":null,"gender":1,"groups":[{"id":94,"name":"3F\u5458\u5de5\u7ec4"}],"house":null,"house_rel_code":1,"id":138,"interviewee":"","interviewee_pinyin":"","is_use":true,"job_number":"","marital_status_code":null,"name":"\u90ed\u674e\u6e05","nation":null,"nationality_code":"","origin":"","password_reseted":false,"people_type":null,"phone":"15074966932","photos":[],"pinyin":"guoliqing","purpose":0,"remark":"","residence_address":"","residence_city_code":null,"residence_district_code":null,"residence_province_code":null,"residence_street_code":null,"source":1,"start_time":0,"subject_type":0,"title":"","village":null,"visit_notify":false,"wg_number":""},{"avatar":"","birthday":null,"building":null,"come_from":"","company_id":1,"create_time":1630545694,"credential_no":"","credential_type":null,"department":"","description":"","domicile_address":"","domicile_city_code":null,"domicile_district_code":null,"domicile_province_code":null,"domicile_street_code":null,"education_code":null,"email":"","end_time":0,"entrance_people_type":null,"entry_date":null,"extra_id":null,"gender":1,"groups":[{"id":88,"name":"1\u697c\u5458\u5de5\u7ec4"}],"house":null,"house_rel_code":1,"id":119,"interviewee":"","interviewee_pinyin":"","is_use":true,"job_number":"","marital_status_code":null,"name":"\u5c0f\u90ed","nation":null,"nationality_code":"","origin":"","password_reseted":false,"people_type":null,"phone":"18170672948","photos":[],"pinyin":"xiaoguo","purpose":0,"remark":"","residence_address":"","residence_city_code":null,"residence_district_code":null,"residence_province_code":null,"residence_street_code":null,"source":1,"start_time":0,"subject_type":0,"title":"","village":null,"visit_notify":false,"wg_number":""},{"avatar":"","birthday":null,"building":null,"come_from":"","company_id":1,"create_time":1630545357,"credential_no":"","credential_type":null,"department":"","description":"","domicile_address":"","domicile_city_code":null,"domicile_district_code":null,"domicile_province_code":null,"domicile_street_code":null,"education_code":null,"email":"","end_time":0,"entrance_people_type":null,"entry_date":null,"extra_id":null,"gender":1,"groups":[{"id":88,"name":"1\u697c\u5458\u5de5\u7ec4"}],"house":null,"house_rel_code":1,"id":118,"interviewee":"","interviewee_pinyin":"","is_use":true,"job_number":"","marital_status_code":null,"name":"\u8d85\u7ea7","nation":null,"nationality_code":"","origin":"","password_reseted":false,"people_type":null,"phone":"13354985621","photos":[],"pinyin":"chaoji","purpose":0,"remark":"","residence_address":"","residence_city_code":null,"residence_district_code":null,"residence_province_code":null,"residence_street_code":null,"source":1,"start_time":0,"subject_type":0,"title":"","village":null,"visit_notify":false,"wg_number":""},{"avatar":"","birthday":null,"building":null,"come_from":"","company_id":1,"create_time":1629783823,"credential_no":"","credential_type":null,"department":"","description":"","domicile_address":"","domicile_city_code":null,"domicile_district_code":null,"domicile_province_code":null,"domicile_street_code":null,"education_code":null,"email":"","end_time":0,"entrance_people_type":null,"entry_date":null,"extra_id":null,"gender":1,"groups":[],"house":null,"house_rel_code":1,"id":88,"interviewee":"","interviewee_pinyin":"","is_use":true,"job_number":"","marital_status_code":null,"name":"\u738b\u529b\u5b8f2","nation":null,"nationality_code":"","origin":"","password_reseted":false,"people_type":null,"phone":"18357036164","photos":[],"pinyin":"wanglihong2","purpose":0,"remark":"","residence_address":"","residence_city_code":null,"residence_district_code":null,"residence_province_code":null,"residence_street_code":null,"source":1,"start_time":0,"subject_type":0,"title":"","village":null,"visit_notify":false,"wg_number":""},{"avatar":"","birthday":null,"building":null,"come_from":"","company_id":1,"create_time":1629702648,"credential_no":"","credential_type":null,"department":"","description":"","domicile_address":"","domicile_city_code":null,"domicile_district_code":null,"domicile_province_code":null,"domicile_street_code":null,"education_code":null,"email":"","end_time":0,"entrance_people_type":null,"entry_date":null,"extra_id":null,"gender":0,"groups":[],"house":null,"house_rel_code":1,"id":86,"interviewee":"","interviewee_pinyin":"","is_use":true,"job_number":"","marital_status_code":null,"name":"w'erwew'er","nation":null,"nationality_code":"","origin":"","password_reseted":false,"people_type":null,"phone":"","photos":[{"company_id":1,"id":106,"origin_url":"/static/upload/origin/2021-08-23/v2_a6caeb3580495f38b8695d1549ad3546a7ca5927.jpg","quality":0.999735,"subject_id":86,"url":"/static/upload/photo/2021-08-23/v2_7a706aeec1dd67da6e497dbd2acb078fd542d49d.jpg","version":7}],"pinyin":"w'erwew'er","purpose":0,"remark":"","residence_address":"","residence_city_code":null,"residence_district_code":null,"residence_province_code":null,"residence_street_code":null,"source":1,"start_time":0,"subject_type":0,"title":"","village":null,"visit_notify":false,"wg_number":""},{"avatar":"","birthday":null,"building":null,"come_from":"","company_id":1,"create_time":1629701357,"credential_no":"","credential_type":null,"department":"","description":"","domicile_address":"","domicile_city_code":null,"domicile_district_code":null,"domicile_province_code":null,"domicile_street_code":null,"education_code":null,"email":"","end_time":0,"entrance_people_type":null,"entry_date":null,"extra_id":null,"gender":0,"groups":[],"house":null,"house_rel_code":1,"id":85,"interviewee":"","interviewee_pinyin":"","is_use":true,"job_number":"","marital_status_code":null,"name":"fsdf","nation":null,"nationality_code":"","origin":"","password_reseted":false,"people_type":null,"phone":"","photos":[{"company_id":1,"id":105,"origin_url":"/static/upload/origin/2021-08-23/v2_de9899d2ae762d0f91033ada87f3e9b5f7f79b1b.jpg","quality":0.99804,"subject_id":85,"url":"/static/upload/photo/2021-08-23/v2_2bc87b2fc7041f41a797f2e42c29016f7f4faff0.jpg","version":7}],"pinyin":"fsdf","purpose":0,"remark":"","residence_address":"","residence_city_code":null,"residence_district_code":null,"residence_province_code":null,"residence_street_code":null,"source":1,"start_time":0,"subject_type":0,"title":"","village":null,"visit_notify":false,"wg_number":""}],"page":{"count":22,"current":1,"size":10,"total":3},"timecost":42}`

	//"page":{"count":22,"current":1,"size":10,"total":3}


	c.JSON(http.StatusOK, gin.H{
		"code":    0,
		"desc": "",
		"data":    temp_map,
		"page":pages,
	})

	//c.JSON(http.StatusOK, json_data)

}
func Copymap(value interface{}) interface{} {
	if valueMap, ok := value.(map[string]interface{}); ok {
		newMap := make(map[string]interface{})
		for k, v := range valueMap {
			newMap[k] = Copymap(v)
		}

		return newMap
	} else if valueSlice, ok := value.([]interface{}); ok {
		newSlice := make([]interface{}, len(valueSlice))
		for k, v := range valueSlice {
			newSlice[k] = Copymap(v)
		}

		return newSlice
	}

	return value
}

func Compare(c *gin.Context){

	_, header, err1 := c.Request.FormFile("image_1")
	if err1 != nil{
		log4go.Error("获取照片文件失败:",err1)
		c.JSON(http.StatusOK, gin.H{
			"code":    -100,
			"desc": "获取照片文件失败",
			"page":map[string]interface{}{},
		})
		return
	}


	_, header2, err1 := c.Request.FormFile("image_2")
	if err1 != nil{
		log4go.Error("获取照片文件失败:",err1)
		c.JSON(http.StatusOK, gin.H{
			"code":    -100,
			"desc": "获取照片文件失败",
			"page":map[string]interface{}{},
		})
		return
	}






	photourl,err := SaveUploadedFile(header,header.Filename)
	if err != nil{
		log4go.Error("获取照片文件失败:",err1)
		c.JSON(http.StatusOK, gin.H{
			"code":    -100,
			"desc": "获取照片文件失败",
			"page":map[string]interface{}{},
		})
		return
	}
	photourl2,err := SaveUploadedFile(header2,header2.Filename)
	if err != nil{
		log4go.Error("获取照片文件失败:",err1)
		c.JSON(http.StatusOK, gin.H{
			"code":    -100,
			"desc": "获取照片文件失败",
			"page":map[string]interface{}{},
		})
		return
	}










	result,score := Getscore(photourl,photourl2)
	json_return := `{"face_info_1":{"rect":{"left":87,"top":109,"right":262,"bottom":283},"quality":0.9955986086279154,"brightness":116.19357429718876,"std_deviation":28.79753933868924},"face_info_2":{"rect":{"left":87,"top":109,"right":262,"bottom":283},"quality":0.9955986086279154,"brightness":116.19357429718876,"std_deviation":28.79753933868924},"same":true,"score":98.16325378417969,"thresholds":{"E3":41.35199737548828,"E4":48.90372848510742,"E5":55.6849365234375,"E6":62.1713752746582,"recognizing":68,"stranger":67,"verify":68,"gate":78}}`
	json_data := make(map[string]interface{},0)
	json.Unmarshal([]byte(json_return),&json_data)
	//log4go.Info(json_data)
	//log4go.Info("result is",result)
	json_data["same"] = result
	json_data["score"] = score
	log4go.Info(json_data)
	c.JSON(http.StatusOK, json_data)

}


func CheckPhoto(photo multipart.File)(error){

	data_return,err := hongtu.AddPhoto(photo)
	if err != nil{
		log4go.Error("添加照片失败:",err)
		return err
	}
	log4go.Info(data_return)
	_,err = data_return.Get("data").Get("uri").String()
	if err != nil{
		log4go.Error("获取照片uri失败:",err)
		log4go.Error("添加照片失败:",err)
		return err
	}
	return nil
}



func Getscore(img1,img2 string)(bool,int){
	score := 0
	score_index := []byte("feautre score :")
	name := "/home/zybox/test/Baidu_Face_Offline_SDK_Linux_ARM_5.0/face-sdk-demo/armlinux/face_compare/run.sh"
	log4go.Info("name value is",name)
	log4go.Info("img1 value is",img1)
	log4go.Info("img2 value is",img2)
	cmd := exec.Command(name,img1,img2)
	output,err := cmd.Output()
	if err != nil{
		log4go.Error("get score err",err)
		//return false,0
	}
	log4go.Info("score_index value is",score_index)
	log4go.Info("output is",string(output))
	located := bytes.LastIndex(output,score_index)
	if located != -1{
		log4go.Info("located is",located)
		score_ := output[located+len(score_index):located+len(score_index)+2]
		log4go.Info("score_ value is",string(score_))
		score,err = strconv.Atoi(string(score_))
		if err != nil{
			log4go.Error("get score err",err)
		}
	}
	if score < 75{
		log4go.Info("score value is false",score)
		return false,score
	}else{
		log4go.Info("score value is true",score)
		return true,score
	}


}



func  SaveUploadedFile(file *multipart.FileHeader, dst string) (string,error) {

	PathExists3()
	uuid := client.UUID4()
	src, err := file.Open()
	if err != nil {
		return "",err
	}
	defer src.Close()
	//创建 dst 文件
	time_str := time.Unix(time.Now().Unix(), 0).Format("2006-01-02")
	out, err := os.Create("/home/zybox/photos/pub/"+time_str+"/"+uuid+"_"+dst)
	if err != nil {
		return "",err
	}
	defer out.Close()
	// 拷贝文件
	_, err = io.Copy(out, src)
	return "/home/zybox/photos/pub/"+time_str+"/"+uuid+"_"+dst,err
}

func SavePhoto(pathurl, dst string)(string,error){

	//log4go.Info("files")
	//log4go.Info(file)

	//presp, err := http.Get(pathurl)
	//if err != nil {
	//	log4go.Error(err.Error())
	//return -1, "", err
	//}
	//if presp.StatusCode != 200 {
	//return -1, "", errors.New("获取照片" + strconv.Itoa(presp.StatusCode) + "错误! 原因：" + presp.Status)
	//}
	//defer presp.Body.Close()

	message,_:=base64.StdEncoding.DecodeString(pathurl)
	//data,_ := base64.StdEncoding.DecodeString(pathurl)
	pix, err := ioutil.ReadAll(bytes.NewReader(message))
	//pix, err := ioutil.ReadAll(bytes.NewReader([]byte(pathurl)))
	if err != nil {
		log4go.Error(err.Error())
		//return -1, "", errors.New("读取图片出错!")
	}
	//创建 dst 文件
	time_str := time.Unix(time.Now().Unix(), 0).Format("2006-01-02")
	//path := "/home/zybox/photos/static/upload/photo/"+time_str+"/"
	out, err := os.Create("/home/zybox/photos/pub/"+time_str+"/"+dst)
	if err != nil {
		return "",err
	}
	defer out.Close()
	// 拷贝文件
	_, err = io.Copy(out, bytes.NewReader(pix))
	return "/pub/"+time_str+"/"+dst,err

}

func PathExists3() (bool) {
	time_str := time.Unix(time.Now().Unix(), 0).Format("2006-01-02")
	path := "/home/zybox/photos/pub/"+time_str+"/"
	_, err := os.Stat(path)
	if err == nil {
		return true
	}
	if os.IsNotExist(err) {
		os.MkdirAll(path,os.ModePerm)
		return PathExists3()
	}
	return true
}



func Recognize(c *gin.Context){


	//json_return := `{"face_info_1":{"rect":{"left":87,"top":109,"right":262,"bottom":283},"quality":0.9955986086279154,"brightness":116.19357429718876,"std_deviation":28.79753933868924},"face_info_2":{"rect":{"left":87,"top":109,"right":262,"bottom":283},"quality":0.9955986086279154,"brightness":116.19357429718876,"std_deviation":28.79753933868924},"same":true,"score":98.16325378417969,"thresholds":{"E3":41.35199737548828,"E4":48.90372848510742,"E5":55.6849365234375,"E6":62.1713752746582,"recognizing":68,"stranger":67,"verify":68,"gate":78}}`
	//json_return := `{"candidates":[{"id":64,"subject_id":52,"photo_id":64,"confidence":33.55346},{"id":106,"subject_id":86,"photo_id":106,"confidence":21.79743},{"id":15,"subject_id":5,"photo_id":15,"confidence":20.482822}],"face_info":{"rect":{"left":146,"top":180,"right":335,"bottom":369},"quality":0.9960721684619784,"brightness":141.17540181691126,"std_deviation":23.86713417529236},"person":{"id":64,"subject_id":52,"photo_id":64,"confidence":33.55346},"recognized":false,"thresholds":{"E3":41.35199737548828,"E4":48.90372848510742,"E5":55.6849365234375,"E6":62.1713752746582,"recognizing":68,"stranger":67,"verify":68,"gate":78}}`
	json_return :=`{"can_door_open":false,"error":7}`
	json_data := make(map[string]interface{},0)
	json.Unmarshal([]byte(json_return),&json_data)
	log4go.Info(json_data)
	c.JSON(http.StatusOK, json_data)

}


func Subjectsgrouplist(c *gin.Context){
	resp4Device := model.Resp4Device{Code: 0, Err_msg: ""}
	data := make([]map[string]interface{},0)
	//json_return := `{"code":0,"data":[{"comment":"floor_15\u8bbf\u5ba2\u7ec4","id":99,"name":"floor_15\u8bbf\u5ba2\u7ec4","subject_count":1,"subject_type":1,"update_by":"admin@91zo.com","update_time":1632282494.0},{"comment":"2\u697c\u8bbf\u5ba2\u7ec4","id":97,"name":"2\u697c\u8bbf\u5ba2\u7ec4","subject_count":0,"subject_type":1,"update_by":"admin@91zo.com","update_time":1630893529.0},{"comment":"3F\u8bbf\u5ba2\u7ec4","id":95,"name":"3F\u8bbf\u5ba2\u7ec4","subject_count":0,"subject_type":1,"update_by":"admin@91zo.com","update_time":1630639023.0},{"comment":"1\u697c\u8bbf\u5ba2\u7ec4","id":89,"name":"1\u697c\u8bbf\u5ba2\u7ec4","subject_count":0,"subject_type":1,"update_by":"admin@91zo.com","update_time":1630544330.0},{"comment":"floor_20\u8bbf\u5ba2\u7ec4","id":17,"name":"floor_20\u8bbf\u5ba2\u7ec4","subject_count":0,"subject_type":1,"update_by":"admin@91zo.com","update_time":1629269644.0}],"page":{"count":5,"current":1,"size":10,"total":1},"timecost":55}`


	//json_data := make(map[string]interface{},0)
	//json.Unmarshal([]byte(json_return),&json_data)

	//log4go.Info(json_data["data"].([]interface{})[0])
	//go func() {
	GetGroupsHongtu()
	//}()

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


	//c.JSON(http.StatusOK, json_data)
	c.JSON(http.StatusOK, gin.H{
		"code":    0,
		"desc": "",
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
	//log4go.Info(c.Request)
	//body, _ := ioutil.ReadAll(c.Request.Body)
	//m := make(map[string]interface{}, 0)
	//json.Unmarshal(body, &m)
	//log4go.Info("m", m)

	json_return := `{"company_id":1,"id":123,"origin_url":"/static/upload/origin/2021-09-26/v2_4b4b6a8f681d68254427e80506426b6aaf68a86c.jpg","quality":null,"subject_id":null,"url":"/static/upload/photo/2021-09-26/v2_a94965c0110044ae16350644f047b517e8f13162.jpg","version":7}`

	json_data := make(map[string]interface{},0)
	json.Unmarshal([]byte(json_return),&json_data)
	_, header, err1 := c.Request.FormFile("photo")
	if err1 != nil{
		log4go.Error("获取照片文件失败:",err1)
		c.JSON(http.StatusOK, gin.H{
			"code":    0,
			"desc": "获取照片文件失败",
			"page":map[string]interface{}{},
		})
		return
	}
	subject_id := c.Request.FormValue("subject_id")
	log4go.Info("get subject_id",subject_id)
	if subject_id == ""{
		photo_file,err := header.Open()
		err = CheckPhoto(photo_file)
		if err != nil{
			log4go.Error("获取照片文件失败:",err1)
			c.JSON(http.StatusOK, gin.H{
				"code":    -100,
				"desc": err.Error(),
				"page":map[string]interface{}{},
			})
			return
		}
		c.JSON(http.StatusOK, gin.H{
			"code":    0,
			"desc": "",
			"data":    json_data,
			"page":map[string]interface{}{},
		})
		return
	}


	photo,err :=header.Open()
	if err != nil{
		log4go.Error("获取照片文件失败:",err)
		c.JSON(http.StatusOK, gin.H{
			"code":    0,
			"desc": "获取照片文件失败",
			"page":map[string]interface{}{},
		})
		return
	}
	data_return,err := hongtu.AddPhoto(photo)
	if err != nil{
		log4go.Error("添加照片失败:",err)
		c.JSON(http.StatusOK, gin.H{
			"code":    -100,
			"desc": "添加照片失败",
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
			"desc": "获取照片uri失败",
			"page":map[string]interface{}{},
		})
		return
	}

	//photoship ,err := server.AddPhotoShip(uri)
	//if err != nil{
	//	log4go.Error("存储uri失败:",err)
	//	c.JSON(http.StatusOK, gin.H{
	//		"code":    -100,
	//		"desc": "存储uri失败",
	//		"page":map[string]interface{}{},
	//	})
	//	return
	//}
	subject_id_int , _ := strconv.Atoi(subject_id)
	subject,err := dao.FindSubjectShip2("","",subject_id_int)
	if err != nil{
			log4go.Error("查询人员失败:",err)
			c.JSON(http.StatusOK, gin.H{
				"code":    -100,
				"desc": "查询人员失败",
				"page":map[string]interface{}{},
			})
			return
	}

	updateparams := map[string]interface{}{
		"imageUri":uri,
		"uuid":subject.Uuid,
	}
	log4go.Info("updateparams",updateparams)
	_,err = hongtu.Personupdate(updateparams)
	if err != nil{
		log4go.Error("更新人员失败:",err)
		c.JSON(http.StatusOK, gin.H{
			"code":    -100,
			"desc": "更新人员失败",
			"page":map[string]interface{}{},
		})
		return
	}

	err = dao.DeleteSubjectShip(subject.Uuid)
	if err != nil{
		log4go.Error("delete subject err",err)
	}

	//json_data["id"] = photoship.ID
	//log4go.Info("json_data is",json_data)
	//c.JSON(http.StatusOK, json_data)
	c.JSON(http.StatusOK, gin.H{
		"code":    0,
		"desc": "",
		"data":    json_data,
		"page":map[string]interface{}{},
	})
}

func DeletePerson(c *gin.Context){
	delete_id,_ := c.Params.Get("id")
	if delete_id == ""{
		c.JSON(http.StatusOK, gin.H{
			"code":    -100,
			"desc": "id错误",
			"data":    map[string]interface{}{},
			"page":map[string]interface{}{},
		})
		return
	}
	id,_ := strconv.Atoi(delete_id)
	subject,err := dao.FindSubjectShip2("","",id)
	if err != nil{
		c.JSON(http.StatusOK, gin.H{
			"code":    -100,
			"desc": "本地查询错误",
			"data":    map[string]interface{}{},
			"page":map[string]interface{}{},
		})
		return
	}

	uuidlist := []string{subject.Uuid}
	params := map[string]interface{}{
		"uuidList":uuidlist,
	}
	_,err = hongtu.DeleteSubject(params)
	if err != nil{
		c.JSON(http.StatusOK, gin.H{
			"code":    -100,
			"desc": "鸿图删除访客失败",
			"data":    map[string]interface{}{},
			"page":map[string]interface{}{},
		})
		return
	}
	dao.DeleteShip2(id)
	c.JSON(http.StatusOK, gin.H{
		"code":    0,
		"desc": "",
		"data":    map[string]interface{}{"id":id},
		"page":map[string]interface{}{},
	})
}