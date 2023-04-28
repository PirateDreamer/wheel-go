package goroutinePool

const (
	defaultScalaThreshold = 1
)

type Config struct {
	//最大协程数量
	ScaleThreshold int32
}

// NewConfig 创建默认的配置
func NewConfig() *Config {
	c := &Config{
		ScaleThreshold: defaultScalaThreshold,
	}
	return c
}
