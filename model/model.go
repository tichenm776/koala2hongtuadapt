package model

type Resp4Device struct {
	Code    int         `json:"code"`
	Err_msg string      `json:"err_msg"`
	Data    interface{} `json:"data"`
	Page    interface{} `json:"Page"`
}

type EmployeeRecords_json struct {
	Name string `json:"name"`
	Page int    `json:"page"`
	Size int    `json:"size"`
}

type EmployeeBinding struct {
	Phone   string `json:"phone"`
	Name   string `json:"name"`
	Unionid int    `json:"unionid"`
}

type Employeedelete struct {
	Proposer   string `json:"proposer"`
	Subject_id int    `json:"subject_id"`
}

type Employeeget struct {
	Proposer   string `json:"proposer"`
	Subject_id int    `json:"subject_id"`
	Page       int    `json:"page"`
	Size       int    `json:"size"`
}

type Personlist struct {
	Subject_type int `json:"subject_type"`
	Page         int `json:"page"`
	Size         int `json:"size"`
}

type Staff struct {
	Code    int    `json:"code"`
	Err_msg string `json:"err_msg"`
	//Photo	string	`json:"photo"`
	Photo_id int `json:"photo_id"`
}

type Photo struct {
	//Code    int         `json:"code"`
	//Err_msg string      `json:"err_msg"`
	Photo string `json:"photo"`
	//Photo_id  int  `json:"photo_id"`
}

type EmployeeInfo struct {
	Subject_type int    `json:"subject_type"`
	Purpose      int    `json:"purpose"`
	Remark       string `json:"remark"`
	Phone        string `json:"phone"`
	Name         string `json:"name"`
	Come_from    string `json:"come_from"`
	//Photo_id     int     `json:"photo_id"` //2021.5.30 注释
	Start_time        int64  `json:"start_time"`
	End_time          int64  `json:"end_time"`
	Interviewee       string `json:"interviewee"`
	Interviewee_phone string `json:"interviewee_phone"`
	//Department  string  `json:"department"` //2021.5.30 注释
	Group_ids int    `json:"group_ids"`
	Proposer  string `json:"proposer"` // 申请人手机号码
	Photo     string `json:"photo"`
	//Retinue     []Retinue	`json:"retinue"`
}

// type Retinue struct{
// 	Name    string    `json:"name"`
// 	Phone   string    `json:"phone"`
// 	Photo   string    `json:"photo"`
// 	//Photo_id   int    `json:"photo_id"` //2021.5.30 注释
// }

//2021.5.30 新增访客查询请求信息
type VisitorReq struct {
	Proposer string `json:"proposer"` //申请人的手机号码
	Phone    string `json:"phone"`    //访客的手机号码
	Name     string `json:"name"`     //访客姓名
	Interviewee_phone string `json:"interviewee_phone"` //2021.6.6被访人手机号码
	Interviewee       string `json:"interviewee"`     //被访人     
	Start_time        int    `json:"start_time"`    //2021.6.6                                       //来访开始时间，时间戳
	End_time          int    `json:"end_time"` //2021.6.6
	Is_auth           string `json:"is_auth"` //2021.6.6授权标志： 0:未授权；1:已授权
	//ID_Card           string `json:"id_card"` //2021.6.6授权标志： 0:未授权；1:已授权

	Page     int    `json:"page"`
	Size     int    `json:"size"`
}

//2021.6.6
type BarCodeReq struct {
	Qrcode string `json:"qrcode"`
}

//2021.5.30 新增微信信息
type WeixinInfo struct {
	WeixinUserId string `json:"weixinUserId"`
	Name         string `json:"name"`
	Phone        string `json:"phone"`
}

//2021.5.30 新增分页信息
type PageInfo struct {
	Count   int `json:"count"`
	Current int `json:"current"`
	Size    int `json:"size"`
	Total   int `json:"total"`
}
//2021.5.30 新增访客查询请求信息
type QRCODE_VisitorReq struct {
	Phone    string `json:"phone" form:"phone"`    //访客的手机号码
	Subject_id        int    `json:"subject_id" form:"subject_id"`    //2021.6.6                                       //该访客在koala服务器内的id
	Start_time        int    `json:"start_time" form:"start_time"`    //2021.6.6                                       //来访开始时间，时间戳
	End_time          int    `json:"end_time" form:"end_time"` //2021.6.6
	Is_auth           string `json:"is_auth" form:"is_auth"` //2021.6.6授权标志： 0:未授权；1:已授权
	Camera_position          string    `json:"camera_position" form:"camera_position"` //2021.6.6
}

type TemplateReq struct{
	Approval_reminder_to_employee_wechat_id	   string	`json:"approval_reminder_to_employee_wechat_id"`
	Approval_reminder_to_employee_sms_id	   string 	`json:"approval_reminder_to_employee_sms_id"`
	Visitor_signin_to_employee_wechat_id  string   `json:"visitor_signin_to_employee_wechat_id"`
	Visitor_signin_to_employee_sms_id	   string	`json:"visitor_signin_to_employee_sms_id"`
	Approval_pass_to_visitor_wechat_id	   string	`json:"approval_pass_to_visitor_wechat_id"`
	Approval_pass_to_visitor_sms_id	   string	`json:"approval_pass_to_visitor_sms_id"`
	Signin_to_visitor_wechat_id	   string	`json:"signin_to_visitor_wechat_id"`
	Signin_to_visitor_sms_id	   string	`json:"signin_to_visitor_sms_id"`
	Signin_to_visitor_place_sms_id	   string	`json:"signin_to_visitor_place_sms_id"`
}

var G_TemplateReq = TemplateReq{}