package entity

// MockConfig はモック設定を表す。
type MockConfig struct {
	AI         bool `yaml:"ai"`
	Repository bool `yaml:"repository"`
}

// Config は設定ファイルの内容を表す。
type Config struct {
	Repo          string     `yaml:"repo"`
	Tone          string     `yaml:"tone"`
	DefaultPrompt string     `yaml:"default_prompt"`
	Model         string     `yaml:"model"`
	Mock          MockConfig `yaml:"mock"`
	LogFile       string     `yaml:"log_file"`
}
