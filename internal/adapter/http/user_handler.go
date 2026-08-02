package httpadapter

import (
	"net/http"

	"github.com/jackc/pgx/v5/pgtype"
)

// GetMe godoc
// @Summary      내 정보 조회
// @Tags         user
// @Produce      json
// @Security     BearerAuth
// @Success      200  {object}  UserResponse
// @Router       /api/me [get]
func (r *Router) getMe(w http.ResponseWriter, req *http.Request) {
	userIDStr, ok := req.Context().Value(contextKeyUserID).(string)
	if !ok || userIDStr == "" {
		writeError(w, http.StatusUnauthorized, "인증 정보가 없습니다.")
		return
	}

	var uid pgtype.UUID
	if err := uid.Scan(userIDStr); err != nil {
		writeError(w, http.StatusBadRequest, "유효하지 않은 사용자 ID입니다.")
		return
	}

	user, err := r.queries.GetUserByID(req.Context(), uid)
	if err != nil {
		writeError(w, http.StatusNotFound, "사용자를 찾을 수 없습니다.")
		return
	}

	writeJSON(w, http.StatusOK, map[string]string{
		"id":       user.ID.String(),
		"email":    user.Email,
		"nickname": user.Nickname,
	})
}
