package httpadapter

import (
	"errors"
	"gachita-api/internal/db"
	"net/http"

	"crypto/rand"
	"encoding/hex"

	"github.com/jackc/pgx/v5/pgconn"
	"github.com/jackc/pgx/v5/pgtype"
)

// CreateRoom godoc
// @Summary      방 생성
// @Tags         room
// @Accept       json
// @Produce      json
// @Security     BearerAuth
// @Param        body  body  CreateRoomRequest  true  "방 정보"
// @Success      201  {object}  RoomResponse
// @Router       /api/rooms [post]
func (r *Router) createRoom(w http.ResponseWriter, req *http.Request) {
	var body CreateRoomRequest
	if err := decodeJSON(req, &body); err != nil {
		writeError(w, http.StatusBadRequest, "invalid json")
		return
	}

	if body.Name == "" {
		writeError(w, http.StatusBadRequest, "제목은 필수 입력 항목입니다.")
		return
	}

	userIDStr, ok := req.Context().Value(contextKeyUserID).(string)
	if !ok || userIDStr == "" {
		writeError(w, http.StatusUnauthorized, "인증 정보가 없습니다.")
		return
	}

	var ownerID pgtype.UUID
	if err := ownerID.Scan(userIDStr); err != nil {
		writeError(w, http.StatusBadRequest, "유효하지 않은 사용자 ID입니다.")
		return
	}

	inviteCode, err := generateInviteCode()
	if err != nil {
		writeError(w, http.StatusInternalServerError, "초대 코드 생성에 실패했습니다.")
		return
	}

	var openchat pgtype.Text
	if body.OpenchatURL != nil && *body.OpenchatURL != "" {
		openchat = pgtype.Text{String: *body.OpenchatURL, Valid: true}
	}
	room, err := r.queries.CreateRoom(req.Context(), db.CreateRoomParams{
		Name:        body.Name,
		InviteCode:  inviteCode,
		OpenchatUrl: openchat,
		OwnerID:     ownerID,
	})
	if err != nil {
		var pgErr *pgconn.PgError
		if errors.As(err, &pgErr) && pgErr.Code == "23505" {
			writeError(w, http.StatusConflict, "초대 코드가 이미 사용 중입니다.")
			return
		}
		writeError(w, http.StatusInternalServerError, "방 생성에 실패했습니다.")
		return
	}

	_, err = r.queries.AddRoomMember(req.Context(), db.AddRoomMemberParams{
		RoomID: room.ID,
		UserID: ownerID,
	})
	if err != nil {
		writeError(w, http.StatusInternalServerError, "방 멤버 추가에 실패했습니다.")
		return
	}

	writeJSON(w, http.StatusCreated, toRoomResponse(room))

}

// JoinRoom godoc
// @Summary      방 입장
// @Tags         room
// @Accept       json
// @Produce      json
// @Security     BearerAuth
// @Param        body  body  JoinRoomRequest  true  "초대 코드"
// @Success      200  {object}  JoinRoomResponse
// @Router       /api/rooms/join [post]
func (r *Router) joinRoom(w http.ResponseWriter, req *http.Request) {
	var body JoinRoomRequest
	if err := decodeJSON(req, &body); err != nil {
		writeError(w, http.StatusBadRequest, "invalid json")
		return
	}
	if body.InviteCode == "" {
		writeError(w, http.StatusBadRequest, "초대 코드는 필수입니다.")
		return
	}

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

	room, err := r.queries.GetRoomByInviteCode(req.Context(), body.InviteCode)
	if err != nil {
		writeError(w, http.StatusNotFound, "유효하지 않은 초대 코드입니다.")
		return
	}

	_, err = r.queries.AddRoomMember(req.Context(), db.AddRoomMemberParams{
		RoomID: room.ID,
		UserID: userID,
	})
	if err != nil {
		var pgErr *pgconn.PgError
		if errors.As(err, &pgErr) && pgErr.Code == "23505" {
			writeError(w, http.StatusConflict, "이미 참여 중인 방입니다.")
			return
		}
		writeError(w, http.StatusInternalServerError, "방 입장에 실패했습니다.")
		return
	}

	writeJSON(w, http.StatusOK, JoinRoomResponse{
		ID:         room.ID.String(),
		Name:       room.Name,
		InviteCode: room.InviteCode,
	})
}

// GetRoom godoc
// @Summary      방 조회
// @Tags         room
// @Produce      json
// @Security     BearerAuth
// @Param        id   path  string  true  "방 ID"
// @Success      200  {object}  RoomDetailResponse
// @Router       /api/rooms/{id} [get]
func (r *Router) getRoom(w http.ResponseWriter, req *http.Request) {
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

	// 멤버만 조회 가능
	_, err := r.queries.GetRoomMember(req.Context(), db.GetRoomMemberParams{
		RoomID: roomID,
		UserID: userID,
	})
	if err != nil {
		writeError(w, http.StatusForbidden, "방에 참여하지 않았습니다.")
		return
	}

	room, err := r.queries.GetRoomByID(req.Context(), roomID)
	if err != nil {
		writeError(w, http.StatusNotFound, "방을 찾을 수 없습니다.")
		return
	}

	members, err := r.queries.ListRoomMembersWithUser(req.Context(), roomID)
	if err != nil {
		writeError(w, http.StatusInternalServerError, "멤버 조회에 실패했습니다.")
		return
	}

	memberList := make([]RoomMemberResponse, 0, len(members))
	for _, m := range members {
		memberList = append(memberList, RoomMemberResponse{
			UserID:   m.UserID.String(),
			Nickname: m.Nickname,
			Email:    m.Email,
		})
	}

	writeJSON(w, http.StatusOK, RoomDetailResponse{
		ID:          room.ID.String(),
		Name:        room.Name,
		InviteCode:  room.InviteCode,
		OpenchatURL: textPtr(room.OpenchatUrl),
		OwnerID:     room.OwnerID.String(),
		Members:     memberList,
	})
}

// ListRooms godoc
// @Summary      내 방 목록
// @Tags         room
// @Produce      json
// @Security     BearerAuth
// @Success      200  {array}   RoomResponse
// @Router       /api/rooms [get]
func (r *Router) listRooms(w http.ResponseWriter, req *http.Request) {
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

	rooms, err := r.queries.ListRoomsByUserID(req.Context(), userID)
	if err != nil {
		writeError(w, http.StatusInternalServerError, "방 목록 조회에 실패했습니다.")
		return
	}

	list := make([]RoomResponse, 0, len(rooms))
	for _, room := range rooms {
		list = append(list, toRoomResponse(room))
	}
	writeJSON(w, http.StatusOK, list)
}
func generateInviteCode() (string, error) {
	b := make([]byte, 4) // 8자 hex
	if _, err := rand.Read(b); err != nil {
		return "", err
	}
	return hex.EncodeToString(b), nil
}
