package main

import (
	"database/sql"
	"encoding/json"
	"fmt"
	"github.com/BurntSushi/toml"
	_ "github.com/go-sql-driver/mysql"
	"github.com/satori/go.uuid"
	"log"
	"math/rand"
	"os"
	"path"
	"time"
)

// 程式設定
type appConfig struct {
	// 主機位址
	Host string
	// 連接埠
	Port int
	// 資料庫名稱
	Database string
	// 登入帳號
	User string
	// 登入密碼
	Password string
	// 主機名稱
	Source string
}

// 時間顯示格式
const timeFormat = "2006-01-02 15:04:05.000000000"

var (
	Config = new(appConfig)
)

func touch(connection *string) {
	var db *sql.DB
	if dbl, err := sql.Open("mysql", *connection); err != nil {
		log.Fatalln(err)
	} else {
		db = dbl

		// 關閉資料庫
		defer func() {
			db.Close()
		}()
	}

	// 產生欄位資料
	//id := uuid.NewV4()
	id, _ := uuid.NewV4()
	actionAt := time.Now().Unix()

	// logging before insert
	fmt.Printf("%v,Inserting data,%v,%v,,%v,\n",
		time.Now().UTC().Format(timeFormat),
		id,
		Config.Source,
		// target
		actionAt)

	if _, err := db.Exec("insert into CrossData(Id,Source,ActionAt) values(?,?,?) ", id, Config.Source, actionAt); err != nil {
		log.Fatalln(err)
	} else {
		// logging after insert
		fmt.Printf("%v,Inserted data,%v,%v,,%v,\n",
			time.Now().UTC().Format(timeFormat),
			id,
			Config.Source,
			// target
			actionAt)
	}
}

func main() {
	// 取得設定檔位置
	configPath := os.Args[1]

	// 處理設定檔路徑
	if !path.IsAbs(configPath) {
		// 取得當前執行目錄
		base, err := os.Getwd()
		if err != nil {
			log.Fatal(err)
		}
		configPath = path.Join(base, configPath)
	}

	// 載入設定
	if _, err := toml.DecodeFile(configPath, &Config); err != nil {
		panic(fmt.Sprintf("%T 讀取設定失敗 %v", err, err))
	} else {
		json.Marshal(Config)
	}

	connection := fmt.Sprintf(
		"%v:%v@tcp(%v:%v)/%v?charset=utf8",
		Config.User,
		Config.Password,
		Config.Host,
		Config.Port,
		Config.Database)

	for {
		touch(&connection)

		// 等待時間
		time.Sleep(time.Second*2 + time.Second*time.Duration(rand.Float64()*5))
	}

}
