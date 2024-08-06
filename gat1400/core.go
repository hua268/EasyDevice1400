package gat1400

import (
	"easyDevice/device"
	"log/slog"
	"math/rand"
	"time"

	"git.lnton.com/lnton/pkg/async"
)

type Core struct {
	cascades *async.Map[string, *CascadeInfo] // 向上级联列表
	cas      Cascade                          // 上级平台信息
	serverID string                           // 本平台的ID
}

func NewCore(id string, cas Cascade) *Core {
	c := Core{
		serverID: id,
		cas:      cas,
		cascades: async.NewMap[string, *CascadeInfo](),
	}

	go c.cascadesRegister()
	go c.sendMsg()
	return &c
}

// AddCascade 添加级联
func (c *Core) AddCascade(cas *device.Cascade) {
	// 添加级联信息
	c.cascades.Store(cas.ID, &CascadeInfo{Cascade: *cas})
}

func (c *Core) sendMsg() {
	ticker := time.NewTicker(5 * time.Second)
	defer ticker.Stop()
	for range ticker.C {
		c.cascades.Range(func(_ string, value *CascadeInfo) bool {
			if value.Status != "Ok" {
				return true
			}

			return true
		})
	}
}

// cascadesRegister 向上级联注册
func (c *Core) cascadesRegister() {
	// 创建定时器，每5秒执行一次注册和心跳操作
	ticker := time.NewTicker(5 * time.Second)
	// 循环执行注册和心跳操作
	defer ticker.Stop()
	for range ticker.C {
		// 遍历级联服务列表
		c.cascades.Range(func(_ string, value *CascadeInfo) bool {
			if !value.Enabled {
				return true
			}

			// 未注册，先注册
			if !value.Registered {
				// 300s 内随机时间，重新注册
				if value.LastRegisteredAt > 0 && time.Now().Unix()-value.LastRegisteredAt < 20 && rand.Intn(10) != 1 { // nolint
					return true
				}
				value.LastRegisteredAt = time.Now().Unix() // 更新最后注册时间
				value.LastHeartbeatAt = time.Now().Unix()  // 更新最后心跳时间
				slog.Info("设备发起注册", "device", c.serverID)
				// 注册级联服务
				if err := c.cas.Register(c.serverID, value.Cascade); err != nil {
					value.Status = err.Error()         // 记录错误状态
					slog.Error("Register", "err", err) // 记录注册错误日志
					return true                        // 继续下一个级联服务
				}
				value.Registered = true // 标记为已注册
				slog.Info("设备注册成功", "device", c.serverID)
				value.Status = "OK" // 更新状态为正常
				//go c.notificationReSend(value.ID) // 异步重新发送通知
			}
			// 发送心跳
			if time.Now().Unix()-value.LastHeartbeatAt >= 55 {
				value.LastHeartbeatAt = time.Now().Unix() // 异步重新发送通知
				slog.Info("发送发起心跳", "device", c.serverID) // 记录日志：心跳
				// 发送心跳
				err := c.cas.Heartbeat(c.serverID, value.Cascade)
				if err != nil {
					value.Registered = false            // 标记为未注册
					value.Status = err.Error()          // 记录错误状态
					slog.Error("Heartbeat", "err", err) // 记录心跳错误日志
					return true                         // 继续下一个级联服务
				}
				slog.Info("发送心跳成功", "device", c.serverID)
				value.Status = "OK" // 更新状态为正常
			}
			return true // 继续下一个级联服务
		})
	}
}
