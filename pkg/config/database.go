package config

type DatabaseConfig struct {
	DataSourceName  string
	Name            string
	MigrationFolder string
	MaxOpen         int
	MaxIdle         int
	MaxLifetime     int
}
