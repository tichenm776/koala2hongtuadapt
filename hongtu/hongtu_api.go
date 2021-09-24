package hongtu

import (
	"bytes"
	"encoding/base64"
	client "go-common/app/service/main/vip/dao/ele-api-client"
	"strings"
	"zhiyuan/koala2hongtuadapt/dao"
	"zhiyuan/koala2hongtuadapt/model"
	//"zhiyuan/zyutil"

	//client "go-common/app/service/main/vip/dao/ele-api-client"
	"io"
	"io/ioutil"
	"net/http"
	"time"
	//"bytes"
	"errors"
	"os"
	//"fmt"
	"github.com/alecthomas/log4go"
	"github.com/bitly/go-simplejson"
	//"io"
	//"mime/multipart"
	//"net/http"
	//"net/http/cookiejar"
	//"net/url"
	//"strconv"
	//"strings"
	hongtu "zhiyuan/koala_api_go/hongtu_api"
)


func PathExists_file(path string) (bool, error) {
	_, err := os.Stat(path)
	if err == nil {
		return true, nil
	}
	if os.IsNotExist(err) {
		return false, nil
	}
	return false, err
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




func SavePhoto4(pathurl, dst string)(string,error){

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

func SavePhoto3(pathurl, dst string)(string,error){

	//log4go.Info("files")
	//log4go.Info(file)

	presp, err := http.Get(pathurl)
	if err != nil {
		log4go.Error(err.Error())
		//return -1, "", err
	}
	if presp.StatusCode != 200 {
		//return -1, "", errors.New("获取照片" + strconv.Itoa(presp.StatusCode) + "错误! 原因：" + presp.Status)
	}
	defer presp.Body.Close()
	pix, err := ioutil.ReadAll(presp.Body)
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




func GetStaffsByNameAndPhone(name,phone string) (map[string]interface{},error) {
	staffs, _, err := GetStaffs2(phone,name, 1, 1000)
	if err != nil {
		return nil, err
	}
	for _, data := range staffs {
		if data["phone"].(string) == phone{
			avatar := ""
			staff := make(map[string]interface{})
			//保存人员照片到本地
			isexists := PathExists3()
			if isexists{
				filename := data["name"].(string)
				if value,ok := data["imageUri"].(string);ok{
					log4go.Info("GET URL",value)
					avatar,err = SavePhoto3(value,filename+".jpg")
					if err != nil{
						log4go.Error("照片存储失败",err)
					}
				}

			}
			uuid := data["uuid"].(string)
			staff["phone"] = data["phone"]
			staff["name"] = data["name"]
			staff["uuid"] = uuid
			staff["avatar"] = avatar
			staff["department"] = data["postion"]

			ship,err := dao.GetSubjectShipUuid(uuid)
			if err != nil{
				if strings.Index(err.Error(),"not found") != -1{
					data := &model.SubjectShip{Uuid: uuid}
					ship,err := dao.CreateSubjectShip(data)
					if err != nil{
						log4go.Error("create subject ship err",err)
						return nil, err
					}
					log4go.Info("create ShipUuid",ship.ID)
					staff["subject_id"] = ship.ID
				}else{
					log4go.Error("get subject ship err",err)
					return nil, err
				}
			}
			if staff["subject_id"] == nil{
				staff["subject_id"] = ship.ID
			}

			return staff,nil
		}
	}
	return nil,errors.New("没有找到员工数据")
}

func Init(host string, koalaProt int) error {
	err := hongtu.Init(host,koalaProt)
	if err !=nil{
		return err
	}
	return nil
}

func KoalaLogin( username string, password string) error {
	return hongtu.KoalaLogin(username,password)
}


func GetStaffs2(phone,name string, page int, size int) ([]map[string]interface{}, *simplejson.Json, error) {

	data := make([]map[string]interface{},0)

	params := map[string]interface{}{
		"name":name,
		"phone":phone,
		"pageNum":page,
		"pageSize":size,
	}
	log4go.Info("params is",params)
	resp_data,err := hongtu.PersonList(params)
	if err != nil{
		log4go.Error("get person list err",err)
	}
	log4go.Info("1111111111111111111111111",resp_data)
	resp_page := simplejson.New()
	resp_page.Set("current",page)
	resp_page.Set("size",size)

	if value,err := resp_data.Get("data").Get("total").Int();err == nil{
		resp_page.Set("total",value)
	}else{
		resp_page.Set("total",0)
	}
	if value,err := resp_data.Get("data").Get("list").Array();err == nil{
		resp_page.Set("count",len(value))
		for _,v := range value{
			if value2,ok := v.(map[string]interface{});ok{
				//
				//value2["phone"] =
				//value2["name"] =
				//value2["subject_id"] =
				//value2["avatar"] =
				//value2["department"] =

				data = append(data, value2)
			}
		}

	}else{
		resp_page.Set("count",0)
	}
	log4go.Info("22222222222222222222222222222222",data)
	log4go.Info("333333333333333333333333333333",resp_page)


	return data, resp_page, nil
}

func GetSubjectId(id int)(map[string]interface{}, error){
	SubjectShip,err := dao.GetSubjectShip(id)
	if err != nil{
		return map[string]interface{}{},err
	}
	params := map[string]interface{}{
		"uuid":SubjectShip.Uuid,
	}
	resp,err := hongtu.PersonQuery(params)
	if err != nil{
		return map[string]interface{}{},err
	}
	staff := make(map[string]interface{})
	staff["phone"],_ = resp.Get("data").Get("phone").String()
	staff["name"],_ = resp.Get("data").Get("name").String()
	staff["subject_id"] = SubjectShip.ID
	staff["avatar"],_ = resp.Get("data").Get("imageUri").String()
	staff["department"],_ = resp.Get("data").Get("postion").String()

	return staff,nil
}

func GetStaffBySubjectId(id int) (map[string]interface{}, error){

	subject,err:=GetSubjectId(id)
	if err!=nil{
		return nil,err
	}
	return subject,err
}
//
//func AddVisitor(params interface{})(*simplejson.Json,error){
//	//name	string	必须		用户名, 长度:[1,40]
//	//type	integer	必须		用户类别, 1 员工 2访客 3 重点人员
//	// imageUri	string	非必须		用户识别照片的uri
//	// identifyNum	string	非必须		身份证号, 格式:允许大小写英文字母,数字, 长度:[1,32]
//	// visitFirm	string	非必须		访客所属单位, 访客可填, 格式:汉字,大小写英文字母,数字, 长度: [1,40]
//	//visitStartTimeStamp	integer	非必须		拜访起始时间(时间戳, 毫秒), 访客必填
//	//visitEndTimeStamp	integer	非必须		拜访结束时间(时间戳, 毫秒), 访客必填
//	//visitReason	string	非必须		拜访原因, 访客可填,格式: 任意字符, 长度: [1, 255]
//	// visitedUuid	string	非必须		受访人的uuid, 访客必填
//	// visitType	integer	非必须		访客类型, 访客可填, 1 普通访客, 2 VIP
//	// phone	string	非必须		手机号, 格式:数字, 长度:[6,18]
//	temp_map := make(map[string]interface{},0)
//
//
//	temp_map["name"] = params.(map[string]interface{})["name"]
//	temp_map["type"] = params.(map[string]interface{})["subject_type"]
//	temp_map["imageUri"] = params.(map[string]interface{})["imageUri"]
//	temp_map["identifyNum"] = params.(map[string]interface{})["imageUri"]
//
//
//
//
//
//
//
//
//
//	temp_map["subject_type"] = visitor.Subject_type
//	temp_map["name"] = visitor.Name
//	temp_map["phone"] = visitor.Phone
//	temp_map["come_from"] = visitor.Come_from
//	temp_map["remark"] = visitor.Remark
//	temp_map["photo_ids"] = []int{visitor.Photo_id}
//	temp_map["start_time"] = visitor.Start_time
//	temp_map["end_time"] = visitor.End_time
//	temp_map["purpose"] = visitor.Purpose
//	temp_map["interviewee"] = visitor.Interviewee
//	//temp_map["categories"]="visitor"
//	//temp_map["visitor_type"]=1
//	log4go.Debug("add koala value", temp_map)
//	if visitor.Group_ids != -1 {
//		temp_map["group_ids"] = []int{visitor.Group_ids}
//	} else {
//		temp_map["group_ids"] = []int{}
//	}
//
//
//
//	hongtu.AddPerson(params)
//
//}
//
func Visitor2Koala(visitor model.Visitor) (model.Visitor, error) {
	var (
		resp4Device model.Resp4Device
	)
	temp_map := make(map[string]interface{})
	resp4Device.Code = 0
	resp4Device.Err_msg = ""

	//defer presp.Body.Close()
	//pix, err := ioutil.ReadAll(presp.Body)
	//if err != nil {
	//	log4go.Error(err.Error())
	//	return errors.New("读取图片出错!")
	//}
	//照片创建
	//photoName, _ := uuid.NewV4()
	//img_name := "./photo" + "/" + photoName.String() + ".jpg"
	//log4go.Info(img_name)
	//out, err := os.Create(img_name)
	//log4go.Info(out)
	//defer out.Close()
	//_, err2 := io.Copy(out, bytes.NewReader([]byte(visitor.Photo)))
	//if err2 != nil {
	//	resp4Device.Code = -200
	//	resp4Device.Err_msg = "photo解析出错"+err.Error()
	//	return 0,errors.New(resp4Device.Err_msg)
	//}
	//defer out.Close()
	//photo_path,err := SavePhoto4(visitor.Photo,visitor.Phone+".jpg")
	//photo_path, err := zyutil.Base64_to_img(visitor.Photo, "./photo/")
	//
	//if err != nil {
	//	resp4Device.Code = -200
	//	resp4Device.Err_msg = "photo解析出错" + err.Error()
	//	return model.Visitor{}, errors.New(resp4Device.Err_msg)
	//}
	photofile, err := os.Open("/home/zybox/photos" + visitor.Photo)
	if err != nil {
		log4go.Error(err)
		resp4Device.Code = -200
		resp4Device.Err_msg = "photo读取出错"
		return model.Visitor{}, errors.New(resp4Device.Err_msg)
	}
	//visitor.Photo = photo_path

	photoresp, err := hongtu.AddPhoto(photofile)
	if err != nil {
		log4go.Error("上传底库失败,失败原因：" + err.Error())
		resp4Device.Code = -200
		resp4Device.Err_msg = "上传底库失败,失败原因：" + err.Error()
		return model.Visitor{}, errors.New(resp4Device.Err_msg)
	}
	uri,err := photoresp.Get("data").Get("uri").String()
	if err != nil{
		log4go.Error("获取人脸图片参数失败：" + err.Error())
		resp4Device.Code = -200
		resp4Device.Err_msg = "获取人脸图片参数失败：" + err.Error()
		return model.Visitor{}, errors.New(resp4Device.Err_msg)
	}
	employee,err := GetStaffsByNameAndPhone(visitor.Interviewee,visitor.Interviewee_phone)
	if err != nil{
		log4go.Error("系统内没有该受访人员：" + err.Error())
		resp4Device.Code = -200
		resp4Device.Err_msg = "系统内没有该受访人员：" + err.Error()
		return model.Visitor{}, errors.New(resp4Device.Err_msg)
	}

	//log4go.Debug(photo_url)
	uuid := client.UUID4()
	visitor.Photo_id = 0
	temp_map["type"] = 2
	temp_map["name"] = visitor.Name
	temp_map["uuid"] = uuid
	temp_map["phone"] = visitor.Phone
	temp_map["imageUri"] = uri
	temp_map["identifyNum"] = visitor.Id_no
	temp_map["visitFirm"] = visitor.Come_from
	temp_map["visitStartTimeStamp"] = visitor.Start_time*1000
	temp_map["visitEndTimeStamp"] = visitor.End_time*1000
	temp_map["visitReason"] = visitor.Purpose_name
	temp_map["visitedUuid"] = employee["uuid"]
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
		return model.Visitor{}, errors.New(resp4Device.Err_msg)
	}
	//res_map, err := res.Map()
	//if err != nil {
	//	log4go.Error(err)
	//	resp4Device.Code = -200
	//	resp4Device.Err_msg = err.Error()
	//	return 0, errors.New(resp4Device.Err_msg)
	//}

	ship,err := dao.GetSubjectShipUuid(uuid)
	if err != nil{
		if strings.Index(err.Error(),"not found") != -1{
			data := &model.SubjectShip{Uuid: uuid}
			ship,err := dao.CreateSubjectShip(data)
			if err != nil{
				log4go.Error("create subject ship err",err)
				return model.Visitor{}, err
			}
			visitor.Subject_id = ship.ID
		}else{
			log4go.Error("get subject ship err",err)
			return model.Visitor{}, err
		}
	}
	visitor.Subject_id = ship.ID
	data := map[string]interface{}{
		"uuid":uuid,
	}
	qrcode,err := hongtu.VisitorCode(data)
	if err != nil {
		log4go.Error("get visitor code err",err)
		resp4Device.Code = -100
		resp4Device.Err_msg = temp_map["name"].(string) + "获取二维码失败"
		log4go.Error(resp4Device.Err_msg)
		return model.Visitor{}, errors.New(resp4Device.Err_msg)
	}
	if value,err := qrcode.Get("data").Get("qrCode").String();err == nil{
		log4go.Debug("qrcode message is",value)
		//filename,err := zyutil.Base64_to_img(value,"./photo/")
		sub := len("data:image/png;base64,")
		PathExists3()
		filename,err := SavePhoto4(value[sub:],visitor.Phone+"_qr.png")
		if err != nil{
			log4go.Error("base64toimg err",err)
			resp4Device.Code = -100
			resp4Device.Err_msg = "保存二维码失败"
			log4go.Error(resp4Device.Err_msg)
			return model.Visitor{}, errors.New(resp4Device.Err_msg)
		}
		photofile, err := ioutil.ReadFile("/home/zybox/photos" + filename)
		if err != nil {
			resp4Device.Code = -200
			resp4Device.Err_msg = "photo读取出错"
			return model.Visitor{}, errors.New(resp4Device.Err_msg)
		}
		log4go.Debug(photofile)
		//visitor.Qrcode_img_url = base64.StdEncoding.EncodeToString(photofile)
		visitor.Qrcode_img_url = filename
		visitor.Qrcode = value
	}

	//visitor.Photo = photo_path
	//if len(res_map["photos"].([]interface{})) == 0 {
	//	resp4Device.Code = -200
	//	resp4Device.Err_msg = "缺少照片"
	//	fmt.Println(res_map["id"].(json.Number).String())
	//	id, _ := strconv.Atoi(res_map["id"].(json.Number).String())
	//	hongtu.DeleteSubject(id)
	//	return 0, errors.New(resp4Device.Err_msg)
	//}
	//subject_id, _ := res_map["id"].(json.Number).Int()
	return visitor, nil
}

func GetEmployeeList(name string)([]map[string]interface{},error){

	temp_map := make(map[string]interface{},0)
	person_map_list := make([]map[string]interface{},0)
	var (
		err  error
	)
	temp_map["type"] = 1
	if name != ""{
		temp_map["name"] = name
	}

	temp_map["pageNum"] = 1
	temp_map["pageSize"] = 10000

	personlist, err := hongtu.PersonList(temp_map)
	if err != nil {
		log4go.Error(temp_map["name"].(string) + "查询失败")
		return nil,err
	}

	if value,err := personlist.Get("code").Int();err == nil{
		if value == 0{
			if array,err := personlist.Get("data").Array();err == nil{
				if len(array) > 0 {
					for k,_ := range array{
						if value,ok := array[k].(map[string]interface{});ok{
							person_map_list = append(person_map_list, value)
						}
					}
				}
			}
		}
	}

	return person_map_list,nil

}

func GetGroupsHongtuList(params interface{})([]map[string]interface{},error){

	person_map_list := make([]map[string]interface{},0)
	var (
		err  error
	)

	personlist, err := hongtu.Passgrouplist(params)
	if err != nil {
		log4go.Error("查询失败",err)
		return nil,err
	}

	if value,err := personlist.Get("code").Int();err == nil{
		if value == 0{
			if array,err := personlist.Get("data").Array();err == nil{
				if len(array) > 0 {
					for k,_ := range array{
						if value,ok := array[k].(map[string]interface{});ok{
							person_map_list = append(person_map_list, value)
						}
					}
				}
			}
		}
	}

	return person_map_list,nil

}
