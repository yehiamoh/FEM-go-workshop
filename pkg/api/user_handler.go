package api

import (
	"encoding/json"
	"errors"
	"log"
	"net/http"
	"regexp"

	"github.com/yehiamoh/go-fem-workshop/pkg/store"
	"github.com/yehiamoh/go-fem-workshop/pkg/utils"
)

type registerUserRequest struct {
	Username string `json:"user_name"`
	Email    string `json:"email"`
	Password string `json:"password"`
	Bio      string `json:"bio"`
}
type UserHandler struct {
	userStore store.UserStore
	logger    *log.Logger
}

func NewUserHandler(userStore store.UserStore, logger *log.Logger) *UserHandler {
	return &UserHandler{
		userStore: userStore,
		logger:    logger,
	}
}

func (h *UserHandler) validateRegisterRequest(req *registerUserRequest) error {
	if req.Username == "" {
		return errors.New("username is required")
	}
	if len(req.Username) > 50 {
		return errors.New("username can't be more than 50 character")
	}
	if req.Email == "" {
		return errors.New("email is required")
	}
	emailRegex := regexp.MustCompile(`^[a-zA-Z0-9._%+-]+@[a-zA-Z0-9.-]+\.[a-zA-Z]{2,}$`)
	if !emailRegex.MatchString(req.Email) {
		return errors.New("invalid email format")
	}
	if req.Password == "" {
		return errors.New("password is Requierd")
	}
	return nil
}
func (h *UserHandler) HandleRegisterUser(w http.ResponseWriter, r *http.Request) {
	var regitserReq registerUserRequest

	err := json.NewDecoder(r.Body).Decode(&regitserReq)
	if err != nil {
		h.logger.Fatalf("Error : Decoding Register Body : %v", err)
		utils.WriteJSON(w, http.StatusBadRequest, utils.Envelope{"error": "Invalid request payload"})
		return
	}
	defer r.Body.Close()

	err = h.validateRegisterRequest(&regitserReq)
	if err != nil {
		h.logger.Fatalf("Error : Decoding Register Body : %v", err)
		utils.WriteJSON(w, http.StatusBadRequest, utils.Envelope{"error": err.Error()})
		return
	}

	user := &store.User{
		Username: regitserReq.Username,
		Email:    regitserReq.Email,
		Bio:      regitserReq.Bio,
	}
	err = user.PasswordHash.Set(regitserReq.Password)
	if err != nil {
		h.logger.Fatalf("Error : Hashing : %v", err)
		utils.WriteJSON(w, http.StatusInternalServerError, utils.Envelope{"error": "Internal server error"})
		return
	}
	err = h.userStore.CreateUser(user)
	if err != nil {
		h.logger.Fatalf("Error : Creating User in the database :%v", err)
		utils.WriteJSON(w, http.StatusInternalServerError, utils.Envelope{"error": "Internal server error"})
		return
	}
	utils.WriteJSON(w, http.StatusCreated, utils.Envelope{"user": user})
}
