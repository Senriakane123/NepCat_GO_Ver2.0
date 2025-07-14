package Handle

import (
	"NepCat_GO/Model"
	"NepCat_GO/NepCatInit/MSGModel"
	"bytes"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"regexp"
	"strconv"
	"strings"
)

func RandomPicManage(rawmsg MSGModel.ResMessage) (isreply bool) {
	fmt.Println("原始消息：", rawmsg.RawMessage)

	var qq string
	var num int = 1 // 默认图数为 1
	var tags []string
	isreply = false

	// 1️⃣ 提取 QQ 号
	qqRe := regexp.MustCompile(`qq=(\d+)`)
	qqMatches := qqRe.FindStringSubmatch(rawmsg.RawMessage)
	if len(qqMatches) >= 2 {
		qq = qqMatches[1]
		fmt.Println("✅ 提取到 QQ:", qq)
	} else {
		fmt.Println("⚠️ 未找到 QQ")
		return
	}

	// 2️⃣ 提取 CQ码后文本部分
	parts := strings.SplitN(rawmsg.RawMessage, "]", 2)
	if len(parts) < 2 {
		fmt.Println("⚠️ 消息格式不正确")
		return
	}
	text := strings.TrimSpace(parts[1])

	// 3️⃣ 判断是否是“随机涩图”开头
	if strings.HasPrefix(text, "随机涩图") {
		text = strings.TrimPrefix(text, "随机涩图")
		if strings.HasPrefix(text, "-") {
			text = strings.TrimPrefix(text, "-")
			segments := strings.SplitN(text, "-", 2)

			// 数量
			numStr := segments[0]
			if n, err := strconv.Atoi(numStr); err == nil {
				num = n
			}

			// 标签
			if len(segments) > 1 {
				tags = strings.Split(segments[1], "，") // 使用中文逗号
			}
		}
	}

	fmt.Println("🎯 QQ号:", qq)
	fmt.Println("🎯 数量:", num)
	fmt.Println("🎯 标签:", tags)

	// 调用 fetchImageURL 函数获取图片 URL
	_, err := fetchImageURL(&Model.ReqParam{
		Num:  num,
		R18:  0,
		Size: []string{"original"},
		Tags: tags,
	})

	if err != nil {
		fmt.Println("⚠️ 获取图片 URL 失败:", err)
		return
	}

	return
}

// 获取图片 URL
func fetchImageURL(reqParams *Model.ReqParam) (*[]Model.PixivImage, error) {
	apiURL := "https://api.lolicon.app/setu/v2"

	// 处理默认参数
	if reqParams.Num <= 0 {
		reqParams.Num = 1
	}
	if reqParams.R18 < 0 || reqParams.R18 > 2 {
		reqParams.R18 = 0
	}
	if reqParams.Size == nil {
		reqParams.Size = []string{"original"}
	}

	// 将参数编码为 JSON
	jsonData, err := json.Marshal(reqParams)
	if err != nil {
		return nil, fmt.Errorf("JSON 序列化失败: %v", err)
	}

	// 发送 POST 请求
	req, err := http.NewRequest("POST", apiURL, bytes.NewBuffer(jsonData))
	if err != nil {
		return nil, fmt.Errorf("创建请求失败: %v", err)
	}
	req.Header.Set("Content-Type", "application/json")

	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		return nil, fmt.Errorf("请求 API 失败: %v", err)
	}
	defer resp.Body.Close()

	// 读取响应体
	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return nil, fmt.Errorf("读取响应数据失败: %v", err)
	}

	// 解析 JSON
	var apiResp Model.APIResponse
	err = json.Unmarshal(body, &apiResp)
	if err != nil {
		return nil, fmt.Errorf("JSON 解析失败: %v", err)
	}

	// 检查错误信息
	if apiResp.Error != "" {
		return nil, fmt.Errorf("API 返回错误: %s", apiResp.Error)
	}

	// 确保返回的数据不为空
	if len(apiResp.Data) == 0 {
		return nil, fmt.Errorf("API 返回空数据")
	}

	// 返回第一张图片
	return &apiResp.Data, nil
}
