package httpadapter

import (
	"time"

	"gachita-api/internal/db"

	"github.com/jackc/pgx/v5/pgtype"
)

// --- requests ---

type SignUpRequest struct {
	Email    string `json:"email" binding:"required" example:"user@example.com"`
	Password string `json:"password" binding:"required" example:"password123"`
	Nickname string `json:"nickname" binding:"required" example:"가지"`
} // @name SignUpRequest

type LoginRequest struct {
	Email    string `json:"email" binding:"required" example:"user@example.com"`
	Password string `json:"password" binding:"required" example:"password123"`
} // @name LoginRequest

type CreateRoomRequest struct {
	Name        string  `json:"name" binding:"required" example:"학교 카풀"`
	OpenchatURL *string `json:"openchat_url" extensions:"x-nullable" example:"https://open.kakao.com/o/xxxxx"`
} // @name CreateRoomRequest

type JoinRoomRequest struct {
	InviteCode string `json:"invite_code" binding:"required" example:"a1b2c3d4"`
} // @name JoinRoomRequest

type CreateHubStopRequest struct {
	Name      string `json:"name" binding:"required" example:"정문"`
	SortOrder *int32 `json:"sort_order" extensions:"x-nullable" example:"0"`
} // @name CreateHubStopRequest

type UpdateHubStopRequest struct {
	Name      string `json:"name" binding:"required" example:"후문"`
	SortOrder *int32 `json:"sort_order" extensions:"x-nullable" example:"1"`
} // @name UpdateHubStopRequest

type CreateQueueEntryRequest struct {
	FromStopID string    `json:"from_stop_id" binding:"required" example:"550e8400-e29b-41d4-a716-446655440000"`
	ToStopID   string    `json:"to_stop_id" binding:"required" example:"550e8400-e29b-41d4-a716-446655440001"`
	TimeStart  time.Time `json:"time_start" binding:"required" example:"2026-08-02T09:00:00Z"`
	TimeEnd    time.Time `json:"time_end" binding:"required" example:"2026-08-02T10:00:00Z"`
	MinSeats   *int32    `json:"min_seats" extensions:"x-nullable" example:"2"`
	MaxSeats   *int32    `json:"max_seats" extensions:"x-nullable" example:"4"`
} // @name CreateQueueEntryRequest

// --- responses ---

type ErrorResponse struct {
	Error string `json:"error" binding:"required" example:"invalid json"`
} // @name ErrorResponse

type MessageResponse struct {
	Message string `json:"message" binding:"required" example:"ok"`
} // @name MessageResponse

type UserResponse struct {
	ID       string `json:"id" binding:"required" example:"550e8400-e29b-41d4-a716-446655440000"`
	Email    string `json:"email" binding:"required" example:"user@example.com"`
	Nickname string `json:"nickname" binding:"required" example:"가지"`
} // @name UserResponse

type LoginResponse struct {
	AccessToken string `json:"access_token" binding:"required"`
	TokenType   string `json:"token_type" binding:"required" example:"Bearer"`
} // @name LoginResponse

type RoomResponse struct {
	ID          string  `json:"id" binding:"required" example:"550e8400-e29b-41d4-a716-446655440000"`
	Name        string  `json:"name" binding:"required" example:"학교 카풀"`
	InviteCode  string  `json:"invite_code" binding:"required" example:"a1b2c3d4"`
	OpenchatURL *string `json:"openchat_url" extensions:"x-nullable"`
	OwnerID     string  `json:"owner_id" binding:"required" example:"550e8400-e29b-41d4-a716-446655440000"`
} // @name RoomResponse

type JoinRoomResponse struct {
	ID         string `json:"id" binding:"required"`
	Name       string `json:"name" binding:"required"`
	InviteCode string `json:"invite_code" binding:"required"`
} // @name JoinRoomResponse

type RoomMemberResponse struct {
	UserID   string `json:"user_id" binding:"required"`
	Nickname string `json:"nickname" binding:"required"`
	Email    string `json:"email" binding:"required"`
} // @name RoomMemberResponse

type RoomDetailResponse struct {
	ID          string               `json:"id" binding:"required"`
	Name        string               `json:"name" binding:"required"`
	InviteCode  string               `json:"invite_code" binding:"required"`
	OpenchatURL *string              `json:"openchat_url" extensions:"x-nullable"`
	OwnerID     string               `json:"owner_id" binding:"required"`
	Members     []RoomMemberResponse `json:"members" binding:"required"`
} // @name RoomDetailResponse

type HubStopResponse struct {
	ID        string `json:"id" binding:"required"`
	RoomID    string `json:"room_id" binding:"required"`
	Name      string `json:"name" binding:"required"`
	SortOrder int32  `json:"sort_order" binding:"required" example:"0"`
} // @name HubStopResponse

type QueueEntryResponse struct {
	ID         string  `json:"id" binding:"required"`
	RoomID     string  `json:"room_id" binding:"required"`
	UserID     string  `json:"user_id" binding:"required"`
	FromStopID string  `json:"from_stop_id" binding:"required"`
	ToStopID   string  `json:"to_stop_id" binding:"required"`
	TimeStart  string  `json:"time_start" binding:"required" example:"2026-08-02T09:00:00Z"`
	TimeEnd    string  `json:"time_end" binding:"required" example:"2026-08-02T10:00:00Z"`
	MinSeats   int32   `json:"min_seats" binding:"required" example:"2"`
	MaxSeats   int32   `json:"max_seats" binding:"required" example:"4"`
	Status     string  `json:"status" binding:"required" enums:"waiting,matched,cancelled"`
	MatchID    *string `json:"match_id" extensions:"x-nullable"`
} // @name QueueEntryResponse

type MatchResponse struct {
	ID        string `json:"id" binding:"required"`
	RoomID    string `json:"room_id" binding:"required"`
	Status    string `json:"status" binding:"required"`
	CreatedAt string `json:"created_at" binding:"required" example:"2026-08-02T09:30:00Z"`
} // @name MatchResponse

type MatchMemberResponse struct {
	UserID       string `json:"user_id" binding:"required"`
	Nickname     string `json:"nickname" binding:"required"`
	QueueEntryID string `json:"queue_entry_id" binding:"required"`
	FromStopID   string `json:"from_stop_id" binding:"required"`
	ToStopID     string `json:"to_stop_id" binding:"required"`
	TimeStart    string `json:"time_start" binding:"required"`
	TimeEnd      string `json:"time_end" binding:"required"`
} // @name MatchMemberResponse

type MatchDetailResponse struct {
	ID        string                `json:"id" binding:"required"`
	RoomID    string                `json:"room_id" binding:"required"`
	Status    string                `json:"status" binding:"required"`
	CreatedAt string                `json:"created_at" binding:"required"`
	Members   []MatchMemberResponse `json:"members" binding:"required"`
} // @name MatchDetailResponse

type NotificationResponse struct {
	ID        string  `json:"id" binding:"required"`
	Type      string  `json:"type" binding:"required"`
	Title     string  `json:"title" binding:"required"`
	Body      string  `json:"body" binding:"required"`
	CreatedAt string  `json:"created_at" binding:"required"`
	ReadAt    *string `json:"read_at" extensions:"x-nullable"`
	RefType   *string `json:"ref_type" extensions:"x-nullable"`
	RefID     *string `json:"ref_id" extensions:"x-nullable"`
} // @name NotificationResponse

// --- mappers ---

func textPtr(t pgtype.Text) *string {
	if !t.Valid {
		return nil
	}
	s := t.String
	return &s
}

func strPtr(s string) *string { return &s }

func toUserResponse(id, email, nickname string) UserResponse {
	return UserResponse{ID: id, Email: email, Nickname: nickname}
}

func toRoomResponse(room db.Room) RoomResponse {
	return RoomResponse{
		ID:          room.ID.String(),
		Name:        room.Name,
		InviteCode:  room.InviteCode,
		OpenchatURL: textPtr(room.OpenchatUrl),
		OwnerID:     room.OwnerID.String(),
	}
}

func toHubStopResponse(s db.HubStop) HubStopResponse {
	return HubStopResponse{
		ID:        s.ID.String(),
		RoomID:    s.RoomID.String(),
		Name:      s.Name,
		SortOrder: s.SortOrder,
	}
}

func toQueueEntryResponse(e db.QueueEntry) QueueEntryResponse {
	return QueueEntryResponse{
		ID:         e.ID.String(),
		RoomID:     e.RoomID.String(),
		UserID:     e.UserID.String(),
		FromStopID: e.FromStopID.String(),
		ToStopID:   e.ToStopID.String(),
		TimeStart:  e.TimeStart.Time.Format(time.RFC3339),
		TimeEnd:    e.TimeEnd.Time.Format(time.RFC3339),
		MinSeats:   e.MinSeats,
		MaxSeats:   e.MaxSeats,
		Status:     e.Status,
	}
}

func queueEntryFromActiveRow(e db.ListMyActiveQueueEntriesRow) QueueEntryResponse {
	resp := QueueEntryResponse{
		ID:         e.ID.String(),
		RoomID:     e.RoomID.String(),
		UserID:     e.UserID.String(),
		FromStopID: e.FromStopID.String(),
		ToStopID:   e.ToStopID.String(),
		TimeStart:  e.TimeStart.Time.Format(time.RFC3339),
		TimeEnd:    e.TimeEnd.Time.Format(time.RFC3339),
		MinSeats:   e.MinSeats,
		MaxSeats:   e.MaxSeats,
		Status:     e.Status,
	}
	if e.MatchID.Valid {
		resp.MatchID = strPtr(e.MatchID.String())
	}
	return resp
}

func queueEntryFromRoomRow(e db.GetQueueEntryByRoomIDRow) QueueEntryResponse {
	resp := QueueEntryResponse{
		ID:         e.ID.String(),
		RoomID:     e.RoomID.String(),
		UserID:     e.UserID.String(),
		FromStopID: e.FromStopID.String(),
		ToStopID:   e.ToStopID.String(),
		TimeStart:  e.TimeStart.Time.Format(time.RFC3339),
		TimeEnd:    e.TimeEnd.Time.Format(time.RFC3339),
		MinSeats:   e.MinSeats,
		MaxSeats:   e.MaxSeats,
		Status:     e.Status,
	}
	if e.MatchID.Valid {
		resp.MatchID = strPtr(e.MatchID.String())
	}
	return resp
}

func toMatchResponse(m db.Match) MatchResponse {
	return MatchResponse{
		ID:        m.ID.String(),
		RoomID:    m.RoomID.String(),
		Status:    m.Status,
		CreatedAt: m.CreatedAt.Time.Format(time.RFC3339),
	}
}

func toNotificationResponse(n db.Notification) NotificationResponse {
	resp := NotificationResponse{
		ID:        n.ID.String(),
		Type:      n.Type,
		Title:     n.Title,
		Body:      n.Body,
		CreatedAt: n.CreatedAt.Time.Format(time.RFC3339),
	}
	if n.ReadAt.Valid {
		s := n.ReadAt.Time.Format(time.RFC3339)
		resp.ReadAt = &s
	}
	if n.RefType.Valid {
		s := n.RefType.String
		resp.RefType = &s
	}
	if n.RefID.Valid {
		s := n.RefID.String()
		resp.RefID = &s
	}
	return resp
}
