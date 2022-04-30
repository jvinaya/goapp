package apihandler

import (
	"fmt"

	"github.com/gin-gonic/gin"
	"github.com/jvinaya/goapp/db"
	"github.com/jvinaya/goapp/token"
	"github.com/jvinaya/goapp/utils"
)

type Server struct {
	config     utils.Config
	store      *db.Store
	tokenMaker token.Maker
	router     *gin.Engine
}

func NewServer(config utils.Config, store *db.Store) (*Server, error) {
	tokMaker, err := token.NewPasetoMaker(config.TokenSymmetricKey)
	if err != nil {
		return nil, fmt.Errorf("connot create token maker : %w", err)
	}
	server := &Server{
		config:     config,
		store:      store,
		tokenMaker: tokMaker,
	}
	server.setupRouter()
	return server, nil
}

func (server *Server) setupRouter() {
	router := gin.Default()
	//for login
	router.POST("/users", server.createUser)
	router.POST("/users/login", server.loginUser)
	//here we check person is authenticated or not
	authRoutes := router.Group("/").Use(authMiddleware(server.tokenMaker))
	authRoutes.GET("/users/:id", server.getUser)
	authRoutes.GET("/users", server.listUsers)
	authRoutes.GET("/users/latest", server.usersDescList)
	authRoutes.PUT("/users", server.updateUser)

	authRoutes.POST("/loans/createLoan", server.createLoan)
	authRoutes.GET("/loans/:id", server.getLoan)
	authRoutes.GET("/loans", server.listLoan)
	authRoutes.GET("/loans/latest", server.loanDescList)
	authRoutes.PUT("/loans/updateApprovalStatus", server.updateLoanStatus)

	authRoutes.POST("/payments", server.createPayment)
	authRoutes.GET("/payments/:id", server.getPayment)
	authRoutes.GET("/payments", server.listPayment)
	authRoutes.GET("/payments/latest", server.paymentDescList)

	authRoutes.GET("/borrowers/:id", server.getBorrower)
	authRoutes.GET("/borrowers", server.listBorrower)
	authRoutes.GET("/borrowers/latest", server.borrowerDescList)

	// router.GET("/requests/:id", server.getRequest)
	// router.GET("/requests", server.listRequest)
	server.router = router
}

// Start runs the HTTP server on a specific address.
func (server *Server) Start(address string) error {
	return server.router.Run(address)
}
func errorResponse(err error) gin.H {
	return gin.H{"error": err.Error()}
}
