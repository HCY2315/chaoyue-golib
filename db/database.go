package db

type Database struct {
	// 数据库类型 mysql，sqlite3， postgres
	// driver: mysql
	Driver string

	// 数据库连接字符串 mysql 缺省信息 charset=utf8&parseTime=True&loc=Local&timeout=1000ms
	Source          string
	ConnMaxIdleTime int
	ConnMaxLifeTime int
	MaxIdleConns    int
	MaxOpenConns    int
	Registers       []DBResolverConfig
}

type DBResolverConfig struct {
	Sources  []string
	Replicas []string
	Policy   string
	Tables   []string
}

var (
	DatabaseConfig  = new(Database)
	DatabasesConfig = make(map[string]*Database)
)
