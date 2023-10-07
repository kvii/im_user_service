package main

import "errors"

// 用户不存在
var ErrUserNotFound = errors.New("user not found")

// 模拟 sql 用户表结构。
type userTable struct {
	id           int    // 主键
	name         string // 用户名
	faceUrl      string // 头像网址
	hasImAccount bool   // 是否有 im 账号
}

// 模拟数据库数据。可以自己添加新用户，模拟系统中已经存在的用户。
var tableUsers = []userTable{
	{1, "a", "", true},
	{2, "b", "", true},
}

// 模拟数据库查询
func findUserByName(name string) (userTable, error) {
	for _, item := range tableUsers {
		if item.name == name {
			return item, nil
		}
	}
	return userTable{}, ErrUserNotFound
}

// 模拟数据库更新
func setHasImAccount(id int) error {
	for i, item := range tableUsers {
		if item.id == id {
			tableUsers[i].hasImAccount = true
			return nil
		}
	}
	return ErrUserNotFound
}
