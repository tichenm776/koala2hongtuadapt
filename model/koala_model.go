package model



type Subject struct {
	//gorm.Model
	Subject_type	int	`json:"subject_type" form:"subject_type"`
	Extra_id string `json:"extra_id" form:"extra_id"`
	Create_time    int64 `json:"create_time" form:"create_time"`
	Job_number   string `json:"job_number" form:"job_number"`	//ID_CARD(身份证)，PASSPORT(护照)，HK_PASS(港澳通行证)，DRIVER_LICENSE(驾照)，TAIWAN_PASS(台湾通行证)，OTHER(其他)
	Visitor_type int    `json:"visitor_type" form:"visitor_type"`
	Title string    `json:"title" form:"title"`
	Entry_date int64    `json:"entry_date" form:"entry_date"`
	Wg_number string    `json:"wg_number" form:"wg_number"`
	Department string    `json:"department" form:"department"`
	Email string    `json:"email" form:"email"`
	Photo_ids []int    `json:"photo_ids" form:"photo_ids"`
	Avatar string    `json:"avatar" form:"avatar"`
	Description int    `json:"description" form:"description"`
	Start_time int64    `json:"start_time" form:"start_time"`
	Interviewee int    `json:"interviewee" form:"interviewee"`
	Phone int    `json:"phone" form:"phone"`
	Birthday int64    `json:"birthday" form:"birthday"`
	Purpose int    `json:"purpose" form:"purpose"`
	Come_from         int    `json:"come_from" form:"come_from"`
	Remark         int    `json:"remark" form:"remark"`
	Group_ids         []int    `json:"group_ids" form:"group_ids"`
	Name         int    `json:"name" form:"name"`
	Gender         int    `json:"gender" form:"gender"`
	End_time         int64    `json:"end_time" form:"end_time"`
}
type Subject2 struct {
	//gorm.Model
	Subject_type	int	`json:"subject_type" form:"subject_type"`
	Extra_id string `json:"extra_id" form:"extra_id"`
	Create_time    int64 `json:"create_time" form:"create_time"`
	Job_number   string `json:"job_number" form:"job_number"`	//ID_CARD(身份证)，PASSPORT(护照)，HK_PASS(港澳通行证)，DRIVER_LICENSE(驾照)，TAIWAN_PASS(台湾通行证)，OTHER(其他)
	Visitor_type int    `json:"visitor_type" form:"visitor_type"`
	Title string    `json:"title" form:"title"`
	Entry_date int64    `json:"entry_date" form:"entry_date"`
	Wg_number string    `json:"wg_number" form:"wg_number"`
	Department string    `json:"department" form:"department"`
	Email string    `json:"email" form:"email"`
	Photo_ids []int    `json:"photo_ids" form:"photo_ids"`
	Avatar string    `json:"avatar" form:"avatar"`
	Description int    `json:"description" form:"description"`
	Start_time int64    `json:"start_time" form:"start_time"`
	Interviewee int    `json:"interviewee" form:"interviewee"`
	Phone int    `json:"phone" form:"phone"`
	Birthday int64    `json:"birthday" form:"birthday"`
	Purpose int    `json:"purpose" form:"purpose"`
	Come_from         int    `json:"come_from" form:"come_from"`
	Remark         int    `json:"remark" form:"remark"`
	Group_ids         []int    `json:"group_ids" form:"group_ids"`
	Name         int    `json:"name" form:"name"`
	Gender         int    `json:"gender" form:"gender"`
	End_time         int64    `json:"end_time" form:"end_time"`
}

type Login struct {
	Username	string	`json:"username"`
	Password	string	`json:"password"`
	Captchas	string	`json:"captchas"`
	Auth_token	bool	`json:"auth_token"`
}
