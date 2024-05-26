package config

type AppConfig struct {
	Port      int    `mapstructure:"port"`
	JwtSecret string `mapstruct:"jwt_secret"`
	BasePath  string `mapstruct:"jwt_secret"`
}
