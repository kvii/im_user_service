package main

import (
	"bytes"
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"net/http"
	"net/url"
	"strconv"
	"strings"
	"time"
)

type UserTokenRequest struct {
	Secret     string `json:"secret"`     // 必填。openIM 密钥，用于注册或者获取token时候的验证，配置在openIM的config/config.yaml中
	PlatformId int    `json:"platformID"` // 必填。平台ID
	UserId     string `json:"userID"`     // 必填。用户ID
}

type UserTokenResponse struct {
	Token             string `json:"token"`             // 获取到的用户 token
	ExpireTimeSeconds int    `json:"expireTimeSeconds"` // token 的过期时间（单位秒）
}

// 向 open im server 发送 http 获取用户 token 请求。
//
// https://doc.rentsoft.cn/restapi/authenticationManagement/getUserToken
func UserToken(ctx context.Context, arg UserTokenRequest) (UserTokenResponse, error) {
	var reply UserTokenResponse
	err := call(ctx, "/auth/user_token", "", arg, &reply)
	if err != nil {
		return UserTokenResponse{}, err
	}
	return reply, nil
}

type UserRegisterRequest struct {
	Secret string             `json:"secret"` // 必填。openIM 密钥，用于注册时候的验证，配置在 openIM 的 config/config.yaml 中
	Users  []UserRegisterUser `json:"users"`  // 必填。用户列表
}

type UserRegisterUser struct {
	UserId   string `json:"userID"`   // 必填。用户ID
	Nickname string `json:"nickname"` // 必填。用户名
	FaceUrl  string `json:"faceURL"`  // 必填。用户头像
}

type UserRegisterResponse struct {
}

// 向 open im sever 发送 http 用户注册请求。
//
// https://doc.rentsoft.cn/restapi/userManagement/userRegister
func UserRegister(ctx context.Context, arg UserRegisterRequest) (UserRegisterResponse, error) {
	var reply UserRegisterResponse
	err := call(ctx, "/user/user_register", "", arg, &reply)
	if err != nil {
		return UserRegisterResponse{}, err
	}
	return reply, nil
}

// 封装 http 请求调用逻辑
func call(ctx context.Context, api, token string, arg, reply any) error {
	// 解析 url
	u, err := url.Parse(config.ImAddr)
	if err != nil {
		return err
	}
	u = u.JoinPath(api)

	// 编码请求 json
	body, err := json.Marshal(arg)
	if err != nil {
		return err
	}

	// 构造请求
	req, err := http.NewRequestWithContext(ctx, http.MethodPost, u.String(), bytes.NewBuffer(body))
	if err != nil {
		return err
	}
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("operationID", operationId())
	if token != "" {
		req.Header.Set("token", token)
	}

	// 执行请求
	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	// 过滤异常响应
	if resp.StatusCode != http.StatusOK {
		return errors.New(resp.Status)
	}
	if !strings.HasPrefix(resp.Header.Get("Content-Type"), "application/json") {
		return errors.New("not json response content type")
	}

	// 解析响应
	data := ImResponse{
		Data: reply,
	}
	err = json.NewDecoder(resp.Body).Decode(&data)
	if err != nil {
		return err
	}
	if data.ErrCode != 0 {
		return ImError{
			ErrCode: data.ErrCode,
			ErrMsg:  data.ErrMsg,
			ErrDlt:  data.ErrDlt,
		}
	}
	return nil
}

type ImResponse struct {
	ErrCode int    `json:"errCode"` // 错误码,0表示成功
	ErrMsg  string `json:"errMsg"`  // 错误简要信息,无错误时为空
	ErrDlt  string `json:"errDlt"`  // 错误详细信息,无错误时为空
	Data    any    `json:"data"`    // 响应
}

type ImError struct {
	ErrCode int    // 错误码,0表示成功
	ErrMsg  string // 错误简要信息,无错误时为空
	ErrDlt  string // 错误详细信息,无错误时为空
}

func (e ImError) Error() string {
	return fmt.Sprintf("im error: code=%d, msg=%s, detail=%s", e.ErrCode, e.ErrMsg, e.ErrDlt)
}

func operationId() string {
	return strconv.FormatInt(time.Now().UnixNano(), 10)
}
