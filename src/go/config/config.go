package config

import (
	"github.com/dyslexic-charactersheets/translator/src/go/log"
	// "fmt"
	toml "github.com/BurntSushi/toml"
	"io/ioutil"
	// "os"
	"database/sql"
	"strconv"
	// _ "github.com/ziutek/mymysql/mysql"
	_ "github.com/ziutek/mymysql/godrv"
)

var Config Configuration

// the types

type Configuration struct {
	Debug       int `toml:"debug"`
	Fail        bool
	Server      serverConfig   `toml:"server"`
	PDF         pdfConfig      `toml:"pdf"`
	Database    databaseConfig `toml:"db"`
	OldDatabase databaseConfig `toml:"old_db"`
	// Github      githubConfig   `toml:"github"`
	Dev         devConfig      `toml:"dev"`
	Live        liveConfig     `toml:"live"`
	Mail        mailConfig     `toml:"mail"`
	Partial		bool           `toml:"partial"`
}

type serverConfig struct {
	Hostname string `toml:"hostname"`
	Port     int    `toml:"port"`
	SSL      bool   `toml:"ssl"`
}

type pdfConfig struct {
	Path string `toml:"path"`
}

type databaseConfig struct {
	Hostname string `toml:"host"`
	Database string `toml:"db"`
	Username string `toml:"user"`
	Password string `toml:"password"`
}

// type githubConfig struct {
// 	AccessToken string `toml:"access_token"`
// }

type devConfig struct {
	DevLoginURL string `toml:"dev_url"`
	SharedSecret string `toml:"shared_secret"`
}

type liveConfig struct {
	LiveLoginURL string `toml:"live_url"`
	SharedSecret string `toml:"shared_secret"`
}

type mailConfig struct {
	Hostname string `toml:"host"`
	Username string `toml:"user"`
	Password string `toml:"password"`
	From string `toml:"from"`
	UseStartTLS bool `toml:"use_start_tls"`
	UseAuth bool `toml:"use_auth"`
}

// loading the config

func init() {
	LoadConfig(true)
}

func LoadConfig(initial bool) {
	config := Configuration{
		Debug: 0,
		Fail:  false,
		Server: serverConfig{
			Hostname: "",
			Port:     9091,
		},
		Mail: mailConfig{
			Hostname: "localhost:25",
			Username: "",
			Password: "",
			From: "noreply@dyslexic-charactersheets.com",
			UseStartTLS: true,
			UseAuth: false,
		},
	}
	if initial {
		Config = config
	}

	configData, err := ioutil.ReadFile("dist/conf/config.toml")
	if err != nil {
		log.Error("config", "Error opening config.toml:", err)
		Config.Fail = true
		return
	}
	if _, err := toml.Decode(string(configData), &config); err != nil {
		// handle error
		log.Error("config", "Error decoding config.toml:", err)
		Config.Fail = true
		return
	}

	config.Fail = false

	// if that worked, swap the config for the new one
	Config = config
	if Config.Debug > 0 {
		DebugConfig()
	}
}

func DebugConfig() {
	log.Log("config", "Config", Config)
}

func (server *serverConfig) Host() string {
	return server.Hostname + ":" + strconv.Itoa(server.Port)
}

func (db *databaseConfig) Open() (*sql.DB, error) {
	conn := db.Database + "/" + db.Username + "/" + db.Password
	if db.Hostname != "localhost" && db.Hostname != "" {
		conn = "tcp:" + db.Hostname + "*" + conn
	}
	if Config.Debug > 0 {
		log.Log("config", "Connecting to", conn)
	}
	sqldb, err := sql.Open("mymysql", conn)
	return sqldb, err
}
