package main

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"time"
)

const BaseURL = "http://localhost:8080/api/v1"

var passed, failed int

type Response struct {
	Code    int             `json:"code"`
	Message string          `json:"message"`
	Data    json.RawMessage `json:"data"`
}

func Test(name, path string) {
	resp, err := http.Get(BaseURL + path)
	if err != nil {
		fmt.Printf("❌ %s\n   Error: %v\n\n", name, err)
		failed++
		return
	}
	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		fmt.Printf("❌ %s\n   Error reading body: %v\n\n", name, err)
		failed++
		return
	}

	if resp.StatusCode != http.StatusOK {
		fmt.Printf("❌ %s\n   Status: %d\n   Body: %s\n\n", name, resp.StatusCode, string(body))
		failed++
		return
	}

	var r Response
	if err := json.Unmarshal(body, &r); err != nil {
		fmt.Printf("❌ %s\n   Error parsing JSON: %v\n   Body: %s\n\n", name, err, string(body))
		failed++
		return
	}

	if r.Code != 200 {
		fmt.Printf("❌ %s\n   Code: %d, Message: %s\n\n", name, r.Code, r.Message)
		failed++
		return
	}

	fmt.Printf("✅ %s\n   %s\n\n", name, string(body))
	passed++
}

func main() {
	fmt.Println("\n========================================")
	fmt.Println("  Auction API Tests")
	fmt.Println("  Time:", time.Now().Format("2006-01-02 15:04:05"))
	fmt.Println("========================================\n")

	// 1. 查询拍卖列表
	fmt.Println("--- List Auctions ---")
	Test("默认查询", "/auctions")
	Test("按状态过滤-active", "/auctions?status=active")
	Test("按状态过滤-ended", "/auctions?status=ended")
	Test("分页", "/auctions?page=1&limit=10")
	Test("按创建时间升序", "/auctions?sort=created_at&order=asc")
	Test("按最高出价降序", "/auctions?sort=highest_bid&order=desc")

	// 汇总
	fmt.Println("========================================")
	fmt.Printf("  Results: ✅ %d passed, ❌ %d failed\n", passed, failed)
	fmt.Println("========================================\n")
}
