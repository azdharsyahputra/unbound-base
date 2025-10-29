package chat

import (
	"fmt"
	"strconv"
	"time"

	"github.com/gofiber/fiber/v2"
)

type ChatHandler struct {
	Service *ChatService
	Hub     *WebSocketHub
}

func NewChatHandler(s *ChatService, hub *WebSocketHub) *ChatHandler {
	return &ChatHandler{Service: s, Hub: hub}
}

func (h *ChatHandler) GetOrCreateChat(c *fiber.Ctx) error {
	var userID uint
	if v := c.Locals("userID"); v != nil {
		userID = v.(uint)
	} else if v := c.Locals("user_id"); v != nil {
		userID = v.(uint)
	} else {
		return fiber.NewError(fiber.StatusUnauthorized, "user id not found in context")
	}

	targetID, err := strconv.Atoi(c.Params("user_id"))
	if err != nil {
		return fiber.NewError(fiber.StatusBadRequest, "invalid target user id")
	}

	chat, err := h.Service.GetOrCreateChat(userID, uint(targetID))
	if err != nil {
		return fiber.NewError(fiber.StatusInternalServerError, err.Error())
	}

	fmt.Printf("[CHAT] User %d open chat with %d -> chat_id=%d\n", userID, targetID, chat.ID)
	return c.JSON(chat)
}

func (h *ChatHandler) GetMessages(c *fiber.Ctx) error {
	chatID, err := strconv.Atoi(c.Params("chat_id"))
	if err != nil {
		return fiber.NewError(fiber.StatusBadRequest, "invalid chat id")
	}

	messages, err := h.Service.GetMessages(uint(chatID))
	if err != nil {
		return fiber.NewError(fiber.StatusInternalServerError, err.Error())
	}

	return c.JSON(messages)
}

func (h *ChatHandler) SendMessage(c *fiber.Ctx) error {
	var userID uint
	if v := c.Locals("userID"); v != nil {
		userID = v.(uint)
	} else if v := c.Locals("user_id"); v != nil {
		userID = v.(uint)
	} else {
		return fiber.NewError(fiber.StatusUnauthorized, "user id not found in context")
	}

	chatID, err := strconv.Atoi(c.Params("chat_id"))
	if err != nil {
		return fiber.NewError(fiber.StatusBadRequest, "invalid chat id")
	}

	var req struct {
		Content string `json:"content"`
	}

	if err := c.BodyParser(&req); err != nil {
		return fiber.NewError(fiber.StatusBadRequest, "invalid JSON body")
	}

	msg, err := h.Service.SendMessage(uint(chatID), userID, req.Content)
	if err != nil {
		return fiber.NewError(fiber.StatusInternalServerError, err.Error())
	}

	fmt.Printf("[MESSAGE] User %d -> chat %d: %s\n", userID, chatID, req.Content)
	return c.JSON(msg)
}

func (h *ChatHandler) MarkAsRead(c *fiber.Ctx) error {
	var userID uint
	if v := c.Locals("userID"); v != nil {
		userID = v.(uint)
	} else if v := c.Locals("user_id"); v != nil {
		userID = v.(uint)
	} else {
		return fiber.NewError(fiber.StatusUnauthorized, "user id not found in context")
	}

	chatID, err := strconv.Atoi(c.Params("chat_id"))
	if err != nil {
		return fiber.NewError(fiber.StatusBadRequest, "invalid chat id")
	}

	now := time.Now()
	if err := h.Service.MarkMessagesAsRead(uint(chatID), userID, now); err != nil {
		return fiber.NewError(fiber.StatusInternalServerError, err.Error())
	}

	h.Hub.BroadcastEvent(uint(chatID), BroadcastPayload{
		Type:      "status_update",
		ChatID:    uint(chatID),
		Status:    "read",
		SenderID:  userID,
		Timestamp: now,
	})

	fmt.Printf("[READ] User %d marked chat %d as read (broadcasted)\n", userID, chatID)
	return c.JSON(fiber.Map{
		"success": true,
		"chat_id": chatID,
		"status":  "read",
	})
}
