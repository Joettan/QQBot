package config

type BotConfig struct {
	AppID uint64 `mapstructure:"AppID"`
	Token string `mapstructure:"Token"`
}
