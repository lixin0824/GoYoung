package model

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"os"
)

const (
	//filePath = "/storage/emulated/0/Android/data/cn.corehub.goyoung/data.json"
	filePath = "./data.json"
)

type User struct {
	UserHard     string `json:"user_hard"`
	UserAccount  string `json:"user_account"`
	PassWord     string `json:"pass_word"`
	LastLoginURL string `json:"last_login_url"`
}

func (info *User) SaveUserInfoJson() {
	data, err := json.Marshal(info)
	if err != nil {
		fmt.Println("JSON 序列化失败" + err.Error())
		return
	}
	err = ioutil.WriteFile(filePath, data, os.ModeAppend)
	if err != nil {
		fmt.Println(err.Error())
		return
	}
}

func (info *User) ReadUserInfoJson() {
	bytes, err := ioutil.ReadFile(filePath)
	if err != nil {
		fmt.Println("data.json 读取失败" + err.Error())
		return
	}
	err = json.Unmarshal(bytes, &info)
	if err != nil {
		fmt.Println("JSON 反序列化失败" + err.Error())
		return
	}
}
