package api

import (
	"database/sql"
	"encoding/json"
	"errors"
	"log"
	"net/http"

	"github.com/yehiamoh/go-fem-workshop/pkg/middleware"
	"github.com/yehiamoh/go-fem-workshop/pkg/store"
	"github.com/yehiamoh/go-fem-workshop/pkg/utils"
)

type WorkoutHandler struct {
	Workoutstore store.WorkoutStore
	Logger       *log.Logger
}

func NewWorkoutHandler(workoutStore store.WorkoutStore, logger *log.Logger) *WorkoutHandler {
	return &WorkoutHandler{
		Workoutstore: workoutStore,
		Logger:       logger,
	}
}

func (wh *WorkoutHandler) HandleGetWorkOutByID(w http.ResponseWriter, r *http.Request) {
	workoutID, err := utils.ReadIDParam(r)
	if err != nil {
		wh.Logger.Printf("ERROR : read ID Param : %v", err)
		utils.WriteJSON(w, http.StatusBadRequest, utils.Envelope{"error": "Invlaid workout Id"})
		return
	}
	workout, err := wh.Workoutstore.GetWorkoutByID(workoutID)
	if err != nil {
		wh.Logger.Printf("ERROR : Get Workout By ID : %v", err)
		utils.WriteJSON(w, http.StatusInternalServerError, utils.Envelope{"error": "Internal Server Error"})
		return
	}
	if workout == nil {
		utils.WriteJSON(w, http.StatusNotFound, utils.Envelope{"message": "Workout not found"})
		return
	}
	utils.WriteJSON(w, http.StatusOK, utils.Envelope{"workout": workout})

}
func (wh *WorkoutHandler) HandleCreateWorkout(w http.ResponseWriter, r *http.Request) {
	var workout store.Workout
	err := json.NewDecoder(r.Body).Decode(&workout)
	if err != nil {
		wh.Logger.Printf("ERROR : Reading RequestBody : %v", err)
		utils.WriteJSON(w, http.StatusBadRequest, utils.Envelope{"error": "Invalid Request payload"})
		return
	}

	currentUser := middleware.GetUser(r)
	if currentUser == nil || currentUser == store.AnonymousUser {
		wh.Logger.Printf("ERROR : Reading RequestBody : %v", err)
		utils.WriteJSON(w, http.StatusBadRequest, utils.Envelope{"error": "User Must be Logged in"})
		return
	}
	workout.UserID = currentUser.ID

	createdWorkout, err := wh.Workoutstore.CreateWorkout(&workout)
	if err != nil {
		wh.Logger.Printf("ERROR : Creating workout : %v", err)
		utils.WriteJSON(w, http.StatusInternalServerError, utils.Envelope{"error": "failed to create workout"})
		return
	}
	utils.WriteJSON(w, http.StatusCreated, utils.Envelope{"workout": createdWorkout})
}
func (wh *WorkoutHandler) HandleUpdateWorkoutByID(w http.ResponseWriter, r *http.Request) {
	workoutID, err := utils.ReadIDParam(r)
	if err != nil {
		wh.Logger.Printf("ERROR : read ID Param : %v", err)
		utils.WriteJSON(w, http.StatusBadRequest, utils.Envelope{"error": "Invlaid workout Id"})
		return
	}
	existingWorkout, err := wh.Workoutstore.GetWorkoutByID(workoutID)
	if err != nil {
		wh.Logger.Printf("ERROR : Get Workout By ID : %v", err)
		utils.WriteJSON(w, http.StatusInternalServerError, utils.Envelope{"error": "Internal Server Error"})
		return
	}
	if existingWorkout == nil {
		http.NotFound(w, r)
		return
	}
	var updateWorkoutRequest struct {
		Title           *string               `json:"title"`
		Description     *string               `json:"description"`
		DurationMinutes *int                  `json:"duration_minutes"`
		CaloriesBurned  *int                  `json:"calories_burned"`
		Entries         *[]store.WorkoutEntry `json:"entries"`
	}
	err = json.NewDecoder(r.Body).Decode(&updateWorkoutRequest)
	if err != nil {
		wh.Logger.Printf("ERROR : Reading RequestBody : %v", err)
		utils.WriteJSON(w, http.StatusBadRequest, utils.Envelope{"error": "Invalid request Payload"})
		return
	}
	if updateWorkoutRequest.Title != nil {
		existingWorkout.Title = *updateWorkoutRequest.Title
	}
	if updateWorkoutRequest.Description != nil {
		existingWorkout.Description = *updateWorkoutRequest.Description
	}
	if updateWorkoutRequest.DurationMinutes != nil {
		existingWorkout.DurationMinutes = *updateWorkoutRequest.DurationMinutes
	}
	if updateWorkoutRequest.CaloriesBurned != nil {
		existingWorkout.CaloriesBurned = *updateWorkoutRequest.CaloriesBurned
	}
	if updateWorkoutRequest.Entries != nil {
		existingWorkout.Entries = *updateWorkoutRequest.Entries
	}
	currentUser := middleware.GetUser(r)
	if currentUser == nil || currentUser == store.AnonymousUser {
		wh.Logger.Printf("ERROR : Reading RequestBody : %v", err)
		utils.WriteJSON(w, http.StatusBadRequest, utils.Envelope{"error": "User Must be Logged in"})
		return
	}
	workoutOwner, err := wh.Workoutstore.GetWorkoutOwner(workoutID)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			utils.WriteJSON(w, http.StatusNotFound, utils.Envelope{"error": "Workout doesn't exists"})
			return
		}
		utils.WriteJSON(w, http.StatusInternalServerError, utils.Envelope{"error": " Internal Server Error"})
		return
	}
	if workoutOwner != currentUser.ID {
		utils.WriteJSON(w, http.StatusForbidden, utils.Envelope{"error": " you aren't authorized to update the workout"})
		return
	}
	if err := wh.Workoutstore.UpdateWorkout(existingWorkout); err != nil {
		wh.Logger.Printf("ERROR : Update Workout : %v", err)
		utils.WriteJSON(w, http.StatusInternalServerError, utils.Envelope{"error": "Internal Server Error"})
		return
	}

	utils.WriteJSON(w, http.StatusOK, utils.Envelope{"workout": existingWorkout})
}
func (wh *WorkoutHandler) HandleDeleteWorkout(w http.ResponseWriter, r *http.Request) {
	workoutID, err := utils.ReadIDParam(r)
	if err != nil {
		wh.Logger.Printf("ERROR : read ID Param : %v", err)
		utils.WriteJSON(w, http.StatusBadRequest, utils.Envelope{"error": "Invlaid workout Id"})
		return
	}
	currentUser := middleware.GetUser(r)
	if currentUser == nil || currentUser == store.AnonymousUser {
		wh.Logger.Printf("ERROR : Reading RequestBody : %v", err)
		utils.WriteJSON(w, http.StatusBadRequest, utils.Envelope{"error": "User Must be Logged in"})
		return
	}
	workoutOwner, err := wh.Workoutstore.GetWorkoutOwner(workoutID)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			utils.WriteJSON(w, http.StatusNotFound, utils.Envelope{"error": "Workout doesn't exists"})
			return
		}
		utils.WriteJSON(w, http.StatusInternalServerError, utils.Envelope{"error": " Internal Server Error"})
		return
	}
	if workoutOwner != currentUser.ID {
		utils.WriteJSON(w, http.StatusForbidden, utils.Envelope{"error": " you aren't authorized to update the workout"})
		return
	}
	err = wh.Workoutstore.DeleteWorkout(workoutID)
	if err == sql.ErrNoRows {
		wh.Logger.Printf("ERROR : read ID Param : %v", err)
		utils.WriteJSON(w, http.StatusNotFound, utils.Envelope{"error": "Workout not found"})
		return
	}
	if err != nil {
		wh.Logger.Printf("ERROR : Deleting Workout : %v", err)
		utils.WriteJSON(w, http.StatusInternalServerError, utils.Envelope{"error": "Internal Server Error"})
		return
	}
	utils.WriteJSON(w, http.StatusOK, utils.Envelope{"message": "Workout Deleted"})
}
