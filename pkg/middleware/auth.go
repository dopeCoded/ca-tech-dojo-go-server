package middleware

import (
	"context"
	"net/http"
	"strings"

	"my-go-project/pkg/auth"
)

// ContextKey はコンテキスト内で使用するキーの型です。
type ContextKey string

const (
	// UserIDKey はコンテキストに格納されるユーザーIDのキーです。
	UserIDKey ContextKey = "userID"
)

// AuthMiddleware はJWTトークンを検証し、認証されたユーザーIDをコンテキストに追加するミドルウェアです。
func AuthMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		// Authorizationヘッダーからトークンを取得
		authHeader := r.Header.Get("Authorization")
		if authHeader == "" {
			http.Error(w, "Authorization header missing", http.StatusUnauthorized)
			return
		}

		// Bearerトークンの形式を確認
		parts := strings.SplitN(authHeader, " ", 2)
		if len(parts) != 2 || strings.ToLower(parts[0]) != "bearer" {
			http.Error(w, "Authorization header format must be Bearer {token}", http.StatusUnauthorized)
			return
		}

		tokenString := parts[1]

		// JWTトークンを検証
		claims, err := auth.ValidateJWT(tokenString)
		if err != nil {
			http.Error(w, "Invalid or expired token", http.StatusUnauthorized)
			return
		}

		// ユーザーIDをコンテキストに追加
		ctx := context.WithValue(r.Context(), UserIDKey, claims.UserID)

		// 次のハンドラーにリクエストを渡す
		next.ServeHTTP(w, r.WithContext(ctx))
	})
}

// GetUserID はコンテキストからユーザーIDを取得します。
func GetUserID(ctx context.Context) (int64, bool) {
	userID, ok := ctx.Value(UserIDKey).(int64)
	return userID, ok
}
