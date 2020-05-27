package initialize

import (
	"fmt"
	"gin-web/models"
	"gin-web/pkg/global"
	_ "github.com/go-sql-driver/mysql" // mysql驱动
	"github.com/jinzhu/gorm"
)

// 初始化mysql数据库
func Mysql() {
	db, err := gorm.Open("mysql", fmt.Sprintf(
		"%s:%s@tcp(%s:%d)/%s?%s",
		global.Conf.Mysql.Username,
		global.Conf.Mysql.Password,
		global.Conf.Mysql.Host,
		global.Conf.Mysql.Port,
		global.Conf.Mysql.Database,
		global.Conf.Mysql.Query,
	))
	if err != nil {
		panic(fmt.Sprintf("初始化mysql异常: %v", err))
	}
	// 打印所有执行的sql
	db.LogMode(global.Conf.Mysql.LogMode)
	global.Mysql = db
	// 表结构
	autoMigrate()
	global.Log.Debug("初始化mysql完成")
	// 初始化数据库日志监听器
	binlog()
}

// 自动迁移表结构
func autoMigrate() {
	global.Mysql.AutoMigrate(
		new(models.SysUser),
		new(models.SysRole),
		new(models.SysMenu),
		new(models.SysApi),
		new(models.SysCasbin),
	)
}

func binlog() {
	MysqlBinlog([]string{
		new(models.SysUser).TableName(),
		new(models.SysRole).TableName(),
		new(models.SysMenu).TableName(),
		new(models.SysApi).TableName(),
		new(models.SysCasbin).TableName(),
		new(models.RelationRoleMenu).TableName(),
	})
}
