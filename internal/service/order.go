package service

import (
	"encoding/json"
	"fmt"
	"log"
	"os"
	"sync"
)

// Order 订单结构体
type Order struct {
	OrderID      string   `json:"order_id"`
	CustomerName string   `json:"customer_name"`
	Items        []string `json:"items"`
	Address      string   `json:"address"`
	Status       string   `json:"status"`
	TotalAmount  float64  `json:"total_amount"`
	CreateDate   string   `json:"create_date"`
}

var (
	orders = []Order{
		{
			OrderID:      "123",
			CustomerName: "张三",
			Items:        []string{"T恤", "牛仔裤"},
			Address:      "北京市西城区金融街1号",
			Status:       "处理中",
			TotalAmount:  299.00,
			CreateDate:   "2025-02-10",
		},
		{
			OrderID:      "456",
			CustomerName: "李四",
			Items:        []string{"连衣裙", "运动鞋"},
			Address:      "北京市海淀区中关村大街2号",
			Status:       "已发货",
			TotalAmount:  588.00,
			CreateDate:   "2025-02-08",
		},
		{
			OrderID:      "789",
			CustomerName: "王五",
			Items:        []string{"夹克", "帽子"},
			Address:      "北京市东城区王府井大街3号",
			Status:       "已送达",
			TotalAmount:  499.00,
			CreateDate:   "2025-02-05",
		},
	}
	mu sync.RWMutex
)

// GetOrder 根据订单ID获取订单
func GetOrder(orderID string) *Order {
	mu.RLock()
	defer mu.RUnlock()
	for _, o := range orders {
		if o.OrderID == orderID {
			return &o
		}
	}
	return nil
}

// UpdateAddress 更新订单收货地址
func UpdateAddress(orderID, newAddr string) bool {
	mu.Lock()
	defer mu.Unlock()
	for i, o := range orders {
		if o.OrderID == orderID {
			orders[i].Address = newAddr
			log.Printf("订单 %s 地址已更新: %s", orderID, newAddr)
			return true
		}
	}
	return false
}

// GetAllOrders 获取所有订单（用于测试）
func GetAllOrders() []Order {
	mu.RLock()
	defer mu.RUnlock()
	return orders
}

// FormatOrder 格式化订单信息为可读文本
func FormatOrder(order *Order) string {
	if order == nil {
		return "订单不存在"
	}
	return fmt.Sprintf(`订单号: %s
客户: %s
商品: %v
金额: ¥%.2f
状态: %s
地址: %s
下单日期: %s`,
		order.OrderID,
		order.CustomerName,
		order.Items,
		order.TotalAmount,
		order.Status,
		order.Address,
		order.CreateDate,
	)
}

// LoadOrdersFromFile 从文件加载订单
func LoadOrdersFromFile(filepath string) error {
	data, err := os.ReadFile(filepath)
	if err != nil {
		return err
	}
	return json.Unmarshal(data, &orders)
}

// SaveOrdersToFile 保存订单到文件
func SaveOrdersToFile(filepath string) error {
	mu.Lock()
	defer mu.Unlock()
	data, err := json.MarshalIndent(orders, "", "  ")
	if err != nil {
		return err
	}
	return os.WriteFile(filepath, data, 0644)
}
