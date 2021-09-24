//author: lei zhenlin  2018-04-15

package util

import (
	//	"bytes"
	"errors"

	"zhiyuan/koala_event_server/src/config"
	//	"io"
	"io/ioutil"
	//	"mime/multipart"
	"net/http"
	//	"net/url"
	//	"strconv"
	"strings"

	//	"net/http/cookiejar"

	"github.com/alecthomas/log4go"
	"github.com/bitly/go-simplejson"
)

func DoResponse(resp *http.Response) (*simplejson.Json, error) {
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
	log4go.Debug(string(body))

	jdata, err := simplejson.NewJson(body)
	if err != nil {
		log4go.Error(err.Error())
		return nil, errors.New("返回报文错误")
	}

	code, _ := jdata.Get("code").Int()
	if code != 0 {
		err_msg, _ := jdata.Get("err_msg").String()
		log4go.Error(err_msg)
		return nil, errors.New(err_msg)
	}
	return jdata, nil
}

func Process_recognized(message string) {
	//log4go.Debug(config.Gconf.Recognized)
	if len(config.Gconf.Recognized) == 0 {
		return
	}

	req, err := http.NewRequest("POST", config.Gconf.Recognized, strings.NewReader(message))
	if err != nil {
		log4go.Error(err)
		return
	}
	log4go.Debug(req.URL)
	log4go.Debug(req.Method)

	req.Header.Set("Content-Type", "application/json")
	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		log4go.Error(err.Error())
		return
	}

	_, err = DoResponse(resp)
	if err != nil {
		log4go.Error(err.Error())
		return
	}

}
