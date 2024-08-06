package main

import (
	"easyDevice/device"
	"easyDevice/gat1400"
	"net/http"
	"time"

	"git.lnton.com/lnton/pkg/orm"
)

func main() {

	ca := gat1400.NewCascade(&http.Client{Timeout: 5 * time.Second})

	c := gat1400.NewCore("34010000001190000001", ca)
	v := device.Cascade{
		ModelWithStrID: orm.NewModelWithStrID("34010000001190000001"),
		Name:           " in.Name",
		Username:       "34010000001190000001",
		Password:       "Aa123456",
		IP:             "124.222.50.22",
		Port:           1400,
		Enabled:        true,
	}

	c.AddCascade(&v)

	v = device.Cascade{
		ModelWithStrID: orm.NewModelWithStrID("34010000001190000002"),
		Name:           "设备2测试平台",
		Username:       "34010000001190000002",
		Password:       "Aa123456",
		IP:             "124.222.50.22",
		Port:           1400,
		Enabled:        true,
	}

	c.AddCascade(&v)

	for {
		time.Sleep(1 * time.Second)
	}

}
