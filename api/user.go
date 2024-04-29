package api

import (
	"fmt"
	"net/http"
	"time"
	"database/sql"
	"github.com/gin-gonic/gin"
	"github.com/lib/pq"
	db "github.com/shoroogAlghamdi/banking_system/db/sqlc"
	"github.com/shoroogAlghamdi/banking_system/util"
)

type createUserRequest struct {
	Username string `json:"username" binding:"required,alphanum"`
	Password string `json:"password" binding:"required,min=6"`
	FullName string `json:"full_name" binding:"required"`
	Email    string `json:"email" binding:"required,email"`
}


type userResponse struct {
	Username          string    `json:"username"`
	FullName          string    `json:"fullName"`
	Email             string    `json:"email"`
	PasswordChangedAt time.Time `json:"passwordChangedAt"`
	CreatedAt         time.Time `json:"createdAt"`
}

func createUserResponse(user db.User) userResponse {
	return userResponse{
		Username: user.Username,
		FullName: user.FullName,
		Email: user.Email,
		CreatedAt: user.CreatedAt,
		PasswordChangedAt: user.PasswordChangedAt,
	}
}
func (server *Server) createUser(ctx *gin.Context) {
	var req createUserRequest
	if err := ctx.ShouldBindJSON(&req); err != nil {
		fmt.Println(err)
		ctx.JSON(http.StatusBadRequest, errResponse(err))
		return
	}

	hashedPassword, err := util.HashedPassword(req.Password)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, errResponse(err))
		return
	}
	arg := db.CreateUserParams{
		Username:       req.Username,
		HashedPassword: hashedPassword,
		Email:          req.Email,
		FullName:       req.FullName,
	}

	user, err := server.store.CreateUser(ctx, arg)
	if err != nil {
		if pqErr, ok := err.(*pq.Error); ok {
			switch pqErr.Code.Name() {
			case "unique_violation":
				ctx.JSON(http.StatusForbidden, errResponse(err))
				return
			}
		}
		ctx.JSON(http.StatusBadRequest, errResponse(err))
		return
	}

	resp := createUserResponse(user)
	ctx.JSON(http.StatusOK, resp)

}

type loginUserRequest struct {
	Username string `json:"username" binding:"required,alphanum"`
	Password string `json:"password" binding:"required,min=6"`
}

type logingUserResponse struct {
	Accesstoken string `json:"access_token"`
	User userResponse `json:"user" binding:"required,min=6"`
}

func (server *Server) loginUser( ctx *gin.Context) {
	var req loginUserRequest
	if err := ctx.ShouldBindJSON(&req); err != nil {
		ctx.JSON(http.StatusBadRequest, errResponse(err))
	}

	user, err := server.store.GetUser(ctx, req.Username)
	if err != nil {
		if err == sql.ErrNoRows {
			ctx.JSON(http.StatusNotFound, errResponse(err))
			return 
		}
		ctx.JSON(http.StatusInternalServerError, errResponse(err))
		return  
	}

	err = util.CheckPassword(req.Password, user.HashedPassword)
	if err != nil {
		ctx.JSON(http.StatusUnauthorized, errResponse(err))
		return 
	}

	access_token, err := server.tokenMaker.CreateToken(req.Username, server.config.AccessTokenDuration)
	if err != nil {
		ctx.JSON(http.StatusUnauthorized, errResponse(err))
		return
	}

	rsp := logingUserResponse{
		Accesstoken: access_token,
		User: createUserResponse(user),
	}

	ctx.JSON(http.StatusOK, rsp)
	 
	

	
}



