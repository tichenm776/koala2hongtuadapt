package dao

import (
	"database/sql"
	"errors"
	"fmt"
	"github.com/alecthomas/log4go"
	"log"
	"runtime"
	// 引入数据库驱动注册及初始化
	"gorm.io/driver/mysql"
	"gorm.io/gorm"

	"zhiyuan/koala2hongtuadapt/model"
	"zhiyuan/zyutil/config"
)

var Db *gorm.DB

var dbcount int

func New() error {
	if runtime.GOOS == "linux" {
		config.Init4ini("./conf.ini")
		fmt.Println("linux find conf.ini")
	} else {
		config.Init4ini("./conf.ini")
		return nil //2021.5.30临时
	}
	errString := Opendb() //2021.5.28 临时注释
	if errString != "" {
		return errors.New(errString)
	}
	return nil
}

func Opendb() string {
	if Db == nil {
		err := Init()
		if err != nil {
			return err.Error()
		}
		return ""
	}
	return ""
}

func Init() error {

	fmt.Println("init configer")
	DBUsername, err := config.Conf4ini.String("db", "db_username")
	DBPassword, err := config.Conf4ini.String("db", "db_password")
	DBIP, err := config.Conf4ini.String("db", "db_ip")
	DBPort, err := config.Conf4ini.String("db", "db_port")
	DBName, err := config.Conf4ini.String("db", "db_name")
	if err != nil {
		log4go.Error(err.Error())
		return errors.New(err.Error())
	}
	if DBUsername == "" || DBPassword == "" || DBIP == "" || DBPort == "" || DBName == "" {
		errString := "缺少 connect db 所需数据！"
		fmt.Println(err.Error())
		log4go.Error(errString)
		return errors.New(errString)
	}
	dsn := DBUsername + ":" + DBPassword + "@tcp(" + DBIP + ":" + DBPort + ")/" + DBName + "?parseTime=true&charset=utf8mb4"
	fmt.Println(dsn)
	Db, err = gorm.Open(mysql.Open(dsn), &gorm.Config{})
	if err != nil {
		log4go.Error("打不开的db，原因：" + err.Error())
	}
	//fmt.Println(err,Db)
	//db.SingularTable(true)
	//defer db.Close()
	//空闲连接和最大打开连接
	dbDB, err := Db.DB()
	if err != nil {
		log4go.Error("db.DB()创建失败，原因：" + err.Error())
		return err
	}
	dbDB.SetMaxIdleConns(20)
	dbDB.SetMaxOpenConns(20)

	if err := dbDB.Ping(); err != nil {
		log4go.Error("ping不通的db，原因：" + err.Error())
		return err
	}
	createTable()
	return nil
}

func createTable() {

	Db.AutoMigrate(&model.Visitor{}, &model.Retinue{}, &model.Employee{},&model.SubjectShip{},&model.PhotoShip{})
	getmigrator := Db.Migrator()
	if !getmigrator.HasTable(model.Template{}) {
		getmigrator.CreateTable(&model.Template{})
		Approval_reminder_to_employee_wechat_id := sql.NullString{String: "", Valid: true}
		Approval_reminder_to_employee_sms_id := sql.NullString{String: "", Valid: true}
		Visitor_signin_to_employee_wechat_id := sql.NullString{String: "", Valid: true}
		Visitor_signin_to_employee_sms_id := sql.NullString{String: "", Valid: true}
		Approval_pass_to_visitor_wechat_id := sql.NullString{String: "", Valid: true}
		Signin_to_visitor_wechat_id := sql.NullString{String: "", Valid: true}
		Signin_to_visitor_sms_id := sql.NullString{String: "", Valid: true}
		sta:= model.Template{
			Approval_reminder_to_employee_wechat_id:Approval_reminder_to_employee_wechat_id,
			Approval_reminder_to_employee_sms_id:Approval_reminder_to_employee_sms_id,
			Visitor_signin_to_employee_wechat_id:Visitor_signin_to_employee_wechat_id,
			Visitor_signin_to_employee_sms_id:Visitor_signin_to_employee_sms_id,
			Approval_pass_to_visitor_wechat_id:Approval_pass_to_visitor_wechat_id,
			Signin_to_visitor_wechat_id:Signin_to_visitor_wechat_id,
			Signin_to_visitor_sms_id:Signin_to_visitor_sms_id,
		}
		Db.Create(&sta)
		log4go.Info("创建信息模板表成功")
	}
	log.Println("end createTable")
}
