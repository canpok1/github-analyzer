package config

import (
	"os"
	"path/filepath"

	"github.com/canpok1/github-analyzer/internal/domain/entity"
	"gopkg.in/yaml.v3"
)

const configFileName = ".github-analyzer.yaml"

// Load はホームディレクトリとカレントディレクトリの設定ファイルを読み込み、
// フィールド単位でマージする。カレントディレクトリの値がホームディレクトリの値を上書きする。
// どちらのファイルも存在しない場合はゼロ値のConfigを返す。
func Load() (entity.Config, error) {
	var paths []string

	if home, err := os.UserHomeDir(); err == nil {
		paths = append(paths, filepath.Join(home, configFileName))
	}

	if cwd, err := os.Getwd(); err == nil {
		paths = append(paths, filepath.Join(cwd, configFileName))
	}

	return loadAndMerge(paths)
}

// loadAndMerge は指定されたパスのリストから設定ファイルを順番に読み込み、
// フィールド単位でマージした結果を返す。後のパスの値が前の値を上書きする。
func loadAndMerge(paths []string) (entity.Config, error) {
	merged := map[string]interface{}{}

	for _, p := range paths {
		raw, err := loadAsMap(p)
		if err != nil {
			return entity.Config{}, err
		}
		deepMerge(merged, raw)
	}

	out, err := yaml.Marshal(merged)
	if err != nil {
		return entity.Config{}, err
	}

	var cfg entity.Config
	if err := yaml.Unmarshal(out, &cfg); err != nil {
		return entity.Config{}, err
	}

	return cfg, nil
}

// loadAsMap はYAMLファイルをmap[string]interface{}として読み込む。
// ファイルが存在しない場合は空のmapを返す。
func loadAsMap(path string) (map[string]interface{}, error) {
	data, err := os.ReadFile(path)
	if err != nil {
		if os.IsNotExist(err) {
			return map[string]interface{}{}, nil
		}
		return nil, err
	}

	var m map[string]interface{}
	if err := yaml.Unmarshal(data, &m); err != nil {
		return nil, err
	}
	if m == nil {
		return map[string]interface{}{}, nil
	}

	return m, nil
}

// deepMerge はsrcのキーをdstにマージする。ネストされたmapはフィールド単位でマージされる。
func deepMerge(dst, src map[string]interface{}) {
	for k, sv := range src {
		dv, exists := dst[k]
		if exists {
			srcMap, srcOK := sv.(map[string]interface{})
			dstMap, dstOK := dv.(map[string]interface{})
			if srcOK && dstOK {
				deepMerge(dstMap, srcMap)
				continue
			}
		}
		dst[k] = sv
	}
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
