# chaoyue-golib

超越专用，golang lib 库

## 一、使用案例

### 6、线程池
```go
package main

import (
	"fmt"

	"github.com/HCY2315/chaoyue-golib/pkg/threads"
)

type name struct {
	i int
}

func main() {
	wn := threads.CreateWorker(10)
	for i := 1; i <= 100; i++ {
		n := new(name)
		n.i = i
		wn.Run(func() {
			fmt.Printf("[%d]:hello threads pool\n", n.i)
		})
	}
	wn.Wait()
}
```

### 5、手机验证码
```go
// 阿里云的
package main

import (
	"context"
	"fmt"
	"net/http"
	"time"

	gin2 "github.com/HCY2315/chaoyue-golib/gin"
	"github.com/HCY2315/chaoyue-golib/log"
	"github.com/HCY2315/chaoyue-golib/pkg/captcha"
	"github.com/HCY2315/chaoyue-golib/pkg/errors"
	"github.com/HCY2315/chaoyue-golib/pkg/limit"
	"github.com/HCY2315/chaoyue-golib/pkg/thirdparty"
	"github.com/gin-gonic/gin"
	"github.com/go-redis/redis/v8"
)

type VerifyController struct {
	smsService captcha.SMSVerifyService
	smsLimiter limit.FixedDurationLimiter
}

func NewVerifyController(smsService captcha.SMSVerifyService, smsLimiter limit.FixedDurationLimiter) *VerifyController {
	return &VerifyController{
		smsService: smsService,
		smsLimiter: smsLimiter,
	}
}

type CreateCaptchaReqVO struct {
	Tag string
	ID  string
}

type GetSMSCaptchaReqVO struct {
	Phone string
}

type GetSMSCaptchaRespVO struct {
	CaptchaID        string `json:"captchaID,omitempty"`
	CaptchaBase64Img string `json:"captchaBase64Img,omitempty"`
}

// @Summary 获取短信验证码
// @Tags 验证码
func (v *VerifyController) GetSMSCode(c *gin.Context) {
	var req GetSMSCaptchaReqVO
	if errBind := c.ShouldBindJSON(&req); errBind != nil {
		c.Error(errors.Wrap(errors.ErrBadRequest, errBind.Error()))
		return
	}
	//todo validate
	// 检查短信频次
	allowSMS, nextRetryTime, errCheckSMSLimit := v.smsLimiter.TryAccess(req.Phone)
	if errCheckSMSLimit != nil {
		log.Errorf("检查SMS频次出错:%s，continue", errCheckSMSLimit.Error())
		//continue
	} else if !allowSMS {
		c.JSON(429, gin2.NewGeneralVO(0, fmt.Sprintf("调用过于频繁，retry at %v", nextRetryTime)))
		return
	}
	// 生成短信验证码
	if errSendSMSCode := v.smsService.SendCodeToPhone(req.Phone); errSendSMSCode != nil {
		c.Error(errors.Wrap(errSendSMSCode, "send sms code to %s", req.Phone))
		return
	}
	var resp GetSMSCaptchaRespVO
	// 没什么好返回的
	c.JSON(http.StatusOK, resp)
}

var (
	// 超时时间
	smsCodeTTLInSeconds int    = 300
	// 五分钟之内容错
	smsCodeLimitMinutes int64  = 5
	// 五分钟之类只允许点击五次
	smsCodeLimitUpper   int64  = 5
	domain              string = ""
	accessKeyID         string = ""
	accessKeySecret     string = ""
	sigName             string = "短信名字"
	tempCode            string = ""
)

func buildRedis() (*redis.Client, error) {
	var opts redis.Options
	opts.Addr = "127.0.0.1:6379"
	opts.Password = ""
	opts.DB = 1
	cli := redis.NewClient(&opts)
	return cli, cli.Ping(context.Background()).Err()
}

func main() {
	redisCli, err := buildRedis()
	if err != nil {
		fmt.Println(err.Error())
		return
	}
	aliSMSSender, errBuildAliSMSSender := thirdparty.NewAliSMSService(domain, accessKeyID,
		accessKeySecret, sigName, tempCode)
	if errBuildAliSMSSender != nil {
		fmt.Println("build sms sender")
		return
	}
	smsService := captcha.NewSmsVerifyImpl("chaoyue:scode", SMSVerifyCodeTTL(), aliSMSSender, redisCli)
	smsLimiter := limit.NewRedisFixedDurationLimiter(redisCli,
		limit.BuildFixedMinuteDurationResolve(smsCodeLimitMinutes),
		smsCodeLimitUpper)

	verifyCtl := NewVerifyController(smsService, smsLimiter)
	r := gin.New()
	r.GET("/api/v1/verify/duanx", verifyCtl.GetSMSCode)
	r.POST("/api/v1/zhengming", func(ctx *gin.Context) {
		phone := ctx.DefaultPostForm("phone", "")
		verifyCode := ctx.DefaultPostForm("code", "")
		match, errVerify := verifyCtl.smsService.VerifyMatch(phone, verifyCode)
		if errVerify != nil {
			ctx.Error(errVerify)
			return
		}
		if !match {
			ctx.Error(errors.Wrap(errors.ErrBadRequest, "手机号和验证码不匹配"))
			ctx.JSON(404, "手机号和验证码不匹配")
			return
		}
		ctx.JSON(200, "success")
	})
	r.Run(":8080")
}

func SMSVerifyCodeTTL() time.Duration {
	return time.Duration(smsCodeTTLInSeconds) * time.Second
}
```

### 4、go prometheus 使用案例
```go
func main() {
	ghf := chaoyueMiddleware.ExcludeByPath("/_metrics")
	prometheusMiddleware := chaoyueMiddleware.NewPrometheusExporter("chaoyue", "127.0.0.1", "127.0.0.1", ghf)
	r.Use(prometheusMiddleware.HandleFunc)
	r.GET("/info", func(c *gin.Context) {
		c.JSON(200, gin.H{
			"message": "ping",
		})
	})
	r.GET("/metrics", prometheusMiddleware.ExportMetricsHandler())
	r.Run(":8490")
}
```

### 3、GORM 使用案例
```go
package main

import (
	"fmt"

	chaoyuedb "github.com/HCY2315/chaoyue-golib/db"
	chaoyueMysql "github.com/HCY2315/chaoyue-golib/db/mysql"
	chaoyueGorm "github.com/HCY2315/chaoyue-golib/gorm"
	"github.com/HCY2315/chaoyue-golib/log"
	chaoyueUtils "github.com/HCY2315/chaoyue-golib/pkg/utils"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
	"gorm.io/gorm/schema"
)

type OpInterface struct {
	Id            string
	InterfaceName string `json:"interfaceName" gorm:"type:varchar(32);comment:接口名"` //
	InterfaceDesc string `json:"interfaceDesc" gorm:"type:blob;comment:接口描述"`       //
	Interface     string `json:"interface" gorm:"type:varchar(128);comment:接口"`     //
	DomainName    string `json:"domainName" gorm:"type:varchar(128);comment:域名"`
}

func (OpInterface) TableName() string {
	return "op_interface"
}

type OpCmdb struct {
	Id string
	Ip string `json:"ip"` //
}

func (OpCmdb) TableName() string {
	return "op_interface"
}

type BothTable struct {
	OpCmdb
	OpInterface
}

func main() {
	// 2、
	fTN := "op_interface"
	sRN := "op_cmdb"
	var op []BothTable
	// SELECT * FROM `op_interface` INNER JOIN `op_cmdb` ON  `op_interface`.`id`=`op_cmdb`.`id`  WHERE op_interface.deleted_at is NULL
	a := chaoyueGorm.NewClauses(1).
		Select(*chaoyueGorm.TableProperties(fTN)).
		FromJoin(chaoyueGorm.JoinBy(fTN, "id", sRN, "id")).
		IsNull(fTN + ".deleted_at").
		Export()
	if err := DBM.Model(&OpInterface{}).Clauses(a...).Debug().Scan(&op).Error; err != nil {
		fmt.Println(err)
	}
	fmt.Println(op)
}

var (
	ChaoyueDatabase = &chaoyuedb.Database{
		Driver:          "mysql",
		Source:          "root:insur132@(127.0.0.1:3306)/md-admin?charset=utf8&parseTime=True&loc=Asia%2FShanghai&timeout=1000ms",
		MaxIdleConns:    10,
		MaxOpenConns:    10,
		ConnMaxIdleTime: 10,
		ConnMaxLifeTime: 10,
		Registers:       []chaoyuedb.DBResolverConfig{},
	}
	DBM *gorm.DB = setup(ChaoyueDatabase)
)

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
	DBM, err := resolverConfig.Init(&gorm.Config{
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
	return DBM
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