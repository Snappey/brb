package config

import "github.com/spf13/viper"

type Config struct {
    AuthToken string
    AppId     string
    GuildId   string
}

func init() {
    viper.AutomaticEnv()
}

func Get() Config {
    return Config{
        AuthToken: viper.GetString("AUTH_TOKEN"),
        AppId:     viper.GetString("APP_ID"),
        GuildId:   viper.GetString("GUILD_ID"),
    }
}
