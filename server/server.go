package main

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"os"
)

var auth, port string

func main() {
	auth, port = initConf()

	http.HandleFunc("/order", handleOrder)
	http.HandleFunc("/transfer", handleTransfer)
	http.ListenAndServe(":"+port, nil)
}

func handleOrder(w http.ResponseWriter, r *http.Request) {
	var req struct {
		Auth string `json:"auth"`
		URL  string `json:"url"`
		Path string `json:"path"`
	}
	var resp struct {
		Code          int    `json:"code"`
		Message       string `json:"msg"`
		ContentLength int64  `json:"contentLength"`
	}
	resp.ContentLength = 0
	err := json.NewDecoder(r.Body).Decode(&req)
	if err != nil {
		resp.Code = 400
		resp.Message = err.Error()
		fmt.Println(err.Error())
		json.NewEncoder(w).Encode(resp)
		return
	}
	if req.Auth != auth {
		resp.Code = 401
		resp.Message = "无权限"
		fmt.Println("[" + r.RemoteAddr + "]: 无权限")
		json.NewEncoder(w).Encode(resp)
		return
	}

	out, err := os.Create("dl/" + req.Path)
	if err != nil {
		resp.Code = 500
		resp.Message = err.Error()
		fmt.Println(err.Error())
		json.NewEncoder(w).Encode(resp)
		return
	}
	defer out.Close()

	data, err := http.Get(req.URL)
	if err != nil {
		resp.Code = 400
		resp.Message = err.Error()
		fmt.Println(err.Error())
		json.NewEncoder(w).Encode(resp)
		return
	}
	defer data.Body.Close()

	_, err = io.Copy(out, data.Body)
	if err != nil {
		resp.Code = 400
		resp.Message = err.Error()
		fmt.Println(err.Error())
		json.NewEncoder(w).Encode(resp)
	} else {
		resp.Code = 200
		resp.Message = "成功"
		fileInfo, _ := os.Stat("dl/" + req.Path)
		resp.ContentLength = fileInfo.Size()
		fmt.Println("[" + r.RemoteAddr + "]: 下载" + req.Path + "成功")
		json.NewEncoder(w).Encode(resp)
	}

}

func handleTransfer(w http.ResponseWriter, r *http.Request) {
	var req struct {
		Auth string `json:"auth"`
		Path string `json:"path"`
	}
	var resp struct {
		Code    int    `json:"code"`
		Message string `json:"msg"`
	}
	err := json.NewDecoder(r.Body).Decode(&req)
	if err != nil {
		resp.Code = 400
		resp.Message = err.Error()
		fmt.Println(err.Error())
		json.NewEncoder(w).Encode(resp)
		return
	}
	if req.Auth != auth {
		resp.Code = 401
		resp.Message = "无权限"
		fmt.Println("[" + r.RemoteAddr + "]: 无权限")
		json.NewEncoder(w).Encode(resp)
		return
	}

	file, err := os.Open("dl/" + req.Path)
	if err != nil {
		resp.Code = 500
		resp.Message = err.Error()
		fmt.Println(err.Error())
		json.NewEncoder(w).Encode(resp)
		return
	}
	defer file.Close()

	_, err = io.Copy(w, file)
	if err != nil {
		fmt.Println(err.Error())
	}
}

/*
* Sample config file:
* {
* 	auth:"fooboo!"
* 	port:"810"
* }
 */

func initConf() (string, string) {
	configFile, err := os.Open("config.json")
	if err != nil {
		if os.IsNotExist(err) {
			fmt.Println("找不到配置文件！")
			os.Exit(0)
		} else {
			fmt.Println(err.Error())
			os.Exit(0)
		}
	}
	defer configFile.Close()

	var config struct {
		Auth string `json:"auth"`
		Port string `json:"port"`
	}
	json.NewDecoder(configFile).Decode(&config)
	fmt.Println("配置文件加载完成")
	return config.Auth, config.Port
}
