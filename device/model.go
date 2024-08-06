package device

import (
	"git.lnton.com/lnton/pkg/orm"
	"github.com/lib/pq"
)

type Cascade struct {
	orm.ModelWithStrID
	Name           string         `gorm:"notNull;default:''" json:"name"`            // 上级平台名称
	Username       string         `gorm:"notNull;default:''" json:"username"`        // 上级视图库名称
	Password       string         `gorm:"notNull;default:''" json:"password"`        // 上级视图库密码
	IP             string         `gorm:"notNull;default:''" json:"ip"`              // 上级 ip
	Port           int            `gorm:"notNull;default:0" json:"port"`             // 上级端口
	Enabled        bool           `gorm:"notNull;default:true" json:"enabled"`       // 是否启用
	VirtualGroupID int            `gorm:"notNull;default:0" json:"virtual_group_id"` // 用于虚拟组织的绑定
	DeviceIDs      pq.StringArray `gorm:"column:device_ids;type:text[];notNull;default:'{}';comment:设备ID数组"  json:"device_ids"`
}
