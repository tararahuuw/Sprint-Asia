package taskcontroller

import (
	"backend/entities"
	"backend/models/taskmodel"
	"encoding/json"
	"net/http"
	"strconv"
	"time"

	"github.com/gorilla/mux"
)

func GetAllDataOnGoing(w http.ResponseWriter, r *http.Request) {
	// Get tasks using taskmodel.GetAllData() or your preferred method
	tasks := taskmodel.GetAllDataOnGoing()

	// Create a custom data structure to represent tasks with subtasks
	var customTasks []struct {
		entities.Task
		SubTasks []entities.SubTask `json:"SubTasks"`
	}

	// Iterate through the tasks and update the 'expired' attribute
	for _, task := range tasks {
		// Check if the deadline is null (time.Time zero value)
		if task.Deadline.IsZero() {
			task.Expired = false
		} else {
			// Compare the deadline with the current time
			if time.Now().After(task.Deadline) {
				task.Expired = true
			} else {
				task.Expired = false
			}
		}

		// Get the subtasks for the current task
		subtasks := taskmodel.GetSubTasksByTaskID(int(task.Id))
		// completeSubTask := 0.0
		// if len(subtasks) > 0 {
		// 	for _, subtask := range subtasks {
		// 		if subtask.Complete == true {
		// 			completeSubTask++
		// 		}
		// 	}
		// 	progress := float64(completeSubTask) / float64(len(subtasks)) * 100
		// 	taskmodel.UpdateTaskProgress(int(task.Id), int(progress))

		// 	if progress == 100 {
		// 		taskmodel.UpdateTaskCompleteByID(int(task.Id), true)
		// 	} else {
		// 		taskmodel.UpdateTaskCompleteByID(int(task.Id), false)
		// 	}
		// }

		// Append the task and its subtasks to the customTasks slice
		customTasks = append(customTasks, struct {
			entities.Task
			SubTasks []entities.SubTask `json:"SubTasks"`
		}{
			Task:     task,
			SubTasks: subtasks,
		})
	}

	// Marshal the customTasks data into JSON
	taskJSON, err := json.Marshal(customTasks)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	// Set the content type to JSON
	w.Header().Set("Content-Type", "application/json")

	// Write the JSON response
	_, err = w.Write(taskJSON)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
}

func GetAllDataComplete(w http.ResponseWriter, r *http.Request) {
	// Get tasks using taskmodel.GetAllData() or your preferred method
	tasks := taskmodel.GetAllDataComplete()

	// Create a custom data structure to represent tasks with subtasks
	var customTasks []struct {
		entities.Task
		SubTasks []entities.SubTask `json:"SubTasks"`
	}

	// Iterate through the tasks and update the 'expired' attribute
	for _, task := range tasks {
		// Check if the deadline is null (time.Time zero value)
		if task.Deadline.IsZero() {
			task.Expired = false
		} else {
			// Compare the deadline with the current time
			if time.Now().After(task.Deadline) {
				task.Expired = true
			} else {
				task.Expired = false
			}
		}

		// Get the subtasks for the current task
		subtasks := taskmodel.GetSubTasksByTaskID(int(task.Id))
		// completeSubTask := 0.0
		// if len(subtasks) > 0 {
		// 	for _, subtask := range subtasks {
		// 		if subtask.Complete == true {
		// 			completeSubTask++
		// 		}
		// 	}
		// 	progress := float64(completeSubTask) / float64(len(subtasks)) * 100
		// 	taskmodel.UpdateTaskProgress(int(task.Id), int(progress))

		// 	if progress == 100 {
		// 		taskmodel.UpdateTaskCompleteByID(int(task.Id), true)
		// 	} else {
		// 		taskmodel.UpdateTaskCompleteByID(int(task.Id), false)
		// 	}
		// }

		// Append the task and its subtasks to the customTasks slice
		customTasks = append(customTasks, struct {
			entities.Task
			SubTasks []entities.SubTask `json:"SubTasks"`
		}{
			Task:     task,
			SubTasks: subtasks,
		})
	}

	// Marshal the customTasks data into JSON
	taskJSON, err := json.Marshal(customTasks)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	// Set the content type to JSON
	w.Header().Set("Content-Type", "application/json")

	// Write the JSON response
	_, err = w.Write(taskJSON)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
}

func CreateTask(w http.ResponseWriter, r *http.Request) {
	// Parse the JSON request body into a Task struct
	var task entities.Task
	err := json.NewDecoder(r.Body).Decode(&task)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	// Create the task in the database
	err = taskmodel.CreateTask(task)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
}

func UpdateTaskNameByID(w http.ResponseWriter, r *http.Request) {
	// Parse the JSON request body into a Task struct
	var updatedTask entities.Task
	err := json.NewDecoder(r.Body).Decode(&updatedTask)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	// Extract the task ID from the URL path parameters
	vars := mux.Vars(r)
	idString := vars["id"]

	// Convert the ID string to an integer
	id, err := strconv.Atoi(idString)
	if err != nil {
		http.Error(w, "Invalid ID", http.StatusBadRequest)
		return
	}

	// Update the task's name in the database
	err = taskmodel.UpdateTaskNameByID(id, updatedTask)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
}

func DeleteTaskById(w http.ResponseWriter, r *http.Request) {
	// Get the task ID from the URL path parameters
	vars := mux.Vars(r)
	idString := vars["id"]

	id, err := strconv.Atoi(idString)
	if err != nil {
		http.Error(w, "Invalid ID", http.StatusBadRequest)
		return
	}

	if err := taskmodel.DeleteTaskById(id); err != nil {
		http.Error(w, "Failed to delete task", http.StatusInternalServerError)
		return
	}

	w.Write([]byte("Task deleted successfully"))
}

func UpdateTaskCompleteByID(w http.ResponseWriter, r *http.Request) {
	// Parse the JSON request body into a Task struct
	var updatedTask entities.Task
	err := json.NewDecoder(r.Body).Decode(&updatedTask)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	// Extract the task ID from the URL path parameters
	vars := mux.Vars(r)
	idString := vars["id"]

	// Convert the ID string to an integer
	id, err := strconv.Atoi(idString)
	if err != nil {
		http.Error(w, "Invalid ID", http.StatusBadRequest)
		return
	}

	// Update the task's 'complete' attribute in the database
	err = taskmodel.UpdateTaskCompleteByID(id, updatedTask.Complete)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
}

func UpdateDeadlineTaskById(w http.ResponseWriter, r *http.Request) {
	// Parse the JSON request body into a map to extract the updated deadline
	var requestBody map[string]string
	err := json.NewDecoder(r.Body).Decode(&requestBody)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	// Extract the updated deadline string from the map (assuming it's in ISO 8601 format)
	updatedDeadlineStr, ok := requestBody["deadline"]
	if !ok {
		http.Error(w, "Missing 'deadline' field in request body", http.StatusBadRequest)
		return
	}

	// Parse the updated deadline string into a time.Time value
	updatedDeadline, err := time.Parse(time.RFC3339, updatedDeadlineStr)
	if err != nil {
		http.Error(w, "Invalid 'deadline' format", http.StatusBadRequest)
		return
	}

	// Extract the task ID from the URL path parameters
	vars := mux.Vars(r)
	idString := vars["id"]

	// Convert the ID string to an integer
	id, err := strconv.Atoi(idString)
	if err != nil {
		http.Error(w, "Invalid ID", http.StatusBadRequest)
		return
	}

	// Update the task's deadline in the database
	err = taskmodel.UpdateTaskDeadlineByID(id, updatedDeadline)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
}
