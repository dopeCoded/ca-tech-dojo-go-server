package main

import (
    "bytes"
    "database/sql"
    "encoding/json"
    "net/http"
    "net/http/httptest"
    "os"
    "testing"

    _ "github.com/go-sql-driver/mysql"
)

var testDB *sql.DB

func TestMain(m *testing.M) {
    // テスト用データベースへの接続設定
    var err error
    dsn := "user:password@tcp(localhost:3306)/dbname?parseTime=true"
    testDB, err = sql.Open("mysql", dsn)
    if err != nil {
        panic(err)
    }
    defer testDB.Close()

    // データベースの初期化（テーブルを作成・クリア）
    setupTestDatabase()

    // グローバルなdbをテスト用に差し替え
    db = testDB

    // テストの実行
    code := m.Run()

    // テスト終了後のクリーンアップ
    teardownTestDatabase()

    os.Exit(code)
}

func setupTestDatabase() {
    // DROP TABLE IF EXISTS users;
    _, err := testDB.Exec(`DROP TABLE IF EXISTS users;`)
    if err != nil {
        panic(err)
    }

    // CREATE TABLE users ...
    _, err = testDB.Exec(`
        CREATE TABLE users (
            id INT AUTO_INCREMENT PRIMARY KEY,
            name VARCHAR(255) NOT NULL,
            created_at DATETIME DEFAULT CURRENT_TIMESTAMP,
            updated_at DATETIME DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP
        ) ENGINE=InnoDB;
    `)
    if err != nil {
        panic(err)
    }
}

func teardownTestDatabase() {
    // テスト終了後にテーブルを削除
    _, err := testDB.Exec(`DROP TABLE IF EXISTS users;`)
    if err != nil {
        panic(err)
    }
}

func TestUserCreateHandler(t *testing.T) {
    // リクエストボディの作成
    requestBody, _ := json.Marshal(UserCreateRequest{Name: "TestUser"})

    // リクエストの作成
    req, err := http.NewRequest("POST", "/user/create", bytes.NewBuffer(requestBody))
    if err != nil {
        t.Fatal(err)
    }
    req.Header.Set("Content-Type", "application/json")

    // レスポンスレコーダーの作成
    rr := httptest.NewRecorder()

    // ハンドラの呼び出し
    handler := http.HandlerFunc(userCreateHandler)
    handler.ServeHTTP(rr, req)

    // ステータスコードのチェック
    if status := rr.Code; status != http.StatusOK {
        t.Errorf("ステータスコードが異なります: got %v want %v", status, http.StatusOK)
    }

    // レスポンスボディのチェック
    var res UserCreateResponse
    err = json.Unmarshal(rr.Body.Bytes(), &res)
    if err != nil {
        t.Fatal("レスポンスのパースに失敗しました:", err)
    }

    if res.Token == "" {
        t.Error("トークンが返されていません")
    }
}

func TestUserGetHandler(t *testing.T) {
    // まずユーザを作成し、トークンを取得
    token := createTestUser(t, "TestUserGet")

    // リクエストの作成
    req, err := http.NewRequest("GET", "/user/get", nil)
    if err != nil {
        t.Fatal(err)
    }
    req.Header.Set("x-token", token)

    // レスポンスレコーダーの作成
    rr := httptest.NewRecorder()

    // ハンドラの呼び出し
    handler := http.HandlerFunc(userGetHandler)
    handler.ServeHTTP(rr, req)

    // ステータスコードのチェック
    if status := rr.Code; status != http.StatusOK {
        t.Errorf("ステータスコードが異なります: got %v want %v", status, http.StatusOK)
    }

    // レスポンスボディのチェック
    var res UserGetResponse
    err = json.Unmarshal(rr.Body.Bytes(), &res)
    if err != nil {
        t.Fatal("レスポンスのパースに失敗しました:", err)
    }

    if res.Name != "TestUserGet" {
        t.Errorf("ユーザ名が異なります: got %v want %v", res.Name, "TestUserGet")
    }
}

func TestUserUpdateHandler(t *testing.T) {
    // ユーザを作成し、トークンを取得
    token := createTestUser(t, "TestUserUpdate")

    // 更新データの作成
    updateData, _ := json.Marshal(UserUpdateRequest{Name: "UpdatedUser"})

    // リクエストの作成
    req, err := http.NewRequest("PUT", "/user/update", bytes.NewBuffer(updateData))
    if err != nil {
        t.Fatal(err)
    }
    req.Header.Set("x-token", token)
    req.Header.Set("Content-Type", "application/json")

    // レスポンスレコーダーの作成
    rr := httptest.NewRecorder()

    // ハンドラの呼び出し
    handler := http.HandlerFunc(userUpdateHandler)
    handler.ServeHTTP(rr, req)

    // ステータスコードのチェック
    if status := rr.Code; status != http.StatusOK {
        t.Errorf("ステータスコードが異なります: got %v want %v", status, http.StatusOK)
    }

    // ユーザ情報の取得と確認
    var name string
    err = testDB.QueryRow("SELECT name FROM users WHERE id = (SELECT id FROM users WHERE name = ?)", "UpdatedUser").Scan(&name)
    if err != nil {
        t.Fatal("ユーザ情報の取得に失敗しました:", err)
    }

    if name != "UpdatedUser" {
        t.Errorf("ユーザ名の更新に失敗しました: got %v want %v", name, "UpdatedUser")
    }
}

// テスト用のユーザを作成し、トークンを取得するヘルパー関数
func createTestUser(t *testing.T, name string) string {
    // リクエストボディの作成
    requestBody, _ := json.Marshal(UserCreateRequest{Name: name})

    // リクエストの作成
    req, err := http.NewRequest("POST", "/user/create", bytes.NewBuffer(requestBody))
    if err != nil {
        t.Fatal(err)
    }
    req.Header.Set("Content-Type", "application/json")

    // レスポンスレコーダーの作成
    rr := httptest.NewRecorder()

    // ハンドラの呼び出し
    handler := http.HandlerFunc(userCreateHandler)
    handler.ServeHTTP(rr, req)

    // ステータスコードのチェック
    if status := rr.Code; status != http.StatusOK {
        t.Fatalf("ユーザ作成に失敗しました: ステータスコード %v", status)
    }

    // レスポンスボディのパース
    var res UserCreateResponse
    err = json.Unmarshal(rr.Body.Bytes(), &res)
    if err != nil {
        t.Fatal("レスポンスのパースに失敗しました:", err)
    }

    if res.Token == "" {
        t.Fatal("トークンが返されていません")
    }

    return res.Token
}
