package apollo

import (
	"github.com/HCY2315/chaoyue-golib/pkg/config/global"
	"github.com/spf13/viper"
)

const (
	MasterDBKey     = "master.config."
	SlaveDBKey      = "slave1.config."
	Slave1EnableKey = "slave1"
)

const (
	HostKey     = "host"
	PasswordKey = "password"
	PortKey     = "port"
	UserNameKey = "username"
	DBNameKey   = "dbname"
)

// 获取数据源数据 master和slave
// 使用slave配置需要先判断slave1Enable是否可用
func GetDataSourceDb() (master, slave global.DataBase, slave1Enable bool) {
	// master
	masterHost := viper.GetString(MasterDBKey + HostKey)
	masterPassword := viper.GetString(MasterDBKey + PasswordKey)
	masterPort := viper.GetString(MasterDBKey + PortKey)
	masterUserName := viper.GetString(MasterDBKey + UserNameKey)
	masterDBName := viper.GetString(MasterDBKey + DBNameKey)
	master = global.DataBase{
		Ip:       masterHost,
		Port:     masterPort,
		User:     masterUserName,
		PassWord: masterPassword,
		DbName:   masterDBName,
	}
	// slave1
	slaveHost := viper.GetString(SlaveDBKey + HostKey)
	slavePassword := viper.GetString(SlaveDBKey + PasswordKey)
	slavePort := viper.GetString(SlaveDBKey + PortKey)
	slaveUserName := viper.GetString(SlaveDBKey + UserNameKey)
	slaveDBName := viper.GetString(SlaveDBKey + DBNameKey)
	slave = global.DataBase{
		Ip:       slaveHost,
		Port:     slavePort,
		User:     slaveUserName,
		PassWord: slavePassword,
		DbName:   slaveDBName,
	}
	slave1Enable = viper.GetBool(Slave1EnableKey)
	return
}
