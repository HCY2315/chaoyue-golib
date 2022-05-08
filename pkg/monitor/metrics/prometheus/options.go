package prometheus

var (
	// DefaultName 默认服务名
	DefaultName = "ccp.microservice.mockservice"
	// DefaultPromePath 默认拉取监控的接口
	DefaultPromePath = CSTDefaultRouter
	// DefaultRegistry 默认服务发现为kubernetes
	DefaultRegistry = "consul"
	// DefaultRegistryAddress 默认服务发现地址为空
	DefaultRegistryAddress = make([]string, 0, 10)
	// DefaultPort 默认exporter端口
	DefaultPort = 50000
)

type Option func(*Options)

// Options prometheus config
type Options struct {
	// PromePath 拉取监控指标接口
	PromePath string // prometheus拉取的路径
	// Port prometheus exporter监听端口
	Port int
	// Registry 值为consul和kubernetes，其它panic
	Registry string
	// Registry 服务发现地址，当服务发现类型为kubernetes时此选项忽略
	RegistryAddress []string
	// Id 在consul中注册的服务服务对应的实例${Name}-uuid
	Id string
	// Name 在consul中注册的服务名
	Name string
}

func newOptions(opts ...Option) Options {
	opt := Options{
		Name:            DefaultName,
		PromePath:       DefaultPromePath,
		Registry:        DefaultRegistry,
		RegistryAddress: DefaultRegistryAddress,
		Port:            DefaultPort,
	}
	for _, o := range opts {
		o(&opt)
	}
	return opt
}

func PromePath(p string) Option {
	return func(o *Options) {
		o.PromePath = p
	}
}

func Registry(r string) Option {
	return func(o *Options) {
		o.Registry = r
	}
}

func RegistryAddress(a []string) Option {
	return func(o *Options) {
		o.RegistryAddress = a
	}
}

func Port(p int) Option {
	return func(o *Options) {
		o.Port = p
	}
}

func ID(id string) Option {
	return func(o *Options) {
		o.Id = id
	}
}

func Name(n string) Option {
	return func(o *Options) {
		o.Name = n
	}
}
