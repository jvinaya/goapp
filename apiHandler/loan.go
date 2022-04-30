package apihandler

import (
	"database/sql"
	"fmt"
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/gin-gonic/gin/binding"
	"github.com/jvinaya/goapp/db"
	"github.com/jvinaya/goapp/token"
)

type createLoanRequest struct {
	Amount float64 `json:"amount" binding:"required,gt=0"`
	Term   int32   `json:"term" binding:"required,min=1,max=6"`
}
type loanResponse struct {
	ID              int64                 `json:"id"`
	Amount          string                `json:"amount"`
	AmountNeedToPay string                `json:"amount_need_to_pay"`
	Term            int32                 `json:"term"`
	ApprovalStatus  db.EnumApprovalStatus `json:"approval_status"`
	IsActive        bool                  `json:"is_active"`
	RepaymentStatus db.EnumPaymentStatus  `json:"repayment_status"`
	CreatedBy       string                `json:"created_by"`
	LastUpdatedBy   string                `json:"last_updated_by"`
	CreatedAt       time.Time             `json:"created_at"`
	LastUpdatedAt   time.Time             `json:"last_updated_at"`
}
type updateLoanStatusRequest struct {
	ID             int64                 `json:"id" binding:"required,min=1"`
	ApprovalStatus db.EnumApprovalStatus `json:"approval_status" binding:"required,min=1,oneof=approved rejected"`
}
type getLoanRequest struct {
	ID int64 `uri:"id" binding:"required,min=1"`
}
type listLoanRequest struct {
	PageID   int32 `form:"page_id" binding:"required,min=1"`
	PageSize int32 `form:"page_size" binding:"required,min=5,max=10"`
}

func newLoanResponse(loan db.Loan) loanResponse {

	return loanResponse{
		ID:              loan.ID,
		Amount:          loan.Amount,
		AmountNeedToPay: loan.AmountNeedToPay,
		Term:            loan.Term,
		ApprovalStatus:  loan.ApprovalStatus,
		IsActive:        loan.IsActive,
		RepaymentStatus: loan.RepaymentStatus,
		CreatedBy:       loan.CreatedBy,
		LastUpdatedBy:   loan.LastUpdatedBy,
		CreatedAt:       loan.CreatedAt,
		LastUpdatedAt:   loan.UpdatedAt,
	}
}

//createLoan create loanand  return the loan with audit details  from the data store
func (server *Server) createLoan(ctx *gin.Context) {
	var req createLoanRequest

	if err := ctx.ShouldBindBodyWith(&req, binding.JSON); err != nil {
		ctx.JSON(http.StatusBadRequest, errorResponse(err))
		return
	}

	authPayload := ctx.MustGet(authorizationPayloadKey).(*token.Payload)
	arg := db.CreateLoanParams{
		Amount:          fmt.Sprint(req.Amount),
		AmountNeedToPay: fmt.Sprint(req.Amount),
		Term:            req.Term,
		ApprovalStatus:  db.EnumApprovalStatusPending,
		RepaymentStatus: db.EnumPaymentStatusUnpaid,
		CreatedBy:       authPayload.Username,
		LastUpdatedBy:   authPayload.Username,
		IpFrom:          ctx.Request.RemoteAddr,
		UserAgent:       ctx.Request.UserAgent(),
	}

	loan, err := server.store.CreateLoanWithBorrower(ctx, arg)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, errorResponse(err))
		return
	}

	ctx.JSON(http.StatusOK, newLoanResponse(loan))
}

//getLoan return the loan from the data store
func (server *Server) getLoan(ctx *gin.Context) {
	var req getLoanRequest
	if err := ctx.ShouldBindUri(&req); err != nil {
		ctx.JSON(http.StatusBadRequest, errorResponse(err))
		return
	}

	loan, err := server.store.GetLoan(ctx, req.ID)
	if err != nil {
		if err == sql.ErrNoRows {
			ctx.JSON(http.StatusNotFound, errorResponse(err))
			return
		}

		ctx.JSON(http.StatusInternalServerError, errorResponse(err))
		return
	}

	ctx.JSON(http.StatusOK, newLoanResponse(loan))
}

// Returns a list of  Loans through pagination
func (server *Server) listLoan(ctx *gin.Context) {
	var req listLoanRequest
	if err := ctx.ShouldBindQuery(&req); err != nil {
		ctx.JSON(http.StatusBadRequest, errorResponse(err))
		return
	}

	arg := db.ListLoanParams{
		Limit:  req.PageSize,
		Offset: (req.PageID - 1) * req.PageSize,
	}

	loans, err := server.store.ListLoan(ctx, arg)

	if err != nil {
		ctx.JSON(http.StatusInternalServerError, errorResponse(err))
		return
	}

	ctx.JSON(http.StatusOK, loans)
}

// Returns a list of  Loans through pagination in Desc order
func (server *Server) loanDescList(ctx *gin.Context) {
	var req listLoanRequest
	if err := ctx.ShouldBindQuery(&req); err != nil {
		ctx.JSON(http.StatusBadRequest, errorResponse(err))
		return
	}

	arg := db.ListDescLoanParams{
		Limit:  req.PageSize,
		Offset: (req.PageID - 1) * req.PageSize,
	}

	loans, err := server.store.ListDescLoan(ctx, arg)

	if err != nil {
		ctx.JSON(http.StatusInternalServerError, errorResponse(err))
		return
	}

	ctx.JSON(http.StatusOK, loans)
}

func (server *Server) updateLoanStatus(ctx *gin.Context) {

	var req updateLoanStatusRequest
	if err := ctx.ShouldBindJSON(&req); err != nil {
		ctx.JSON(http.StatusBadRequest, errorResponse(err))
		return
	}

	authPayload := ctx.MustGet(authorizationPayloadKey).(*token.Payload)

	arg := db.UpdateLoanStatusParams{
		ID:             req.ID,
		ApprovalStatus: req.ApprovalStatus,
		LastUpdatedBy:  authPayload.Username,
		UpdatedAt:      time.Now(),
	}

	loan, err := server.store.UpdateLoanStatus(ctx, arg)

	if err != nil {
		ctx.JSON(http.StatusInternalServerError, errorResponse(err))
		return
	}

	ctx.JSON(http.StatusOK, loan)
}
