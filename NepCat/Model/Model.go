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
