package DatabaseControl

import (
	"fmt"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
)

type DatabaseControl struct {
	db *gorm.DB
}

var dBHandele DatabaseControl

func GetDBHandle() *DatabaseControl {
	return &dBHandele
}

func (DB *DatabaseControl) GormDBInit(username, password, address, dbname string, port int) error {
	var err error
	//config := ConfigManage.GetConfig()

	dsn := fmt.Sprintf("%s:%s@tcp(%s:%d)/%s?charset=utf8mb4&parseTime=True&loc=Local",
		username, password, address, port, dbname)

	DB.db, err = gorm.Open(mysql.Open(dsn), &gorm.Config{})
	if err != nil {
		fmt.Println("连接 MySQL 数据库失败:", err)
		return nil
	}

	return err
}
