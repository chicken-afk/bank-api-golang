package api

import (
	db "bank-api/db/sqlc"
	"database/sql"
	"net/http"

	"github.com/gin-gonic/gin"
)

type CreateAccountRequest struct {
	Owner    string `json:"Owner" binding:"required"`
	Currency string `json:"currency" binding:"required"`
}

type response struct {
	status_code int32      `json:"status_code"`
	message     string     `json:"message"`
	data        db.Account `json:"data"`
}

func (server *Server) CreateAccount(ctx *gin.Context) {
	var req CreateAccountRequest
	if err := ctx.ShouldBindJSON(&req); err != nil {
		ctx.JSON(http.StatusBadRequest, errorResponse(err))
		return
	}

	arg := db.CreateAccountParams{
		Owner:    req.Owner,
		Currency: req.Currency,
		Balance:  0,
	}

	account, err := server.store.CreateAccount(ctx, arg)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, errorResponse(err))
		return
	}
	// res.status_code = 200
	// res.message = "Successfully created new account"
	// res.data = account

	ctx.JSON(http.StatusOK, account)
}

type GetAccountRequest struct {
	ID int64 `uri:"id" binding:"required"`
}

func (server *Server) GetAccount(ctx *gin.Context) {
	var req GetAccountRequest

	if err := ctx.ShouldBindUri(&req); err != nil {
		ctx.JSON(http.StatusBadRequest, errorResponse(err))
		return
	}

	account, err := server.store.GetAccount(ctx, req.ID)
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

type ListAccountRequest struct {
	Page     int32 `form:"page" binding:"required,min=1"`
	PageSize int32 `form:"page_size" binding:"required,min=5,max=20"`
}

func (server *Server) ListAccount(ctx *gin.Context) {
	var req ListAccountRequest
	if err := ctx.ShouldBindQuery(&req); err != nil {
		ctx.JSON(http.StatusBadRequest, errorResponse(err))
		return
	}

	arg := db.ListAccountsParams{
		Limit:  int32(req.PageSize),
		Offset: (req.Page - 1) * req.PageSize,
	}

	accounts, err := server.store.ListAccounts(ctx, arg)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, errorResponse(err))
		return
	}

	ctx.JSON(http.StatusOK, accounts)
}

type DeleteAccountRequest struct {
	ID int64 `uri:"id" binding:"required"`
}

type DeleteAccountResponse struct {
	Message string     `json:"message"`
	Data    db.Account `json:"account"`
}

func (server *Server) DeleteAccount(ctx *gin.Context) {
	var req DeleteAccountRequest
	if err := ctx.ShouldBindUri(&req); err != nil {
		ctx.JSON(http.StatusBadRequest, errorResponse(err))
		return
	}
	account, err := server.store.GetAccount(ctx, req.ID)
	if err != nil {
		if err == sql.ErrNoRows {
			ctx.JSON(http.StatusNotFound, errorResponse(err))
			return
		}
		ctx.JSON(http.StatusInternalServerError, errorResponse(err))
		return
	}

	err = server.store.DeleteAccount(ctx, req.ID)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, errorResponse(err))
		return
	}
	res := DeleteAccountResponse{
		Message: "Berhasil Menghapus Data",
		Data:    account,
	}

	ctx.JSON(http.StatusOK, res)

}

// type EditBalanceRequest struct {
// 	ID      int64 `uri:"id" binding:"required"`
// 	Balance int64 `json:"balance" binding:"required"`
// }
type EditBalanceRequest struct {
	ID      int64 `uri:"id"`
	Balance int64 `json:"balance"`
}
type EditBalanceResponse struct {
	Message string     `json:"message"`
	Data    db.Account `json:"account"`
}

func (server *Server) EditBalance(ctx *gin.Context) {
	var req EditBalanceRequest

	// Bind data dari URI
	// if err := ctx.ShouldBindUri(&req); err != nil {
	// 	ctx.JSON(http.StatusBadRequest, errorResponse(err))
	// 	return
	// }

	// // Bind data dari body JSON
	// if err := ctx.ShouldBindJSON(&req); err != nil {
	// 	ctx.JSON(http.StatusBadRequest, errorResponse(err))
	// 	return
	// }
	ctx.ShouldBindUri(&req)
	ctx.ShouldBindJSON(&req)

	arg := db.AddAccountBalanceParams{
		ID:     req.ID,
		Amount: req.Balance,
	}

	account, err := server.store.AddAccountBalance(ctx, arg)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, errorResponse(err))
		return
	}

	res := EditBalanceResponse{
		Message: "Berhasil Merubah Data",
		Data:    account,
	}

	ctx.JSON(http.StatusOK, res)
}

type EditAccountRequest struct {
	Balance  int64  `json:"balance" binding:"required"`
	Owner    string `json:"owner" binding:"required"`
	Currency string `json:"currency" binding:"required"`
	ID       int64  `json:"id" binding:"required"`
}

func (server *Server) EditAccount(ctx *gin.Context) {
	// // Bind data dari body JSON
	var req EditAccountRequest
	if err := ctx.ShouldBindJSON(&req); err != nil {
		ctx.JSON(http.StatusBadRequest, errorResponse(err))
		return
	}

	arg := db.UpdateAccountParams{
		ID:       req.ID,
		Owner:    req.Owner,
		Currency: req.Currency,
		Balance:  req.Balance,
	}

	account, err := server.store.UpdateAccount(ctx, arg)
	if err != nil {
		ctx.JSON(500, errorResponse(err))
	}

	res := EditBalanceResponse{
		Message: "Berhasil Merubah Data",
		Data:    account,
	}

	ctx.JSON(http.StatusOK, res)
}
