package config

// Config 定义整体配置结构
type Config struct {
	Server    ServerConfig    `mapstructure:"server"`
	Database  DatabaseConfig  `mapstructure:"database"`
	Redis     RedisConfig     `mapstructure:"redis"`
	JWT       JWTConfig       `mapstructure:"jwt"`
	MinIO     MinIOConfig     `mapstructure:"minio"`
	RateLimit RateLimitConfig `mapstructure:"rateLimit"`
	Telemetry TelemetryConfig `mapstructure:"telemetry"`
}

type RedisConfig struct {
	Host     string `mapstructure:"host"`
	Port     int    `mapstructure:"port"`
	Password string `mapstructure:"password"`
	DB       int    `mapstructure:"db"`
}

type ServerConfig struct {
	Port int `mapstructure:"port"`
}

type DatabaseConfig struct {
	DSN             string `mapstructure:"dsn"`
	MaxOpenConns    int    `mapstructure:"maxOpenConns"`
	MaxIdleConns    int    `mapstructure:"maxIdleConns"`
	ConnMaxLifetime int    `mapstructure:"connMaxLifetime"`
}

type JWTConfig struct {
	SecretKey  string `mapstructure:"secretKey"`
	ExpireTime string `mapstructure:"expireTime"`
}

type MinIOConfig struct {
	Endpoint  string `mapstructure:"endpoint"`
	AccessKey string `mapstructure:"accessKey"`
	SecretKey string `mapstructure:"secretKey"`
	Bucket    string `mapstructure:"bucket"`
	UseSSL    bool   `mapstructure:"useSSL"`
}

type TelemetryConfig struct {
	ZipkinEndpoint string `mapstructure:"zipkinEndpoint"`
	ServiceName    string `mapstructure:"serviceName"`
}

// RateLimitConfig 限流配置
type RateLimitConfig struct {
	Enable  bool   `mapstructure:"enable"`
	Method  string `mapstructure:"method"`  // "token_bucket", "sliding_window"
	QPS     int    `mapstructure:"qps"`     // 每秒请求数
	Burst   int    `mapstructure:"burst"`   // 突发容量
	KeyFunc string `mapstructure:"keyFunc"` // 限流Key来源: "ip"(默认), "user_id"
}

var GlobalConfig *Config

// InitConfig 保留给 main() 调用，实际由 database.DB() 的 ensureConfig() 懒加载
func InitConfig() error {
	return nil
}
