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

	res,err:=h.User.ConfirmationRegister(c,&req)
	if err!=nil{
		h.Log.Error("ConfirmationRegister funksiyasiga yuborishda xatolik.","error",err)
		c.AbortWithStatusJSON(400,gin.H{
			"error":err.Error(),
		})
		return
	}

	c.JSON(200,res)
}

func (h Handler) GetUSerByEmail(c *gin.Context){
	req:=pb.GetUSerByEmailReq{
		Email: c.Param("email"),
	}

	res,err:=h.User.GetUSerByEmail(c,&req)
	if err!=nil{
		h.Log.Error("GetUSerByEmail funksiyasida xatolik.","error",err.Error())
		c.AbortWithStatusJSON(500,gin.H{
			"error":err.Error(),
		})
		return
	}

	c.JSON(200,res)
}

func (h Handler) UpdatePassword(c *gin.Context){
	req:=pb.UpdatePasswordReq{
		OldPassword: c.Param("old password"),
		NewPassword: c.Param("new password"),
		Email: c.Param("email"),
	}

	res,err:=h.User.UpdatePassword(c,&req)
	if err!=nil{
		h.Log.Error("UpdatePassword funksiyasiga malumot yuboishda xatolik.","error",err.Error())
		c.AbortWithStatusJSON(500,gin.H{
			"error":err.Error(),
		})
		return
	}

	c.JSON(200,res)
}

func (h Handler) ResetPassword(c *gin.Context){
	req:=pb.ResetPasswordReq{
		Email: c.Param("email"),
	}

	res,err:=email.Email(req.Email)
	if err!=nil{
		h.Log.Error("Email ga xabar yuuborishda xatolik.","error",err.Error())
		c.AbortWithStatusJSON(400,gin.H{
			"error":err.Error(),
		})
		return
	}

	c.JSON(200,res)
}

func (h Handler) ConfirmationPassword(c *gin.Context){
	req:=pb.ConfirmationReq{}

	if err := json.NewDecoder(c.Request.Body).Decode(&req); err != nil {
		c.AbortWithStatusJSON(http.StatusBadRequest, err)
		return
	}

	res,err:=h.User.ConfirmationPassword(c,&req)
	if err!=nil{
		h.Log.Error("ConfirmationPassword funksiyasiga malumot yuborishda xatolik","error",err.Error())
		c.AbortWithStatusJSON(500,gin.H{
			"error":err.Error(),
		})
		return
	}

	c.JSON(200,res)
}

func (h Handler) UpdateUser(c *gin.Context){
	req:=pb.UpdateUserRequest{}

	if err := json.NewDecoder(c.Request.Body).Decode(&req); err != nil {
		c.AbortWithStatusJSON(http.StatusBadRequest, err)
		return
	}
	
	res,err:=h.User.UpdateUser(c,&req)
	if err!=nil{
		h.Log.Error("UpdateUser funksiyasoga xabar yuborishda xatolik","error",err.Error())
		c.AbortWithStatusJSON(500,gin.H{
			"error":err.Error(),
		})
		return
	}

	c.JSON(200,res)
}

func (h Handler) DeleteUser(c *gin.Context){
	req:=pb.UserId{
		Id: c.Param("id"),
	}

	res,err:=h.User.DeleteUser(c,&req)
	if err!=nil{
		h.Log.Error("DeleteUserga malumot yuborishda xatolik","error",err.Error())
		c.AbortWithStatusJSON(500,gin.H{
			"error":err.Error(),
		})
		return
	}

	c.JSON(200,res)
}

func (h Handler) UpdateRole(c *gin.Context){
	
}