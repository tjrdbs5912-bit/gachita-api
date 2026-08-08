package main

import (
	"fmt"
	"os"
	"path/filepath"
	"regexp"
	"strings"
)

func main() {
	normalizeHandlers()
	fixAuth()
	fixUser()
	fixRoom()
	fixHub()
	fixQueue()
	fixMatch()
	fixNotif()
	fixResponse()
	removeSwaggerModels()
	fmt.Println("restored")
}

func normalizeHandlers() {
	dir := "internal/adapter/http"
	entries, err := os.ReadDir(dir)
	if err != nil {
		panic(err)
	}
	for _, e := range entries {
		if !strings.HasSuffix(e.Name(), "_handler.go") {
			continue
		}
		p := filepath.Join(dir, e.Name())
		b, err := os.ReadFile(p)
		if err != nil {
			panic(err)
		}
		s := strings.ReplaceAll(string(b), "\r\n", "\n")
		if err := os.WriteFile(p, []byte(s), 0o644); err != nil {
			panic(err)
		}
	}
}

func rw(path string, fn func(string) string) {
	b, err := os.ReadFile(path)
	if err != nil {
		panic(err)
	}
	s := fn(string(b))
	if err := os.WriteFile(path, []byte(s), 0o644); err != nil {
		panic(err)
	}
	fmt.Println("updated", path)
}

func stripType(s, name string) string {
	re := regexp.MustCompile(`(?s)type ` + name + ` struct \{.*?\}\n\n`)
	return re.ReplaceAllString(s, "")
}

func fixResponse() {
	rw("internal/adapter/http/response.go", func(s string) string {
		return strings.ReplaceAll(s, `writeJSON(w, status, map[string]string{"error": msg})`, `writeJSON(w, status, ErrorResponse{Error: msg})`)
	})
}

func removeSwaggerModels() {
	_ = os.Remove("internal/adapter/http/swagger_models.go")
}

func fixAuth() {
	rw("internal/adapter/http/auth_handler.go", func(s string) string {
		s = stripType(s, "signUpRequest")
		s = stripType(s, "loginRequest")
		s = strings.ReplaceAll(s, "signUpRequest", "SignUpRequest")
		s = strings.ReplaceAll(s, "loginRequest", "LoginRequest")
		s = strings.ReplaceAll(s, `writeJSON(w, http.StatusCreated, map[string]string{
		"id":       user.ID.String(),
		"email":    user.Email,
		"nickname": user.Nickname,
	})`, `writeJSON(w, http.StatusCreated, toUserResponse(user.ID.String(), user.Email, user.Nickname))`)
		s = strings.ReplaceAll(s, `writeJSON(w, http.StatusOK, map[string]string{
		"access_token": accessToken,
		"token_type":   "Bearer",
	})`, `writeJSON(w, http.StatusOK, LoginResponse{
		AccessToken: accessToken,
		TokenType:   "Bearer",
	})`)
		return s
	})
}

func fixUser() {
	rw("internal/adapter/http/user_handler.go", func(s string) string {
		return strings.ReplaceAll(s, `writeJSON(w, http.StatusOK, map[string]string{
		"id":       user.ID.String(),
		"email":    user.Email,
		"nickname": user.Nickname,
	})`, `writeJSON(w, http.StatusOK, toUserResponse(user.ID.String(), user.Email, user.Nickname))`)
	})
}

func fixRoom() {
	rw("internal/adapter/http/room_handler.go", func(s string) string {
		s = stripType(s, "createRoomRequest")
		s = stripType(s, "joinRoomRequest")
		s = strings.ReplaceAll(s, "createRoomRequest", "CreateRoomRequest")
		s = strings.ReplaceAll(s, "joinRoomRequest", "JoinRoomRequest")

		oldCreate := `	room, err := r.queries.CreateRoom(req.Context(), db.CreateRoomParams{
		Name:       body.Name,
		InviteCode: inviteCode,
		OpenchatUrl: pgtype.Text{
			String: body.OpenchatUrl,
			Valid:  body.OpenchatUrl != "",
		},
		OwnerID: ownerID,
	})`
		newCreate := `	var openchat pgtype.Text
	if body.OpenchatURL != nil && *body.OpenchatURL != "" {
		openchat = pgtype.Text{String: *body.OpenchatURL, Valid: true}
	}
	room, err := r.queries.CreateRoom(req.Context(), db.CreateRoomParams{
		Name:        body.Name,
		InviteCode:  inviteCode,
		OpenchatUrl: openchat,
		OwnerID:     ownerID,
	})`
		s = strings.ReplaceAll(s, oldCreate, newCreate)

		s = strings.ReplaceAll(s, `writeJSON(w, http.StatusCreated, map[string]string{
		"id":           room.ID.String(),
		"name":         room.Name,
		"invite_code":  room.InviteCode,
		"openchat_url": room.OpenchatUrl.String,
		"owner_id":     room.OwnerID.String(),
	})`, "writeJSON(w, http.StatusCreated, toRoomResponse(room))")

		s = strings.ReplaceAll(s, `writeJSON(w, http.StatusOK, map[string]string{
		"id":          room.ID.String(),
		"name":        room.Name,
		"invite_code": room.InviteCode,
	})`, `writeJSON(w, http.StatusOK, JoinRoomResponse{
		ID:         room.ID.String(),
		Name:       room.Name,
		InviteCode: room.InviteCode,
	})`)

		oldGet := `	memberList := make([]map[string]string, 0, len(members))
	for _, m := range members {
		memberList = append(memberList, map[string]string{
			"user_id":  m.UserID.String(),
			"nickname": m.Nickname,
			"email":    m.Email,
		})
	}

	writeJSON(w, http.StatusOK, map[string]interface{}{
		"id":           room.ID.String(),
		"name":         room.Name,
		"invite_code":  room.InviteCode,
		"openchat_url": room.OpenchatUrl.String,
		"owner_id":     room.OwnerID.String(),
		"members":      memberList,
	})`
		newGet := `	memberList := make([]RoomMemberResponse, 0, len(members))
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
	})`
		s = strings.ReplaceAll(s, oldGet, newGet)

		oldList := `	list := make([]map[string]string, 0, len(rooms))
	for _, room := range rooms {
		list = append(list, map[string]string{
			"id":           room.ID.String(),
			"name":         room.Name,
			"invite_code":  room.InviteCode,
			"openchat_url": room.OpenchatUrl.String,
			"owner_id":     room.OwnerID.String(),
		})
	}
	writeJSON(w, http.StatusOK, list)`
		newList := `	list := make([]RoomResponse, 0, len(rooms))
	for _, room := range rooms {
		list = append(list, toRoomResponse(room))
	}
	writeJSON(w, http.StatusOK, list)`
		s = strings.ReplaceAll(s, oldList, newList)
		return s
	})
}

func fixHub() {
	rw("internal/adapter/http/hub_stop_handler.go", func(s string) string {
		s = stripType(s, "createHubStopRequest")
		s = stripType(s, "updateHubStopRequest")
		s = strings.ReplaceAll(s, "createHubStopRequest", "CreateHubStopRequest")
		s = strings.ReplaceAll(s, "updateHubStopRequest", "UpdateHubStopRequest")
		s = strings.ReplaceAll(s, "\t\"strconv\"\n\n", "\n")
		s = strings.ReplaceAll(s, `writeJSON(w, http.StatusCreated, map[string]string{
		"id":         stop.ID.String(),
		"room_id":    stop.RoomID.String(),
		"name":       stop.Name,
		"sort_order": strconv.Itoa(int(stop.SortOrder)),
	})`, "writeJSON(w, http.StatusCreated, toHubStopResponse(stop))")
		s = strings.ReplaceAll(s, `	list := make([]map[string]string, 0, len(stops))
	for _, s := range stops {
		list = append(list, map[string]string{
			"id":         s.ID.String(),
			"room_id":    s.RoomID.String(),
			"name":       s.Name,
			"sort_order": strconv.Itoa(int(s.SortOrder)),
		})
	}
	writeJSON(w, http.StatusOK, list)`, `	list := make([]HubStopResponse, 0, len(stops))
	for _, stop := range stops {
		list = append(list, toHubStopResponse(stop))
	}
	writeJSON(w, http.StatusOK, list)`)
		s = strings.ReplaceAll(s, `writeJSON(w, http.StatusOK, map[string]string{
		"id":         stop.ID.String(),
		"room_id":    stop.RoomID.String(),
		"name":       stop.Name,
		"sort_order": strconv.Itoa(int(stop.SortOrder)),
	})`, "writeJSON(w, http.StatusOK, toHubStopResponse(stop))")
		re := regexp.MustCompile(`writeJSON\(w, http\.StatusOK, map\[string\]string\{\s*"message":\s*"([^"]*)",?\s*\}\)`)
		s = re.ReplaceAllString(s, `writeJSON(w, http.StatusOK, MessageResponse{Message: "$1"})`)
		return s
	})
}

func fixQueue() {
	rw("internal/adapter/http/queue_entry_handler.go", func(s string) string {
		s = stripType(s, "createQueueEntryRequest")
		s = strings.ReplaceAll(s, "createQueueEntryRequest", "CreateQueueEntryRequest")
		s = strings.ReplaceAll(s, "\t\"strconv\"\n", "")
		reMap := regexp.MustCompile(`(?s)func queueEntryToMap\(e db\.QueueEntry\) map\[string\]string \{.*?\n\}\n\n`)
		s = reMap.ReplaceAllString(s, "")
		s = strings.ReplaceAll(s, `	resp := queueEntryToMap(entry)
	if matched {
		resp["status"] = "matched"
		resp["match_id"] = match.ID.String()
	}
	writeJSON(w, http.StatusCreated, resp)`, `	resp := toQueueEntryResponse(entry)
	if matched {
		resp.Status = "matched"
		resp.MatchID = strPtr(match.ID.String())
	}
	writeJSON(w, http.StatusCreated, resp)`)
		s = strings.ReplaceAll(s, `	list := make([]map[string]string, 0, len(entries))
	for _, e := range entries {
		list = append(list, queueEntryToMap(e))
	}
	writeJSON(w, http.StatusOK, list)`, `	list := make([]QueueEntryResponse, 0, len(entries))
	for _, e := range entries {
		list = append(list, toQueueEntryResponse(e))
	}
	writeJSON(w, http.StatusOK, list)`)
		s = strings.ReplaceAll(s, "writeJSON(w, http.StatusOK, queueEntryToMap(entry))", "writeJSON(w, http.StatusOK, toQueueEntryResponse(entry))")
		if !strings.Contains(s, "time.") {
			s = strings.ReplaceAll(s, "\t\"time\"\n\n", "\n")
			s = strings.ReplaceAll(s, "\t\"time\"\n", "")
		}
		return s
	})
}

func fixMatch() {
	rw("internal/adapter/http/match_handler.go", func(s string) string {
		s = strings.ReplaceAll(s, `	list := make([]map[string]string, 0, len(matches))
	for _, m := range matches {
		list = append(list, map[string]string{
			"id":         m.ID.String(),
			"room_id":    m.RoomID.String(),
			"status":     m.Status,
			"created_at": m.CreatedAt.Time.Format(time.RFC3339),
		})
	}
	writeJSON(w, http.StatusOK, list)`, `	list := make([]MatchResponse, 0, len(matches))
	for _, m := range matches {
		list = append(list, toMatchResponse(m))
	}
	writeJSON(w, http.StatusOK, list)`)
		s = strings.ReplaceAll(s, `	memberList := make([]map[string]string, 0, len(members))
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
	})`, `	memberList := make([]MatchMemberResponse, 0, len(members))
	for _, m := range members {
		memberList = append(memberList, MatchMemberResponse{
			UserID:       m.UserID.String(),
			Nickname:     m.Nickname,
			QueueEntryID: m.QueueEntryID.String(),
			FromStopID:   m.FromStopID.String(),
			ToStopID:     m.ToStopID.String(),
			TimeStart:    m.TimeStart.Time.Format(time.RFC3339),
			TimeEnd:      m.TimeEnd.Time.Format(time.RFC3339),
		})
	}

	writeJSON(w, http.StatusOK, MatchDetailResponse{
		ID:        match.ID.String(),
		RoomID:    match.RoomID.String(),
		Status:    match.Status,
		CreatedAt: match.CreatedAt.Time.Format(time.RFC3339),
		Members:   memberList,
	})`)
		return s
	})
}

func fixNotif() {
	rw("internal/adapter/http/notification_handler.go", func(s string) string {
		reMap := regexp.MustCompile(`(?s)func notificationToMap\(n db\.Notification\) map\[string\]string \{.*?\n\}\n\n`)
		s = reMap.ReplaceAllString(s, "")
		s = strings.ReplaceAll(s, "\t\"time\"\n\n", "\n")
		s = strings.ReplaceAll(s, `	list := make([]map[string]string, 0, len(items))
	for _, n := range items {
		list = append(list, notificationToMap(n))
	}
	writeJSON(w, http.StatusOK, list)`, `	list := make([]NotificationResponse, 0, len(items))
	for _, n := range items {
		list = append(list, toNotificationResponse(n))
	}
	writeJSON(w, http.StatusOK, list)`)
		s = strings.ReplaceAll(s, "writeJSON(w, http.StatusOK, notificationToMap(n))", "writeJSON(w, http.StatusOK, toNotificationResponse(n))")
		re := regexp.MustCompile(`writeJSON\(w, http\.StatusOK, map\[string\]string\{\s*"message":\s*"([^"]*)",?\s*\}\)`)
		s = re.ReplaceAllString(s, `writeJSON(w, http.StatusOK, MessageResponse{Message: "$1"})`)
		return s
	})
}
