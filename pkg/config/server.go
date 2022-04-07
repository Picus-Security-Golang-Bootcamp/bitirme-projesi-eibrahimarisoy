package config

type ServerConfig struct {
	AppVersion       string
	Mode             string
	RoutePrefix      string
	Debug            bool
	Port             string
	TimeoutSecs      int64
	ReadTimeoutSecs  int64
	WriteTimeoutSecs int64
}
