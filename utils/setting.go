package utils

import (
	"fmt"
	"gopkg.in/ini.v1"
)

var (
	HttpPort	string

	Db			string
	DbHost		string
	DbPort		string
	DbUser		string
	DbPassWord 	string
	DbName		string

	RoleUser		int
	RoleAdmin		int
	RoleSuperAdmin	int
)

func init() {
	file , err := ini.Load("config/config.ini")
	if err != nil {
		fmt.Println("配置文件读取错误，请检查文件路径：", err)
	}
	LoadServer(file)
	LoadDateBase(file)
	LoadUser(file)
}

// 抽离参数

func LoadServer(file *ini.File) {
	HttpPort = file.Section("server").Key("HttpPort").MustString(":2333")
}

func LoadDateBase(file *ini.File) {
	Db = file.Section("database").Key("Db").MustString("mysql")
	DbHost = file.Section("database").Key("DbHost").MustString("gz-cdb-8eb5lnxf.sql.tencentcdb.com")
	DbPort = file.Section("database").Key("DbPort").MustString("58893")
	DbUser = file.Section("database").Key("DbUser").MustString("smsklxq1")
	DbPassWord = file.Section("database").Key("DbPassWord").MustString("smsklxq1")
	DbName = file.Section("database").Key("DbName").MustString("smsklxq1")
}

func LoadUser(file *ini.File) {
	RoleUser = file.Section("user").Key("RoleUser").MustInt(1)
	RoleAdmin = file.Section("user").Key("RoleAdmin").MustInt(2)
	RoleSuperAdmin = file.Section("user").Key("RoleSuperAdmin").MustInt(3)
}