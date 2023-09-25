package subtaskmodel

import (
	"backend/config"
	"backend/entities"
	"fmt"
	"time"
)

func scanTime(dest *time.Time) interface{} {
	return (*MySQLTime)(dest)
}

type MySQLTime time.Time

func (t *MySQLTime) Scan(v interface{}) error {
	val, ok := v.([]uint8)
	if !ok {
		return fmt.Errorf("Expected []uint8, got %T", v)
	}

	if string(val) == "0000-00-00 00:00:00" {
		*t = MySQLTime(time.Time{}) // Set to zero time
		return nil
	}

	tm, err := time.Parse("2006-01-02 15:04:05", string(val))
	if err != nil {
		return err
	}

	*t = MySQLTime(tm)
	return nil
}

func GetSubTasksByID(id int) []entities.SubTask {
	// Define the SQL query to select subtasks by task ID
	query := `
        SELECT id, id_task, name, complete, created_at, updated_at
        FROM subtask
        WHERE id = ?
    `

	// Execute the SQL query with the task ID
	rows, err := config.DB.Query(query, id)
	if err != nil {
		// Handle the error if needed
		// You can log the error or return an empty slice
		return nil
	}
	defer rows.Close()

	var subtasks []entities.SubTask

	for rows.Next() {
		var subtask entities.SubTask

		// Scan the row data into a SubTask struct
		err := rows.Scan(&subtask.Id, &subtask.Id_Task, &subtask.Name, &subtask.Complete, scanTime(&subtask.CreatedAt), scanTime(&subtask.UpdatedAt))

		if err != nil {
			// Handle the error as needed, e.g., log it or skip the subtask
			fmt.Println("Error:", err)
			continue
		}

		subtasks = append(subtasks, subtask)
	}

	return subtasks
}

func CreateSubTask(subtask entities.SubTask) (int, error) {
	// Define the SQL query to insert a new subtask
	query := `
		INSERT INTO subtask (id_task, name, complete, created_at, updated_at)
		VALUES (?, ?, ?, ?, ?)
	`

	// Execute the SQL query with the provided subtask data
	result, err := config.DB.Exec(query,
		subtask.Id_Task,
		subtask.Name,
		subtask.Complete,
		subtask.CreatedAt.Format("2006-01-02 15:04:05"), // Format the time as a string
		subtask.UpdatedAt.Format("2006-01-02 15:04:05"), // Format the time as a string
	)

	if err != nil {
		return 0, err
	}

	// Get the ID of the newly created subtask
	lastInsertID, err := result.LastInsertId()
	if err != nil {
		return 0, err
	}

	return int(lastInsertID), nil
}

// Function to update the name of a subtask by its ID
func UpdateSubTaskNameByID(id int, updatedSubTask entities.SubTask) error {
	// Define the SQL query to update the subtask's name
	query := `
		UPDATE subtask
		SET name = ?
		WHERE id = ?
	`

	// Execute the SQL query with the updated subtask name and ID
	result, err := config.DB.Exec(query, updatedSubTask.Name, id)
	if err != nil {
		return err
	}

	// Check the number of rows affected to verify if the update was successful
	rowsAffected, err := result.RowsAffected()
	if err != nil {
		return err
	}

	if rowsAffected == 0 {
		return fmt.Errorf("SubTask with ID %d not found", id)
	}

	return nil
}

func DeleteSubTaskById(id int) error {
	// Define the SQL query to delete the subtask by ID
	query := "DELETE FROM subtask WHERE id = ?"

	// Execute the SQL query with the task's ID
	result, err := config.DB.Exec(query, id)
	if err != nil {
		return err
	}

	// Check the number of rows affected to verify if the delete was successful
	rowsAffected, err := result.RowsAffected()
	if err != nil {
		return err
	}

	if rowsAffected == 0 {
		return fmt.Errorf("SubTask with ID %d not found", id)
	}

	return nil
}

func UpdateSubTaskCompleteByID(id int, complete bool) error {
	// Define the SQL query to update the subtask's 'complete' attribute by ID
	query := `
		UPDATE subtask
		SET complete = ?
		WHERE id = ?
	`

	// Execute the SQL query with the updated 'complete' value and task ID
	_, err := config.DB.Exec(query, complete, id)
	if err != nil {
		return err
	}

	return nil
}

// func UpdateTaskProgress(taskID int) error {
// 	// Calculate the progress for the task with the given ID
// 	var completedSubtasks, totalSubtasks int

// 	// Query the database to count the completed and total subtasks for the task
// 	err := config.DB.QueryRow(`
// 		SELECT
// 			COUNT(CASE WHEN complete = true THEN 1 ELSE NULL END) AS completed_subtasks,
// 			COUNT(*) AS total_subtasks
// 		FROM subtask
// 		WHERE id_task = ?`, taskID).Scan(&completedSubtasks, &totalSubtasks)

// 	if err != nil {
// 		return err
// 	}

// 	// Calculate the progress percentage
// 	var progress int
// 	if totalSubtasks > 0 {
// 		progress = (completedSubtasks * 100) / totalSubtasks
// 	}

// 	// Update the task's progress in the database
// 	_, err = config.DB.Exec(`
// 		UPDATE task
// 		SET progress = ?
// 		WHERE id = ?`, progress, taskID)

// 	if err != nil {
// 		return err
// 	}

// 	return nil
// }

// // CountCompletedSubtasksByTaskID counts the number of completed subtasks for a given task ID.
// func CountCompletedSubtasksByTaskID(taskID int) (int, error) {
// 	// Define the SQL query to count completed subtasks for the given task ID
// 	query := `
// 		SELECT COUNT(*) FROM subtask
// 		WHERE id_task = ? AND complete = true
// 	`

// 	// Execute the SQL query and store the result in the `count` variable
// 	var count int
// 	err := config.DB.QueryRow(query, taskID).Scan(&count)
// 	if err != nil {
// 		return 0, err
// 	}

// 	return count, nil
// }

// // CountTotalSubtasksByTaskID counts the total number of subtasks for a given task ID.
// func CountTotalSubtasksByTaskID(taskID int) (int, error) {
// 	// Define the SQL query to count all subtasks for the given task ID
// 	query := `
// 		SELECT COUNT(*) FROM subtask
// 		WHERE id_task = ?
// 	`

// 	// Execute the SQL query and store the result in the `count` variable
// 	var count int
// 	err := config.DB.QueryRow(query, taskID).Scan(&count)
// 	if err != nil {
// 		return 0, err
// 	}

// 	return count, nil
// }
