package gat1400

import (
	"bytes"
	"easyDevice/device"
	"easyDevice/pkg/digest"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"strings"
)

// Cascade 级联
type Cascade struct {
	cli *http.Client
}

// NewCascade 级联对象创建
func NewCascade(cli *http.Client) Cascade {
	return Cascade{cli: cli}
}

// Register 向上级图库注册
func (c *Cascade) Register(deviceID string, cas device.Cascade) error {
	// 构建注册请求的路径
	const path = `/VIID/System/Register`
	// 构建注册请求的JSON数据
	uri := fmt.Sprintf("http://%s:%d%s", cas.IP, cas.Port, path)
	// 构建注册请求的JSON数据
	var in struct {
		RegisterObject struct {
			DeviceID string // 注册设备的 ID，指采集设备、采集系统、应用平台、分析系统等
		}
	}
	in.RegisterObject.DeviceID = deviceID
	body, _ := json.Marshal(in)
	// 发送注册请求
	return registerAndUnRegister(uri, cas.Username, cas.Password, body)
}

// registerAndUnRegister 向指定的服务进行注册或注销操作
func registerAndUnRegister(uri, name, pass string, body []byte) error {
	// 创建带Digest认证的HTTP客户端传输
	ts := digest.NewTransport(name, pass)

	// 构建POST请求，是哦那个application/VIID+JSON类型的Content-Type，请求体为Body
	req, _ := http.NewRequest(http.MethodPost, uri, bytes.NewReader(body))
	req.Header.Set("Content-Type", "application/VIID+JSON")

	// 发送请求并接收响应
	resp, err := ts.RoundTrip(req)
	if err != nil {
		// 网络或发送请求过程中发生错误，直接返回该错误。
		return err
	}
	// 确保响应体被关闭，即使在读取体内容后也确保释放资源
	defer resp.Body.Close() // nolint

	// 检查响应状态码，如果不是200，则返回错误信息
	if resp.StatusCode != 200 {
		b, _ := io.ReadAll(resp.Body)
		// 返回响应状态码和错误信息
		return fmt.Errorf("status[%s] err[%s]", resp.Status, strings.ReplaceAll(strings.ReplaceAll(string(b), " ", ""), "\"", `"`))
	}

	// 解析响应体，检查响应状态码
	var out ResponseStatus
	if err := json.NewDecoder(resp.Body).Decode(&out); err != nil {
		// 解析响应体出错，返回错误。
		return err
	}
	// 检查响应体中的状态码，如果为0，则代表操作成功，否则返回错误信息
	if out.ResponseStatusObject.StatusCode == 0 {
		// 操作成功，返回nil。
		return nil
	}
	// 返回包含HTTP状态码和响应状态字符串的错误信息。
	return fmt.Errorf("%d:%s", resp.StatusCode, out.ResponseStatusObject.StatusString)
}

// UnRegister 向上级图库取消注册
func (c *Cascade) UnRegister(deviceID string, cas device.Cascade) error {
	// 构建注销请求的路径
	const path = `/VIID/System/UnRegister`
	// 构建注销请求的JSON数据
	uri := fmt.Sprintf("http://%s:%d%s", cas.IP, cas.Port, path)
	var in struct {
		UnRegisterObject struct {
			DeviceID string
		}
	}
	// 设置注销请求中的设备ID
	in.UnRegisterObject.DeviceID = deviceID
	// 将注销请求序列化为JSON格式
	body, _ := json.Marshal(in)
	// 发送注销请求
	return registerAndUnRegister(uri, cas.Username, cas.Password, body)
}

// Heartbeat 发送心跳
func (c *Cascade) Heartbeat(deviceID string, cas device.Cascade) error {
	// 定义心跳信息请求的路径
	const path = `/VIID/System/Keepalive`

	// 根据设备提供的IP和端口构造完整的请求URI
	uri := fmt.Sprintf("http://%s:%d%s", cas.IP, cas.Port, path)

	// 创建输入结构，用于序列化成JSON格式的心跳信息
	var in struct {
		KeepaliveObject struct {
			DeviceID string
		}
	}
	// 设置心跳信息中的设备ID
	in.KeepaliveObject.DeviceID = deviceID
	// 将心跳信息序列化为JSON。
	b, _ := json.Marshal(in)

	// 调用do方法发送心跳请求，并获取响应状态、状态码和可能的错误。
	resp, statusCode, err := c.do(deviceID, uri, b, nil)
	// 出现网络请求错误，则直接返回该错误
	if err != nil {
		return err
	}
	// 状态码为200且响应状态为0表示成功接收心跳信息，无需返回错误
	if statusCode == 200 && resp.ResponseStatusObject.StatusCode == 0 {
		return nil
	}
	// 状态码为401表示权限问题，返回相应的错误信息。
	if statusCode == 401 {
		return fmt.Errorf("没有权限")
	}
	// 其他情况下，返回包含状态码和状态字符串的错误信息。
	return fmt.Errorf("%d:%s", statusCode, resp.ResponseStatusObject.StatusString)
}

// do 方法向指定的路径发送POST请求，携带JSON格式的数据和自定义的HTTP头信息，并返回解析后的响应状态及HTTP状态码。
func (c *Cascade) do(deviceID, path string, data []byte, header map[string]string) (*ResponseStatus, int, error) {
	// 根据给定的path和data创建POST请求
	req, _ := http.NewRequest(http.MethodPost, path, bytes.NewReader(data))

	// 设置“Content-Type”为“application/VIID+JSON”和“User-Identify”为提供的deviceID。
	req.Header.Add("Content-Type", "application/VIID+JSON")
	req.Header.Add("User-Identify", deviceID)

	// 迭代提供的header map，并将每一对key-value作为请求头添加到请求中。
	for k, v := range header {
		req.Header.Add(k, v)
	}
	// 使用Cascade结构的HTTP客户端执行请求，并接收响应或错误
	resp, err := c.cli.Do(req)
	if err != nil {
		// 如果请求发送失败，返回nil响应状态、0状态码和对应的错误。
		return nil, 0, err
	}

	// 确保响应体在函数返回时关闭，以释放网络资源。
	defer resp.Body.Close() // nolint
	var out ResponseStatus
	// 将响应体的内容解析为ResponseStatus结构，并忽略解析过程中可能出现的错误。
	_ = json.NewDecoder(resp.Body).Decode(&out)
	// 返回解析后的响应状态、HTTP状态码和nil表示的无错误。
	return &out, resp.StatusCode, nil
}
