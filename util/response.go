package util

// Response 基础序列化器
type Response struct {
	Code  int         `json:"code"`
	//Data  interface{} `json:"data,omitempty"`
	Data  interface{} `json:"data"`
	Msg   string      `json:"message"`
	//Desc string      `json:"desc,omitempty"`
	Desc string      `json:"desc"`
}
func Err(errCode int, msg string, err error) Response {
	res := Response{
		Code: errCode,
		Msg:  msg,
		Desc:  err.Error(),
	}
	// 生产环境隐藏底层报错
	//if err != nil && gin.Mode() != gin.ReleaseMode {
	//	res.Desc = err.Error()
	//}
	return res
}
