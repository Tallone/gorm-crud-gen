package entity

import "time"

// Order 订单
type Order struct {
	ID        uint      `gorm:"primarykey" json:"-"`
	CreatedAt time.Time `json:"-"`
	UpdatedAt time.Time `json:"-"`

	// Platform 支付平台. @see cons.OpenPlatformWechat
	Platform string `gorm:"not null;index:idx_p_o" json:"-"`
	// PlatformOrderNo 支付平台订单号
	PlatformOrderNo string `gorm:"not null;index:idx_p_o" json:"-"`

	// OrderNo 订单号
	OrderNo string `gorm:"not null;uniqueIndex" json:"orderNo"`
	// UserID 用户ID
	UserID uint `gorm:"not null;index" json:"userID"`
	// Amount 订单金额，单位为分
	Amount int64 `gorm:"not null" json:"amount"`
	// Status 订单状态. @see cons.OrderStatusPending
	Status string `gorm:"not null;default:pending" json:"status"`
	// Note 订单备注
	Note string `json:"note"`
}
