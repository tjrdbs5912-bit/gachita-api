package httpadapter

import (
	"net/http"
	"time"

	"gachita-api/internal/db"

	"github.com/jackc/pgx/v5/pgtype"
)

func notificationToMap(n db.Notification) map[string]string {
	m := map[string]string{
		"id":         n.ID.String(),
		"type":       n.Type,
		"title":      n.Title,
		"body":       n.Body,
		"created_at": n.CreatedAt.Time.Format(time.RFC3339),
		"read_at":    "",
		"ref_type":   "",
		"ref_id":     "",
	}
	if n.ReadAt.Valid {
		m["read_at"] = n.ReadAt.Time.Format(time.RFC3339)
	}
	if n.RefType.Valid {
		m["ref_type"] = n.RefType.String
	}
	if n.RefID.Valid {
		m["ref_id"] = n.RefID.String()
	}
	return m
}

// ListNotifications godoc
// @Summary      내 알림 목록
// @Tags         notification
// @Produce      json
// @Security     BearerAuth
// @Success      200  {array}  map[string]string
// @Router       /api/notifications [get]
func (r *Router) listNotifications(w http.ResponseWriter, req *http.Request) {
	userIDStr, ok := req.Context().Value(contextKeyUserID).(string)
	if !ok || userIDStr == "" {
		writeError(w, http.StatusUnauthorized, "인증 정보가 없습니다.")
		return
	}
	var userID pgtype.UUID
	if err := userID.Scan(userIDStr); err != nil {
		writeError(w, http.StatusBadRequest, "유효하지 않은 사용자 ID입니다.")
		return
	}

	items, err := r.queries.ListNotificationsByUserID(req.Context(), userID)
	if err != nil {
		writeError(w, http.StatusInternalServerError, "알림 조회에 실패했습니다.")
		return
	}

	list := make([]map[string]string, 0, len(items))
	for _, n := range items {
		list = append(list, notificationToMap(n))
	}
	writeJSON(w, http.StatusOK, list)
}

// ReadNotification godoc
// @Summary      알림 읽음 처리
// @Tags         notification
// @Produce      json
// @Security     BearerAuth
// @Param        notificationId  path  string  true  "알림 ID"
// @Success      200  {object}  map[string]string
// @Router       /api/notifications/{notificationId}/read [put]
func (r *Router) readNotification(w http.ResponseWriter, req *http.Request) {
	userIDStr, ok := req.Context().Value(contextKeyUserID).(string)
	if !ok || userIDStr == "" {
		writeError(w, http.StatusUnauthorized, "인증 정보가 없습니다.")
		return
	}
	var userID pgtype.UUID
	if err := userID.Scan(userIDStr); err != nil {
		writeError(w, http.StatusBadRequest, "유효하지 않은 사용자 ID입니다.")
		return
	}

	var notificationID pgtype.UUID
	if err := notificationID.Scan(req.PathValue("notificationId")); err != nil {
		writeError(w, http.StatusBadRequest, "유효하지 않은 알림 ID입니다.")
		return
	}

	n, err := r.queries.MarkNotificationRead(req.Context(), db.MarkNotificationReadParams{
		ID:     notificationID,
		UserID: userID,
	})
	if err != nil {
		writeError(w, http.StatusNotFound, "읽음 처리할 알림이 없습니다.")
		return
	}

	writeJSON(w, http.StatusOK, notificationToMap(n))
}

// ReadAllNotifications godoc
// @Summary      알림 전체 읽음 처리
// @Tags         notification
// @Produce      json
// @Security     BearerAuth
// @Success      200  {object}  map[string]string
// @Router       /api/notifications/read-all [put]
func (r *Router) readAllNotifications(w http.ResponseWriter, req *http.Request) {
	userIDStr, ok := req.Context().Value(contextKeyUserID).(string)
	if !ok || userIDStr == "" {
		writeError(w, http.StatusUnauthorized, "인증 정보가 없습니다.")
		return
	}
	var userID pgtype.UUID
	if err := userID.Scan(userIDStr); err != nil {
		writeError(w, http.StatusBadRequest, "유효하지 않은 사용자 ID입니다.")
		return
	}

	if err := r.queries.MarkAllNotificationsRead(req.Context(), userID); err != nil {
		writeError(w, http.StatusInternalServerError, "읽음 처리에 실패했습니다.")
		return
	}

	writeJSON(w, http.StatusOK, map[string]string{
		"message": "모든 알림을 읽음 처리했습니다.",
	})
}
