package config

import (
	"os"
	"path/filepath"

	"github.com/canpok1/github-analyzer/internal/domain/entity"
	"gopkg.in/yaml.v3"
)

const configFileName = ".github-analyzer.yaml"

// Load はホームディレクトリの設定ファイル (~/.github-analyzer.yaml) を読み込む。
// ファイルが存在しない場合はゼロ値のConfigを返す。
func Load() (entity.Config, error) {
	home, err := os.UserHomeDir()
	if err != nil {
		return entity.Config{}, nil
	}
	return LoadFromPath(filepath.Join(home, configFileName))
}

// LoadFromPath は指定パスからYAML設定ファイルを読み込む。
// ファイルが存在しない場合はゼロ値のConfigを返す。
func LoadFromPath(path string) (entity.Config, error) {
	data, err := os.ReadFile(path)
	if err != nil {
		if os.IsNotExist(err) {
			return entity.Config{}, nil
		}
		return entity.Config{}, err
	}

	var cfg entity.Config
	if err := yaml.Unmarshal(data, &cfg); err != nil {
		return entity.Config{}, err
	}

	return cfg, nil
}
