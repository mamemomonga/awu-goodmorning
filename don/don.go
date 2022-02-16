package don

import (
	"context"
	"fmt"
	"log"

	"github.com/mattn/go-mastodon"
	"github.com/schollz/jsonstore"
)

// UserLogin ユーザのログイン情報
type UserLogin struct {
	Domain   string `yaml:"domain"`
	Email    string `yaml:"email"`
	Password string `yaml:"password"`
}

// ClientKeys クライアントキーペア
type ClientKeys struct {
	ClientID     string
	ClientSecret string
}

// Don Donオブジェクト
type Don struct {
	ctx    context.Context
	client *mastodon.Client
	cnf    Config
}

// Config Donコンストラクタ設定
type Config struct {
	Store      *jsonstore.JSONStore
	ClientName string
	UserLogin  UserLogin
}

// NewDon Donコンストラクタ
func NewDon(cnf Config) (t *Don) {
	t = new(Don)
	t.ctx = context.Background()
	t.cnf = cnf
	return t
}

// Register アプリケーション登録(true=新規登録)
func (t *Don) Register() (bool, error) {
	domain := t.cnf.UserLogin.Domain

	clientKeys := ClientKeys{}
	t.cnf.Store.Get(domain, &clientKeys)
	newReg := false

	if (clientKeys.ClientID == "") || (clientKeys.ClientSecret == "") {
		// アプリケーション登録
		app, err := mastodon.RegisterApp(t.ctx, &mastodon.AppConfig{
			Server:     fmt.Sprintf("https://%s/", domain),
			ClientName: t.cnf.ClientName,
			Scopes:     "read write follow",
		})
		if err != nil {
			return false, err
		}
		clientKeys = ClientKeys{
			ClientID:     app.ClientID,
			ClientSecret: app.ClientSecret,
		}
		t.cnf.Store.Set(domain, clientKeys)
		log.Printf("info: Register App %s", domain)
		newReg = true
	}

	// クライアント
	t.client = mastodon.NewClient(&mastodon.Config{
		Server:       fmt.Sprintf("https://%s/", domain),
		ClientID:     clientKeys.ClientID,
		ClientSecret: clientKeys.ClientSecret,
	})

	// 認証
	err := t.client.Authenticate(t.ctx, t.cnf.UserLogin.Email, t.cnf.UserLogin.Password)
	if err != nil {
		return newReg, err
	}
	log.Printf("info: Mastodon Client Start %s", domain)
	return newReg, nil
}

// Toot トゥート
func (t *Don) Toot(message string) (err error) {
	toot := mastodon.Toot{Status: message}
	status, err := t.client.PostStatus(t.ctx, &toot)
	if err != nil {
		return err
	}
	_ = status
	// spew.Dump(status)
	log.Printf("info: Toot %s", message)
	return nil
}
