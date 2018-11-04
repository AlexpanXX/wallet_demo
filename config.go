package main

// 钱包的配置参数
type config struct {
	Magic uint32
	SeedList []string
	DefaultPort uint16
}