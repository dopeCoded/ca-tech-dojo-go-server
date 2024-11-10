#!/bin/bash

# ユーザー作成
echo "Creating a new user..."
create_response=$(curl -s -X POST -H "Content-Type: application/json" -d '{"name": "TestUser"}' http://localhost:8080/user/create)
echo "Response from /user/create:"
echo $create_response

# 取得したJWTトークンを抽出
token=$(echo $create_response | jq -r '.token')
echo "JWT Token: $token"

# ユーザー情報の取得 (/user/get)
echo "Getting user information before update..."
get_response_before=$(curl -s -X GET -H "x-token: $token" http://localhost:8080/user/get)
echo "Response from /user/get (before update):"
echo $get_response_before

# ユーザー情報の更新 (/user/update)
echo "Updating user name to 'UpdatedUser'..."
update_response=$(curl -s -X PUT -H "Content-Type: application/json" -H "x-token: $token" -d '{"name": "UpdatedUser"}' http://localhost:8080/user/update)
echo "Response from /user/update:"
echo $update_response

# ユーザー情報の取得 (/user/get) - 更新後
echo "Getting user information after update..."
get_response_after=$(curl -s -X GET -H "x-token: $token" http://localhost:8080/user/get)
echo "Response from /user/get (after update):"
echo $get_response_after

# ガチャ実行 (/gacha/draw)
echo "Drawing gacha 3 times..."
gacha_response=$(curl -s -X POST -H "Content-Type: application/json" -H "x-token: $token" -d '{"times": 3}' http://localhost:8080/gacha/draw)
echo "Response from /gacha/draw:"
echo $gacha_response

# ユーザー所持キャラクター一覧取得 (/character/list)
echo "Listing user characters..."
character_list_response=$(curl -s -X GET -H "x-token: $token" http://localhost:8080/character/list)
echo "Response from /character/list:"
echo $character_list_response
