package agent

import (
	"context"
	"fast-fashion-agent/internal/service"
	"fmt"
)

type OrderAgent struct {
	*BaseAgent
}

func NewOrderAgent(base *BaseAgent) *OrderAgent {
	return &OrderAgent{BaseAgent: base}
}

// Handle 处理订单相关问题
func (oa *OrderAgent) Handle(ctx context.Context, question, history, orderInfo string) (string, error) {
	if orderInfo == "" {
		return "抱歉，未找到订单信息。请提供正确的订单号。", nil
	}

	prompt := fmt.Sprintf(`你是一个快时尚电商客服AI助手，负责处理订单相关问题。

请严格按照以下决策树进行回复:

%s

对话历史:
%s

订单信息:
%s

客户问题: %s

回复要求:
1. 严格按照决策树流程处理
2. 回答简洁明了，避免长篇大论
3. 使用友好、专业的语气
4. 不要使用正式的结束语（如"祝您生活愉快"等）
5. 回复应该直接针对客户问题，提供明确答案或行动方案
6. 如果需要客户提供更多信息，礼貌地询问

现在请根据以上信息回复客户的问题。`,
		service.OrderSOP, history, orderInfo, question)

	return oa.Generate(ctx, prompt)
}
