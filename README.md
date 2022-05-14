# chaoyue-golib

超越专用，golang lib 库

## 一、使用案例
### 1、反向代理
```go
func main() {
	r := gin.New()
	r.Use(core())
	r.GET("/api/v1/info", func(ctx *gin.Context) {
		http.RedirectHandler(ctx, urlAddress, "/api/v1")
	})
	r.Run(":8080")
}
```

### 2、连接数据库案例（mysql）
```go
package main

import (
	"fmt"

	chaoyuedb "github.com/HCY2315/chaoyue-golib/db"
	chaoyueMysql "github.com/HCY2315/chaoyue-golib/db/mysql"
	"github.com/HCY2315/chaoyue-golib/log"
	chaoyueUtils "github.com/HCY2315/chaoyue-golib/pkg/utils"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
	"gorm.io/gorm/schema"
)

type OpInterface struct {
	gorm.Model
	// models.ControlBy

	InterfaceName string `json:"interfaceName" gorm:"type:varchar(32);comment:接口名"` //
	InterfaceDesc string `json:"interfaceDesc" gorm:"type:blob;comment:接口描述"`       //
	Interface     string `json:"interface" gorm:"type:varchar(128);comment:接口"`     //
	DomainName    string `json:"domainName" gorm:"type:varchar(128);comment:域名"`
}

func (OpInterface) TableName() string {
	return "op_interface"
}

func main() {
	var chaoyueDatabase = &chaoyuedb.Database{
		Driver:          "mysql",
		Source:          "root:insur132@(127.0.0.1:3306)/md-admin?charset=utf8&parseTime=True&loc=Asia%2FShanghai&timeout=1000ms",
		MaxIdleConns:    10,
		MaxOpenConns:    10,
		ConnMaxIdleTime: 10,
		ConnMaxLifeTime: 10,
		Registers:       []chaoyuedb.DBResolverConfig{},
	}
	db := setup(chaoyueDatabase)
	var op []OpInterface
	if err := db.Model(&OpInterface{}).Scan(&op).Error; err != nil {
		fmt.Print(err)
	}
	fmt.Println(op)
}

func setup(c *chaoyuedb.Database) *gorm.DB {
	registers := make([]chaoyueMysql.ResolverConfigure, len(c.Registers))
	for i := range c.Registers {
		registers[i] = chaoyueMysql.NewResolverConfigure(
			c.Registers[i].Sources,
			c.Registers[i].Replicas,
			c.Registers[i].Policy,
			c.Registers[i].Tables)
	}
	resolverConfig := chaoyueMysql.NewConfigure(c.Source, c.MaxIdleConns, c.MaxOpenConns, c.ConnMaxIdleTime, c.ConnMaxLifeTime, registers)
	db, err := resolverConfig.Init(&gorm.Config{
		NamingStrategy: schema.NamingStrategy{
			SingularTable: true,
		},
		// DisableForeignKeyConstraintWhenMigrating: true,
		Logger: log.NewGormLogger(),
	}, mysql.Open)

	if err != nil {
		log.Debugf(chaoyueUtils.Red(c.Driver+" connect error :"), err)
	} else {
		log.Debugf(chaoyueUtils.Green(c.Driver + " connect success !"))
	}
	return db
}

```