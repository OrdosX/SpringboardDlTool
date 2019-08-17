package main

import (
	"bytes"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"net/http"
	"os"
	"strconv"
)

func main() {
	auth, host, port := initConf()
	address := "http://" + host + ":" + port
	url, path := initArg()

	_, err := os.Stat("dl")
	if os.IsNotExist(err) {
		os.Mkdir("dl", os.ModeDir)
	}

	contentLength, err := placeOrder(auth, address, url, path)
	if err != nil {
		fmt.Println(err.Error())
		os.Exit(0)
	}
	fmt.Println("服务器下载成功，等待取回")

	err = getFile(auth, address, path, contentLength)
	if err != nil {
		fmt.Println(err.Error())
		os.Exit(0)
	}
	fmt.Fprintf(os.Stdout, "任务成功完成                                                                                                          \n")
}

/*
* Sample config file:
* {
* 	auth:"fooboo!"
* 	host:"192.168.114.514"
* 	port:"810"
* }
 */

func initConf() (string, string, string) {
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
		Host string `json:"host"`
		Port string `json:"port"`
	}
	json.NewDecoder(configFile).Decode(&config)
	fmt.Println("配置文件加载完成")
	return config.Auth, config.Host, config.Port
}

/*
* Sample argument format:
* ./client.exe https://www.example.com/example.zip ex.zip
 */

func initArg() (string, string) {
	if len(os.Args) != 3 {
		fmt.Println("参数错误！")
		os.Exit(0)
	}
	return os.Args[1], os.Args[2]
}

func placeOrder(auth string, address string, url string, path string) (int64, error) {
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
	req.Auth = auth
	req.URL = url
	req.Path = path

	b := new(bytes.Buffer)
	json.NewEncoder(b).Encode(req)
	data, err := http.Post(address+"/order", "application/json", b)
	if err != nil {
		return -1, err
	}

	json.NewDecoder(data.Body).Decode(&resp)
	if resp.Code != 200 {
		return -1, errors.New(`服务端错误！\n状态码: ` + strconv.Itoa(resp.Code) + `\n报错信息： ` + resp.Message)
	}
	return resp.ContentLength, nil
}

func getFile(auth string, address string, path string, contentLength int64) error {
	var req struct {
		Auth string `json:"auth"`
		Path string `json:"path"`
	}
	req.Auth = auth
	req.Path = path
	b := new(bytes.Buffer)
	json.NewEncoder(b).Encode(req)
	data, err := http.Post(address+"/transfer", "application/json", b)
	if err != nil {
		return err
	}
	defer data.Body.Close()

	out, err := os.Create("dl/" + path)
	if err != nil {
		return err
	}
	defer out.Close()

	fileSize := int64(0)
	buffer := make([]byte, 32*1024)
	for {
		bytesRead, readErr := data.Body.Read(buffer)
		if bytesRead > 0 {
			bytesWritten, writeErr := out.Write(buffer[0:bytesRead])
			if bytesWritten > 0 {
				fileSize += int64(bytesWritten)
			}
			if writeErr != nil {
				return writeErr
			}
			if bytesRead != bytesWritten {
				return io.ErrShortWrite
			}
		}
		if readErr != nil {
			if readErr != io.EOF {
				return readErr
			}
			break
		}

		progress := fileSize * 100 / contentLength
		fmt.Fprintf(os.Stdout, "%d%% [%s]\r", progress, getS(progress, "#")+getS(100-progress, " "))
	}
	return nil
}

func getS(n int64, char string) (s string) {
	for i := int64(0); i < n; i++ {
		s += char
	}
	return
}
