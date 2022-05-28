package controllers

import (
	"net/http"
	"strconv"

	"example.com/simple-api/models"
	"example.com/simple-api/services"
	"github.com/gin-gonic/gin"
)

type MailController struct {
	MailService services.MailService
}

func New(mailService services.MailService) MailController {
	return MailController{
		MailService: mailService,
	}
}

func (mc *MailController) CreateMail(ctx *gin.Context) {
	var mail models.Mail
	if err := ctx.ShouldBindJSON(&mail); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"message": err.Error()})
		return
	}
	err := mc.MailService.CreateMail(&mail)
	if err != nil {
		ctx.JSON(http.StatusBadGateway, gin.H{"message": err.Error()})
		return
	}
	ctx.JSON(http.StatusOK, gin.H{"message": "success"})
}

func (mc *MailController) GetMail(ctx *gin.Context) {
	id := ctx.Param("id")
	idc, _ := strconv.Atoi(id)
	mail, err := mc.MailService.GetMail(&idc)
	if err != nil {
		ctx.JSON(http.StatusBadGateway, gin.H{"message": err.Error()})
		return
	}
	ctx.JSON(http.StatusOK, mail)
}

func (mc *MailController) GetAll(ctx *gin.Context) {
	mails, err := mc.MailService.GetAll()
	if err != nil {
		ctx.JSON(http.StatusBadGateway, gin.H{"message": err.Error()})
		return
	}
	ctx.JSON(http.StatusOK, mails)
}

func (mc *MailController) UpdateMail(ctx *gin.Context) {
	var mail models.Mail
	id := ctx.Param("id")
	idc, _ := strconv.Atoi(id)
	if err := ctx.ShouldBindJSON(&mail); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"message": err.Error()})
		return
	}
	err := mc.MailService.UpdateMail(&idc, &mail)
	if err != nil {
		ctx.JSON(http.StatusBadGateway, gin.H{"message": err.Error()})
		return
	}
	ctx.JSON(http.StatusOK, gin.H{"message": "success"})
}

func (mc *MailController) DeleteMail(ctx *gin.Context) {
	id := ctx.Param("id")
	idc, _ := strconv.Atoi(id)
	err := mc.MailService.DeleteMail(&idc)
	if err != nil {
		ctx.JSON(http.StatusBadGateway, gin.H{"message": err.Error()})
		return
	}
	ctx.JSON(http.StatusOK, gin.H{"message": "success"})
}

func (mc *MailController) MailRoutes(rg *gin.RouterGroup) {
	route := rg.Group("/mail")
	route.POST("", mc.CreateMail)
	route.GET("", mc.GetAll)
	route.GET("/:id", mc.GetMail)
	route.DELETE("/:id", mc.DeleteMail)
	route.PUT("/:id", mc.UpdateMail)
}
