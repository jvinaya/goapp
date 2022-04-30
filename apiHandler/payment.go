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

type createPaymentRequest struct {
	Amount float64 `json:"amount" binding:"required,gt=0"`
	LoanID int64   `json:"loan_id" binding:"min=1" `
}
type paymentResponse struct {
	ID            int64     `json:"id"`
	LoanID        int64     `json:"loan_id"`
	UserID        int64     `json:"user_id"`
	Amount        string    `json:"amount"`
	CreatedBy     string    `json:"created_by"`
	LastUpdatedBy string    `json:"last_updated_by"`
	CreatedAt     time.Time `json:"created_at"`
	LastUpdatedAt time.Time `json:"last_updated_at"`
}
type txnResponse struct {
	ID                int64                `json:"payment_id"`
	LoanID            int64                `json:"loan_id"`
	UserID            int64                `json:"user_id"`
	Amount            string               `json:"payment_amount"`
	PendingLoanAmount string               `json:"pending_loan_amount"`
	LoanPaymentStatus db.EnumPaymentStatus `json:"loan_payment_status"`
	LoanAmount        string               `json:"loan_amount"`
	CreatedBy         string               `json:"created_by"`
	LastUpdatedBy     string               `json:"last_updated_by"`
	CreatedAt         time.Time            `json:"created_at"`
	LastUpdatedAt     time.Time            `json:"last_updated_at"`
}

func newPaymentResponse(loan db.Payment) paymentResponse {

	return paymentResponse{
		ID:            loan.ID,
		LoanID:        loan.LoanID,
		UserID:        loan.UserID,
		Amount:        loan.Amount,
		CreatedBy:     loan.CreatedBy,
		CreatedAt:     loan.CreatedAt,
		LastUpdatedBy: loan.LastUpdatedBy,
		LastUpdatedAt: loan.UpdatedAt,
	}

}
func newTxnResponse(txn db.TransactionDetail) txnResponse {

	return txnResponse{

		ID:                txn.CurrentPayment.ID,
		LoanID:            txn.CurrentPayment.LoanID,
		UserID:            txn.CurrentPayment.UserID,
		Amount:            txn.CurrentPayment.Amount,
		PendingLoanAmount: txn.LoanDetails.AmountNeedToPay,
		LoanPaymentStatus: txn.LoanDetails.RepaymentStatus,
		LoanAmount:        txn.LoanDetails.Amount,
		CreatedBy:         txn.CurrentPayment.CreatedBy,
		LastUpdatedBy:     txn.CurrentPayment.LastUpdatedBy,
		CreatedAt:         txn.CurrentPayment.CreatedAt,
		LastUpdatedAt:     txn.CurrentPayment.UpdatedAt,
	}

}

//createPayment create Payment  return the payment with audit details  from the data store

func (server *Server) createPayment(ctx *gin.Context) {
	var req createPaymentRequest

	if err := ctx.ShouldBindBodyWith(&req, binding.JSON); err != nil {
		ctx.JSON(http.StatusBadRequest, errorResponse(err))
		return
	}

	authPayload := ctx.MustGet(authorizationPayloadKey).(*token.Payload)
	arg := db.CreatePaymentParams{
		Amount:        fmt.Sprint(req.Amount),
		LoanID:        req.LoanID,
		CreatedBy:     authPayload.Username,
		LastUpdatedBy: authPayload.Username,
		IpFrom:        ctx.Request.RemoteAddr,
		UserAgent:     ctx.Request.UserAgent(),
	}

	txn, err := server.store.CreatePaymentTerms(ctx, arg)
	if err != nil {
		if err.Error() == db.LaonIsPaid {
			ctx.JSON(http.StatusBadRequest, errorResponse(err))
			return
		}
		ctx.JSON(http.StatusInternalServerError, errorResponse(err))
		return
	}

	ctx.JSON(http.StatusOK, newTxnResponse(txn))
}

type getPaymentRequest struct {
	ID int64 `uri:"id" binding:"required,min=1"`
}

//getPayment return the payment from the data store
func (server *Server) getPayment(ctx *gin.Context) {
	var req getPaymentRequest
	if err := ctx.ShouldBindUri(&req); err != nil {
		ctx.JSON(http.StatusBadRequest, errorResponse(err))
		return
	}

	payment, err := server.store.GetPayment(ctx, req.ID)
	if err != nil {
		if err == sql.ErrNoRows {
			ctx.JSON(http.StatusNotFound, errorResponse(err))
			return
		}

		ctx.JSON(http.StatusInternalServerError, errorResponse(err))
		return
	}

	ctx.JSON(http.StatusOK, newPaymentResponse(payment))
}

type listPaymentRequest struct {
	PageID   int32 `form:"page_id" binding:"required,min=1"`
	PageSize int32 `form:"page_size" binding:"required,min=5,max=10"`
}

// Returns a list of  Payments through pagination
func (server *Server) listPayment(ctx *gin.Context) {
	var req listPaymentRequest
	if err := ctx.ShouldBindQuery(&req); err != nil {
		ctx.JSON(http.StatusBadRequest, errorResponse(err))
		return
	}

	arg := db.ListPaymentParams{
		Limit:  req.PageSize,
		Offset: (req.PageID - 1) * req.PageSize,
	}

	payments, err := server.store.ListPayment(ctx, arg)

	if err != nil {
		ctx.JSON(http.StatusInternalServerError, errorResponse(err))
		return
	}

	ctx.JSON(http.StatusOK, payments)
}

// Returns a list of  Payments through pagination in Desc Order

func (server *Server) paymentDescList(ctx *gin.Context) {
	var req listPaymentRequest
	if err := ctx.ShouldBindQuery(&req); err != nil {
		ctx.JSON(http.StatusBadRequest, errorResponse(err))
		return
	}

	arg := db.ListDescPaymentParams{
		Limit:  req.PageSize,
		Offset: (req.PageID - 1) * req.PageSize,
	}

	payments, err := server.store.ListDescPayment(ctx, arg)

	if err != nil {
		ctx.JSON(http.StatusInternalServerError, errorResponse(err))
		return
	}

	ctx.JSON(http.StatusOK, payments)
}
