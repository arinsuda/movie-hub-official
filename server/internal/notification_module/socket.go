package notification_module

import (
	"net/http"
	"strconv"
	"sync"

	"github.com/zishang520/socket.io/servers/socket/v3"
	"github.com/zishang520/socket.io/v3/pkg/types"
)

type TokenVerifier interface {
	VerifyToken(token string) (uint, error)
}

type Hub struct {
	io       *socket.Server
	verifier TokenVerifier
	mu       sync.RWMutex
}

func NewHub(verifier TokenVerifier, allowedOrigin string) *Hub {
	opts := socket.DefaultServerOptions()
	opts.SetCors(&types.Cors{
		Origin:      allowedOrigin,
		Credentials: true,
	})

	io := socket.NewServer(nil, opts)

	h := &Hub{
		io:       io,
		verifier: verifier,
	}

	io.On("connection", func(clients ...any) {
		client := clients[0].(*socket.Socket)

		auth := client.Handshake().Auth
		token, _ := auth["token"].(string)

		userID, err := verifier.VerifyToken(token)
		if err != nil {
			_ = client.Emit("error", "unauthorized")
			client.Disconnect(true)
			return
		}

		client.Join(roomFor(userID))
	})

	return h
}

func roomFor(userID uint) socket.Room {
	return socket.Room("user:" + strconv.FormatUint(uint64(userID), 10))
}

func (h *Hub) Handler() http.Handler {
	return h.io.ServeHandler(nil)
}

func (h *Hub) EmitNew(userID uint, n NotificationResponse) {
	_ = h.io.To(roomFor(userID)).Emit("notification:new", n)
}

func (h *Hub) EmitRead(userID uint, ids []uint) {
	_ = h.io.To(roomFor(userID)).Emit("notification:read", map[string]any{"ids": ids})
}

func (h *Hub) EmitDeleted(userID uint, ids []uint) {
	_ = h.io.To(roomFor(userID)).Emit("notification:deleted", map[string]any{"ids": ids})
}

func (h *Hub) EmitUnreadCount(userID uint, count int64) {
	_ = h.io.To(roomFor(userID)).Emit("notification:unread_count", map[string]any{
		"unread_count": count,
	})
}
