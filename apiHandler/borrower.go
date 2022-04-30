package apihandler

import (
	"database/sql"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/jvinaya/goapp/db"
)

type getBorrowerRequest struct {
	ID int64 `uri:"id" binding:"required,min=1"`
}

//getBorrower return the borrower from the data store
func (server *Server) getBorrower(ctx *gin.Context) {
	var req getBorrowerRequest
	if err := ctx.ShouldBindUri(&req); err != nil {
		ctx.JSON(http.StatusBadRequest, errorResponse(err))
		return
	}

	borrower, err := server.store.GetBorrower(ctx, req.ID)
	if err != nil {
		if err == sql.ErrNoRows {
			ctx.JSON(http.StatusNotFound, errorResponse(err))
			return
		}

		ctx.JSON(http.StatusInternalServerError, errorResponse(err))
		return
	}

	ctx.JSON(http.StatusOK, borrower)
}

type listBorrowerRequest struct {
	PageID   int32 `form:"page_id" binding:"required,min=1"`
	PageSize int32 `form:"page_size" binding:"required,min=5,max=10"`
}

// Returns a list of  Borrowers through pagination
func (server *Server) listBorrower(ctx *gin.Context) {
	var req listBorrowerRequest
	if err := ctx.ShouldBindQuery(&req); err != nil {
		ctx.JSON(http.StatusBadRequest, errorResponse(err))
		return
	}

	// authPayload := ctx.MustGet(authorizationPayloadKey).(*token.Payload)

	arg := db.ListBorrowerParams{
		// Owner:  authPayload.Username,
		Limit:  req.PageSize,
		Offset: (req.PageID - 1) * req.PageSize,
	}

	borrower, err := server.store.ListBorrower(ctx, arg)

	if err != nil {
		ctx.JSON(http.StatusInternalServerError, errorResponse(err))
		return
	}

	ctx.JSON(http.StatusOK, borrower)
}

// Returns a list of  Borrowers in Desc Order through pagination
func (server *Server) borrowerDescList(ctx *gin.Context) {
	var req listBorrowerRequest
	if err := ctx.ShouldBindQuery(&req); err != nil {
		ctx.JSON(http.StatusBadRequest, errorResponse(err))
		return
	}

	arg := db.ListDescBorrowerParams{
		Limit:  req.PageSize,
		Offset: (req.PageID - 1) * req.PageSize,
	}

	borrower, err := server.store.ListDescBorrower(ctx, arg)

	if err != nil {
		ctx.JSON(http.StatusInternalServerError, errorResponse(err))
		return
	}

	ctx.JSON(http.StatusOK, borrower)
}
