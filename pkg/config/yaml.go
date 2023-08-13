package config

type BotConfig struct {
	AppID uint64 `mapstructure:"AppID"`
	Token string `mapstructure:"Token"`
}

type GPTConfig struct {
	GPTToken string `mapstructure:"GPTToken"`
}

type RedisConfig struct {
	Host     string `mapstructure:"Host"`
	Port     string `mapstructure:"Port"`
	Password string `mapstructure:"Password"`
}
