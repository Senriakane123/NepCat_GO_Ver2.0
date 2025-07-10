package NepCatInit

import (
	ConfigManage "NepCat_GO/ConfigModule"
	"bytes"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"strconv"
)

type ReqHandler func(ReqMessage string, Resp map[string]interface{})

var ReqApiMap = map[string]ReqHandler{}

func HttpReqInit(ReqType string, handler ReqHandler) {
	if handler != nil {
		handler = handler
	} else {
		handler = nil
	}
	ReqApiMap[ReqType] = handler
}

func HandleMessage(ReqMessage string, Resp map[string]interface{}) {
	url := "http://" + ConfigManage.GetWebConfig().NepcatInfo.LocalAddress + ":" + strconv.Itoa(ConfigManage.GetWebConfig().NepcatInfo.Port) + "/" + ReqMessage
	//url := "http://127.0.0.1:3000/" + ReqMessage
	// 将 Go 的 map 转换为 JSON
	jsonData, err := json.Marshal(Resp)
	if err != nil {
		fmt.Println("JSON 序列化失败:", err)
		return
	}

	// 创建 HTTP 请求
	req, err := http.NewRequest("POST", url, bytes.NewBuffer(jsonData))
	if err != nil {
		fmt.Println("创建请求失败:", err)
		return
	}

	// 设置请求头
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Authorization", "Bearer your_access_token") // 如果需要身份验证

	// 发送请求
	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		fmt.Println("请求失败:", err)
		return
	}
	defer resp.Body.Close()

	// 读取响应数据
	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		fmt.Println("读取响应失败:", err)
		return
	}

	// 输出 API 返回结果
	fmt.Println("状态码:", resp.Status)
	fmt.Println("响应数据:", string(body))

	//WebSocket 服务器地址（注意修改为你的 go-cqhttp 监听地址）
}
