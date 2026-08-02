package httpadapter

import (
	"net/http"
	"time"

	"gachita-api/internal/db"

	"github.com/jackc/pgx/v5/pgtype"
)

// ListMatches godoc
// @Summary      매칭 목록
// @Tags         match
// @Produce      json
// @Security     BearerAuth
// @Param        id   path  string  true  "방 ID"
// @Success      200  {array}   MatchResponse
// @Router       /api/rooms/{id}/matches [get]
func (r *Router) listMatches(w http.ResponseWriter, req *http.Request) {
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

	matches, err := r.queries.ListMatchesByRoomID(req.Context(), roomID)
	if err != nil {
		writeError(w, http.StatusInternalServerError, "매칭 목록 조회에 실패했습니다.")
		return
	}

	list := make([]map[string]string, 0, len(matches))
	for _, m := range matches {
		list = append(list, map[string]string{
			"id":         m.ID.String(),
			"room_id":    m.RoomID.String(),
			"status":     m.Status,
			"created_at": m.CreatedAt.Time.Format(time.RFC3339),
		})
	}
	writeJSON(w, http.StatusOK, list)
}

// GetMatch godoc
// @Summary      매칭 상세 (멤버 포함)
// @Tags         match
// @Produce      json
// @Security     BearerAuth
// @Param        id       path  string  true  "방 ID"
// @Param        matchId  path  string  true  "매칭 ID"
// @Success      200  {object}  MatchDetailResponse
// @Router       /api/rooms/{id}/matches/{matchId} [get]
func (r *Router) getMatch(w http.ResponseWriter, req *http.Request) {
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

	var matchID pgtype.UUID
	if err := matchID.Scan(req.PathValue("matchId")); err != nil {
		writeError(w, http.StatusBadRequest, "유효하지 않은 매칭 ID입니다.")
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

	match, err := r.queries.GetMatchByID(req.Context(), matchID)
	if err != nil || match.RoomID != roomID {
		writeError(w, http.StatusNotFound, "매칭을 찾을 수 없습니다.")
		return
	}

	members, err := r.queries.ListMatchMembersByMatchID(req.Context(), matchID)
	if err != nil {
		writeError(w, http.StatusInternalServerError, "매칭 멤버 조회에 실패했습니다.")
		return
	}

	memberList := make([]map[string]string, 0, len(members))
	for _, m := range members {
		memberList = append(memberList, map[string]string{
			"user_id":        m.UserID.String(),
			"nickname":       m.Nickname,
			"queue_entry_id": m.QueueEntryID.String(),
			"from_stop_id":   m.FromStopID.String(),
			"to_stop_id":     m.ToStopID.String(),
			"time_start":     m.TimeStart.Time.Format(time.RFC3339),
			"time_end":       m.TimeEnd.Time.Format(time.RFC3339),
		})
	}

	writeJSON(w, http.StatusOK, map[string]interface{}{
		"id":         match.ID.String(),
		"room_id":    match.RoomID.String(),
		"status":     match.Status,
		"created_at": match.CreatedAt.Time.Format(time.RFC3339),
		"members":    memberList,
	})
}
