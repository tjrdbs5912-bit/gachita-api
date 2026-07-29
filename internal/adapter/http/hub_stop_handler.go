package httpadapter

import (
	"net/http"
	"strconv"

	"gachita-api/internal/db"

	"github.com/jackc/pgx/v5/pgtype"
)

type createHubStopRequest struct {
	Name      string `json:"name"`
	SortOrder *int32 `json:"sort_order"`
}

// CreateHubStop godoc
// @Summary      거점 추가
// @Tags         hub_stop
// @Accept       json
// @Produce      json
// @Security     BearerAuth
// @Param        id    path  string  true  "방 ID"
// @Param        body  body  createHubStopRequest  true  "거점 정보"
// @Success      201  {object}  map[string]string
// @Router       /api/rooms/{id}/stops [post]
func (r *Router) createHubStop(w http.ResponseWriter, req *http.Request) {
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

	var roomID pgtype.UUID
	if err := roomID.Scan(req.PathValue("id")); err != nil {
		writeError(w, http.StatusBadRequest, "유효하지 않은 방 ID입니다.")
		return
	}

	// 멤버만 추가 가능
	_, err := r.queries.GetRoomMember(req.Context(), db.GetRoomMemberParams{
		RoomID: roomID,
		UserID: userID,
	})
	if err != nil {
		writeError(w, http.StatusForbidden, "방에 참여하지 않았습니다.")
		return
	}

	var body createHubStopRequest
	if err := decodeJSON(req, &body); err != nil {
		writeError(w, http.StatusBadRequest, "invalid json")
		return
	}
	if body.Name == "" {
		writeError(w, http.StatusBadRequest, "거점 이름은 필수입니다.")
		return
	}

	var sortOrder int32
	if body.SortOrder != nil {
		sortOrder = *body.SortOrder
	}

	stop, err := r.queries.CreateHubStop(req.Context(), db.CreateHubStopParams{
		RoomID:    roomID,
		Name:      body.Name,
		SortOrder: sortOrder,
	})
	if err != nil {
		writeError(w, http.StatusInternalServerError, "거점 추가에 실패했습니다.")
		return
	}

	writeJSON(w, http.StatusCreated, map[string]string{
		"id":         stop.ID.String(),
		"room_id":    stop.RoomID.String(),
		"name":       stop.Name,
		"sort_order": strconv.Itoa(int(stop.SortOrder)),
	})
}

// ListHubStops godoc
// @Summary      거점 목록
// @Tags         hub_stop
// @Produce      json
// @Security     BearerAuth
// @Param        id   path  string  true  "방 ID"
// @Success      200  {array}  map[string]string
// @Router       /api/rooms/{id}/stops [get]
func (r *Router) listHubStops(w http.ResponseWriter, req *http.Request) {
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

	var roomID pgtype.UUID
	if err := roomID.Scan(req.PathValue("id")); err != nil {
		writeError(w, http.StatusBadRequest, "유효하지 않은 방 ID입니다.")
		return
	}

	_, err := r.queries.GetRoomMember(req.Context(), db.GetRoomMemberParams{
		RoomID: roomID,
		UserID: userID,
	})
	if err != nil {
		writeError(w, http.StatusForbidden, "방에 참여하지 않았습니다.")
		return
	}

	stops, err := r.queries.ListHubStopsByRoomID(req.Context(), roomID)
	if err != nil {
		writeError(w, http.StatusInternalServerError, "거점 조회에 실패했습니다.")
		return
	}

	list := make([]map[string]string, 0, len(stops))
	for _, s := range stops {
		list = append(list, map[string]string{
			"id":         s.ID.String(),
			"room_id":    s.RoomID.String(),
			"name":       s.Name,
			"sort_order": strconv.Itoa(int(s.SortOrder)),
		})
	}
	writeJSON(w, http.StatusOK, list)
}

type updateHubStopRequest struct {
	Name      string `json:"name"`
	SortOrder *int32 `json:"sort_order"`
}

// UpdateHubStop godoc
// @Summary      거점 수정
// @Tags         hub_stop
// @Accept       json
// @Produce      json
// @Security     BearerAuth
// @Param        id      path  string  true  "방 ID"
// @Param        stopId  path  string  true  "거점 ID"
// @Param        body    body  updateHubStopRequest  true  "거점 정보"
// @Success      200  {object}  map[string]string
// @Router       /api/rooms/{id}/stops/{stopId} [put]
func (r *Router) updateHubStop(w http.ResponseWriter, req *http.Request) {
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

	var roomID pgtype.UUID
	if err := roomID.Scan(req.PathValue("id")); err != nil {
		writeError(w, http.StatusBadRequest, "유효하지 않은 방 ID입니다.")
		return
	}

	var stopID pgtype.UUID
	if err := stopID.Scan(req.PathValue("stopId")); err != nil {
		writeError(w, http.StatusBadRequest, "유효하지 않은 거점 ID입니다.")
		return
	}

	_, err := r.queries.GetRoomMember(req.Context(), db.GetRoomMemberParams{
		RoomID: roomID,
		UserID: userID,
	})
	if err != nil {
		writeError(w, http.StatusForbidden, "방에 참여하지 않았습니다.")
		return
	}

	var body updateHubStopRequest
	if err := decodeJSON(req, &body); err != nil {
		writeError(w, http.StatusBadRequest, "invalid json")
		return
	}
	if body.Name == "" {
		writeError(w, http.StatusBadRequest, "거점 이름은 필수입니다.")
		return
	}

	var sortOrder int32
	if body.SortOrder != nil {
		sortOrder = *body.SortOrder
	}

	stop, err := r.queries.UpdateHubStop(req.Context(), db.UpdateHubStopParams{
		ID:        stopID,
		RoomID:    roomID,
		Name:      body.Name,
		SortOrder: sortOrder,
	})
	if err != nil {
		writeError(w, http.StatusNotFound, "거점을 찾을 수 없습니다.")
		return
	}

	writeJSON(w, http.StatusOK, map[string]string{
		"id":         stop.ID.String(),
		"room_id":    stop.RoomID.String(),
		"name":       stop.Name,
		"sort_order": strconv.Itoa(int(stop.SortOrder)),
	})
}

// DeleteHubStop godoc
// @Summary      거점 삭제
// @Tags         hub_stop
// @Produce      json
// @Security     BearerAuth
// @Param        id      path  string  true  "방 ID"
// @Param        stopId  path  string  true  "거점 ID"
// @Success      200  {object}  map[string]string
// @Router       /api/rooms/{id}/stops/{stopId} [delete]
func (r *Router) deleteHubStop(w http.ResponseWriter, req *http.Request) {
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

	var roomID pgtype.UUID
	if err := roomID.Scan(req.PathValue("id")); err != nil {
		writeError(w, http.StatusBadRequest, "유효하지 않은 방 ID입니다.")
		return
	}

	var stopID pgtype.UUID
	if err := stopID.Scan(req.PathValue("stopId")); err != nil {
		writeError(w, http.StatusBadRequest, "유효하지 않은 거점 ID입니다.")
		return
	}

	_, err := r.queries.GetRoomMember(req.Context(), db.GetRoomMemberParams{
		RoomID: roomID,
		UserID: userID,
	})
	if err != nil {
		writeError(w, http.StatusForbidden, "방에 참여하지 않았습니다.")
		return
	}

	_, err = r.queries.DeleteHubStop(req.Context(), db.DeleteHubStopParams{
		ID:     stopID,
		RoomID: roomID,
	})
	if err != nil {
		writeError(w, http.StatusNotFound, "거점을 찾을 수 없습니다.")
		return
	}

	writeJSON(w, http.StatusOK, map[string]string{
		"message": "거점이 삭제되었습니다.",
	})
}
