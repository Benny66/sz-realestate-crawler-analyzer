package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"
	"time"
)

// 企业微信机器人配置
const WechatWebhookURL = "https://qyapi.weixin.qq.com/cgi-bin/webhook/send?key=f5e8ee74-efbd-4e02-9011-20ab8decb475"

// WechatMessage 企业微信消息结构
type WechatMessage struct {
	MsgType string `json:"msgtype"`
	Text    struct {
		Content string `json:"content"`
	} `json:"text"`
}

// WechatMarkdownMessage 企业微信 Markdown 消息
type WechatMarkdownMessage struct {
	MsgType  string `json:"msgtype"`
	Markdown struct {
		Content string `json:"content"`
	} `json:"markdown"`
}

// SendWechatMessage 发送企业微信消息
func SendWechatMessage(content string) error {
	msg := WechatMessage{
		MsgType: "text",
	}
	msg.Text.Content = content

	msgBytes, err := json.Marshal(msg)
	if err != nil {
		return fmt.Errorf("序列化消息失败: %v", err)
	}

	resp, err := http.Post(WechatWebhookURL, "application/json", bytes.NewBuffer(msgBytes))
	if err != nil {
		return fmt.Errorf("发送消息失败: %v", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return fmt.Errorf("企业微信接口返回异常: %s", resp.Status)
	}

	return nil
}

// SendWechatMarkdownMessage 发送企业微信 Markdown 消息
func SendWechatMarkdownMessage(content string) error {
	msg := WechatMarkdownMessage{
		MsgType: "markdown",
	}
	msg.Markdown.Content = content

	msgBytes, err := json.Marshal(msg)
	if err != nil {
		return fmt.Errorf("序列化消息失败: %v", err)
	}

	resp, err := http.Post(WechatWebhookURL, "application/json", bytes.NewBuffer(msgBytes))
	if err != nil {
		return fmt.Errorf("发送消息失败: %v", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return fmt.Errorf("企业微信接口返回异常: %s", resp.Status)
	}

	return nil
}

// 推送消息模板
func BuildPriceAlertMessage(projectName, buildingName, houseType string, currentPrice, lastPrice float64, changePct float64) string {
	direction := "上涨"
	if changePct < 0 {
		direction = "下跌"
	}
	
	return fmt.Sprintf("🏠 **价格变化提醒**\n"+
		"**楼盘**: %s %s %s\n"+
		"**当前价格**: %.2f 元/㎡\n"+
		"**上次价格**: %.2f 元/㎡\n"+
		"**变化幅度**: %s %.2f%%\n"+
		"**时间**: %s",
		projectName, buildingName, houseType,
		currentPrice, lastPrice,
		direction, changePct,
		time.Now().Format("2006-01-02 15:04:05"))
}

func BuildSaleAlertMessage(projectName, buildingName, houseType string, currentRate, lastRate float64, changePct float64) string {
	return fmt.Sprintf("📊 **销售状态提醒**\n"+
		"**楼盘**: %s %s %s\n"+
		"**当前在售率**: %.2f%%\n"+
		"**上次在售率**: %.2f%%\n"+
		"**变化幅度**: %.2f%%\n"+
		"**时间**: %s",
		projectName, buildingName, houseType,
		currentRate, lastRate,
		changePct,
		time.Now().Format("2006-01-02 15:04:05"))
}

func BuildNewFavoriteMessage(projectName, buildingName, houseType, zone string) string {
	return fmt.Sprintf("⭐ **新收藏楼盘**\n"+
		"**楼盘**: %s\n"+
		"**楼栋**: %s\n"+
		"**户型**: %s\n"+
		"**区域**: %s\n"+
		"**时间**: %s",
		projectName, buildingName, houseType, zone,
		time.Now().Format("2006-01-02 15:04:05"))
}
