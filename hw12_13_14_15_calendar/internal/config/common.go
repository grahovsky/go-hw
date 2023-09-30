package config

type Storage struct {
	Type string `mapstructure:"type" env:"STORAGE_TYPE"`
	DB   DB     `mapstructure:"db"`
}

type DB struct {
	Name     string `mapstructure:"name" env:"DB_NAME"`
	Host     string `mapstructure:"host" env:"DB_HOST"`
	User     string `mapstructure:"user" env:"DB_USER"`
	Password string `mapstructure:"password" env:"DB_PASSWORD"`
}

type RMQ struct {
	Host     string `mapstructure:"host" env:"RMQ_HOST"`
	Port     string `mapstructure:"port" env:"RMQ_PORT"`
	User     string `mapstructure:"user" env:"RMQ_USER"`
	Password string `mapstructure:"password" env:"RMQ_PASSWORD"`
	Queue    string `mapstructure:"queue" env:"RMQ_QUEUE"`
}
