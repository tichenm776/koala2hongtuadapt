package dao

import (
	"errors"
	"fmt"
	"github.com/alecthomas/log4go"
	"github.com/satori/go.uuid"
	"strings"
	"time"
	"zhiyuan/koala2hongtuadapt/model"
	"zhiyuan/zyutil"
)

//新增访客记录 黄俊add
func CreateVisitor(data *model.Visitor) (Visitor_obj model.Visitor, err error) {
	if err := Db.Create(data); err.Error != nil {

		log4go.Error(err)
		return model.Visitor{}, errors.New("Visitor insert db err")
	}
	if err := Db.Last(&Visitor_obj); err.Error != nil {
		log4go.Error(err)
		return model.Visitor{}, errors.New("select visitor in last time err")
	}
	//2021.6.8 黄俊新增
	retinues := make([]model.Retinue, 10)
	Db.Where("parent_id=?", Visitor_obj.ID).Find(&retinues)
	Visitor_obj.Retinues = retinues
	return Visitor_obj, nil
}

//根据ID查询访客数据  黄俊add
func GetVisitor(id int) (visitor model.Visitor,isFind bool) {
	visitor=model.Visitor{}
	tx := Db.First(&visitor, id)
	if tx.RowsAffected !=0 {
		retinues := make([]model.Retinue, 10)
		Db.Where("parent_id=?", visitor.ID).Find(&retinues)
		visitor.Retinues = retinues
	}
	isFind = tx.RowsAffected == 1
	return
}
//根据ID查询访客数据  黄俊add
func GetVisitorBySubjectid(subject_id int) (visitor model.Visitor,isFind bool) {
	visitor=model.Visitor{}
	tx := Db.Debug().Where("subject_id = ?",subject_id).Last(&visitor)
	//if tx.RowsAffected !=0 {
	//	retinues := make([]model.Retinue, 10)
	//	Db.Where("parent_id=?", visitor.ID).Find(&retinues)
	//	visitor.Retinues = retinues
	//}
	isFind = tx.RowsAffected == 1
	return
}
//保存更新  黄俊add
func UpdateVisitor(visitorInfo model.Visitor) {
	Db.Debug().Save(&visitorInfo)
	return
}

//更新二维码，2021.6.8 新增signin_time
func UpdateVisitorBarCode(visitor model.Visitor) {
	Db.Exec("update visitors set qrcode=?,signin_time=now() where id=?",visitor.Qrcode,visitor.ID)
	return
}
//更新二维码，2021.6.24 修改subjectid
func UpdateVisitorSubjectId(visitor model.Visitor) (model.Visitor){
	Db.Exec("update visitors set subject_id=? where phone=?",visitor.Subject_id,visitor.Phone)
	Db.Exec("update retinues set subject_id=? where phone=?",visitor.Subject_id,visitor.Phone)
	return visitor
}
//更新二维码，2021.6.8 新增signin_time
func UpdateRetinueBarCode(retinue model.Retinue) {
	Db.Exec("update retinues set qrcode=?,signin_time=now() where id=?",retinue.Qrcode,retinue.ID)
	return
}
func UpdateRetinue(retinue model.Retinue) {
	Db.Debug().Save(&retinue)
	return
}
//访客记录删除  黄俊add 2021.6.6
func DeleteVisitor(id int) {
	Db.Unscoped().Where("parent_id = ?", id).Delete(&model.Retinue{})
	Db.Unscoped().Where("id = ?", id).Delete(&model.Visitor{})
	return
}

// 追加SQL where条件 黄俊add
func appendCondition(sqlWhere, condition string) (sqlWhere1 string) {
	if sqlWhere != "" {
		sqlWhere1 = sqlWhere + " and " + condition
	} else {
		sqlWhere1 = condition
	}
	return
}

//查询列表  黄俊add
//func GetVisitors(req *model.VisitorReq) (visitors []model.Visitor, err error, page model.PageInfo) {
//	sqlWhere := ""
//	if req.Proposer != "" {
//		sqlWhere = appendCondition(sqlWhere, fmt.Sprintf("proposer = '%s' ", req.Proposer))
//	}
//	if req.Phone != "" {
//		sqlWhere = appendCondition(sqlWhere, fmt.Sprintf("phone = '%s' ", req.Phone))
//	}
//	if req.Name != "" {//访客姓名
//		sqlWhere = appendCondition(sqlWhere, fmt.Sprintf("name = '%s' ", req.Name))
//	}
//	if req.Interviewee != "" {//被访人
//		sqlWhere = appendCondition(sqlWhere, fmt.Sprintf("interviewee = '%s' ", req.Interviewee))
//	}
//	if req.Interviewee_phone != "" {//被访人手机号码
//		sqlWhere = appendCondition(sqlWhere, fmt.Sprintf("interviewee_phone = '%s' ", req.Interviewee_phone))
//	}
//	if req.Is_auth != "" { //授权标志
//		sqlWhere = appendCondition(sqlWhere, fmt.Sprintf("is_auth = %s ", req.Is_auth))
//	}
//	if req.Page == 0 {
//		req.Page = 1
//	}
//	if req.Size == 0 {
//		req.Size = 5
//	}
//	fmt.Println(sqlWhere)
//	offset := (req.Page - 1) * req.Size
//	count := 0
//	var visitor model.Visitor
//	Db.Model(&visitor).Where(sqlWhere).Order("ID desc").Offset(offset).Limit(req.Size).Find(&visitors)
//	fmt.Println(visitors)
//	Db.Raw("select count(*) from visitors where " + sqlWhere).Scan(&count)
//	total := zyutil.GetTotal(count, req.Size)
//	page = model.PageInfo{Count: count, Current: req.Page, Size: req.Size, Total: total}
//
//	for i := 0; i < len(visitors); i++ {
//		retinues := make([]model.Retinue, req.Size)
//		Db.Where("parent_id=?", visitors[i].ID).Find(&retinues)
//		visitors[i].Retinues = retinues
//	}
//	return visitors, nil, page
//}
//查询列表  Tim add
func GetVisitors(req *model.VisitorReq)(visitors []model.Visitor, err error, page model.PageInfo){

	log4go.Debug("params Proposer is:",req.Proposer)

	DBdate := Db
	DBCount := Db
	count := int64(0)
	if req.Start_time != 0 {
		DBdate = DBdate.Where("start_time >= ?", req.Start_time)
		DBCount = DBdate.Where("start_time >= ?", req.Start_time)
	}
	if req.End_time != 0 {
		DBdate = DBdate.Where("end_time <= ?", req.End_time)
		DBCount = DBdate.Where("end_time <= ?", req.End_time)
	}
	if req.Proposer != "" {
		DBdate = DBdate.Where("proposer = ?", req.Proposer)
		DBCount = DBdate.Where("proposer = ?", req.Proposer)
	}
	if req.Phone != "" {
		DBdate = DBdate.Where("phone = ?", req.Phone)
		DBCount = DBdate.Where("phone = ?", req.Phone)
	}
	if req.Name != "" {//访客姓名
		DBdate = DBdate.Where("name like ?", "%" +req.Name+"%")
		DBCount = DBdate.Where("name like ?", "%" +req.Name+"%")
	}
	if req.Interviewee != "" {//被访人
		DBdate = DBdate.Where("interviewee like ?", "%" +req.Interviewee+"%")
		DBCount = DBdate.Where("interviewee like ?", "%" +req.Interviewee+"%")
	}
	if req.Interviewee_phone != "" {//被访人手机号码
		DBdate = DBdate.Where("interviewee_phone = ?", req.Interviewee_phone)
		DBCount = DBdate.Where("interviewee_phone = ?", req.Interviewee_phone)
	}
	if req.Is_auth != "" { //授权标志
		DBdate = DBdate.Where("is_auth = ?", req.Is_auth)
		DBCount = DBdate.Where("is_auth = ?", req.Is_auth)
	}
	if req.Page == 0 {
		req.Page = 1
	}
	if req.Size == 0 {
		req.Size = 5
	}
	DBdate = DBdate.Order("id desc")
	if req.Page  > 0 {
		DBdate = DBdate.Limit(req.Size).Offset((req.Page - 1) * req.Size)
	}
	if err := DBdate.Model(model.Visitor{}).Find(&visitors).Error;err!= nil {

		return visitors, err, page
	}

	if err := DBCount.Model(model.Visitor{}).Count(&count).Error;err != nil {
		return visitors, err, page
	}
	count2 := int(count)
	total := zyutil.GetTotal(count2, req.Size)
	page.Size = req.Size
	page.Count = count2
	page.Total = total
	page.Current = req.Page

	for i := 0; i < len(visitors); i++ {
		retinues := make([]model.Retinue, req.Size)
		Db.Where("parent_id=?", visitors[i].ID).Find(&retinues)
		visitors[i].Retinues = retinues
	}



	return visitors, nil, page
}
//访客签到
func VisitorSign(phone string) (visitor model.Visitor,err error) {
	visitor=model.Visitor{}
	retinue:=model.Retinue{}
	result:=Db.Where("phone = ? and FROM_UNIXTIME(start_time,'%Y-%m-%d') = CURDATE() and is_auth=1",phone).Last(&visitor)
	uid, _ := uuid.NewV4()
	if result.RowsAffected == 0{//随行人员
		result=Db.Raw("select * from retinues r where phone = ? and exists( select id from visitors where FROM_UNIXTIME(start_time,'%Y-%m-%d') = CURDATE() and is_auth=1 and id=r.parent_id) ",phone).Last(&retinue)
		if result.RowsAffected == 0{//没有找到记录
			err = errors.New("没有找到今日已授权的访客记录")
		}else{
			barcode:= fmt.Sprintf(`{"id":%d,"qrcode":"%s","phone":"%s","type":"retinue"}`,retinue.ID,uid,phone)
			retinue.Qrcode = barcode
			//retinue.Qrcode_img_url = uid + ".jpg"
			UpdateRetinueBarCode(retinue)
			visitor,_ = GetVisitor(retinue.Parent_id)
		}
	}else{//访客
		barcode:= fmt.Sprintf(`{"id":%d,"qrcode":"%s","phone":"%s","type":"visitor"}`,visitor.ID,uid,phone)
		visitor.Qrcode = barcode
		//visitor.Qrcode_img_url = config.Gconf.CloudUrl + uid + ".jpg"
		UpdateVisitorBarCode(visitor)
		visitor,_ = GetVisitor(visitor.ID)
	}
	return
}

//新增员工记录 黄俊add
func CreateEmployee(data *model.Employee) (employee model.Employee, err error) {
	if err := Db.Create(data); err.Error != nil {
		log4go.Error(err)
		return model.Employee{}, errors.New("Employee insert db err")
	}
	if err := Db.Last(&employee); err.Error != nil {
		log4go.Error(err)
		return model.Employee{}, errors.New("select Employee in last time err")
	}
	return employee, nil
}

//根据电话号码查询员工数据  黄俊add
func GetEmployee(phone string) (employee model.Employee,isFind bool ) {
	result := Db.Debug().Limit(1).Last(&employee, "phone = ?", phone)
	isFind = result.RowsAffected == 1
	return
}
func GetEmployee2(phone string) (employee model.Employee,isFind bool ) {
	result := Db.Debug().Limit(1).Last(&employee, "phone = ?", phone)
	isFind = result.RowsAffected == 1
	return
}
//保存更新  黄俊add 2021.6.5
func UpdateEmployee(employee model.Employee) {
	Db.Save(&employee)
	return
}

//根据微信号删除用户信息  黄俊add 2021.6.6
func DeleteEmployeeByWeixinUserId(weixinUserId string) {
	Db.Unscoped().Where("weixin_user_id = ?", weixinUserId).Delete(&model.Employee{})
	return
}

//根据id查询随行人员数据  黄俊add 2021.6.6
func GetRetinue(id int) (retinue model.Retinue,isFind bool ) {
	tx := Db.Last(&retinue, id)
	isFind = tx.RowsAffected == 1
	return
}

func Createtimes(data *model.Visitor) (Visitor_obj model.Visitor, err error) {
	for index := 1; index <= 100; index++ {
		data.Subject_id = index
		Db.Create(data)
		data.ID++
		continue
	}
	return Visitor_obj, nil
}

func Delete(subject_id int, proposer string) (err error) {
	if err := Db.Where("proposer = ? and subject_id = ?", proposer, subject_id).Delete(&model.Visitor{}); err.Error != nil {
		log4go.Error(err)
		return errors.New("Visitor delete db err")
	}
	return nil
}
func Delete_Retinue(subject_id int) (err error) {
	if err := Db.Where("parent_id = ?", subject_id).Delete(&model.Retinue{}); err.Error != nil {
		log4go.Error(err)
		return errors.New("Visitor delete db err")
	}
	return nil
}
func Get_Retinue(subject_id, page, size int) (result []model.Retinue, err error, count int, total int) {

	DBdate := Db
	if subject_id != 0 {
		DBdate = DBdate.Where("parent_id = ?", subject_id)
	}

	DBdate = DBdate.Order("id desc")
	if page > 0 {
		DBdate = DBdate.Limit(size).Offset((page - 1) * size)
	}
	if err := DBdate.Model(model.Retinue{}).Find(&result); err.Error != nil {
		return result, err.Error, 0, 0
	}
	DBcount := Db
	if subject_id != 0 {
		DBcount = DBcount.Where("parent_id = ?", subject_id)
	}
	// if err := DBcount.Model(model.Retinue{}).Count(&count); err.Error != nil {
	// 	return result, err.Error, 0,0
	// }
	total = zyutil.GetTotal(count, size)
	return result, nil, count, total
}

func Get(proposer string, subject_id, page, size int) (result []model.Visitor, err error, count int, total int) {

	DBdate := Db
	if subject_id != 0 {
		DBdate = DBdate.Where("subject_id = ?", subject_id)
	}
	if proposer != "" {
		DBdate = DBdate.Where("proposer = ?", proposer)
	}
	DBdate = DBdate.Order("id desc")
	if page > 0 {
		DBdate = DBdate.Limit(size).Offset((page - 1) * size)
	}
	if err := DBdate.Model(model.Visitor{}).Find(&result); err.Error != nil {
		return result, err.Error, 0, 0
	}
	DBcount := Db
	if subject_id != 0 {
		DBcount = DBcount.Where("subject_id = ?", subject_id)
	}
	if proposer != "" {
		DBcount = DBcount.Where("proposer = ?", proposer)
	}
	// if err := DBcount.Model(model.Visitor{}).Count(&count); err.Error != nil {
	// 	return result, err.Error, 0,0
	// }
	total = zyutil.GetTotal(count, size)
	return result, nil, count, total
}

//func Update(subject_id int)(result []model.Visitor,err error){
//	if err := Db.Model(&result).Where("subject_id = ?",subject_id).Update("photo","");err.Error!=nil{
//		log4go.Error(err)
//		return []model.Visitor{},errors.New("Visitor updata db err")
//	}
//	return []model.Visitor{},nil
//}

func Gettime(proposer string, subject_id int, timenow int) (result []model.Visitor, err error) {
	DBdate := Db
	if subject_id != 0 {
		DBdate = DBdate.Where("subject_id = ?", subject_id)
	}
	if proposer != "" {

		DBdate = DBdate.Where("proposer = ?", proposer)
	}
	DBdate = DBdate.Order("id desc")
	if err := DBdate.Where("end_time < ?", timenow).Find(&result); err.Error != nil {
		return result, err.Error
	}
	return result, nil
}

func GetAccount() (Account_obj model.Account, err error) {

	if err := Db.Model(&model.Account{}).First(&Account_obj); err.Error != nil {
		log4go.Error("select updated account in last time err :", err)
		return model.Account{}, err.Error
	}
	return Account_obj, nil
}

func CreateRetinue(data *model.Retinue) (Retinue_obj model.Retinue, err error) {
	if err := Db.Create(data); err.Error != nil {

		log4go.Error(err)
		return Retinue_obj, errors.New("Visitor insert db err")
	}
	if err := Db.Last(&Retinue_obj); err.Error != nil {
		log4go.Error(err)
		return Retinue_obj, errors.New("select visitor in last time err")
	}
	return Retinue_obj, nil
}

func GetRetinues(subject_id, page, size int) (result []model.Retinue, err error, count int, total int) {
	//page = 1
	//size = 10000
	DBdate := Db
	if subject_id != 0 {
		DBdate = DBdate.Where("parent_id = ?", subject_id)
	}
	DBdate = DBdate.Order("id desc")
	if page > 0 {
		DBdate = DBdate.Limit(size).Offset((page - 1) * size)
	}
	if err := DBdate.Model(model.Retinue{}).Find(&result); err.Error != nil {
		return result, err.Error, 0, 0
	}
	DBcount := Db
	if subject_id != 0 {
		DBcount = DBcount.Where("parent_id = ?", subject_id)
	}
	// if err := DBcount.Model(model.Retinue{}).Count(&count); err.Error != nil {
	// 	return result, err.Error, 0,0
	// }
	total = zyutil.GetTotal(count, size)
	return result, nil, count, total
}

func GetTemplate() (result model.Template, err error) {
	DBdate := Db

	if err := DBdate.Model(model.Template{}).Last(&result).Error; err != nil {
		log4go.Error(err.Error())
		return result, errors.New("查询信息模板失败!")
	}
	return result, nil
}

func UpdateTemplate(template_params model.Template) (err error) {

	if err := Db.Debug().Table("templates").Where("id = ?", 1).Updates(&template_params).Error; err != nil {
		log4go.Error(err.Error())
		if strings.Index(err.Error(),"not found") != -1{
			template_params.Id = 1
			if err := Db.Debug().Table("templates").Create(&template_params).Error; err != nil {
				if err != nil{
					return err
				}
			}
		}else{
			return errors.New("更新信息模板失败!")
		}
	}
	return nil
}
//验证开门  Tim add
func GetVisitor_Sign(req *model.QRCODE_VisitorReq)(open bool, err error){

	//log4go.Debug("params Proposer is:",req.Proposer)
	visitors := model.Visitor{}
	retinues := model.Retinue{}
	DBdate := Db
	DBdata_Retinue := Db
	if req.Start_time != 0 {
		DBdate = DBdate.Where("start_time >= ?", req.Start_time)
	}
	//if req.End_time != 0 {
	DBdate = DBdate.Where("end_time >= ?", time.Now().Unix())
	//}
	if req.Phone != "" {
		DBdate = DBdate.Where("phone = ?", req.Phone)
		DBdata_Retinue = DBdata_Retinue.Where("phone = ?", req.Phone)
	}
	if req.Is_auth != "" { //授权标志
		DBdate = DBdate.Where("is_auth = ?", req.Is_auth)
	}

	DBdate = DBdate.Order("id desc")

	if err := DBdate.Debug().Model(model.Visitor{}).Last(&visitors).Error;err!= nil {

		if strings.Compare(err.Error(),"record not found") == 0 {

			if err := DBdata_Retinue.Debug().Model(model.Retinue{}).Last(&retinues).Error;err != nil{

				//if strings.Compare(err.Error(),"Record not find") != -1 {
				return false, nil
				//}
			}
			return false, nil
		}
		return false, nil
	}
	return true, nil
}

func GetSubjectShip(id int) (result model.SubjectShip, err error) {

	if err := Db.Model(model.SubjectShip{}).Where("id = ?",id).Last(&result); err.Error != nil {
		return result, err.Error
	}
	return result, nil
}

func GetSubjectShipUuid(uuid string) (result model.SubjectShip, err error) {

	if err := Db.Model(model.SubjectShip{}).Where("uuid = ?",uuid).Last(&result); err.Error != nil {
		return result, err.Error
	}
	return result, nil
}

func GetPhotoShip(id int) (result model.PhotoShip, err error) {

	if err := Db.Model(model.PhotoShip{}).Where("id = ?",id).Last(&result); err.Error != nil {
		return result, err.Error
	}
	return result, nil
}

func CreateSubjectShip(data *model.SubjectShip) (SubjectShip model.SubjectShip, err error) {
	if err := Db.Create(data); err.Error != nil {
		log4go.Error(err)
		return SubjectShip, errors.New("Employee insert db err")
	}
	if err := Db.Last(&SubjectShip); err.Error != nil {
		log4go.Error(err)
		return SubjectShip, errors.New("select Employee in last time err")
	}
	return SubjectShip, nil
}

func CreatePhotoShip(data *model.PhotoShip) (PhotoShip model.PhotoShip, err error) {
	if err := Db.Create(data); err.Error != nil {
		log4go.Error(err)
		return PhotoShip, errors.New("Employee insert db err")
	}
	if err := Db.Last(&PhotoShip); err.Error != nil {
		log4go.Error(err)
		return PhotoShip, errors.New("select Employee in last time err")
	}
	return PhotoShip, nil
}


func CreateGroupShip(data *model.GroupShip) (GroupShip model.GroupShip, err error) {
	if err := Db.Create(data); err.Error != nil {
		log4go.Error(err)
		return GroupShip, errors.New("Employee insert db err")
	}
	if err := Db.Last(&GroupShip); err.Error != nil {
		log4go.Error(err)
		return GroupShip, errors.New("select Employee in last time err")
	}
	return GroupShip, nil
}

func FindSubjectShip(uuid string) (SubjectShip []model.SubjectShip, err error) {

	Dbdata := Db.Debug()

	if uuid != ""{
		Dbdata.Where("uuid = ?",uuid)
	}
	if err := Dbdata.Find(&SubjectShip); err.Error != nil {
		log4go.Error(err)
		return SubjectShip, errors.New("select all employee  err")
	}
	return SubjectShip, nil
}

func FindGroupShip(uuid string) (GroupShip []model.GroupShip, err error) {

	Dbdata := Db.Debug()

	if uuid != ""{
		Dbdata.Where("subject_uuid = ?",uuid)
	}
	if err := Dbdata.Find(&GroupShip); err.Error != nil {
		log4go.Error(err)
		return GroupShip, errors.New("select all employee  err")
	}
	return GroupShip, nil
}

func FindPhotoShip(uuid string) (PhotoShip []model.PhotoShip, err error) {

	Dbdata := Db.Debug()

	if uuid != ""{
		Dbdata.Where("subject_uuid = ?",uuid)
	}
	if err := Dbdata.Find(&PhotoShip); err.Error != nil {
		log4go.Error(err)
		return PhotoShip, errors.New("select all employee  err")
	}
	return PhotoShip, nil
}

func DeleteSubjectShip(uuid string) ( err error) {
	SubjectShip := model.SubjectShip{}
	if err := Db.Debug().Where("uuid = ?",uuid).Delete(&SubjectShip); err.Error != nil {
		log4go.Error(err)
		return  errors.New("select all employee  err")
	}
	return  nil
}

func DeleteGroupShip(uuid string) ( err error) {
	GroupShip := model.GroupShip{}
	if err := Db.Debug().Where("subject_uuid = ?",uuid).Delete(&GroupShip); err.Error != nil {
		log4go.Error(err)
		return  errors.New("select all employee  err")
	}
	return  nil
}

func DeletePhotoShip(uuid string) ( err error) {
	PhotoShip := model.PhotoShip{}
	if err := Db.Debug().Where("subject_uuid = ?",uuid).Delete(&PhotoShip); err.Error != nil {
		log4go.Error(err)
		return  errors.New("select all employee  err")
	}
	return  nil
}