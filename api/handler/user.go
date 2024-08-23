package handler

import (
	"context"
	"encoding/json"
	"fmt"
	"log"
	"math/rand"
	"net/http"
	"strconv"
	"time"

	"google_docs_user/api/email"
	"google_docs_user/api/token"
	pb "google_docs_user/genproto/user"
	"google_docs_user/models"

	"github.com/gin-gonic/gin"
	"golang.org/x/crypto/bcrypt"
)

// @Summary      Register a new user
// @Description  This endpoint registers a new user by taking user details, hashing the password, and generating a confirmation code.
// @Tags         auth
// @Accept       json
// @Produce      json
// @Param        user  body      models.Register  true  "User Registration Data"
// @Success      202   {object}  user.RegisterResp
// @Failure      400   {object}  string
// @Failure      500   {object}  string
// @Router       /register [post]
func (h Handler) Register(c *gin.Context) {
	req := models.Register{}

	err := json.NewDecoder(c.Request.Body).Decode(&req)
	if err != nil {
		h.Log.Error(fmt.Sprintf("bodydan malumotlarni olishda xatolik: %v", err))
		c.AbortWithStatusJSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	hashpassword, err := bcrypt.GenerateFromPassword([]byte(req.Password), bcrypt.DefaultCost)
	if err != nil {
		h.Log.Error(fmt.Sprintf("Pasworni hashlashda xatolik: %v", err))
		c.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	req.Password = string(hashpassword)

	rand.Seed(time.Now().UnixNano())
	randomCode := rand.Intn(900000) + 100000

	req1 := pb.RegisterReq{
		Email:     req.Email,
		FirstName: req.FirstName,
		LastName:  req.LastName,
		Password:  req.Password,
		Code:      int64(randomCode),
	}

	resp, err := h.User.Register(c, &req1)
	if err != nil {
		h.Log.Error(fmt.Sprintf("Foydalanuvchi malumotlarni bazga yuborishda xatolik: %v", err))
		c.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	email.SendCode(req1.Email, string(randomCode))

	c.JSON(http.StatusAccepted, resp)
}

// @Summary      Login a user
// @Description  This endpoint logs in a user by checking the credentials and generating JWT tokens.
// @Tags         auth
// @Accept       json
// @Produce      json
// @Param        credentials  body      user.LoginReq  true  "User Login Data"
// @Success      202   {object}  token.Tokens
// @Failure      400   {object}  string
// @Failure      401   {object}  string
// @Failure      500   {object}  string
// @Router       /login [post]
func (h Handler) LoginUser(c *gin.Context) {
	req := pb.LoginReq{}

	if err := json.NewDecoder(c.Request.Body).Decode(&req); err != nil {
		c.AbortWithStatusJSON(http.StatusBadRequest, err)
		return
	}

	req1 := pb.GetUSerByEmailReq{
		Email: req.Email,
	}

	user, err := h.User.GetUSerByEmail(c, &req1)
	if err != nil {
		h.Log.Error(fmt.Sprintf("GetbyUserda xatolik: %v", err))
		c.AbortWithStatusJSON(http.StatusInternalServerError, err)
		return
	}

	if err = bcrypt.CompareHashAndPassword([]byte(user.User.Password), []byte(req.Password)); err != nil {
		log.Printf("Password comparison failed: %v", err)
		c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"error": "Invalid credentials"})
		return
	}

	token := token.GenerateJWT(&pb.User{
		Id:       user.User.Id,
		Email:    req.Email,
		Password: req.Password,
		Role:     user.User.Role,
	})

	_, err = h.User.StoreRefreshToken(context.Background(), &pb.StoreRefreshTokenReq{
		UserId:  user.User.Id,
		Refresh: token.Refresh,
	})

	if err != nil {
		h.Log.Error(fmt.Sprintf("storefreshtokenda xatolik: %v", err))
		c.AbortWithStatusJSON(http.StatusInternalServerError, err)
		return
	}

	c.JSON(http.StatusAccepted, token)
}

// ConfirmationRegister godoc
// @Summary      Confirm user registration
// @Description  This endpoint confirms user registration by verifying the email and confirmation code.
// @Tags         auth
// @Accept       json
// @Produce      json
// @Param        email  path      string  true  "User Email"
// @Param        code   path      string  true  "Confirmation Code"
// @Success      200    {object}  user.ConfirmationRegisterResp
// @Failure      400    {object}  string
// @Failure      500    {object}  string
// @Router       /confirm/{email}/{code} [get]
func (h Handler) ConfirmationRegister(c *gin.Context) {
	codestr := c.Param("code")

	code, err := strconv.ParseInt(codestr, 10, 64)
	if err != nil {
		h.Log.Error(fmt.Sprintf("Kod stringdan int64 ga o'tkazishda xatolik: %v", err))
		c.AbortWithStatusJSON(http.StatusBadRequest, gin.H{"error": "Kod noto'g'ri formatda"})
		return
	}

	req := pb.ConfirmationRegisterReq{
		Email: c.Param("email"),
		Code:  code,
	}

	res, err := h.User.ConfirmationRegister(c, &req)
	if err != nil {
		h.Log.Error("ConfirmationRegister funksiyasiga yuborishda xatolik.", "error", err)
		c.AbortWithStatusJSON(400, gin.H{
			"error": err.Error(),
		})
		return
	}

	c.JSON(200, res)
}

// @Summary      Get user by email
// @Description  This endpoint retrieves user details by email.
// @Tags         user
// @Accept       json
// @Produce      json
// @Param        email  path      string  true  "User Email"
// @Success      200    {object}  user.GetUSerByEmailResp
// @Failure      500    {object}  string
// @Router       /user/{email} [get]
func (h Handler) GetUSerByEmail(c *gin.Context) {
	req := pb.GetUSerByEmailReq{
		Email: c.Param("email"),
	}

	res, err := h.User.GetUSerByEmail(c, &req)
	if err != nil {
		h.Log.Error("GetUSerByEmail funksiyasida xatolik.", "error", err.Error())
		c.AbortWithStatusJSON(500, gin.H{
			"error": err.Error(),
		})
		return
	}

	c.JSON(200, res)
}

// @Summary      Update user password
// @Description  This endpoint updates the user password after validating the old password.
// @Tags         user
// @Accept       json
// @Produce      json
// @Param        old_password  path      string  true  "Old Password"
// @Param        new_password  path      string  true  "New Password"
// @Param        email         path      string  true  "User Email"
// @Success      200    {object}  user.UpdatePasswordResp
// @Failure      401    {object}  string
// @Failure      500    {object}  string
// @Router       /user/update_password/{email}/{old_password}/{new_password} [put]
func (h Handler) UpdatePassword(c *gin.Context) {
	req := pb.UpdatePasswordReq{
		OldPassword: c.Param("old password"),
		NewPassword: c.Param("new password"),
		Email:       c.Param("email"),
	}

	req1 := pb.GetUSerByEmailReq{
		Email: req.Email,
	}

	resemail, err := h.User.GetUSerByEmail(c, &req1)
	if err != nil {
		h.Log.Error("Email buyicha malumotlarni olishda xatolik", "error", err.Error())
		c.AbortWithStatusJSON(500, gin.H{
			"error": err.Error(),
		})
		return
	}

	if err = bcrypt.CompareHashAndPassword([]byte(resemail.User.Password), []byte(req.OldPassword)); err != nil {
		log.Printf("Password comparison failed: %v", err)
		c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"error": "Invalid credentials"})
		return
	}

	hashpassword, err := bcrypt.GenerateFromPassword([]byte(req.NewPassword), bcrypt.DefaultCost)
	if err != nil {
		h.Log.Error(fmt.Sprintf("Pasworni hashlashda xatolik: %v", err))
		c.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	req.NewPassword = string(hashpassword)

	res, err := h.User.UpdatePassword(c, &req)
	if err != nil {
		h.Log.Error("UpdatePassword funksiyasiga malumot yuboishda xatolik.", "error", err.Error())
		c.AbortWithStatusJSON(500, gin.H{
			"error": err.Error(),
		})
		return
	}

	c.JSON(200, res)
}

// @Summary      Reset user password
// @Description  This endpoint sends a password reset email to the specified email address.
// @Tags         auth
// @Accept       json
// @Produce      json
// @Param        email  path      string  true  "User Email"
// @Success      200    {object}  user.ResetPasswordResp
// @Failure      400    {object}  string
// @Failure      500    {object}  string
// @Router       /reset_password/{email} [get]
func (h Handler) ResetPassword(c *gin.Context) {
	req := pb.ResetPasswordReq{
		Email: c.Param("email"),
	}

	res, err := email.Email(req.Email)
	if err != nil {
		h.Log.Error("Email ga xabar yuuborishda xatolik.", "error", err.Error())
		c.AbortWithStatusJSON(400, gin.H{
			"error": err.Error(),
		})
		return
	}

	c.JSON(200, res)
}

// @Summary      Confirm new password
// @Description  This endpoint confirms the new password by updating the user's password after validation.
// @Tags         auth
// @Accept       json
// @Produce      json
// @Param        confirmation  body      user.ConfirmationReq  true  "Confirmation Data"
// @Success      200    {object}  user.ConfirmationPasswordResp
// @Failure      400    {object}  string
// @Failure      500    {object}  string
// @Router       /confirmation_password [post]
func (h Handler) ConfirmationPassword(c *gin.Context) {
	req := pb.ConfirmationReq{}

	if err := json.NewDecoder(c.Request.Body).Decode(&req); err != nil {
		c.AbortWithStatusJSON(http.StatusBadRequest, err)
		return
	}

	hashpassword, err := bcrypt.GenerateFromPassword([]byte(req.NewPassword), bcrypt.DefaultCost)
	if err != nil {
		h.Log.Error(fmt.Sprintf("Pasworni hashlashda xatolik: %v", err))
		c.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	req.NewPassword = string(hashpassword)

	res, err := h.User.ConfirmationPassword(c, &req)
	if err != nil {
		h.Log.Error("ConfirmationPassword funksiyasiga malumot yuborishda xatolik", "error", err.Error())
		c.AbortWithStatusJSON(500, gin.H{
			"error": err.Error(),
		})
		return
	}

	c.JSON(200, res)
}

// @Summary      Update user details
// @Description  This endpoint updates the user's details based on the provided information.
// @Tags         user
// @Accept       json
// @Produce      json
// @Param        user  body      user.UpdateUserRequest  true  "User Update Data"
// @Success      200    {object}  user.UpdateUserResponse
// @Failure      400    {object}  string
// @Failure      500    {object}  string
// @Router       /update_user [put]
func (h Handler) UpdateUser(c *gin.Context) {
	req := pb.UpdateUserRequest{}

	if err := json.NewDecoder(c.Request.Body).Decode(&req); err != nil {
		c.AbortWithStatusJSON(http.StatusBadRequest, err)
		return
	}

	res, err := h.User.UpdateUser(c, &req)
	if err != nil {
		h.Log.Error("UpdateUser funksiyasoga xabar yuborishda xatolik", "error", err.Error())
		c.AbortWithStatusJSON(500, gin.H{
			"error": err.Error(),
		})
		return
	}

	c.JSON(200, res)
}

// @Summary      Delete user
// @Description  This endpoint deletes a user based on the provided user ID.
// @Tags         user
// @Accept       json
// @Produce      json
// @Param        id  path      string  true  "User ID"
// @Success      200    {object}  user.DeleteUserResponse
// @Failure      400    {object}  string
// @Failure      500    {object}  string
// @Router       /delete_user/{id} [delete]
func (h Handler) DeleteUser(c *gin.Context) {
	req := pb.UserId{
		Id: c.Param("id"),
	}

	res, err := h.User.DeleteUser(c, &req)
	if err != nil {
		h.Log.Error("DeleteUserga malumot yuborishda xatolik", "error", err.Error())
		c.AbortWithStatusJSON(500, gin.H{
			"error": err.Error(),
		})
		return
	}

	c.JSON(200, res)
}

// @Summary      Update user role
// @Description  This endpoint updates the user's role based on the provided email and new role.
// @Tags         user
// @Accept       json
// @Produce      json
// @Param        email  path      string  true  "User Email"
// @Param        role   path      string  true  "New User Role"
// @Success      200    {object}  user.UpdateRoleResponse
// @Failure      400    {object}  string
// @Failure      500    {object}  string
// @Router       /update_role/{email}/{role} [put]
func (h Handler) UpdateRole(c *gin.Context) {
	req := pb.UpdateRoleReq{
		Email: c.Param("email"),
		Role:  c.Param("role"),
	}

	res, err := h.User.UpdateRole(c, &req)
	if err != nil {
		h.Log.Error("Update role ga malumot yuborishda xatolik", "error", err.Error())
		c.AbortWithStatusJSON(500, gin.H{
			"error": err.Error(),
		})
		return
	}

	c.JSON(200, res)
}
