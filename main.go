package main

import (
	"context"
	"os"

	"fast-fashion-agent/internal/agent"
	"fast-fashion-agent/internal/logger"
	"fast-fashion-agent/internal/mcp"
	"fast-fashion-agent/internal/middleware"

	"github.com/gin-gonic/gin"
	"github.com/tmc/langchaingo/llms/ollama"
)

const (
	DefaultOllamaBaseURL = "http://localhost:11434"
	DefaultModel         = "qwen3:14b" // 兼容 qwen3:14b
	DefaultPort          = "8080"
)

func main() {
	// 初始化日志
	logger.InitLogger(true)

	// 从环境变量读取配置，使用默认值
	ollamaURL := getEnv("OLLAMA_BASE_URL", DefaultOllamaBaseURL)
	modelName := getEnv("OLLAMA_MODEL", DefaultModel)
	port := getEnv("SERVER_PORT", DefaultPort)

	logger.Infof("🚀 Fast Fashion Agent 初始化中...")
	logger.Infof("   - Ollama URL: %s", ollamaURL)
	logger.Infof("   - Model: %s", modelName)

	// 初始化 Ollama LLM
	llm, err := ollama.New(
		ollama.WithModel(modelName),
		ollama.WithServerURL(ollamaURL),
	)
	if err != nil {
		logger.Fatalf("❌ 初始化 Ollama 失败: %v", err)
	}

	logger.Info("✅ LLM 连接成功")

	// 初始化 Agents
	baseAgent := agent.NewBaseAgent(llm)
	intentAgent := agent.NewIntentAgent(baseAgent)
	orderAgent := agent.NewOrderAgent(baseAgent)

	// 初始化 MCP Server (工具定义)
	mcpSrv := mcp.NewMCPServer()
	_ = mcpSrv // MCP 工具已注册，可用于外部 MCP 客户端
	logger.Info("✅ MCP Server 初始化成功")

	// Gin 路由
	r := gin.New()
	r.Use(middleware.Recovery())
	r.Use(middleware.GinZapLogger())

	// 测试接口
	r.GET("/test", func(c *gin.Context) {
		resp, err := intentAgent.Recognize(context.Background(), "我的订单在哪里？订单号123", "")
		if err != nil {
			logger.Errorf("Intent recognition failed: %v", err)
			c.JSON(500, gin.H{"error": err.Error()})
			return
		}
		c.JSON(200, gin.H{"intent": resp})
	})

	r.POST("/chat", func(c *gin.Context) {
		var req struct {
			Message string `json:"message"`
			History string `json:"history"`
		}
		if err := c.ShouldBindJSON(&req); err != nil {
			logger.Warnf("Invalid request: %v", err)
			c.JSON(400, gin.H{"error": "无效请求"})
			return
		}

		logger.Infof("Received message: %s", req.Message)

		// 意图识别
		intent, err := intentAgent.Recognize(context.Background(), req.Message, req.History)
		if err != nil {
			logger.Errorf("Intent recognition failed: %v", err)
			c.JSON(500, gin.H{"error": err.Error()})
			return
		}

		logger.Infof("Detected intent: %s", intent)

		var response string
		switch intent {
		case "ORDER":
			// 提取订单号 (简化版)
			orderID := "123" // 默认测试订单
			orderInfo := mcp.GetOrderForTesting(orderID)
			if orderInfo == "" {
				response = "抱歉，未找到该订单信息。"
			} else {
				response, err = orderAgent.Handle(context.Background(), req.Message, req.History, orderInfo)
				if err != nil {
					logger.Errorf("Order agent failed: %v", err)
					c.JSON(500, gin.H{"error": err.Error()})
					return
				}
			}
		case "LOGISTICS":
			response = "物流问题：请联系客服处理配送相关问题。"
		default:
			response = "抱歉，我不太理解您的问题。请问是关于订单还是配送？"
		}

		logger.Infof("Generated response: %s", response)

		c.JSON(200, gin.H{
			"intent":   intent,
			"response": response,
		})
	})

	logger.Info("")
	logger.Info("🎉 服务启动成功！")
	logger.Info("─────────────────────────────────")
	logger.Infof("   🧪 Test intent:  http://localhost:%s/test", port)
	logger.Infof("   💬 Chat API:     http://localhost:%s/chat", port)
	logger.Info("─────────────────────────────────")
	logger.Info("")

	if err := r.Run(":" + port); err != nil {
		logger.Fatalf("❌ 服务启动失败: %v", err)
	}
}

func getEnv(key, defaultValue string) string {
	if value := os.Getenv(key); value != "" {
		return value
	}
	return defaultValue
}
