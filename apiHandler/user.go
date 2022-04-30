package apihandler

import (
	"database/sql"
	"fmt"
	"net/http"
	"strings"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/gin-gonic/gin/binding"
	"github.com/jvinaya/goapp/db"
	"github.com/jvinaya/goapp/token"
	"github.com/jvinaya/goapp/utils"
	"github.com/lib/pq"
)

type createUserRequest struct {
	Name     string `json:"name" binding:"required,min=3"`
	Mobile   string `json:"mobile"`
	Email    string `json:"email" binding:"email,required"`
	Password string `json:"password" binding:"required,min=6"`
	Address  string `json:"address"`
}
type updateUserRequest struct {
	ID       int64  `json:"id" binding:"required,min=1"`
	Name     string `json:"name" binding:"required,min=3"`
	Mobile   string `json:"mobile"`
	Email    string `json:"email" binding:"email,required"`
	Password string `json:"password" binding:"required,min=6"`
	Address  string `json:"address"`
}
type userResponse struct {
	ID                int64          `json:"id"`
	Name              string         `json:"name"`
	Mobile            sql.NullString `json:"mobile"`
	Address           sql.NullString `json:"address"`
	Email             string         `json:"email"`
	PasswordChangedAt time.Time      `json:"password_changed_at"`
	CreatedBy         string         `json:"created_by"`
	LastUpdatedBy     string         `json:"last_updated_by"`
	CreatedAt         time.Time      `json:"created_at"`
	LastUpdatedAt     time.Time      `json:"last_updated_at"`
}

func newUserResponse(user db.User) userResponse {
	return userResponse{
		ID:                user.ID,
		Name:              user.Name,
		Mobile:            user.Mobile,
		Address:           user.Address,
		Email:             user.Email,
		PasswordChangedAt: user.PasswordChangedAt,
		CreatedBy:         user.CreatedBy,
		LastUpdatedBy:     user.LastUpdatedBy,
		CreatedAt:         user.CreatedAt,
		LastUpdatedAt:     user.UpdatedAt,
	}
}

//createUser create User
func (server *Server) createUser(ctx *gin.Context) {
	var req createUserRequest
	if err := ctx.ShouldBindBodyWith(&req, binding.JSON); err != nil {
		ctx.JSON(http.StatusBadRequest, errorResponse(err))
		return
	}

	hashedPassword, err := utils.HashPassword(req.Password)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, errorResponse(err))
		return
	}
	arg := db.CreateUserParams{
		Name:           req.Name,
		Mobile:         sql.NullString{String: req.Mobile, Valid: len(strings.TrimSpace(req.Mobile)) > 0},
		Email:          req.Email,
		CreatedBy:      req.Email,
		LastUpdatedBy:  req.Email,
		IpFrom:         ctx.Request.RemoteAddr,
		UserAgent:      ctx.Request.UserAgent(),
		HashedPassword: hashedPassword,
		Address:        sql.NullString{String: req.Address, Valid: len(strings.TrimSpace(req.Address)) > 0},
	}

	user, err := server.store.CreateUser(ctx, arg)
	if err != nil {
		if pqErr, ok := err.(*pq.Error); ok {
			switch pqErr.Code.Name() {
			case "foreign_key_violation", "unique_violation":
				ctx.JSON(http.StatusForbidden, errorResponse(err))
				return
			}
		}
		ctx.JSON(http.StatusInternalServerError, errorResponse(err))
		return
	}

	ctx.JSON(http.StatusOK, newUserResponse(user))
}

//updateUser update User Details
func (server *Server) updateUser(ctx *gin.Context) {
	var req updateUserRequest
	if err := ctx.ShouldBindBodyWith(&req, binding.JSON); err != nil {
		ctx.JSON(http.StatusBadRequest, errorResponse(err))
		return
	}
	authPayload := ctx.MustGet(authorizationPayloadKey).(*token.Payload)

	account, err := server.store.GetUserByEmail(ctx, authPayload.Username)
	if err != nil {
		if err == sql.ErrNoRows {
			ctx.JSON(http.StatusNotFound, errorResponse(err))
			return
		}

		ctx.JSON(http.StatusInternalServerError, errorResponse(err))
		return
	}
	if account.ID != req.ID {
		ctx.JSON(http.StatusUnauthorized, errorResponse(fmt.Errorf("you can't update other details")))
		return
	}

	hashedPassword, err := utils.HashPassword(req.Password)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, errorResponse(err))
		return
	}
	arg := db.UpdateUserParams{
		ID:             req.ID,
		Name:           req.Name,
		Mobile:         sql.NullString{String: req.Mobile, Valid: len(strings.TrimSpace(req.Mobile)) > 0},
		Email:          req.Email,
		LastUpdatedBy:  req.Email,
		HashedPassword: hashedPassword,
		UpdatedAt:      time.Now(),
		Address:        sql.NullString{String: req.Address, Valid: len(strings.TrimSpace(req.Address)) > 0},
	}

	user, err := server.store.UpdateUser(ctx, arg)
	if err != nil {

		ctx.JSON(http.StatusInternalServerError, errorResponse(err))
		return
	}

	ctx.JSON(http.StatusOK, newUserResponse(user))
}

type getUserRequest struct {
	ID int64 `uri:"id" binding:"required,min=1"`
}

//getUser return the User from the data store
func (server *Server) getUser(ctx *gin.Context) {
	var req getUserRequest
	if err := ctx.ShouldBindUri(&req); err != nil {
		ctx.JSON(http.StatusBadRequest, errorResponse(err))
		return
	}

	account, err := server.store.GetUser(ctx, req.ID)
	if err != nil {
		if err == sql.ErrNoRows {
			ctx.JSON(http.StatusNotFound, errorResponse(err))
			return
		}

		ctx.JSON(http.StatusInternalServerError, errorResponse(err))
		return
	}

	ctx.JSON(http.StatusOK, account)
}

type listUserRequest struct {
	PageID   int32 `form:"page_id" binding:"required,min=1"`
	PageSize int32 `form:"page_size" binding:"required,min=5,max=10"`
}

// Returns a list of  Users through pagination

func (server *Server) listUsers(ctx *gin.Context) {
	var req listUserRequest
	if err := ctx.ShouldBindQuery(&req); err != nil {
		ctx.JSON(http.StatusBadRequest, errorResponse(err))
		return
	}

	// authPayload := ctx.MustGet(authorizationPayloadKey).(*token.Payload)
	arg := db.ListUserParams{
		// Owner:  authPayload.Username,
		Limit:  req.PageSize,
		Offset: (req.PageID - 1) * req.PageSize,
	}

	accounts, err := server.store.ListUser(ctx, arg)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, errorResponse(err))
		return
	}

	ctx.JSON(http.StatusOK, accounts)
}

// Returns a list of  Users through pagination in Desc Order
func (server *Server) usersDescList(ctx *gin.Context) {
	var req listUserRequest
	if err := ctx.ShouldBindQuery(&req); err != nil {
		ctx.JSON(http.StatusBadRequest, errorResponse(err))
		return
	}

	arg := db.ListDescUserParams{
		Limit:  req.PageSize,
		Offset: (req.PageID - 1) * req.PageSize,
	}

	accounts, err := server.store.ListDescUser(ctx, arg)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, errorResponse(err))
		return
	}

	ctx.JSON(http.StatusOK, accounts)
}

type loginUserRequest struct {
	Username string `json:"username" binding:"email,required"`
	Password string `json:"password" binding:"required,min=6"`
}
type loginUserResponse struct {
	AccessToken string       `json:"access_token"`
	User        userResponse `json:"user"`
}

//loginUser function is user login and return data with  token
func (server *Server) loginUser(ctx *gin.Context) {
	var req loginUserRequest
	if err := ctx.ShouldBindJSON(&req); err != nil {
		ctx.JSON(http.StatusBadRequest, errorResponse(err))
		return
	}
	user, err := server.store.GetUserByEmail(ctx, req.Username)
	if err != nil {
		if err == sql.ErrNoRows {
			ctx.JSON(http.StatusNotFound, errorResponse(err))
			return
		}
		ctx.JSON(http.StatusInternalServerError, errorResponse(err))
		return
	}
	err = utils.CheckPassword(req.Password, user.HashedPassword)
	if err != nil {
		ctx.JSON(http.StatusUnauthorized, errorResponse(err))
		return
	}

	accessToken, _, err := server.tokenMaker.CreateToken(
		user.Email,
		server.config.AccessTokenDuration,
	)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, errorResponse(err))
		return
	}
	rsp := loginUserResponse{
		AccessToken: accessToken,
		User:        newUserResponse(user),
	}

	ctx.JSON(http.StatusOK, rsp)
}
