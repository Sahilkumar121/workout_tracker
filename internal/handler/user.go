package handler

import (
	"encoding/json"
	"errors"
	"fmt"
	"log/slog"
	"net/http"
	"strings"

	"github.com/Sahilkumar121/workout_tracker/internal/helper"
	"github.com/Sahilkumar121/workout_tracker/internal/httpx"
	"github.com/Sahilkumar121/workout_tracker/internal/middleware"
	"github.com/go-playground/validator/v10"
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgconn"
	"github.com/jackc/pgx/v5/pgxpool"
)

type UserHandler struct {
	DB         *pgxpool.Pool
	Logger     *slog.Logger
	Validate   *validator.Validate
	SecreteKey string
}

func NewUserHandler(db *pgxpool.Pool, logger *slog.Logger, secreteKey string) *UserHandler {
	return &UserHandler{
		DB:         db,
		Logger:     logger,
		Validate:   validator.New(),
		SecreteKey: secreteKey,
	}
}

// register user

func (register *UserHandler) RegisterUser(w http.ResponseWriter, r *http.Request) {
	// context request id
	requestID := middleware.RequestIDFromContext(r.Context())
	// store in struct
	var registerData RegisterUser
	err := json.NewDecoder(r.Body).Decode(&registerData)
	if err != nil {
		register.Logger.Error("failed decode data", "err", err, "requestID", requestID)
		httpx.Error(w, http.StatusInternalServerError, "something went error", httpx.CodeInternalServer)
		return
	}
	// validate the input data
	err = register.Validate.Struct(registerData)
	if err != nil {
		register.Logger.Error("failed validate data", "err", err, "requestID", requestID)
		httpx.Error(w, http.StatusBadRequest, "invalid field, please provide correct fields", httpx.CodeInvalidInput)
		return
	}
	// convert plain password to hash password
	hashPassword, err := helper.GetHashPassword(registerData.Password)
	if err != nil {
		register.Logger.Error("failed hash password", "err", err, "requestID", requestID)
		httpx.Error(w, http.StatusInternalServerError, "server is not responding, try again later", httpx.CodeInternalServer)
		return
	}
	registerData.Password = hashPassword
	// stored the hash password
	// insert in database
	_, err = register.DB.Exec(
		r.Context(),
		`insert into users (first_name, last_name, email, password) values ($1, $2, $3, $4);`,
		registerData.FirstName, registerData.LastName, registerData.Email, registerData.Password)

	if err != nil {
		var pgxErr *pgconn.PgError
		if errors.As(err, &pgxErr) && pgxErr.Code == "23505" {
			register.Logger.Error("failed user resgister", "err", pgxErr, "requestID", requestID)
			httpx.Error(w, http.StatusConflict, "user already exist", httpx.CodeConflict)
			return
		}
		register.Logger.Error("failed user register", "err", err, "requestID", requestID)
		httpx.Error(w, http.StatusInternalServerError, "server is not responding, try again later", httpx.CodeInternalServer)
		return
	}
	helper.SetResponse(w, http.StatusCreated)
}

// login user

func (login *UserHandler) LoginUser(w http.ResponseWriter, r *http.Request) {
	// request id from context
	requestID := middleware.RequestIDFromContext(r.Context())
	// login data from request body
	var loginData LoginUser
	err := json.NewDecoder(r.Body).Decode(&loginData)
	// decode error
	if err != nil {
		login.Logger.Error("failed decode data", "err", err, "requestID", requestID)
		httpx.Error(w, http.StatusInternalServerError, "something went wrong", httpx.CodeInternalServer)
		return
	}
	// validate request data
	err = login.Validate.Struct(loginData)
	// validate error
	if err != nil {
		login.Logger.Error("bad request", "err", err, "requestID", requestID)
		httpx.Error(w, http.StatusBadRequest, "please  provided all information", httpx.CodeInvalidInput)
		return
	}
	// check for user existence in DB
	row := login.DB.QueryRow(r.Context(), `select id, password from users where email=$1`, loginData.Email)
	var (
		userID       string
		userPassword string
	)
	// scan to select first row
	err = row.Scan(&userID, &userPassword)
	// if no row found then user is not registered
	// if not atleast on row is reterieve user not found
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			httpx.Error(w, http.StatusUnauthorized, "login failed", httpx.CodeUnauthorized)
			return
		}
		login.Logger.Error("database query error during login", "err", err, "requestID", requestID)
		httpx.Error(w, http.StatusUnauthorized, "login failed", httpx.CodeUnauthorized)
		return
	}
	// check for password
	err = helper.CheckHashPassword(userPassword, loginData.Password)
	if err != nil {
		httpx.Error(w, http.StatusUnauthorized, "login failed", httpx.CodeUnauthorized)
		return
	}
	token, err := helper.CreateAccessToken(userID, []byte(login.SecreteKey))
	if err != nil {
		login.Logger.Error("failed creation token", "err", err, "requestID", requestID)
		httpx.Error(w, http.StatusInternalServerError, "something weent wrong, try again later", httpx.CodeInternalServer)
		return
	}
	helper.SetResponse(w, http.StatusOK)
	_ = json.NewEncoder(w).Encode(LoginUserResponse{Token: token, TokenType: "Bearer"})
}

// user profile

func (me *UserHandler) GetProfile(w http.ResponseWriter, r *http.Request) {

	requestID := middleware.RequestIDFromContext(r.Context())
	user_id := middleware.AuthUserIDFromContext(r.Context())

	if user_id == "" {
		me.Logger.Error("unauthorized attempt", "requestID", requestID)
		httpx.Error(w, http.StatusUnauthorized, "unauthorized", httpx.CodeUnauthorized)
		return
	}
	row := me.DB.QueryRow(r.Context(), `SELECT id, first_name, last_name, email, created_at FROM users WHERE id = $1`, user_id)
	var userData User
	err := row.Scan(&userData.ID, &userData.FirstName, &userData.LastName, &userData.Email, &userData.CreatedAt)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			me.Logger.Warn("user profile not found", "userID", user_id, "requestID", requestID)
			httpx.Error(w, http.StatusNotFound, "user profile not found", httpx.CodeNotFound)
			return
		}
		me.Logger.Error("failed scan database", "err", err, "requesID", requestID)
		httpx.Error(w, http.StatusInternalServerError, "something went wrong", httpx.CodeInternalServer)
		return
	}
	helper.SetResponse(w, http.StatusOK)
	_ = json.NewEncoder(w).Encode(userData)
}

// user update there details
func (update *UserHandler) UpdateUser(w http.ResponseWriter, r *http.Request) {

	requestID := middleware.RequestIDFromContext(r.Context())
	user_id := middleware.AuthUserIDFromContext(r.Context())

	var updateData UpdateUserData
	err := json.NewDecoder(r.Body).Decode(&updateData)
	if err != nil {
		update.Logger.Error("decode error", "err", err, "requestID", requestID)
		httpx.Error(w, http.StatusInternalServerError, "something went wrong", httpx.CodeInternalServer)
		return
	}

	err = update.Validate.Struct(updateData)
	if err != nil {
		update.Logger.Error("incorrect body request", "err", err, "requestID", requestID)
		httpx.Error(w, http.StatusBadRequest, "please provide proper data", httpx.CodeInvalidInput)
		return
	}

	if updateData.Password != "" {
		hash, err := helper.GetHashPassword(updateData.Password)
		if err != nil {
			update.Logger.Error("hash password failed", "err", err, "requestID", requestID)
			httpx.Error(w, http.StatusInternalServerError, "something went wrong", httpx.CodeInternalServer)
			return
		}

		updateData.Password = hash
	}

	setParts := []string{}
	args := []any{}
	argID := 1

	if updateData.FirstName != "" {
		setParts = append(setParts, fmt.Sprintf("first_name=$%d", argID))
		args = append(args, updateData.FirstName)
		argID += 1
	}

	if updateData.LastName != "" {
		setParts = append(setParts, fmt.Sprintf("last_name=$%d", argID))
		args = append(args, updateData.LastName)
		argID += 1
	}

	if updateData.Email != "" {
		setParts = append(setParts, fmt.Sprintf("email=$%d", argID))
		args = append(args, updateData.Email)
		argID += 1
	}

	if updateData.Password != "" {
		setParts = append(setParts, fmt.Sprintf("password=$%d", argID))
		args = append(args, updateData.Password)
		argID += 1
	}

	if len(setParts) == 0 {
		httpx.Error(w, http.StatusBadRequest, "please provide at least one field to update", httpx.CodeInternalServer)
		return
	}

	args = append(args, user_id)

	query := fmt.Sprintf("update users set %s where id=$%d", strings.Join(setParts, ", "), argID)

	_, err = update.DB.Exec(r.Context(), query, args...)
	if err != nil {
		update.Logger.Error("database update failed", "err", err, "requestID", requestID)
		httpx.Error(w, http.StatusInternalServerError, "failed to update user", httpx.CodeInternalServer)
		return
	}

	helper.SetResponse(w, http.StatusOK)
	_ = json.NewEncoder(w).Encode(map[string]string{"message": "update successfully"})
}

// delete user
func (delete *UserHandler) DeleteUser(w http.ResponseWriter, r *http.Request) {

	requestID := middleware.RequestIDFromContext(r.Context())
	user_id := middleware.AuthUserIDFromContext(r.Context())

	commandtag, err := delete.DB.Exec(r.Context(), `delete from user where id=$1`, user_id)
	if err != nil {
		delete.Logger.Error("failed database query", "err", err, "requestID", requestID, "userID", user_id)
		httpx.Error(w, http.StatusInternalServerError, "something went wrong", httpx.CodeInternalServer)
		return
	}

	if commandtag.RowsAffected() == 0 {
		httpx.Error(w, http.StatusBadRequest, "no user id found", httpx.CodeInvalidInput)
		return
	}

	helper.SetResponse(w, http.StatusOK)
	_ = json.NewEncoder(w).Encode(map[string]string{"message": "delete successfully"})
}
