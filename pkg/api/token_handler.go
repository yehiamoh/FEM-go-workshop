package api

import (
	"encoding/json"
	"log"
	"net/http"
	"time"

	"github.com/yehiamoh/go-fem-workshop/pkg/store"
	"github.com/yehiamoh/go-fem-workshop/pkg/tokens"
	"github.com/yehiamoh/go-fem-workshop/pkg/utils"
)

type TokenHandler struct {
	tokenStore store.TokenStore
	userStore  store.UserStore
	logger     *log.Logger
}
type createtokenRequest struct {
	Username string `json:"user_name"`
	Password string `json:"password"`
}

func NewTokenHandler(tokenStore store.TokenStore, userStore store.UserStore, logger *log.Logger) *TokenHandler {
	return &TokenHandler{
		tokenStore: tokenStore,
		userStore:  userStore,
		logger:     logger,
	}
}
func (h *TokenHandler) HandleCreateToken(w http.ResponseWriter, r *http.Request) {
	var req createtokenRequest
	err := json.NewDecoder(r.Body).Decode(&req)
	if err != nil {
		h.logger.Printf("ERROR : createTokenRequest :%v", err)
		utils.WriteJSON(w, http.StatusBadRequest, utils.Envelope{"error": "inaavalid request payload"})
		return
	}

	user, err := h.userStore.GetUserByUserName(req.Username)
	if err != nil || user == nil {
		h.logger.Printf("ERROR: GetUserByUserName: %v ", err)
		utils.WriteJSON(w, http.StatusInternalServerError, utils.Envelope{"error": "internal server error"})
		return
	}
	DoesPassMatch, err := user.PasswordHash.IsMatch(req.Password)
	if err != nil {
		h.logger.Printf("ERROR: passoword hash matching: %v ", err)
		utils.WriteJSON(w, http.StatusInternalServerError, utils.Envelope{"error": "internal server error"})
		return
	}
	if !DoesPassMatch {
		utils.WriteJSON(w, http.StatusUnauthorized, utils.Envelope{"error": "invalid credintials"})
		return
	}
	token, err := h.tokenStore.CreateNewToken(user.ID, 24*time.Hour, tokens.ScopeAuth)
	if err != nil {
		h.logger.Printf("ERROR: Creating Token: %v ", err)
		utils.WriteJSON(w, http.StatusInternalServerError, utils.Envelope{"error": "internal server error"})
		return
	}
	utils.WriteJSON(w, http.StatusCreated, utils.Envelope{"auth_token": token})
}
