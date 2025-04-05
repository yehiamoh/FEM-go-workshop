package api

import (
	"database/sql"
	"encoding/json"
	"fmt"
	"net/http"
	"strconv"

	"github.com/go-chi/chi/v5"
	"github.com/yehiamoh/go-fem-workshop/internal/store"
)

type WorkoutHandler struct {
	Workoutstore store.WorkoutStore
}

func NewWorkoutHandler(workoutStore store.WorkoutStore) *WorkoutHandler {
	return &WorkoutHandler{
		Workoutstore: workoutStore,
	}
}

func (wh *WorkoutHandler) HandleGetWorkOutByID(w http.ResponseWriter, r *http.Request) {
	paramsWorkoutID := chi.URLParam(r, "id")
	if paramsWorkoutID == "" {
		http.NotFound(w, r)
		return
	}
	workoutID, err := strconv.Atoi(paramsWorkoutID)
	if err != nil {
		http.NotFound(w, r)
		return
	}
	workout, err := wh.Workoutstore.GetWorkoutByID(workoutID)
	if err != nil {
		http.Error(w, "Failed to retrieve the workout", http.StatusNotFound)
		return
	}
	w.WriteHeader(http.StatusOK)
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(workout)

	//fmt.Fprintf(w, "this is the workout id %d\n", workoutID)
}
func (wh *WorkoutHandler) HandleCreateWorkout(w http.ResponseWriter, r *http.Request) {
	var workout store.Workout
	err := json.NewDecoder(r.Body).Decode(&workout)
	if err != nil {
		fmt.Println(err)
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}
	createdWorkout, err := wh.Workoutstore.CreateWorkout(&workout)
	if err != nil {
		fmt.Println(err)
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(createdWorkout)
}
func (wh *WorkoutHandler) HandleUpdateWorkoutByID(w http.ResponseWriter, r *http.Request) {
	paramsWorkoutID := chi.URLParam(r, "id")
	if paramsWorkoutID == "" {
		http.NotFound(w, r)
		return
	}
	workoutID, err := strconv.Atoi(paramsWorkoutID)
	if err != nil {
		http.NotFound(w, r)
		return
	}
	existingWorkout, err := wh.Workoutstore.GetWorkoutByID(workoutID)
	if err != nil {
		http.Error(w, "Failed to Fetch workout", http.StatusInternalServerError)
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
		http.Error(w, err.Error(), http.StatusBadRequest)
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
	if err := wh.Workoutstore.UpdateWorkout(existingWorkout); err != nil {
		http.Error(w, "Failed to update workout", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(existingWorkout)
}
func (wh *WorkoutHandler) HandleDeleteWorkout(w http.ResponseWriter, r *http.Request) {
	paramWorkoutId := chi.URLParam(r, "id")
	if paramWorkoutId == "" {
		http.NotFound(w, r)
		return
	}
	workoutID, err := strconv.Atoi(paramWorkoutId)
	if err != nil {
		http.NotFound(w, r)
		return
	}
	err = wh.Workoutstore.DeleteWorkout(workoutID)
	if err == sql.ErrNoRows {
		http.Error(w, "Workout not found", http.StatusNotFound)
		return
	}
	if err != nil {
		http.Error(w, "Failed To Delete the workout", http.StatusInternalServerError)
		return
	}
	w.WriteHeader(http.StatusOK)
	response := map[string]interface{}{
		"message": fmt.Sprintf("Workout with ID %d has been successfully delted", workoutID),
		"status":  "success",
	}
	json.NewEncoder(w).Encode(response)
}
