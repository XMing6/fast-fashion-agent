package agent

import (
	"context"
	"fmt"
	"strings"
)

type IntentAgent struct {
	*BaseAgent
}

func NewIntentAgent(base *BaseAgent) *IntentAgent {
	return &IntentAgent{BaseAgent: base}
}

// Recognize 识别用户意图
// 返回: ORDER (订单相关), LOGISTICS (物流相关), UNKNOWN (未知)
func (ia *IntentAgent) Recognize(ctx context.Context, question, history string) (string, error) {
	prompt := fmt.Sprintf(`你是一个快时尚电商客服的意图识别系统。

对话历史:
%s

当前问题: %s

请分析用户意图，只回复以下之一:
- ORDER (订单相关: 订单状态、修改、取消、支付等)
- LOGISTICS (物流相关: 配送、地址、签收等)
- UNKNOWN (无法识别或其他问题)

例如:
"我的订单在哪" -> ORDER
"什么时候发货" -> ORDER
"我想改地址" -> LOGISTICS
"包裹显示已签收但我没收到" -> LOGISTICS

只回复 ORDER 或 LOGISTICS 或 UNKNOWN，不要包含其他内容。`,
		history, question)

	resp, err := ia.Generate(ctx, prompt)
	if err != nil {
		return "UNKNOWN", err
	}

	resp = strings.TrimSpace(strings.ToUpper(resp))
	if strings.Contains(resp, "ORDER") {
		return "ORDER", nil
	} else if strings.Contains(resp, "LOGISTICS") {
		return "LOGISTICS", nil
	}
	return "UNKNOWN", nil
}
