package model

import (
	"database/sql"
	"gorm.io/gorm"
)

type SubjectShip struct {
	ID   int    `gorm:"primary_key:AUTO_INCREMENT;column:id;not null" json:"subject_id"`
	Uuid string `gorm:"column:uuid;index:idx_uuid" json:"uuid"` //2021.5.30被访人手机号码
	Uri  string `gorm:"column:uri;index:idx_uri" json:"uri"` //2021.5.30被访人手机号码
	Name  string `gorm:"column:name;index:idx_name" json:"name"` //2021.5.30被访人手机号码
	Phone  string `gorm:"column:phone;index:idx_phone" json:"phone"` //2021.5.30被访人手机号码
	VisitType  string `gorm:"column:visitType;index:idx_visitType" json:"visitType"` //2021.5.30被访人手机号码
	IdentifyNum  string `gorm:"column:identifyNum;index:idx_identifyNum" json:"identifyNum"` //2021.5.30被访人手机号码
}

type PhotoShip struct {
	ID   int    `gorm:"primary_key:AUTO_INCREMENT;column:id;not null" json:"id"`
	Uri string `gorm:"column:uri;index:idx_uri" json:"uri"` //2021.5.30被访人手机号码
	SubjectUuid string `gorm:"column:subject_uuid;index:idx_subject_uuid" json:"-"`
}

type GroupShip struct {
	ID   int    `gorm:"primary_key:AUTO_INCREMENT;column:id;not null" json:"subject_id"`
	Uuid string `gorm:"column:uuid;index:idx_uuid" json:"uuid"` //2021.5.30被访人手机号码
	Type string `gorm:"column:type;index:idx_type" json:"type"` //2021.5.30被访人手机号码
	Name string `gorm:"column:name;index:idx_name" json:"name"` //2021.5.30被访人手机号码
	SubjectUuid string `gorm:"column:subject_uuid;index:idx_subject_uuid" json:"-"`
}

type Visitor struct {
	//gorm.Model
	ID                int	 `gorm:"column:id" json:"id"`
	Name              string `json:"name"`                                                 //访客姓名
	Phone             string `gorm:"index:idx_phone" json:"phone" `                        //访客手机号码
	Come_from         string `json:"come_from"`                                            //访客单位
	Purpose           int    `json:"purpose"`                                              //来访目的id
	Purpose_name      string `json:"purpose_name"`                                         //来访目的描述
	Interviewee       string `json:"interviewee"`     //被访人                                      //被访人
	Interviewee_phone string `gorm:"index:idx_interviewee_phone" json:"interviewee_phone"` //2021.5.30被访人手机号码
	Proposer          string `gorm:"index:idx_proposer" json:"proposer"`                   //申请人的手机号码
	Start_time        int    `json:"start_time"`                                           //来访开始时间，时间戳
	End_time          int    `json:"end_time"`
	Is_auth           int    `json:"is_auth"` //2021.5.26授权标志： 0:未授权；1:已授权
	Remark            string `json:"remark"`
	Group_ids         int    `json:"group_ids"`
	Group_Name        string `json:"group_Name"`
	Department        string `json:"department"`
	Id_no           string `json:"id_no"` //2021.7.1访客身份证
	Subject_type int       `json:"subject_type"` //考拉系统对应的访客类型，固定为1
	Subject_id   int       `json:"subject_id"`   //考拉系统对应的访客id，初始为空，只有通过接口写入考拉里了才会返回有
	Photo        string    `json:"photo"`        //考拉系统对应的访客照片url，只有通过现场访客机拍照录入考拉里了才会返回有
	Photo_id     int       `json:"photo_id"`     //考拉系统对应的访客照片id，只有通过现场访客机拍照录入考拉里了才会返回有
	Retinues     []Retinue `gorm:"foreignKey:Parent_id" json:"retinues" `
	Qrcode	        string `json:"qrcode"` //2021.6.5 二维码
	Qrcode_img_url  string `json:"qrcode_img_url"` //2021.6.5
	Signin_time  int       `json:"signin_time"`     //2021.6.8 访客签到时间
}

type NettyObj struct {

	SceneImg         string	 `json:"sceneImg" form:"sceneImg"`
	IdImg         string	 `json:"idImg" form:"idImg"`
	Fvisitormobile         string	 `json:"fvisitormobile" form:"fvisitormobile"`
	//Fvisitdate         string	 `json:"fvisitdate"`
}

type Retinue struct {
	//gorm.Model
	ID         int	 `json:"id"`
	Parent_id  int    //关联visitor表的id字段
	Name       string `json:"name"`
	Phone      string `json:"phone"`
	Subject_id int    `json:"subject_id"` //考拉系统对应的访客id，初始为空，只有通过接口写入考拉里了才会返回有
	Photo      string `json:"photo"`      //考拉系统对应的访客照片url，只有通过现场访客机拍照录入考拉里了才会返回有
	Photo_id   int    `json:"photo_id"`   //考拉系统对应的访客照片id，只有通过现场访客机拍照录入考拉里了才会返回有
	Id_no           string `json:"id_no"` //2021.7.1访客身份证
	Qrcode	        string `json:"qrcode"` //2021.6.5
	Qrcode_img_url  string `json:"qrcode_img_url"` //2021.6.5
	Signin_time  int       `json:"signin_time"`     //2021.6.8 访客签到时间
}

//2021.6.2 新增，用户信息，便于快捷取得考拉数据
type Employee struct {
	gorm.Model
	Phone      string `json:"phone"`
	Subject_id int    `json:"subject_id"` //考拉系统对应的访客id，初始为空，只有通过接口写入考拉里了才会返回有
	Uuid string    `json:"uuid"` //考拉系统对应的访客id，初始为空，只有通过接口写入考拉里了才会返回有
	WeixinUserId string `json:"weixinUserId"` //2021.6.6 新增微信id
	//Name string `json:"name"` //2021.7.10 新增微信名称
}

type Account struct {
	ID         int    `gorm:"primary_key:AUTO_INCREMENT;column:id;not null" json:"id"`
	Ip_address string `gorm:"column:ip_address" json:"ip_address"`
	Account    string `gorm:"column:account" json:"account"`
	Password   string `gorm:"column:password" json:"password"`
	Activation int    `gorm:"column:activation" json:"activation"`
}

type Template struct{
	Id  int 	`gorm:"primary_key;AUTO_INCREMENT;not null" json:"id" toml:"id"`
	Approval_reminder_to_employee_wechat_id	   sql.NullString	`gorm:"approval_reminder_to_employee_wechat_id" json:"approval_reminder_to_employee_wechat_id" toml:"approval_reminder_to_employee_wechat_id"`
	Approval_reminder_to_employee_sms_id	   sql.NullString 	`gorm:"approval_reminder_to_employee_sms_id" json:"approval_reminder_to_employee_sms_id" toml:"approval_reminder_to_employee_sms_id"`
	Visitor_signin_to_employee_wechat_id   sql.NullString   `gorm:"visitor_signin_to_employee_wechat_id" json:"visitor_signin_to_employee_wechat_id" toml:"visitor_signin_to_employee_wechat_id"`
	Visitor_signin_to_employee_sms_id	   sql.NullString	`gorm:"visitor_signin_to_employee_sms_id" json:"visitor_signin_to_employee_sms_id" toml:"visitor_signin_to_employee_sms_id"`
	Approval_pass_to_visitor_wechat_id	   sql.NullString	`gorm:"approval_pass_to_visitor_wechat_id" json:"approval_pass_to_visitor_wechat_id" toml:"approval_pass_to_visitor_wechat_id"`
	Approval_pass_to_visitor_sms_id	   sql.NullString	`gorm:"approval_pass_to_visitor_sms_id" json:"approval_pass_to_visitor_sms_id" toml:"approval_pass_to_visitor_sms_id"`
	Signin_to_visitor_wechat_id	   sql.NullString	`gorm:"signin_to_visitor_wechat_id" json:"signin_to_visitor_wechat_id" toml:"signin_to_visitor_wechat_id"`
	Signin_to_visitor_sms_id	   sql.NullString	`gorm:"signin_to_visitor_sms_id" json:"signin_to_visitor_sms_id" toml:"signin_to_visitor_sms_id"`
	Signin_to_visitor_place_sms_id	   sql.NullString	`gorm:"signin_to_visitor_place_sms_id" json:"signin_to_visitor_place_sms_id" toml:"signin_to_visitor_place_sms_id"`
}




