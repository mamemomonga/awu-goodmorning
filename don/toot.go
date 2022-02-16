package don

import (
	"io/ioutil"
	"log"
	"os"

	"github.com/schollz/jsonstore"
	"gopkg.in/yaml.v2"
)

// UserConfig
type UserConfig struct {
	ClientName string    `yaml:"client_name"`
	Mastodon   UserLogin `yaml:"mastodon"`
}

// 設定ファイルの読み込み
func userConfigLoad(filename string) (cnf UserConfig, err error) {
	buf, err := ioutil.ReadFile(filename)
	if err != nil {
		return cnf, err
	}
	err = yaml.Unmarshal(buf, &cnf)
	if err != nil {
		return cnf, err
	}
	log.Printf("trace: Load %s", filename)
	return cnf, nil
}

func TootIt(message string) (err error) {
	// カレントディレクトリの mastodon.yaml を読込
	cfg, err := userConfigLoad("mastodon.yaml")
	if err != nil {
		return err
	}

	// serversファイルのロードもしくは新規作成
	serversFilename := "mastodon.json"
	var servers *jsonstore.JSONStore
	if _, err := os.Stat(serversFilename); err == nil {
		servers, err = jsonstore.Open(serversFilename)
		if err != nil {
			return err
		}
	} else {
		servers = new(jsonstore.JSONStore)
	}

	// donインスタンス作成
	don := NewDon(Config{
		ClientName: cfg.ClientName,
		UserLogin:  cfg.Mastodon,
		Store:      servers,
	})

	// アプリ登録
	if save, err := don.Register(); err == nil {
		if save {
			jsonstore.Save(servers, serversFilename)
		}
	} else {
		return err
	}

	// トゥート
	don.Toot(message)
	return nil
}
