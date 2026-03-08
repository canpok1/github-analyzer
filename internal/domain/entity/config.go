package entity

// Config は設定ファイルの内容を表す。
type Config struct {
	Repo          string `yaml:"repo"`
	Tone          string `yaml:"tone"`
	DefaultPrompt string `yaml:"default_prompt"`
	Model         string `yaml:"model"`
}
