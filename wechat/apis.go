package wechat

import (
	"fmt"
	"net/http"
	"zhiyuan/zyutil/config"
	"io/ioutil"
	"encoding/json"
	"net/url"
	"github.com/alecthomas/log4go"
	"strconv"
	"bytes"
)

//2.1 查询微信是否已绑定 /v2/wechat/binding?weixinUserId=
// {
//     "code": 0,
//     "data": {
//         "avatar": "http://thirdwx.qlogo.cn/mmopen/YLt86Zd0qic11xR087H4qAlFibEiaBdGlicoJGU9JOfF8zrqovAoO8jDPRfMyXt8gI7sa3gbtibhoick2b2fW9BjlJyP3Wd6gtQN5A/132",
//         "nick": "雷振林",
//         "phone": "13355780752",
//         "weixinUserId": "oJajH6SCem9MSNa_gUgVsqABFPR4"
//     },
//     "err_msg": ""
// }
func QueryWechatBinding(weixinUserId string)(m map[string]interface{}){
	baseUrl := config.Gconf.WechatServiceUrl
	urlStr:=fmt.Sprintf("%s/v2/wechat/binding?weixinUserId=%s" ,baseUrl, weixinUserId)
	fmt.Println("QueryWechatBing",urlStr)
	resp, err := http.Get( urlStr )
	defer resp.Body.Close()
	m = ParseResponse(resp, err)
	return
}

//2. 微信绑定 /v2/wechat/binding 
//form data: weixinUserId, phone
func WechatBinding(weixinUserId,phone string)(m map[string]interface{}){
	baseUrl := config.Gconf.WechatServiceUrl
	urlStr:=fmt.Sprintf("%s/v2/wechat/binding" ,baseUrl)
	fmt.Println("WechatBing",urlStr)
	reqData:=url.Values{
		"weixinUserId":{weixinUserId},
		"phone":{phone},
	}
	resp, err := http.PostForm(urlStr,reqData)
	defer resp.Body.Close()
	m = ParseResponse(resp, err)
	return
}

//处理其他接口返回的数据
func ParseResponse(resp *http.Response,err error)(m map[string]interface{}){
	if err != nil{
		log4go.Error(err.Error())
		m = map[string]interface{}{ 
			"code":resp.StatusCode,
			 "err_msg": err.Error(),
		}
		return	
	}
	if resp.StatusCode!=200{
		data, err1 := ioutil.ReadAll(resp.Body)
		if err1 != nil{
			log4go.Error(err1.Error())
			m = map[string]interface{}{ 
				"code":resp.StatusCode,
				"err_msg": err1.Error(),
			}
			return	
		}else{
			m = map[string]interface{}{ 
				"code":resp.StatusCode,
				"err_msg": "微信接口异常",
			}
			log4go.Error(strconv.Itoa(resp.StatusCode)  + string(data) )
		}	
		return
	}
	data, err1 := ioutil.ReadAll(resp.Body)
	if err1 != nil{
		log4go.Error(err1.Error())
		m = map[string]interface{}{ 
			"code":resp.StatusCode,
			 "err_msg": err1.Error(),
		}
		return	
	}
	m = make(map[string]interface{})
	err1 = json.Unmarshal(data, &m)
	if err1 != nil{
		log4go.Error(err1.Error())
		m = map[string]interface{}{ 
			"code":resp.StatusCode,
			"err_msg": err1.Error(),
		}
	}
	return
}

//3 微信解绑 /v2/wechat/binding/:weixinUserId
func WechatUnBinding(weixinUserId string)(m map[string]interface{}){
	baseUrl := config.Gconf.WechatServiceUrl
	urlStr:=fmt.Sprintf("%s/v2/wechat/binding/%s" ,baseUrl, weixinUserId)
	log4go.Info("WechatUnBinding:" + urlStr)
	req, _ := http.NewRequest("DELETE", urlStr, nil)
	resp, err := http.DefaultClient.Do(req)
	defer resp.Body.Close()
	m = ParseResponse(resp, err)
	return
}

// 4 微信发送 POST /v2/message  2021.6.2
// | 参数         | 类型   | 必选 | 描述                     |
// | phone        | string |  是  | 用户手机号码(不能为空)   |
// | tplMessageId | string |  是  | 消息模板编号(不能为空)   |
// | url          | string |  非  | 消息详情URL(允许为空)    |
// | data         | string |  是  | 消息内容参数表(不能为空) |
// * data消息内容参数表说明
// | 字段     | 类型   | 必选 | 允许为空 | 说明               |
// | first    | string | 非   | 是       | first参数值        |
// | keyword1 | string | 非   | 是       | keyword1参数值     |
// | keyword2 | string | 非   | 是       | keyword2参数值     |
// | keyword3 | string | 非   | 是       | keyword3参数值     |
// | keyword4 | string | 非   | 是       | keyword4参数值     |
// | keyword5 | string | 非   | 是       | keyword5参数值     |
// | ...      | string | 非   | 是       | keyword参数值      |
// | remark   | string | 非   | 是       | remark参数值参数值 |
func SendWechatMessage(phone,tplMessageId,url string,data map[string]interface{})(m map[string]interface{}){
	baseUrl := config.Gconf.WechatServiceUrl
	urlStr:=fmt.Sprintf("%s/v2/message" ,baseUrl)
	log4go.Info("SendWechatMessage:" + urlStr)
	databytes, _ := json.Marshal(data)
	jsonData := map[string]interface{}{ 
		"phone":phone,
		"tplMessageId": tplMessageId,
		"url": url,
		"data":string(databytes),
	}
	jsonbytes, err := json.Marshal(jsonData)
	if err != nil{
		log4go.Error(err.Error())
		m = map[string]interface{}{ 
			"code":-100,
			"err_msg": err.Error(),
		}
		return	
	}
 	resp,err:=http.Post(urlStr,"application/json;charset=utf-8",bytes.NewBuffer(jsonbytes))
	defer resp.Body.Close()
	m = ParseResponse(resp, err)
	return
}