package config

//Composing structs for unmarshalling. Writer is lumberjack's setup struct.
//It's annotated for JSON out-of-the-box.
//Logrus is for logging level and log output settings.

type AppenderType string

const (
	RollingFile AppenderType = "rolling_file"
	Console                  = "console"
)

/*type GlobalConfig struct {
	Logging Logging `json:"logging"`
}*/

type Destination struct {
	File       AppenderType `json:"file"`
	BufferSize int          `json:"buffer_size"`
	MaxSize    int          `json:"max_size"`
	MaxAge     int          `json:"max_age"`
	MaxBackups int          `json:"max_backups"`
	LocalTime  bool         `json:"local_time"`
}

type Appender struct {
	Type        AppenderType `json:"type"`
	Enabled     bool         `json:"enabled"`
	Level       *string      `json:"level"`
	Destination *Destination `json:"destination"`
}

type Logging struct {
	Level          *string    `json:"level"`
	Appenders      []Appender `json:"appenders"`
	ConsoleOutHTTP bool       `json:"console_out_http"`
	LogHTTP        bool       `json:"http_log"`
}
