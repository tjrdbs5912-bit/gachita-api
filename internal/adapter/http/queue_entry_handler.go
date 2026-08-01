package httpadapter

import (
	"context"
	"net/http"
	"strconv"
	"time"

	"gachita-api/internal/db"

	"github.com/jackc/pgx/v5/pgtype"
)

type createQueueEntryRequest struct {
	FromStopID string    `json:"from_stop_id"`
	ToStopID   string    `json:"to_stop_id"`
	TimeStart  time.Time `json:"time_start"`
	TimeEnd    time.Time `json:"time_end"`
	MinSeats   *int32    `json:"min_seats"`
	MaxSeats   *int32    `json:"max_seats"`
}

func queueEntryToMap(e db.QueueEntry) map[string]string {
	return map[string]string{
		"id":           e.ID.String(),
		"room_id":      e.RoomID.String(),
		"user_id":      e.UserID.String(),
		"from_stop_id": e.FromStopID.String(),
		"to_stop_id":   e.ToStopID.String(),
		"time_start":   e.TimeStart.Time.Format(time.RFC3339),
		"time_end":     e.TimeEnd.Time.Format(time.RFC3339),
		"min_seats":    strconv.Itoa(int(e.MinSeats)),
		"max_seats":    strconv.Itoa(int(e.MaxSeats)),
		"status":       e.Status,
	}
}

// CreateQueueEntry godoc
// @Summary      매칭 대기 등록
// @Tags         queue
// @Accept       json
// @Produce      json
// @Security     BearerAuth
// @Param        id    path  string  true  "방 ID"
// @Param        body  body  createQueueEntryRequest  true  "대기 정보"
// @Success      201  {object}  map[string]string
// @Router       /api/rooms/{id}/queue [post]
func (r *Router) createQueueEntry(w http.ResponseWriter, req *http.Request) {
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

	var body createQueueEntryRequest
	if err := decodeJSON(req, &body); err != nil {
		writeError(w, http.StatusBadRequest, "invalid json")
		return
	}

	var fromStopID, toStopID pgtype.UUID
	if err := fromStopID.Scan(body.FromStopID); err != nil {
		writeError(w, http.StatusBadRequest, "유효하지 않은 출발 거점 ID입니다.")
		return
	}
	if err := toStopID.Scan(body.ToStopID); err != nil {
		writeError(w, http.StatusBadRequest, "유효하지 않은 도착 거점 ID입니다.")
		return
	}
	if body.FromStopID == body.ToStopID {
		writeError(w, http.StatusBadRequest, "출발 거점과 도착 거점이 같을 수 없습니다.")
		return
	}

	if body.TimeStart.IsZero() || body.TimeEnd.IsZero() {
		writeError(w, http.StatusBadRequest, "time_start, time_end는 필수입니다.")
		return
	}
	if !body.TimeStart.Before(body.TimeEnd) {
		writeError(w, http.StatusBadRequest, "time_start는 time_end보다 빨라야 합니다.")
		return
	}

	minSeats := int32(2)
	maxSeats := int32(4)
	if body.MinSeats != nil {
		minSeats = *body.MinSeats
	}
	if body.MaxSeats != nil {
		maxSeats = *body.MaxSeats
	}
	if minSeats < 1 || maxSeats < minSeats {
		writeError(w, http.StatusBadRequest, "인원 설정이 올바르지 않습니다. (1 <= min <= max)")
		return
	}

	// 거점이 이 방 소속인지 확인
	fromStop, err := r.queries.GetHubStopByID(req.Context(), fromStopID)
	if err != nil || fromStop.RoomID != roomID {
		writeError(w, http.StatusBadRequest, "출발 거점이 이 방에 없습니다.")
		return
	}
	toStop, err := r.queries.GetHubStopByID(req.Context(), toStopID)
	if err != nil || toStop.RoomID != roomID {
		writeError(w, http.StatusBadRequest, "도착 거점이 이 방에 없습니다.")
		return
	}

	ctx := req.Context()
	tx, err := r.pool.Begin(ctx)
	if err != nil {
		writeError(w, http.StatusInternalServerError, "대기 등록에 실패했습니다.")
		return
	}
	defer tx.Rollback(ctx)
	qtx := r.queries.WithTx(tx)

	entry, err := qtx.CreateQueueEntry(ctx, db.CreateQueueEntryParams{
		RoomID:     roomID,
		UserID:     userID,
		FromStopID: fromStopID,
		ToStopID:   toStopID,
		TimeStart:  pgtype.Timestamptz{Time: body.TimeStart, Valid: true},
		TimeEnd:    pgtype.Timestamptz{Time: body.TimeEnd, Valid: true},
		MinSeats:   minSeats,
		MaxSeats:   maxSeats,
	})
	if err != nil {
		writeError(w, http.StatusInternalServerError, "대기 등록에 실패했습니다.")
		return
	}

	match, matched, err := tryMatch(ctx, qtx, entry)
	if err != nil {
		writeError(w, http.StatusInternalServerError, "매칭 처리에 실패했습니다.")
		return
	}

	if err := tx.Commit(ctx); err != nil {
		writeError(w, http.StatusInternalServerError, "대기 등록에 실패했습니다.")
		return
	}

	resp := queueEntryToMap(entry)
	if matched {
		resp["status"] = "matched"
		resp["match_id"] = match.ID.String()
	}
	writeJSON(w, http.StatusCreated, resp)
}

// tryMatch는 새 대기(entry)와 같은 구간·겹치는 시간의 waiting 대기들을 모아,
// 인원 조건(min~max)이 맞으면 매칭을 생성한다.
func tryMatch(ctx context.Context, q *db.Queries, entry db.QueueEntry) (db.Match, bool, error) {
	candidates, err := q.ListMatchCandidates(ctx, db.ListMatchCandidatesParams{
		RoomID:     entry.RoomID,
		FromStopID: entry.FromStopID,
		ToStopID:   entry.ToStopID,
		TimeStart:  entry.TimeEnd,
		TimeEnd:    entry.TimeStart,
	})
	if err != nil {
		return db.Match{}, false, err
	}

	if len(candidates) < 2 {
		return db.Match{}, false, nil
	}

	//후보들의 max_seats 중 가장 작은 값
	capSize := candidates[0].MaxSeats
	for _, c := range candidates {
		if c.MaxSeats < capSize {
			capSize = c.MaxSeats
		}
	}
	if int32(len(candidates)) < capSize {
		capSize = int32(len(candidates))
	}

	group := candidates[:capSize]

	// 새 대기가 그룹에 포함되지 않으면 이번엔 매칭하지 않음
	included := false
	for _, g := range group {
		if g.ID == entry.ID {
			included = true
			break
		}
	}
	if !included {
		return db.Match{}, false, nil
	}

	need := group[0].MinSeats
	for _, g := range group {
		if g.MinSeats > need {
			need = g.MinSeats
		}
	}
	if int32(len(group)) < need {
		return db.Match{}, false, nil
	}

	match, err := q.CreateMatch(ctx, entry.RoomID)
	if err != nil {
		return db.Match{}, false, err
	}

	entryIDs := make([]pgtype.UUID, 0, len(group))
	for _, g := range group {
		if err := q.AddMatchMember(ctx, db.AddMatchMemberParams{
			MatchID:      match.ID,
			UserID:       g.UserID,
			QueueEntryID: g.ID,
		}); err != nil {
			return db.Match{}, false, err
		}
		entryIDs = append(entryIDs, g.ID)
	}

	if err := q.MarkQueueEntriesMatched(ctx, entryIDs); err != nil {
		return db.Match{}, false, err
	}

	// 매칭된 멤버 전원에게 알림 생성
	for _, g := range group {
		if _, err := q.CreateNotification(ctx, db.CreateNotificationParams{
			UserID:  g.UserID,
			Type:    "match_confirmed",
			Title:   "매칭이 성사되었습니다!",
			Body:    "같이 탈 사람이 모였어요. 매칭 상세에서 확인하세요.",
			RefType: pgtype.Text{String: "match", Valid: true},
			RefID:   match.ID,
		}); err != nil {
			return db.Match{}, false, err
		}
	}

	return match, true, nil
}

// ListQueueEntries godoc
// @Summary      매칭 대기 목록
// @Tags         queue
// @Produce      json
// @Security     BearerAuth
// @Param        id   path  string  true  "방 ID"
// @Success      200  {array}  map[string]string
// @Router       /api/rooms/{id}/queue [get]
func (r *Router) listQueueEntries(w http.ResponseWriter, req *http.Request) {
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

	entries, err := r.queries.ListWaitingQueueEntriesByRoomID(req.Context(), roomID)
	if err != nil {
		writeError(w, http.StatusInternalServerError, "대기 목록 조회에 실패했습니다.")
		return
	}

	list := make([]map[string]string, 0, len(entries))
	for _, e := range entries {
		list = append(list, queueEntryToMap(e))
	}
	writeJSON(w, http.StatusOK, list)
}

// CancelQueueEntry godoc
// @Summary      매칭 대기 취소
// @Tags         queue
// @Produce      json
// @Security     BearerAuth
// @Param        id       path  string  true  "방 ID"
// @Param        entryId  path  string  true  "대기 ID"
// @Success      200  {object}  map[string]string
// @Router       /api/rooms/{id}/queue/{entryId} [delete]
func (r *Router) cancelQueueEntry(w http.ResponseWriter, req *http.Request) {
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

	var entryID pgtype.UUID
	if err := entryID.Scan(req.PathValue("entryId")); err != nil {
		writeError(w, http.StatusBadRequest, "유효하지 않은 대기 ID입니다.")
		return
	}

	entry, err := r.queries.CancelQueueEntry(req.Context(), db.CancelQueueEntryParams{
		ID:     entryID,
		RoomID: roomID,
		UserID: userID,
	})
	if err != nil {
		writeError(w, http.StatusNotFound, "취소할 대기가 없습니다.")
		return
	}

	writeJSON(w, http.StatusOK, queueEntryToMap(entry))
}
