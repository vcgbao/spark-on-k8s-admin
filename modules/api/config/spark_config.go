package config

type SparkConfig struct {
	Namespace         string
	SparkUIServiceUrl string `yaml:"sparkUiServiceUrl" mapstructure:"sparkUiServiceUrl"`
	ModifyRedirectUrl bool   `yaml:"modifyRedirectUrl" mapstructure:"sparkUiServiceUrl"`
}
