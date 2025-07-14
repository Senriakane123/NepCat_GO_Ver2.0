package Model

type SysMonitorStatus struct {
	TimeStamp string       `json:"time_stamp"` //时间戳
	CPU       float64      `json:"CPU"`        //CPU 使用率
	Memory    float64      `json:"Memory"`     // 内存使用率
	Network   NetworkSpeed `json:"Network"`
	Disk      DiskSpeed    `json:"Disk"` //

}

type DiskSpeed struct {
	ReadBPS  float64 `json:"read_bps"`  // 磁盘每秒读取字节数
	WriteBPS float64 `json:"write_bps"` // 磁盘每秒写入字节数
}

type NetworkSpeed struct {
	RecvBPS float64 `json:"recv_bps"` // 网络每秒接收字节数
	SendBPS float64 `json:"send_bps"` // 网络每秒发送字节数
}

type ReqParam struct {
	Num     int      `json:"num,omitempty"`
	R18     int      `json:"r18,omitempty"`
	Keyword string   `json:"keyword,omitempty"`
	Tags    []string `json:"tag,omitempty"`
	Size    []string `json:"size,omitempty"`
	Proxy   string   `json:"proxy,omitempty"`
}

// api获取图片
type APIResponse struct {
	Error string       `json:"error"`
	Data  []PixivImage `json:"data"`
} //type Data struct {

// PixivImage 结构体表示单个图片信息
type PixivImage struct {
	ID     uint     `gorm:"primaryKey;autoIncrement" json:"id"`      // 设为主键，自增
	PID    int      `gorm:"uniqueIndex" json:"pid"`                  // 唯一索引
	Title  string   `gorm:"type:varchar(255);not null" json:"title"` // 设定字符串长度
	Author string   `gorm:"type:varchar(255);not null" json:"author"`
	Tags   []string `gorm:"serializer:json" json:"tags"` // JSON 序列化存储

	URLs struct {
		Original string `json:"original"`
	} `gorm:"embedded" json:"urls"` // 使用 embedded 方式存储
}
