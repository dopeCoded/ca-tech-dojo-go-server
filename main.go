package main

import (
    "database/sql"
    "fmt"
    "log"
    "net/http"
    "os"
    "time"

    _ "github.com/go-sql-driver/mysql"
)

var db *sql.DB
var jwtKey []byte

func main() {
    // 環境変数からデータベース接続情報を取得
    dsn := os.Getenv("DB_DSN")
    if dsn == "" {
        log.Fatal("DB_DSN 環境変数が設定されていません")
    }

    // 環境変数からJWTキーを取得
    key := os.Getenv("JWT_KEY")
    if key == "" {
        log.Fatal("JWT_KEY 環境変数が設定されていません")
    }
    jwtKey = []byte(key)

    // データベース接続の再試行
    var err error
    for i := 0; i < 5; i++ {
        db, err = sql.Open("mysql", dsn)
        if err == nil {
            if pingErr := db.Ping(); pingErr == nil {
                break
            }
        }
        log.Printf("データベースへの接続に失敗しました。再試行します... (%d/5)", i+1)
        time.Sleep(5 * time.Second)
    }

    if err != nil {
        log.Fatalf("データベースへの接続に失敗しました: %v", err)
    }
    defer db.Close()

    // ルーティング
    http.HandleFunc("/user/create", userCreateHandler)
    http.HandleFunc("/user/get", userGetHandler)
    http.HandleFunc("/user/update", userUpdateHandler)
    http.HandleFunc("/gacha/draw", gachaDrawHandler)
    http.HandleFunc("/character/list", characterListHandler)

    // サーバーの起動
    fmt.Println("Server is running on port 8080")
    if err := http.ListenAndServe(":8080", nil); err != nil {
        log.Fatal("サーバーの起動に失敗しました:", err)
    }
}
