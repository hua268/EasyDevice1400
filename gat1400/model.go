package gat1400

import "easyDevice/device"

// ResponseStatus 响应
type ResponseStatus struct {
	ResponseStatusObject ResponseStatusObject
}

// ResponseStatusObject 响应对象
type ResponseStatusObject struct {
	Id           string // nolint  ，请求注册的 DeviceID
	StatusCode   int    // 本次注册的操作响应码
	RequestURL   string // 本次请求的资源定位符
	StatusString string // 本次注册的操作响应说明
	LocalTime    string // 被注册方系统时间，例如 20180116141531
}

// CascadeInfo 级联信息
type CascadeInfo struct {
	device.Cascade

	LastHeartbeatAt  int64 // 上次心跳时间
	LastRegisteredAt int64 // 上次注册时间
	Registered       bool
	Status           string
}
