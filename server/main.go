package main

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"sync"
	"time"

	"github.com/gorilla/websocket"
)

type LoginRequest struct {
	Username string `json:"username"`
	Password string `json:"password"`
}

type LoginResponse struct {
	Success    bool        `json:"success"`
	Message    string      `json:"message"`
	Token      string      `json:"token,omitempty"`
	ClientInfo *ClientInfo `json:"clientInfo,omitempty"`
}

func handleLogin(w http.ResponseWriter, r *http.Request) {
	// 添加 CORS 头部
	w.Header().Set("Access-Control-Allow-Origin", "*")
	w.Header().Set("Access-Control-Allow-Methods", "POST, GET, OPTIONS")
	w.Header().Set("Access-Control-Allow-Headers", "Content-Type")

	// 处理 OPTIONS 请求
	if r.Method == http.MethodOptions {
		w.WriteHeader(http.StatusOK)
		return
	}

	fmt.Println("Received login request")
	if r.Method != http.MethodPost {
		fmt.Println("Method not allowed")
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	var req LoginRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		fmt.Println("Invalid request body:", err)
		http.Error(w, "Invalid request body", http.StatusBadRequest)
		return
	}

	fmt.Printf("Login request: username=%s, password=%s\n", req.Username, req.Password)

	if req.Username == "" || req.Password == "" {
		response := LoginResponse{
			Success: false,
			Message: "Username and password are required",
		}
		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(response)
		return
	}

	// 生成客户端ID
	clientID := "client_" + time.Now().Format("20060102150405") + "_" + req.Username
	machineID := "machine_" + time.Now().Format("20060102150405") + "_" + req.Username

	// 获取客户端IP
	clientIP := r.RemoteAddr
	if forwardedFor := r.Header.Get("X-Forwarded-For"); forwardedFor != "" {
		clientIP = forwardedFor
	}

	// 保存客户端信息
	clientInfo := &ClientInfo{
		ClientID:  clientID,
		MachineID: machineID,
		Username:  req.Username,
		IP:        clientIP,
		Online:    true,
		LastSeen:  time.Now(),
		LoginTime: time.Now(),
	}

	clientInfosMutex.Lock()
	clientInfos[clientID] = clientInfo
	clientInfosMutex.Unlock()

	fmt.Printf("Client logged in: %s (username: %s, ip: %s)\n", clientID, req.Username, clientIP)

	response := LoginResponse{
		Success:    true,
		Message:    "Login successful",
		Token:      "sample-token-" + req.Username,
		ClientInfo: clientInfo,
	}

	// 添加客户端信息到响应
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(response)
	fmt.Println("Login response sent")
}

type Machine struct {
	ID       string `json:"id"`
	Name     string `json:"name"`
	IP       string `json:"ip"`
	Online   bool   `json:"online"`
	LastSeen string `json:"lastSeen"`
}

// WebSocket升级器
var upgrader = websocket.Upgrader{
	CheckOrigin: func(r *http.Request) bool {
		return true // 允许所有来源的连接
	},
}

// 信令消息结构
type SignalingMessage struct {
	Type       string      `json:"type"`       // offer, answer, ice-candidate
	SenderID   string      `json:"senderId"`   // 发送者ID
	ReceiverID string      `json:"receiverId"` // 接收者ID
	Data       interface{} `json:"data"`       // 消息数据
}

// 客户端信息结构
type ClientInfo struct {
	ClientID  string    `json:"clientId"`
	MachineID string    `json:"machineId"`
	Username  string    `json:"username"`
	IP        string    `json:"ip"`
	Online    bool      `json:"online"`
	LastSeen  time.Time `json:"lastSeen"`
	LoginTime time.Time `json:"loginTime"`
}

// 客户端连接管理
type Client struct {
	conn      *websocket.Conn
	clientID  string
	machineID string
}

// 客户端连接映射
var (
	clients          = make(map[string]*Client)
	clientsMutex     sync.Mutex
	clientInfos      = make(map[string]*ClientInfo)
	clientInfosMutex sync.Mutex
)

func handleMachines(w http.ResponseWriter, r *http.Request) {
	// 添加 CORS 头部
	w.Header().Set("Access-Control-Allow-Origin", "*")
	w.Header().Set("Access-Control-Allow-Methods", "GET, POST, OPTIONS")
	w.Header().Set("Access-Control-Allow-Headers", "Content-Type")

	// 处理 OPTIONS 请求
	if r.Method == http.MethodOptions {
		w.WriteHeader(http.StatusOK)
		return
	}

	fmt.Println("Received machines request")

	// 从客户端信息生成机器列表
	var machines []Machine

	clientInfosMutex.Lock()
	fmt.Printf("Number of clients: %d\n", len(clientInfos))
	for _, clientInfo := range clientInfos {
		machine := Machine{
			ID:       clientInfo.MachineID,
			Name:     clientInfo.Username + "的机器",
			IP:       clientInfo.IP,
			Online:   clientInfo.Online,
			LastSeen: clientInfo.LastSeen.Format("2006-01-02 15:04:05"),
		}
		machines = append(machines, machine)
		fmt.Printf("Added machine: %s (ID: %s, IP: %s)\n", machine.Name, machine.ID, machine.IP)
	}
	clientInfosMutex.Unlock()

	fmt.Printf("Returning %d machines\n", len(machines))

	response := map[string]interface{}{
		"success": true,
		"message": "Machines retrieved successfully",
		"data":    machines,
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(response)
	fmt.Println("Machines response sent")
}

// 处理WebSocket连接
func handleWebSocket(w http.ResponseWriter, r *http.Request) {
	// 升级HTTP连接为WebSocket
	conn, err := upgrader.Upgrade(w, r, nil)
	if err != nil {
		fmt.Println("WebSocket升级失败:", err)
		return
	}
	defer conn.Close()

	// 从查询参数中获取客户端ID和机器ID
	clientID := r.URL.Query().Get("clientId")
	machineID := r.URL.Query().Get("machineId")

	if clientID == "" {
		fmt.Println("客户端ID不能为空")
		return
	}

	// 注册客户端
	clientsMutex.Lock()
	clients[clientID] = &Client{
		conn:      conn,
		clientID:  clientID,
		machineID: machineID,
	}
	clientsMutex.Unlock()

	fmt.Printf("客户端 %s (机器ID: %s) %v 已连接\n", clientID, machineID, clients[clientID])

	// 处理消息
	for {
		var msg SignalingMessage
		err := conn.ReadJSON(&msg)
		if err != nil {
			fmt.Printf("读取消息失败: %v\n", err)
			break
		}

		fmt.Printf("收到信令消息: %s 从 %s 到 %s\n", msg.Type, msg.SenderID, msg.ReceiverID)

		// 转发消息给目标客户端
		clientsMutex.Lock()
		receiver, exists := clients[msg.ReceiverID]
		clientsMutex.Unlock()

		if exists {
			err := receiver.conn.WriteJSON(msg)
			if err != nil {
				fmt.Printf("发送消息失败: %v\n", err)
			}
		} else {
			fmt.Printf("接收者 %s 不存在\n", msg.ReceiverID)
		}
	}

	// 连接关闭后移除客户端
	clientsMutex.Lock()
	delete(clients, clientID)
	clientsMutex.Unlock()

	fmt.Printf("客户端 %s 已断开连接\n", clientID)
}

func main() {
	http.HandleFunc("/api/login", handleLogin)
	http.HandleFunc("/api/machines", handleMachines)
	http.HandleFunc("/ws", handleWebSocket)

	port := ":8081"
	fmt.Printf("Server starting on port %s...\n", port)
	log.Fatal(http.ListenAndServe(port, nil))
}
