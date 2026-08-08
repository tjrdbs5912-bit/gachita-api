package httpadapter

import (
	"errors"
	"net/http"

	"github.com/jackc/pgx/v5/pgconn"
	"golang.org/x/crypto/bcrypt"

	"gachita-api/internal/db"
)

// SignUp godoc
// @Summary      회원가입
// @Tags         auth
// @Accept       json
// @Produce      json
// @Param        body  body  SignUpRequest  true  "가입 정보"
// @Success      201  {object}  UserResponse
// @Router       /api/auth/signup [post]
func (r *Router) signUp(w http.ResponseWriter, req *http.Request) {
	var body SignUpRequest
	if err := decodeJSON(req, &body); err != nil {
		writeError(w, http.StatusBadRequest, "invalid json")
		return
	}

	if body.Email == "" || body.Password == "" || body.Nickname == "" {
		writeError(w, http.StatusBadRequest, "이메일, 비밀번호, 닉네임은 필수 입력 항목입니다.")
		return
	}

	hash, err := bcrypt.GenerateFromPassword([]byte(body.Password), bcrypt.DefaultCost)
	if err != nil {
		writeError(w, http.StatusInternalServerError, "비밀번호 해싱에 실패했습니다.")
		return
	}

	user, err := r.queries.CreateUser(req.Context(), db.CreateUserParams{
		Email:        body.Email,
		PasswordHash: string(hash),
		Nickname:     body.Nickname,
	})
	if err != nil {
		var pgErr *pgconn.PgError
		if errors.As(err, &pgErr) && pgErr.Code == "23505" {
			writeError(w, http.StatusConflict, "이메일 또는 닉네임이 이미 사용 중입니다.")
			return
		}
		writeError(w, http.StatusInternalServerError, "사용자 생성에 실패했습니다.")
		return
	}

	writeJSON(w, http.StatusCreated, toUserResponse(user.ID.String(), user.Email, user.Nickname))
}

// Login godoc
// @Summary      로그인
// @Tags         auth
// @Accept       json
// @Produce      json
// @Param        body  body  LoginRequest  true  "로그인 정보"
// @Success      200  {object}  LoginResponse
// @Router       /api/auth/login [post]
func (r *Router) login(w http.ResponseWriter, req *http.Request) {
	var body LoginRequest
	if err := decodeJSON(req, &body); err != nil {
		writeError(w, http.StatusBadRequest, "invalid json")
		return
	}

	if body.Email == "" || body.Password == "" {
		writeError(w, http.StatusBadRequest, "이메일과 비밀번호는 필수 입력 항목입니다.")
		return
	}

	user, err := r.queries.GetUserByEmail(req.Context(), body.Email)
	if err != nil {
		writeError(w, http.StatusUnauthorized, "이메일 또는 비밀번호가 올바르지 않습니다.")
		return
	}

	if err := bcrypt.CompareHashAndPassword([]byte(user.PasswordHash), []byte(body.Password)); err != nil {
		writeError(w, http.StatusUnauthorized, "이메일 또는 비밀번호가 올바르지 않습니다.")
		return
	}

	accessToken, err := r.tokens.IssueAccessToken(user.ID.String(), user.Email)
	if err != nil {
		writeError(w, http.StatusInternalServerError, "토큰 발급에 실패했습니다.")
		return
	}

	writeJSON(w, http.StatusOK, LoginResponse{
		AccessToken: accessToken,
		TokenType:   "Bearer",
	})
}
