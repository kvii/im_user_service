package main

import (
	"encoding/json"
	"log/slog"
	"net/http"
	"strconv"
)

func init() {
	initConfig()
}

func main() {
	http.HandleFunc("/api/login", loginHandler)
	http.HandleFunc("/api/logout", logoutHandler)
	http.ListenAndServe(":9090", nil)
}

type LoginRequest struct {
	UserName string `json:"userName"`
	Platform int    `json:"platform"`
}

type LoginResponse struct {
	Id      int    `json:"id"`
	Name    string `json:"name"`
	FaceUrl string `json:"faceUrl"`
	ImToken string `json:"imToken"`
}

func loginHandler(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()

	// 解析参数
	var arg LoginRequest
	err := json.NewDecoder(r.Body).Decode(&arg)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	// 一个已上线的项目，登录逻辑对接 open im server 的示意

	// 1. 从数据库中查询用户
	user, err := findUserByName(arg.UserName)
	if err != nil {
		http.Error(w, err.Error(), http.StatusNotFound)
		return
	}

	// 2. 若数据库字段表明没有 im 账户就注册一个
	if !user.hasImAccount {
		logger.InfoContext(ctx, "register im account",
			slog.String("url", "/user/user_register"),
		)
		_, err := UserRegister(ctx, UserRegisterRequest{
			Secret: "openIM123",
			Users: []UserRegisterUser{{
				UserId:   strconv.Itoa(user.id),
				Nickname: user.name,
				FaceUrl:  user.faceUrl,
			}},
		})
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}

		// 3. 将数据库字段更新为“已注册”
		setHasImAccount(user.id)
	}

	// 4. 获取 im token
	logger.InfoContext(ctx, "get user token",
		slog.String("url", "/auth/user_token"),
	)
	reply, err := UserToken(ctx, UserTokenRequest{
		Secret:     "openIM123",
		PlatformId: arg.Platform,
		UserId:     strconv.Itoa(user.id),
	})
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	// 5. 响应用户信息与 im token
	data := LoginResponse{
		Id:      user.id,
		Name:    user.name,
		FaceUrl: user.faceUrl,
		ImToken: reply.Token,
	}
	renderJson(w, http.StatusOK, data)
}

type LogoutRequest struct {
	Id int `json:"id,omitempty"`
}

type LogoutResponse struct {
}

// 登出。无逻辑。
func logoutHandler(w http.ResponseWriter, r *http.Request) {
	// 解析参数
	var arg LogoutRequest
	err := json.NewDecoder(r.Body).Decode(&arg)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	// 响应
	renderJson(w, http.StatusOK, LogoutResponse{})
}

func renderJson(w http.ResponseWriter, code int, a any) error {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(code)

	bs, err := json.Marshal(a)
	if err != nil {
		return err
	}
	_, err = w.Write(bs)
	return err
}
