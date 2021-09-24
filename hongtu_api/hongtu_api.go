package hongtu_api

import (
	"bytes"
	"os"

	//"bytes"
	"crypto/md5"
	"encoding/hex"
	"encoding/json"
	"errors"
	"fmt"
	"github.com/alecthomas/log4go"
	"github.com/bitly/go-simplejson"
	"io"
	"mime/multipart"

	//"io"
	"io/ioutil"
	"math/rand"
	//"mime/multipart"
	"net/http"
	"strconv"
	"strings"
	"time"
)

//header
// ctimestamp	客户端本地unix毫秒时间戳
// cnonce	6位随机字符串，可含字母，客户端自行生成
// cappkey	标识第三方系统身份，由鸿图分配
// csign	按签名方案计算的签名
//​ MD5，32位小写
// csign = MD5(A-B-C-D-E-ctimestamp-cnonce-cappkey)

//A：请求URI，以/开头，如/v1/api/device/list
//B：request的请求method，大写，如GET、POST、PUT、DELETE
//C：请求参数按照参数key正序排序后连接的字符串：如请求参数中有akey=1,ckey=2,bkey=3,则连接的参数字符串为：akey=1&bkey=3&ckey=2(注意：参数值采用url编码前的原值)，目前，我们大部分接口都是是有body传参，所以这里是""
//D：请求的requestbody中的内容的MD5(32位小写，下同)，如无内容，则为""【注，对于文件上传类url，该字段用空字符串替代】
//E： cappkey对应的秘钥，如：appkey1
//F：header中的ctimestamp、cnonce、cappkey

const secret  = "sdfajk3242324fa!djq7"


var G_HT = HT{}

type HT struct {
	Url_A 		string
	Method_B 	string
	Qury_Msg_C 	string
	Req_Body_D	string
	Secret_E	string
	Other_F		string
	timestamp	string
	vcode	string
}

func (h *HT)GetMD5code(encodestr string)(output string){
	//str,_ := GetMacAddrs()
	md5_c := md5.New()
	md5_c.Write([]byte(encodestr))
	return hex.EncodeToString(md5_c.Sum(nil))
}

func (h *HT)Getcnonce()string{
	nowtime := time.Now().UnixNano()
	rnd := rand.New(rand.NewSource(nowtime))
	vcode := fmt.Sprintf("%06v",rnd.Int31n(1000000))
	h.vcode = vcode
	return h.vcode
}
func (h *HT)Getctimestamp(){
	nowtime := time.Now().UnixNano()
	//rnd := rand.New(rand.NewSource(nowtime))
	//vcode := fmt.Sprintf("%06v",rnd.Int31n(1000000))
	h.timestamp = strconv.FormatInt(nowtime, 10)[0:13]
	//return h.timestamp
}

func (h *HT)GetOther_F(){

	//nowtime := time.Now().UnixNano()
	//rnd := rand.New(rand.NewSource(nowtime))
	//vcode := fmt.Sprintf("%06v",rnd.Int31n(1000000))
	//s_nowtime := strconv.FormatInt(nowtime, 10)
	h.Other_F = h.timestamp+"-"+h.vcode+"-"+"appkey1"
	//return s_nowtime+"-"+vcode+"-"+secret
}
func (h *HT)BodyMD5(){
	if h.Req_Body_D == ""{
		return
	}
	temp:= h.GetMD5code(h.Req_Body_D)
	h.Req_Body_D = temp
}
func(h *HT)GetSign()(string){
	h.BodyMD5()
	h.GetOther_F()
	cSign := h.GetMD5code(h.Url_A+"-"+h.Method_B+"-"+h.Qury_Msg_C+"-"+h.Req_Body_D+"-"+secret+"-"+h.Other_F)
	return cSign
}

func (h *HT)BuildHeard(r *http.Request)(){
	h.Getctimestamp()
	h.Getcnonce()
	r.Header.Set("ctimestamp",h.timestamp)
	r.Header.Set("cnonce",h.vcode)
	r.Header.Set("cappkey","appkey1")
	r.Header.Set("csign",h.GetSign())
}


//func AddPhotoByFile(photo multipart.File,hongtuHost string)(int, string, error){
//	client := &http.Client{
//		//Jar: jar,
//	}
//	METHOD := "POST"
//	url := "/v1/api/person/uploadImage"
//
//	ht := HT{}
//	ht.Req_Body_D = ""
//	ht.Method_B = METHOD
//	ht.Url_A = url
//
//	body := new(bytes.Buffer)
//	writer := multipart.NewWriter(body)
//	part, err := writer.CreateFormFile("photo", "photo.jpg")
//	if err != nil {
//		return -1, "", err
//	}
//	_, err = io.Copy(part, photo)
//	if err != nil {
//		log4go.Error(err.Error())
//		return -1, "", err
//	}
//	err = writer.Close()
//	if err != nil {
//		return -1, "", err
//	}
//	request, err := http.NewRequest(METHOD, hongtuHost+url, body)
//	G_HT.BuildHeard(request)
//	request.Header.Add("Content-Type", writer.FormDataContentType())
//	log4go.Debug(request.URL)
//	log4go.Debug(request.Method)
//	log4go.Debug(request.Header)
//	resp, err := client.Do(request)
//	if err != nil {
//		log4go.Error(err.Error())
//		return -1, "", err
//	}
//	resp_json, err := doResponse(resp)
//	if err != nil {
//		log4go.Error(err.Error())
//		return -1, "", err
//	}
//	//photo_id, _ := resp_json.Get("data").Get("id").Int()
//	//photo_url, _ := resp_json.Get("data").Get("url").String()
//
//	//if ishave := strings.HasPrefix(photo_url, "http://") &&  len(photo_url) > 0; ishave {
//	//} else {
//	//	photo_url =  photo_url
//	//}
//	//log4go.Debug(photo_url)
//	return photo_id, photo_url, nil
//}

var baseUrl = ""


func Init(host string, koalaProt int) error {
	log4go.Info("koala init -----------------------------------------")
	//baseUrl = "http://" + host + ":" + strconv.Itoa(koalaProt)
	baseUrl = "http://" + host
	//baseUrl = "http://" + host
	log4go.Info("koala url: " + baseUrl)
	fmt.Println("koala url: " + baseUrl)
	// All users of cookiejar should import "golang.org/x/net/publicsuffix"
	//jar, err := cookiejar.New(&cookiejar.Options{PublicSuffixList: publicsuffix.List})
	//var err error
	//jar, err = cookiejar.New(nil)
	//if err != nil {
	//	//log4go.Crash(err)
	//	log4go.Error(err)
	//	return err
	//}
	log4go.Info("init success -----------------------------------------")
	return nil
}
func KoalaLogin( username string, password string) error {
	return nil

}

func doResponse(resp *http.Response) (*simplejson.Json, error) {
	log4go.Debug(resp.Status)
	if resp.StatusCode != 200 {
		return nil, errors.New(resp.Status)
	}

	defer resp.Body.Close()

	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		log4go.Error(err.Error())
		return nil, errors.New("Read response body error")
	}

	jdata, err := simplejson.NewJson(body)
	if err != nil {
		log4go.Error(string(body))
		log4go.Error(err.Error())
		if strings.Index(string(body), "<span>记住我</span>") == -1 {
			return nil, errors.New("Face++返回报文错误")
		} else {
			return nil, errors.New("登陆失效，请重新登录!")
		}

	}
	code, _ := jdata.Get("code").Int()
	if code != 0 {
		//log4go.Error(string(body))
		desc, _ := jdata.Get("msg").String()
		log4go.Error(desc)
		return nil, errors.New(desc)
	}
	return jdata, nil
}


func GetDeviceList(params interface{})(*simplejson.Json, error){
	client := &http.Client{}
	METHOD := "POST"
	url := "/v1/api/device/list"

	ht := HT{}

	ht.Method_B = METHOD
	ht.Url_A = url
	ht.Qury_Msg_C = ""

	jsonBytes, err := json.Marshal(params)
	if err != nil {
		return nil, err
	}
	log4go.Info(string(jsonBytes))
	ht.Req_Body_D = string(jsonBytes)
	request, err := http.NewRequest(METHOD, baseUrl+url, strings.NewReader(string(jsonBytes)))
	ht.BuildHeard(request)
	request.Header.Add("Content-Type", "application/json")
	log4go.Debug(request.URL)
	log4go.Debug(request.Method)
	log4go.Debug(request.Header)
	resp, err := client.Do(request)
	if err != nil {
		log4go.Error(err.Error())
		return nil, err
	}
	resp_json, err := doResponse(resp)
	if err != nil {
		log4go.Error(err.Error())
		return nil, err
	}
	return resp_json, nil
}


func ReadFile(path string)  ([]byte){
	f, err := os.Open(path)
	if err != nil {
		fmt.Println("read file fail", err)
		return nil
	}
	defer f.Close()

	fd, err := ioutil.ReadAll(f)
	if err != nil {
		fmt.Println("read to fd fail", err)
		return nil
	}

	return fd
}


//用于photo为File形式
func AddPhoto(photo multipart.File ) (*simplejson.Json, error) {
	client := &http.Client{}

	METHOD := "POST"
	url := "/v1/api/person/uploadImage"

	ht := HT{}

	ht.Method_B = METHOD
	ht.Url_A = url
	ht.Qury_Msg_C = ""


	body := new(bytes.Buffer)
	writer := multipart.NewWriter(body)
	part, err := writer.CreateFormFile("file", "photo.jpg")

	if err != nil {
		return nil, err
	}
	_, err = io.Copy(part, photo)
	if err != nil {
		log4go.Error(err.Error())
		return nil, err
	}
	err = writer.Close()
	if err != nil {
		return nil, err
	}
	ht.Req_Body_D = ""
	//ht.Req_Body_D = `name=\"file\"; filename=\"photo.jpg\"`
	request, err := http.NewRequest(METHOD,baseUrl+url, body)
	request.Header.Add("Content-Type", writer.FormDataContentType())
	ht.BuildHeard(request)
	log4go.Debug(request.URL)
	log4go.Debug(request.Method)
	log4go.Debug(request.Header)
	resp, err := client.Do(request)
	if err != nil {
		log4go.Error(err.Error())
		return nil, err
	}
	resp_json, err := doResponse(resp)
	if err != nil {
		log4go.Error(err.Error())
		return nil, err
	}
	//&{map[code:0 data:map[uri:_ZzEwMF9mb3JldmVyQnVja2V0_ffc2834c9fea442da253f6c9dc56a4c9] msg:成功]}
	return resp_json, nil
}


func AddPhotoTest(photo []byte ) (*simplejson.Json, error) {
	client := &http.Client{}

	METHOD := "POST"
	url := "/v1/api/person/uploadImage"

	ht := HT{}

	ht.Method_B = METHOD
	ht.Url_A = url
	ht.Qury_Msg_C = ""


	body := new(bytes.Buffer)
	writer := multipart.NewWriter(body)
	part, err := writer.CreateFormFile("file", "photo.jpg")

	if err != nil {
		return nil, err
	}
	_, err = io.Copy(part, bytes.NewReader(photo))
	if err != nil {
		log4go.Error(err.Error())
		return nil, err
	}
	err = writer.Close()
	if err != nil {
		return nil, err
	}
	ht.Req_Body_D = ""
	//ht.Req_Body_D = `name=\"file\"; filename=\"photo.jpg\"`
	request, err := http.NewRequest(METHOD,baseUrl+url, body)
	request.Header.Add("Content-Type", writer.FormDataContentType())
	ht.BuildHeard(request)
	log4go.Debug(request.URL)
	log4go.Debug(request.Method)
	log4go.Debug(request.Header)
	resp, err := client.Do(request)
	if err != nil {
		log4go.Error(err.Error())
		return nil, err
	}
	resp_json, err := doResponse(resp)
	if err != nil {
		log4go.Error(err.Error())
		return nil, err
	}
	//&{map[code:0 data:map[uri:_ZzEwMF9mb3JldmVyQnVja2V0_ffc2834c9fea442da253f6c9dc56a4c9] msg:成功]}
	return resp_json, nil
}


func AddPerson(params interface{} )(*simplejson.Json,error){
	client := &http.Client{}
	METHOD := "POST"
	url := "/v1/api/person/batchAdd"

	ht := HT{}

	ht.Method_B = METHOD
	ht.Url_A = url
	ht.Qury_Msg_C = ""

	jsonBytes, err := json.Marshal(params)
	if err != nil {
		return nil, err
	}
	//log4go.Info(string(jsonBytes))
	ht.Req_Body_D = string(jsonBytes)
	request, err := http.NewRequest(METHOD,baseUrl+url, strings.NewReader(string(jsonBytes)))
	if err != nil {
		log4go.Error(err)
		return  nil,errors.New("New request error")
	}
	log4go.Debug(request.URL)
	log4go.Debug(request.Method)
	ht.BuildHeard(request)
	request.Header.Set("Content-Type", "application/json")
	resp, err := client.Do(request)
	if err != nil {
		log4go.Error(err.Error())
		return  nil,err
	}

	resp_json, err := doResponse(resp)
	fmt.Println(resp_json)
	if err != nil {
		log4go.Error(err.Error())
		return  nil,err
	}

	return resp_json,nil

}


func DeleteSubject(params interface{}) (*simplejson.Json,error) {
	client := &http.Client{}

	METHOD := "POST"
	url := "/v1/api/person/batchDelete"

	ht := HT{}

	ht.Method_B = METHOD
	ht.Url_A = url
	ht.Qury_Msg_C = ""

	jsonBytes, err := json.Marshal(params)
	if err != nil {
		return nil, err
	}
	//log4go.Info(string(jsonBytes))
	ht.Req_Body_D = string(jsonBytes)
	request, err := http.NewRequest(METHOD,baseUrl+url, strings.NewReader(string(jsonBytes)))
	if err != nil {
		log4go.Error(err)
		return nil,errors.New("New request error")
	}
	log4go.Debug(request.URL)
	log4go.Debug(request.Method)
	ht.BuildHeard(request)
	request.Header.Set("Content-Type", "application/json")
	resp, err := client.Do(request)
	if err != nil {
		log4go.Error(err.Error())
		return nil,err
	}

	resp_json, err := doResponse(resp)
	if err != nil {
		log4go.Error(err.Error())
		return nil,err
	}

	return resp_json,nil
}

func PersonList(params interface{}) (*simplejson.Json,error) {
	client := &http.Client{}

	METHOD := "POST"
	url := "/v1/api/person/list"

	ht := HT{}

	ht.Method_B = METHOD
	ht.Url_A = url
	ht.Qury_Msg_C = ""

	jsonBytes, err := json.Marshal(params)
	if err != nil {
		return nil, err
	}
	//log4go.Info(string(jsonBytes))
	ht.Req_Body_D = string(jsonBytes)
	request, err := http.NewRequest(METHOD,baseUrl+url, strings.NewReader(string(jsonBytes)))
	if err != nil {
		log4go.Error(err)
		return nil,errors.New("New request error")
	}
	log4go.Debug(request.URL)
	log4go.Debug(request.Method)
	ht.BuildHeard(request)
	request.Header.Set("Content-Type", "application/json")
	//req.Header.Set("Content-Type", "application/json")
	resp, err := client.Do(request)
	if err != nil {
		log4go.Error(err.Error())
		return nil,err
	}

	resp_json, err := doResponse(resp)
	if err != nil {
		log4go.Error(err.Error())
		return nil,err
	}

	return resp_json,nil
}

func PersonQuery(params interface{}) (*simplejson.Json,error) {
	client := &http.Client{}

	METHOD := "POST"
	url := "/v1/api/person/query"

	ht := HT{}

	ht.Method_B = METHOD
	ht.Url_A = url
	ht.Qury_Msg_C = ""

	jsonBytes, err := json.Marshal(params)
	if err != nil {
		return nil, err
	}
	//log4go.Info(string(jsonBytes))
	ht.Req_Body_D = string(jsonBytes)
	request, err := http.NewRequest(METHOD,baseUrl+url, strings.NewReader(string(jsonBytes)))
	if err != nil {
		log4go.Error(err)
		return nil,errors.New("New request error")
	}
	log4go.Debug(request.URL)
	log4go.Debug(request.Method)
	ht.BuildHeard(request)
	request.Header.Set("Content-Type", "application/json")
	resp, err := client.Do(request)
	if err != nil {
		log4go.Error(err.Error())
		return nil,err
	}

	resp_json, err := doResponse(resp)
	if err != nil {
		log4go.Error(err.Error())
		return nil,err
	}

	return resp_json,nil
}
func VisitorCode(params interface{}) (*simplejson.Json,error) {
	client := &http.Client{}

	METHOD := "POST"
	url := "/v1/api/person/visitorCode"

	ht := HT{}

	ht.Method_B = METHOD
	ht.Url_A = url
	ht.Qury_Msg_C = ""

	jsonBytes, err := json.Marshal(params)
	if err != nil {
		return nil, err
	}
	//log4go.Info(string(jsonBytes))
	ht.Req_Body_D = string(jsonBytes)
	request, err := http.NewRequest(METHOD,baseUrl+url, strings.NewReader(string(jsonBytes)))
	if err != nil {
		log4go.Error(err)
		return nil,errors.New("New request error")
	}
	log4go.Debug(request.URL)
	log4go.Debug(request.Method)
	ht.BuildHeard(request)
	request.Header.Set("Content-Type", "application/json")
	resp, err := client.Do(request)
	if err != nil {
		log4go.Error(err.Error())
		return nil,err
	}

	resp_json, err := doResponse(resp)
	if err != nil {
		log4go.Error(err.Error())
		return nil,err
	}

	return resp_json,nil
}


func Passgrouplist_WEB(params interface{}) (*simplejson.Json,error) {
	client := &http.Client{}

	METHOD := "POST"
	url := "/v1/web/pass/group/page"

	ht := HT{}

	ht.Method_B = METHOD
	ht.Url_A = url
	ht.Qury_Msg_C = ""

	jsonBytes, err := json.Marshal(params)
	if err != nil {
		return nil, err
	}
	//log4go.Info(string(jsonBytes))
	ht.Req_Body_D = string(jsonBytes)
	request, err := http.NewRequest(METHOD,baseUrl+url, strings.NewReader(string(jsonBytes)))
	if err != nil {
		log4go.Error(err)
		return nil,errors.New("New request error")
	}
	log4go.Debug(request.URL)
	log4go.Debug(request.Method)
	ht.BuildHeard(request)
	request.Header.Set("Content-Type", "application/json")
	//req.Header.Set("Content-Type", "application/json")
	resp, err := client.Do(request)
	if err != nil {
		log4go.Error(err.Error())
		return nil,err
	}

	resp_json, err := doResponse(resp)
	if err != nil {
		log4go.Error(err.Error())
		return nil,err
	}

	return resp_json,nil
}
func Passgrouplist(params interface{}) (*simplejson.Json,error) {
	client := &http.Client{}

	METHOD := "POST"
	url := "/v1/api/pass/group/list"

	ht := HT{}

	ht.Method_B = METHOD
	ht.Url_A = url
	ht.Qury_Msg_C = ""

	jsonBytes, err := json.Marshal(params)
	if err != nil {
		return nil, err
	}
	//log4go.Info(string(jsonBytes))
	ht.Req_Body_D = string(jsonBytes)
	request, err := http.NewRequest(METHOD,baseUrl+url, strings.NewReader(string(jsonBytes)))
	if err != nil {
		log4go.Error(err)
		return nil,errors.New("New request error")
	}
	log4go.Debug(request.URL)
	log4go.Debug(request.Method)
	ht.BuildHeard(request)
	request.Header.Set("Content-Type", "application/json")
	//req.Header.Set("Content-Type", "application/json")
	resp, err := client.Do(request)
	if err != nil {
		log4go.Error(err.Error())
		return nil,err
	}

	resp_json, err := doResponse(resp)
	if err != nil {
		log4go.Error(err.Error())
		return nil,err
	}

	return resp_json,nil
}
func Personupdate(params interface{}) (*simplejson.Json,error) {
	client := &http.Client{}

	METHOD := "POST"
	url := "/v1/api/person/update"

	ht := HT{}

	ht.Method_B = METHOD
	ht.Url_A = url
	ht.Qury_Msg_C = ""

	jsonBytes, err := json.Marshal(params)
	if err != nil {
		return nil, err
	}
	//log4go.Info(string(jsonBytes))
	ht.Req_Body_D = string(jsonBytes)
	request, err := http.NewRequest(METHOD,baseUrl+url, strings.NewReader(string(jsonBytes)))
	if err != nil {
		log4go.Error(err)
		return nil,errors.New("New request error")
	}
	log4go.Debug(request.URL)
	log4go.Debug(request.Method)
	ht.BuildHeard(request)
	request.Header.Set("Content-Type", "application/json")
	//req.Header.Set("Content-Type", "application/json")
	resp, err := client.Do(request)
	if err != nil {
		log4go.Error(err.Error())
		return nil,err
	}

	resp_json, err := doResponse(resp)
	if err != nil {
		log4go.Error(err.Error())
		return nil,err
	}

	return resp_json,nil
}


func PhotoSearch(params interface{}) (*simplejson.Json,error) {
	//uri	string	必须		图片寻址uri，通过<上传人员图片接口>获取uri
	//groupTypes	integer []	必须		底库来源 1-员工; 2-访客; 3-重点人员; 4-陌生人(备注：鸿图3003不支持4-陌生人);
	//threshold	integer	必须		阈值大小, 数值为正整数，阈值区间(0, 100)
	client := &http.Client{}

	METHOD := "POST"
	url := "/v1/api/photo/search"

	ht := HT{}

	ht.Method_B = METHOD
	ht.Url_A = url
	ht.Qury_Msg_C = ""

	jsonBytes, err := json.Marshal(params)
	if err != nil {
		return nil, err
	}
	//log4go.Info(string(jsonBytes))
	ht.Req_Body_D = string(jsonBytes)
	request, err := http.NewRequest(METHOD,baseUrl+url, strings.NewReader(string(jsonBytes)))
	if err != nil {
		log4go.Error(err)
		return nil,errors.New("New request error")
	}
	log4go.Debug(request.URL)
	log4go.Debug(request.Method)
	ht.BuildHeard(request)
	request.Header.Set("Content-Type", "application/json")
	//req.Header.Set("Content-Type", "application/json")
	resp, err := client.Do(request)
	if err != nil {
		log4go.Error(err.Error())
		return nil,err
	}

	resp_json, err := doResponse(resp)
	if err != nil {
		log4go.Error(err.Error())
		return nil,err
	}

	return resp_json,nil
}