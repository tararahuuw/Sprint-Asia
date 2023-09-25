package subtaskcontroller

import (
	"backend/entities"
	"backend/models/subtaskmodel"
	"backend/models/taskmodel"
	"encoding/json"
	"fmt"
	"net/http"
	"strconv"

	"github.com/gorilla/mux"
)

func CreateSubTask(w http.ResponseWriter, r *http.Request) {
	// Parse the JSON request body into a SubTask struct
	var subtask entities.SubTask
	err := json.NewDecoder(r.Body).Decode(&subtask)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	// Create the subtask in the database and get the ID of the newly created subtask
	subtaskID, err := subtaskmodel.CreateSubTask(subtask)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	fmt.Printf("id: %+v\n", subtaskID)
	subtask2 := subtaskmodel.GetSubTasksByID(subtaskID)
	taskID := int(subtask2[0].Id_Task)
	fmt.Printf("id task: %+v\n", taskID)
	subtasks := taskmodel.GetSubTasksByTaskID(int(subtask2[0].Id_Task))
	fmt.Printf("subtask: %+v\n", subtask) // Assume this function returns a single SubTask
	if subtasks != nil {
		totalSubtasks := len(subtasks)
		fmt.Printf("totalSubtasks: %+v\n", totalSubtasks)
		// Retrieve the total number of completed subtasks for the task
		completedSubtasks := taskmodel.GetCompletedSubTasksByTaskID(taskID)
		fmt.Printf("completedSubtask: %+v\n", completedSubtasks)
		// Calculate progress
		progress := 0
		if totalSubtasks > 0 {
			progress = (completedSubtasks * 100) / totalSubtasks
		}
		fmt.Printf("progress: %+v\n", progress)
		// Update task progress and completion status
		taskmodel.UpdateTaskProgress(taskID, progress)
		taskmodel.UpdateTaskCompleteByID(taskID, progress == 100)
	}
}

func UpdateSubTaskNameByID(w http.ResponseWriter, r *http.Request) {
	// Parse the JSON request body into a SubTask struct
	var updatedSubTask entities.SubTask
	err := json.NewDecoder(r.Body).Decode(&updatedSubTask)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	// Extract the subtask ID from the URL path parameters
	vars := mux.Vars(r)
	idString := vars["id"]

	// Convert the ID string to an integer
	id, err := strconv.Atoi(idString)
	if err != nil {
		http.Error(w, "Invalid ID", http.StatusBadRequest)
		return
	}

	fmt.Printf("updatedSubTask: %+v\n", updatedSubTask)
	// Update the subtask's name in the database
	err = subtaskmodel.UpdateSubTaskNameByID(id, updatedSubTask)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
}

func DeleteSubTaskById(w http.ResponseWriter, r *http.Request) {
	// Get the task ID from the URL path parameters
	vars := mux.Vars(r)
	idString := vars["id"]

	id, err := strconv.Atoi(idString)
	if err != nil {
		http.Error(w, "Invalid ID", http.StatusBadRequest)
		return
	}

	fmt.Printf("id: %+v\n", id)
	subtask := subtaskmodel.GetSubTasksByID(id)
	taskID := int(subtask[0].Id_Task)
	fmt.Printf("id task: %+v\n", taskID)

	if err := subtaskmodel.DeleteSubTaskById(id); err != nil {
		http.Error(w, "Failed to delete subtask", http.StatusInternalServerError)
		return
	}
	subtasks := taskmodel.GetSubTasksByTaskID(int(subtask[0].Id_Task))
	fmt.Printf("subtask: %+v\n", subtask) // Assume this function returns a single SubTask
	if subtasks != nil {
		totalSubtasks := len(subtasks)
		fmt.Printf("totalSubtasks: %+v\n", totalSubtasks)
		// Retrieve the total number of completed subtasks for the task
		completedSubtasks := taskmodel.GetCompletedSubTasksByTaskID(taskID)
		fmt.Printf("completedSubtask: %+v\n", completedSubtasks)
		// Calculate progress
		progress := 0
		if totalSubtasks > 0 {
			progress = (completedSubtasks * 100) / totalSubtasks
		}
		fmt.Printf("progress: %+v\n", progress)
		// Update task progress and completion status
		taskmodel.UpdateTaskProgress(taskID, progress)
		taskmodel.UpdateTaskCompleteByID(taskID, progress == 100)
	}
}

// UpdateSubTaskCompleteById updates the 'complete' attribute of a subtask by its ID.
func UpdateSubTaskCompleteById(w http.ResponseWriter, r *http.Request) {
	// Parse the JSON request body into a SubTask struct
	var updatedSubTask entities.SubTask
	err := json.NewDecoder(r.Body).Decode(&updatedSubTask)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	// Extract the subtask ID from the URL path parameters
	vars := mux.Vars(r)
	idString := vars["id"]

	// Convert the ID string to an integer
	id, err := strconv.Atoi(idString)
	if err != nil {
		http.Error(w, "Invalid ID", http.StatusBadRequest)
		return
	}

	// Update the subtask's 'complete' attribute in the database
	err = subtaskmodel.UpdateSubTaskCompleteByID(id, updatedSubTask.Complete)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	// Calculate task progress and update it in the database
	subtask := subtaskmodel.GetSubTasksByID(id) // Assume this function returns a single SubTask
	if subtask != nil {
		taskID := int(subtask[0].Id_Task)
		// Retrieve the total number of subtasks for the task
		totalSubtasks := taskmodel.GetTotalSubTasksByTaskID(taskID)
		// Retrieve the total number of completed subtasks for the task
		completedSubtasks := taskmodel.GetCompletedSubTasksByTaskID(taskID)
		// Calculate progress
		progress := 0
		if totalSubtasks > 0 {
			progress = (completedSubtasks * 100) / totalSubtasks
		}
		// Update task progress and completion status
		taskmodel.UpdateTaskProgress(taskID, progress)
		taskmodel.UpdateTaskCompleteByID(taskID, progress == 100)
	}
}
