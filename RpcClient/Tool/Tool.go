package Tool

import (
	"fmt"
	"math/rand"
	"time"
)

func Init() {
	rand.Seed(time.Now().UnixNano())
}

// GenerateSN 生成唯一的 SN（序列号）
func GenerateSN() string {
	timestamp := time.Now().Format("20060102150405") // 格式：yyyyMMddHHmmss
	randomPart := rand.Intn(1000000)                 // 6位随机数
	return fmt.Sprintf("SN%s%06d", timestamp, randomPart)
}

func GenerateNumSN() int {
	return rand.Intn(900000) + 100000 // 生成 100000 ~ 999999 之间的随机数
}
