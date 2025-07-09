package ConfigManage

import (
	"fmt"
	"gopkg.in/yaml.v2"
	"os"
)

type WebConfig struct {
	// 数据库配置
	Database struct {
		Type     string `yaml:"db_type"`
		Address  string `yaml:"db_address"`
		Name     string `yaml:"db_name"`
		User     string `yaml:"db_user"`
		Password string `yaml:"db_password"`
		Port     int    `yaml:"port"`
	} `yaml:"database"`

	// 日志配置
	Logging struct {
		Level  string `yaml:"log_level"`
		LogDir string `yaml:"log_dir"`
	} `yaml:"logging"`

	// 服务配置（新增）
	Server struct {
		HTTPEnabled  bool   `yaml:"http_enabled"`
		HTTPPort     int    `yaml:"http_port"`
		HTTPSEnabled bool   `yaml:"https_enabled"`
		HTTPSPort    int    `yaml:"https_port"`
		RpcServer    string `yaml:"rpcserver"`
	} `yaml:"server"`

	// Nepcat配置
	NepcatInfo struct {
		LocalAddress string `yaml:"localaddress"`
		Port         int    `yaml:"port"`
	} `yaml:"nepcatinfo"`
	WebsocketInfo struct {
		Host     string `yaml:"host"`
		Port     int    `yaml:"port"`
		Scheme   string `yaml:"scheme"`
		Path     string `yaml:"path"`
		Rawquery string `yaml:"rawquery"`
	} `yaml:"websocket"`
	Mode struct {
		ReplyMode string `yaml:"replymode"`
	} `yaml:"mode"`
}

var webconf WebConfig

func GetWebConfig() *WebConfig {
	return &webconf
}

// 加载配置文件
func ConfigInit(ConfigFileName string) error {
	fmt.Println("------------------------------------------------------------------------加载本地yaml配置文件------------------------------------------------------------------------")
	file, err := os.ReadFile(ConfigFileName)
	if err != nil {
		return fmt.Errorf("failed to read config file: %v", err)
	}
	//yaml.Decoder{}
	if err = yaml.Unmarshal(file, &webconf); err != nil {
		return fmt.Errorf("failed to unmarshal config: %v", err)
	}

	return nil
}
