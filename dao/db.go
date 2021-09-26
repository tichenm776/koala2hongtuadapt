package dao

import (
	"errors"
	"fmt"
	"github.com/alecthomas/log4go"
	"log"
	"strconv"

	// 引入数据库驱动注册及初始化
	"gorm.io/driver/mysql"
	"gorm.io/gorm"

	"zhiyuan/koala2hongtuadapt/model"
	"zhiyuan/zyutil/config"
)

var Db *gorm.DB

var dbcount int

func New() error {
	//if runtime.GOOS == "linux" {
	//	config.Init4ini("./conf.ini")
	//	fmt.Println("linux find conf.ini")
	//} else {
	//	config.Init4ini("./conf.ini")
	//	return nil //2021.5.30临时
	//}
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
	//DBUsername, err := config.Conf4ini.String("db", "db_username")
	//DBPassword, err := config.Conf4ini.String("db", "db_password")
	//DBIP, err := config.Conf4ini.String("db", "db_ip")
	//DBPort, err := config.Conf4ini.String("db", "db_port")
	//DBName, err := config.Conf4ini.String("db", "db_name")
	//if err != nil {
	//	log4go.Error(err.Error())
	//	return errors.New(err.Error())
	//}
	var(
		err error
	)

	DBUsername := config.Gconf.DBUsername
	DBPassword := config.Gconf.DBPassword
	DBIP := config.Gconf.DBIP
	DBPort := config.Gconf.DBPort
	DBName := config.Gconf.DBName


	if DBUsername == "" || DBPassword == "" || DBIP == "" || DBPort == 0 || DBName == "" {
		errString := "缺少 connect db 所需数据！"
		//fmt.Println(err.Error())
		log4go.Error(errString)
		return errors.New(errString)
	}
	dsn := DBUsername + ":" + DBPassword + "@tcp(" + DBIP + ":" + strconv.Itoa(DBPort) + ")/" + DBName + "?parseTime=true&charset=utf8mb4"
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

	Db.AutoMigrate(&model.SubjectShip{},&model.PhotoShip{},&model.GroupShip{})
	log.Println("end createTable")
}
