package controllers

import (
	"encoding/json"
	"net/http"
	"strconv"

	"example.com/simple-api/models"
	"example.com/simple-api/services"
	"github.com/go-chi/chi/v5"
)

type MailController struct {
	MailService services.MailService
}

func NewMail(mailService services.MailService) MailController {
	return MailController{
		MailService: mailService,
	}
}

func (mc *MailController) CreateMail(w http.ResponseWriter, r *http.Request) {
	var mail models.Mail
	if err := json.NewDecoder(r.Body).Decode(&mail); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}
	err := mc.MailService.CreateMail(&mail)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadGateway)
		return
	}
	w.Write([]byte("success"))
	w.WriteHeader(201)
}

func (mc *MailController) GetMail(w http.ResponseWriter, r *http.Request) {
	id := chi.URLParam(r, "id")
	idc, _ := strconv.Atoi(id)
	mail, err := mc.MailService.GetMail(&idc)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadGateway)
		return
	}
	json.NewEncoder(w).Encode(mail)
	w.WriteHeader(201)
}

func (mc *MailController) GetAll(w http.ResponseWriter, r *http.Request) {
	mails, err := mc.MailService.GetAll()
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadGateway)
		return
	}
	json.NewEncoder(w).Encode(mails)
	w.WriteHeader(201)
}

func (mc *MailController) UpdateMail(w http.ResponseWriter, r *http.Request) {
	var mail models.Mail
	id := chi.URLParam(r, "id")
	idc, _ := strconv.Atoi(id)
	if err := json.NewDecoder(r.Body).Decode(&mail); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}
	err := mc.MailService.UpdateMail(&idc, &mail)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadGateway)
		return
	}
	w.Write([]byte("success"))
	w.WriteHeader(201)
}

func (mc *MailController) DeleteMail(w http.ResponseWriter, r *http.Request) {
	id := chi.URLParam(r, "id")
	idc, _ := strconv.Atoi(id)
	err := mc.MailService.DeleteMail(&idc)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadGateway)
		return
	}
	w.Write([]byte("success"))
	w.WriteHeader(201)
}
