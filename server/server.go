package server

import (
	"bufio"
	"bytes"
	"encoding/base64"
	"encoding/json"
	"errors"
	"fmt"
	"github.com/alecthomas/log4go"
	"github.com/bitly/go-simplejson"
	"github.com/satori/go.uuid"
	"io"
	"io/ioutil"
	"net/http"
	"os"
	"reflect"
	"strconv"
	"strings"
	"time"
	"zhiyuan/koala2hongtuadapt/dao"
	koala "zhiyuan/koala_api_go/koala_api"
	"zhiyuan/zyutil/config"
	//"zhiyuan/koala2hongtuadapt/koala"
	"zhiyuan/koala2hongtuadapt/model"
)

func GetPurposeMapList() []map[string]interface{} {
	Purpose_map_list := make([]map[string]interface{}, 0)
	label := []string{"其它", "面试", "商务", "亲友", "快递送货"}

	for k, v := range label {
		Purpose_map := map[string]interface{}{
			"label": v,
			"value": k,
		}
		Purpose_map_list = append(Purpose_map_list, Purpose_map)
	}

	return Purpose_map_list
}

func GetPersonGroupList(subject_type, page, size int) ([]map[string]interface{}, error) {

	//xxxx:=make([]map[string]interface{},0)
	//koala.Init("192.168.18.51",80)
	//koala.KoalaLogin(config.Gconf.KoalaUsername,config.Gconf.KoalaPassword)
	resp, err := koala.GetPersonGroupList(subject_type, page, size)
	if err != nil {
		return nil, err
	}
	return resp, err
}

func GetStaffs(name string, page int, size int) ([]map[string]interface{}, *simplejson.Json, error) {

	//koala.KoalaLogin(config.Gconf.KoalaUsername,config.Gconf.KoalaPassword)
	data, staff_page, err := koala.GetStaffs2(name, page, size)
	if err != nil {
		return nil, nil, err
	}
	//for _,v1:= range staff_page {
	//	fmt.Println(v1)
	//}
	temp_map_list := make([]map[string]interface{}, 0)
	for _, v := range data {
		temp_map := make(map[string]interface{})
		temp_map["name"] = v["name"]
		temp_map["subject_id"] = v["id"]
		temp_map["avatar"] = v["avatar"]
		temp_map["department"] = v["department"]
		temp_map_list = append(temp_map_list, temp_map)
	}

	return temp_map_list, staff_page, err
}

func DeleteSubject(subject_id int, proposer string) (err error) {

	err = dao.Delete(subject_id, proposer)
	if err != nil {
		return err
	}
	people, err, _, _ := dao.Get_Retinue(subject_id, 1, 1000)
	if err != nil {
		return err
	}
	for _, v := range people {
		err = koala.DeleteSubject(v.Subject_id)
		if err != nil {
			continue
		}
	}
	err = dao.Delete_Retinue(subject_id)
	if err != nil {
		return err
	}
	return nil
}

func GetSubject(proposer string, subject_id, page, size int) (result []model.Visitor, result_data map[string]interface{}, err error, count int, total int) {
	//result := make([]model.Visitor,0)
	//fmt.Println("------------------------------------3")
	//fmt.Println(proposer,subject_id,page,size)
	//fmt.Println("------------------------------------3")
	result, err, count, total = dao.Get(proposer, subject_id, page, size)
	fmt.Println(result)
	if err != nil {
		return result, result_data, err, 0, 0
	}
	for k, v := range result {
		photoUrl := "http://" + config.Gconf.KoalaHost + ":80" + v.Photo
		photobs, _ := photoUrltobase64(photoUrl, "./photo/")
		//fmt.Println(photobs)
		result[k].Photo = photobs
	}
	if len(result) == 1 && subject_id != 0 {
		Retinues, _, _, _ := dao.GetRetinues(subject_id, 1, 10000)
		for ki, vi := range Retinues {
			photoUrl := "http://" + config.Gconf.KoalaHost + ":80" + vi.Photo
			photobs, _ := photoUrltobase64(photoUrl, "./photo/")
			//fmt.Println(photobs)
			Retinues[ki].Photo = photobs
		}
		if len(result) == 1 {
			purpose := map[int]interface{}{
				0: "其它",
				1: "面试",
				2: "商务",
				3: "亲友",
				4: "快递送货",
			}
			for k, v1 := range purpose {
				if result[0].Purpose == k {
					result[0].Purpose_name = v1.(string)
				}
			}
			resp, _ := GetPersonGroupList(1, 1, 1000)
			for _, v := range resp {
				data, _ := strconv.Atoi(v["id"].(json.Number).String())
				purpose := map[int]interface{}{
					data: v["name"],
				}
				for k, v1 := range purpose {
					if result[0].Group_ids == k {
						result[0].Group_Name = v1.(string)
					}
				}
			}
		}
		resultmap := Struct2Map(result[0])

		resultmap["retinue"] = Retinues
		resultmap["purpose"] = resultmap["purpose_name"]
		resultmap["group_ids"] = resultmap["group_name"]

		return nil, resultmap, nil, count, total
	}
	if len(result) == 0 {
		//return result, errors.New("未找到对应访客信息!"),0,0
		return result, result_data, nil, 0, 0
	}
	//result_data = result
	return result, nil, nil, count, total
}
func Struct2Map(obj interface{}) map[string]interface{} {
	t := reflect.TypeOf(obj)
	v := reflect.ValueOf(obj)

	var data = make(map[string]interface{})
	for i := 0; i < t.NumField(); i++ {
		data[strings.ToLower(t.Field(i).Name)] = v.Field(i).Interface()
	}
	return data
}

/**
根据姓名和手机号码，查询员工信息
*/
func StaffsBinding(name, phone string, unionid int) (map[string]interface{}, error) {

	//koala.KoalaLogin(config.Gconf.KoalaUsername,config.Gconf.KoalaPassword)
	//data,err := koala.BindingStaffs(phone,unionid)
	// 根据姓名找到员工
	//data, _, err := koala.GetStaffs("", 1, 10000)
    _, data, err := koala.GetSubjectsByCondition("employee", name, "", 1, 10000, "")

	if err != nil {
		return nil, err
	}
    if len(data)==0 { // 没有匹配的员工
		return nil, errors.New("绑定失败，失败原因：未找到对应手机号!")
    }

	temp_map_list := make([]map[string]interface{}, 0)
	for _, v := range data {
		if v["phone"] == "" || v["phone"] != phone {
			continue
		}
		//fmt.Println(v["phone"])
		temp_map := make(map[string]interface{})
		temp_map["name"] = v["name"]
		temp_map["subject_id"] = v["id"]
		temp_map["department"] = v["department"]

		if len(v["avatar"].(string)) != 0 {
			url := "http://" + config.Gconf.KoalaHost + ":80" + v["avatar"].(string)
			photoUrl, _ := photoUrltobase64(url, "./photo/")
			temp_map["avatar"] = photoUrl
			temp_map_list = append(temp_map_list, temp_map)
		} else {
			if len(v["photos"].([]interface{})) == 0 {
				temp_map["avatar"] = ""
				temp_map_list = append(temp_map_list, temp_map)
			} else {
				url := "http://" + config.Gconf.KoalaHost + ":80" + v["photos"].([]interface{})[0].(map[string]interface{})["url"].(string)
				photoUrl, _ := photoUrltobase64(url, "./photo/")
				temp_map["avatar"] = photoUrl
				temp_map_list = append(temp_map_list, temp_map)
			}
		}
	}
	//fmt.Println(temp_map_list)
	if len(temp_map_list) == 0 {
		return nil, errors.New("绑定失败，失败原因：未找到对应手机号!")
	}

	return temp_map_list[0], err
}

func photoUrltobase64(photoUrl string, base_path string) (string, error) {
	log4go.Info(photoUrl)
	presp, err := http.Get(photoUrl)
	if err != nil {
		log4go.Error(err.Error())
		return "", errors.New("2.0api-DealModifyPersons: 访问照片地址失败！")
	}
	if presp.StatusCode != 200 {
		return "", errors.New("获取照片" + strconv.Itoa(presp.StatusCode) + "错误! 原因：" + presp.Status)
	}
	defer presp.Body.Close()
	pix, err := ioutil.ReadAll(presp.Body)
	if err != nil {
		log4go.Error(err.Error())
		return "", errors.New("读取图片出错!")
	}
	//照片创建
	photoName, _ := uuid.NewV4()
	img_name := base_path + photoName.String() + ".jpg"
	log4go.Info(img_name)
	out, err := os.Create(img_name)
	log4go.Info(out)
	_, err2 := io.Copy(out, bytes.NewReader(pix))
	if err2 != nil {
		log4go.Error(err2.Error())
		return "", errors.New("下载图片出错!")
	}
	out.Close()
	imgFile, err := os.Open(img_name)
	fInfo, _ := imgFile.Stat()
	var size int64 = fInfo.Size()
	bufstore := make([]byte, size)
	fReader := bufio.NewReader(imgFile)
	fReader.Read(bufstore)
	photobase64 := "data:image/jpeg;base64," + base64.StdEncoding.EncodeToString(bufstore)
	//fmt.Println(photobase64)
	return photobase64, err
}

func DeleteNextTime() (err error) {
	//koala.Init(config.Gconf.KoalaHost,config.Gconf.KoalaPort)
	//koala.KoalaLogin(config.Gconf.KoalaUsername, config.Gconf.KoalaPassword)

	timenow := int(time.Now().Unix())
	//fmt.Println(timenow)
	result, err := dao.Gettime("", 0, timenow)
	if err != nil {
		log4go.Error("Visitor get db err:", err.Error())
	}
	for _, v := range result {
		err = koala.ModSubject(v.Subject_id, 1)
		if err != nil {
			log4go.Error("Visitor mod photo_ids err:", err.Error())
		}
	}
	return err
}


func AddSubjectShip(uuid string)(model.SubjectShip,error){
	data := model.SubjectShip{Uuid: uuid}
	SubjectShip,err := dao.CreateSubjectShip((*model.SubjectShip)(&data))
	if err != nil{
		return model.SubjectShip{}, err
	}
	return SubjectShip,nil
}

func AddPhotoShip(uri string)(model.PhotoShip,error){
	data := model.PhotoShip{Uri: uri}
	PhotoShip,err := dao.CreatePhotoShip(&data)
	if err != nil{
		return model.PhotoShip{}, err
	}
	return PhotoShip,nil
}
