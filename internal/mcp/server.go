package mcp

import (
	"context"
	"encoding/json"
	"fast-fashion-agent/internal/service"

	mcp "github.com/mark3labs/mcp-go/mcp"
	"github.com/mark3labs/mcp-go/server"
)

type MCPServer struct {
	srv *server.MCPServer
}

func NewMCPServer() *MCPServer {
	s := server.NewMCPServer("fast-fashion-customer-service", "1.0.0")

	// Tool: get_order_info - 获取订单详情
	getOrderTool := mcp.NewTool("get_order_info",
		mcp.WithDescription("根据订单ID获取订单详情"),
		mcp.WithString("order_id", mcp.Required(), mcp.Description("订单ID")),
	)
	s.AddTool(getOrderTool, func(ctx context.Context, req mcp.CallToolRequest) (*mcp.CallToolResult, error) {
		orderID, _ := req.RequireString("order_id")
		order := service.GetOrder(orderID)
		if order == nil {
			return mcp.NewToolResultError("订单未找到"), nil
		}
		data, _ := json.Marshal(order)
		result, _ := mcp.NewToolResultJSON(data)
		return result, nil
	})

	// Tool: update_order_address - 更新收货地址
	updateAddrTool := mcp.NewTool("update_order_address",
		mcp.WithDescription("更新订单的收货地址"),
		mcp.WithString("order_id", mcp.Required(), mcp.Description("订单ID")),
		mcp.WithString("new_address", mcp.Required(), mcp.Description("新地址")),
	)
	s.AddTool(updateAddrTool, func(ctx context.Context, req mcp.CallToolRequest) (*mcp.CallToolResult, error) {
		orderID, _ := req.RequireString("order_id")
		newAddr, _ := req.RequireString("new_address")
		if !service.UpdateAddress(orderID, newAddr) {
			return mcp.NewToolResultError("更新地址失败"), nil
		}
		return mcp.NewToolResultText("地址更新成功"), nil
	})

	// Tool: get_sop_tree - 获取SOP决策树
	getSOPTool := mcp.NewTool("get_sop_tree",
		mcp.WithDescription("获取SOP决策树内容"),
		mcp.WithString("sop_type", mcp.Required(), mcp.Description("SOP类型: order 或 logistics")),
	)
	s.AddTool(getSOPTool, func(ctx context.Context, req mcp.CallToolRequest) (*mcp.CallToolResult, error) {
		typ, _ := req.RequireString("sop_type")
		var tree string
		switch typ {
		case "order":
			tree = service.OrderSOP
		case "logistics":
			tree = service.LogisticsSOP
		default:
			return mcp.NewToolResultError("无效的SOP类型"), nil
		}
		return mcp.NewToolResultText(tree), nil
	})

	return &MCPServer{srv: s}
}

// GetOrderForTesting 供内部测试使用
func GetOrderForTesting(orderID string) string {
	order := service.GetOrder(orderID)
	if order == nil {
		return ""
	}
	data, _ := json.MarshalIndent(order, "", "  ")
	return string(data)
}


