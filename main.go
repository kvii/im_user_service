package main

import (
	"encoding/json"
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

	// 查询用户
	user, err := findUserByName(arg.UserName)
	if err != nil {
		http.Error(w, err.Error(), http.StatusNotFound)
		return
	}

	// 若没有 im 账户就注册一个
	if !user.hasImAccount {
		logger.InfoContext(ctx, "register im account")
		_, err := UserRegister(ctx, UserRegisterRequest{
			Secret: "openIM123",
			Users: []UserRegisterUser{{
				UserId:   strconv.Itoa(user.id),
				Nickname: user.name,
				FaceUrl:  "",
			}},
		})
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
	}
	setImAccount(user.id)

	// 获取 im token
	logger.InfoContext(ctx, "get user token")
	reply, err := UserToken(ctx, UserTokenRequest{
		Secret:     "openIM123",
		PlatformId: arg.Platform,
		UserId:     strconv.Itoa(user.id),
	})
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	// 响应
	data := LoginResponse{
		Id:      user.id,
		Name:    user.name,
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
