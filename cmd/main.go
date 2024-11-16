package main

import (
	"database/sql"
	"log"
	"net/http"
	"os"

	"my-go-project/internal/handler"
	"my-go-project/internal/repository"
	"my-go-project/internal/service"
	"my-go-project/pkg/auth"
	"my-go-project/pkg/middleware"

	_ "github.com/go-sql-driver/mysql"
)

func main() {
	// 環境変数からJWTKeyを取得
	jwtKey := os.Getenv("JWT_KEY")
	if jwtKey == "" {
		log.Fatal("JWT_KEY environment variable is not set")
	}
	// authパッケージにJWTKeyを設定
	auth.JWTKey = []byte(jwtKey)

	// データベース接続
	db, err := sql.Open("mysql", os.Getenv("DB_DSN"))
	if err != nil {
		log.Fatalf("Failed to connect to database: %v", err)
	}
	defer db.Close()

	// リポジトリの初期化
	userRepo := repository.NewUserRepository(db)
	gachaRepo := repository.NewGachaRepository(db)

	// サービスの初期化
	userService := service.NewUserService(userRepo)
	gachaService := service.NewGachaService(gachaRepo)

	// ハンドラーの初期化
	userHandler := handler.NewUserHandler(userService)
	gachaHandler := handler.NewGachaHandler(gachaService)

	// ルーターの設定
	mux := http.NewServeMux()

	// 認証不要なルート
	mux.HandleFunc("/user/create", userHandler.CreateUser)

	// 認証が必要なルート
	authenticatedMux := http.NewServeMux()
	authenticatedMux.HandleFunc("/user/get", userHandler.GetUser)
	authenticatedMux.HandleFunc("/user/update", userHandler.UpdateUser)
	authenticatedMux.HandleFunc("/gacha/draw", gachaHandler.DrawGacha)
	authenticatedMux.HandleFunc("/character/list", gachaHandler.ListCharacters)

	// ミドルウェアを適用
	mux.Handle("/auth/", middleware.AuthMiddleware(authenticatedMux))

	// サーバーの起動
	log.Println("Server is running on port 8080")
	if err := http.ListenAndServe(":8080", mux); err != nil {
		log.Fatalf("Failed to start server: %v", err)
	}
}
