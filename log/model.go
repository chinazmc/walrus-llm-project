package log

type LogConfig struct {
	Level        string `mapstructure:"level"`
	Name         string `mapstructure:"name"`
	Path         string `mapstructure:"path"`
	MaxAge       int    `mapstructure:"max_age"`
	MaxSize      int    `mapstructure:"max_size"`
	MaxBackups   int    `mapstructure:"max_backups"`
	Compress     bool   `mapstructure:"compress"`
	LocalTime    bool   `mapstructure:"local_time"`
	LogInConsole bool   `mapstructure:"log_in_console"`
}
