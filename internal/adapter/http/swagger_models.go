package httpadapter

// ErrorResponse 공통 에러 응답
type ErrorResponse struct {
	Error string `json:"error"`
} // @name ErrorResponse

// MessageResponse 단순 메시지 응답
type MessageResponse struct {
	Message string `json:"message"`
} // @name MessageResponse

// UserResponse 사용자 정보
type UserResponse struct {
	ID       string `json:"id"`
	Email    string `json:"email"`
	Nickname string `json:"nickname"`
} // @name UserResponse

// LoginResponse 로그인 응답
type LoginResponse struct {
	AccessToken string `json:"access_token"`
	TokenType   string `json:"token_type"`
} // @name LoginResponse

// RoomResponse 방 정보
type RoomResponse struct {
	ID          string `json:"id"`
	Name        string `json:"name"`
	InviteCode  string `json:"invite_code"`
	OpenchatURL string `json:"openchat_url"`
	OwnerID     string `json:"owner_id"`
} // @name RoomResponse

// JoinRoomResponse 방 입장 응답
type JoinRoomResponse struct {
	ID         string `json:"id"`
	Name       string `json:"name"`
	InviteCode string `json:"invite_code"`
} // @name JoinRoomResponse

// RoomMemberResponse 방 멤버
type RoomMemberResponse struct {
	UserID   string `json:"user_id"`
	Nickname string `json:"nickname"`
	Email    string `json:"email"`
} // @name RoomMemberResponse

// RoomDetailResponse 방 상세 (멤버 포함)
type RoomDetailResponse struct {
	ID          string               `json:"id"`
	Name        string               `json:"name"`
	InviteCode  string               `json:"invite_code"`
	OpenchatURL string               `json:"openchat_url"`
	OwnerID     string               `json:"owner_id"`
	Members     []RoomMemberResponse `json:"members"`
} // @name RoomDetailResponse

// HubStopResponse 거점
type HubStopResponse struct {
	ID        string `json:"id"`
	RoomID    string `json:"room_id"`
	Name      string `json:"name"`
	SortOrder string `json:"sort_order"`
} // @name HubStopResponse

// QueueEntryResponse 매칭 대기
type QueueEntryResponse struct {
	ID         string `json:"id"`
	RoomID     string `json:"room_id"`
	UserID     string `json:"user_id"`
	FromStopID string `json:"from_stop_id"`
	ToStopID   string `json:"to_stop_id"`
	TimeStart  string `json:"time_start"`
	TimeEnd    string `json:"time_end"`
	MinSeats   string `json:"min_seats"`
	MaxSeats   string `json:"max_seats"`
	Status     string `json:"status"`
	MatchID    string `json:"match_id,omitempty"`
} // @name QueueEntryResponse

// MatchResponse 매칭 요약
type MatchResponse struct {
	ID        string `json:"id"`
	RoomID    string `json:"room_id"`
	Status    string `json:"status"`
	CreatedAt string `json:"created_at"`
} // @name MatchResponse

// MatchMemberResponse 매칭 멤버
type MatchMemberResponse struct {
	UserID       string `json:"user_id"`
	Nickname     string `json:"nickname"`
	QueueEntryID string `json:"queue_entry_id"`
	FromStopID   string `json:"from_stop_id"`
	ToStopID     string `json:"to_stop_id"`
	TimeStart    string `json:"time_start"`
	TimeEnd      string `json:"time_end"`
} // @name MatchMemberResponse

// MatchDetailResponse 매칭 상세
type MatchDetailResponse struct {
	ID        string                `json:"id"`
	RoomID    string                `json:"room_id"`
	Status    string                `json:"status"`
	CreatedAt string                `json:"created_at"`
	Members   []MatchMemberResponse `json:"members"`
} // @name MatchDetailResponse

// NotificationResponse 알림
type NotificationResponse struct {
	ID        string `json:"id"`
	Type      string `json:"type"`
	Title     string `json:"title"`
	Body      string `json:"body"`
	CreatedAt string `json:"created_at"`
	ReadAt    string `json:"read_at"`
	RefType   string `json:"ref_type"`
	RefID     string `json:"ref_id"`
} // @name NotificationResponse
