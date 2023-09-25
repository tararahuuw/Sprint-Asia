package taskmodel

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

func GetAllDataOnGoing() []entities.Task {
	rows, err := config.DB.Query(`SELECT * FROM task where complete = false`)
	if err != nil {
		panic(err)
	}

	defer rows.Close()

	var tasks []entities.Task

	for rows.Next() {
		var task entities.Task

		if err := rows.Scan(&task.Id, &task.Name, scanTime(&task.Deadline), &task.Complete, &task.Progress, &task.Expired, scanTime(&task.CreatedAt), scanTime(&task.UpdatedAt)); err != nil {
			panic(err)
		}

		tasks = append(tasks, task)
	}

	return tasks
}

func GetAllDataComplete() []entities.Task {
	rows, err := config.DB.Query(`SELECT * FROM task where complete = true`)
	if err != nil {
		panic(err)
	}

	defer rows.Close()

	var tasks []entities.Task

	for rows.Next() {
		var task entities.Task

		if err := rows.Scan(&task.Id, &task.Name, scanTime(&task.Deadline), &task.Complete, &task.Progress, &task.Expired, scanTime(&task.CreatedAt), scanTime(&task.UpdatedAt)); err != nil {
			panic(err)
		}

		tasks = append(tasks, task)
	}

	return tasks
}

func CreateTask(task entities.Task) error {
	// Define the SQL query to insert a new task
	query := `
		INSERT INTO task (name, deadline, complete, progress, expired, created_at, updated_at)
		VALUES (?, ?, ?, ?, ?, ?, ?)
	`

	// Execute the SQL query with the provided task data
	_, err := config.DB.Exec(query,
		task.Name,
		task.Deadline.Format("2006-01-02 15:04:05"), // Format the time as a string
		task.Complete,
		task.Progress,
		task.Expired,
		task.CreatedAt.Format("2006-01-02 15:04:05"), // Format the time as a string
		task.UpdatedAt.Format("2006-01-02 15:04:05"), // Format the time as a string
	)

	if err != nil {
		return err
	}

	return nil
}

// Function to update the name of a task by its ID
func UpdateTaskNameByID(id int, updatedTask entities.Task) error {
	// Define the SQL query to update the task's name
	query := `
		UPDATE task
		SET name = ?
		WHERE id = ?
	`

	// Execute the SQL query with the updated task name and ID
	result, err := config.DB.Exec(query, updatedTask.Name, id)
	if err != nil {
		return err
	}

	// Check the number of rows affected to verify if the update was successful
	rowsAffected, err := result.RowsAffected()
	if err != nil {
		return err
	}

	if rowsAffected == 0 {
		return fmt.Errorf("Task with ID %d not found", id)
	}

	return nil
}

func DeleteTaskById(id int) error {
	// Define the SQL query to delete the task by ID
	query := "DELETE FROM task WHERE id = ?"

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
		return fmt.Errorf("Task with ID %d not found", id)
	}

	return nil
}

func UpdateTaskCompleteByID(id int, complete bool) error {
	// Define the SQL query to update the task's 'complete' attribute by ID
	query := `
		UPDATE task
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

// GetSubTasksByTaskID retrieves subtasks associated with a task by its ID
func GetSubTasksByTaskID(taskID int) []entities.SubTask {
	// Define the SQL query to select subtasks by task ID
	query := "SELECT id, id_task, name, complete FROM subtask WHERE id_task = ?"

	// Execute the query and retrieve subtasks
	rows, err := config.DB.Query(query, taskID)
	if err != nil {
		// Handle the error as needed, e.g., log it or return an empty slice
		fmt.Println("Error:", err)
		return nil
	}
	defer rows.Close()

	// Create a slice to hold subtasks
	var subtasks []entities.SubTask

	// Iterate through the result set and populate subtasks
	for rows.Next() {
		var subtask entities.SubTask
		err := rows.Scan(&subtask.Id, &subtask.Id_Task, &subtask.Name, &subtask.Complete)
		if err != nil {
			// Handle the error as needed, e.g., log it or skip the subtask
			fmt.Println("Error:", err)
			continue
		}
		subtasks = append(subtasks, subtask)
	}

	if err := rows.Err(); err != nil {
		// Handle the error as needed, e.g., log it
		fmt.Println("Error:", err)
	}

	if subtasks == nil {
		subtasks = []entities.SubTask{}
	}

	return subtasks
}

// UpdateTaskDeadlineByID updates the deadline of a task by its ID in the database.
func UpdateTaskDeadlineByID(id int, updatedDeadline time.Time) error {
	// Prepare the SQL query to update the task's deadline
	query := `
		UPDATE task
		SET deadline = ?
		WHERE id = ?
	`

	// Execute the SQL query with the updated deadline and task ID
	_, err := config.DB.Exec(query, updatedDeadline, id)
	if err != nil {
		return err
	}

	return nil
}

// UpdateTaskProgress updates the progress of a task by its ID.
func UpdateTaskProgress(taskID int, progress int) error {
	// Define the SQL query to update the task's progress by ID
	query := `
		UPDATE task
		SET progress = ?
		WHERE id = ?
	`

	// Execute the SQL query with the updated `progress` and `taskID`
	_, err := config.DB.Exec(query, progress, taskID)
	if err != nil {
		return err
	}

	return nil
}

// GetTotalSubTasksByTaskID mengembalikan total subtask berdasarkan ID task.
func GetTotalSubTasksByTaskID(taskID int) int {
	// Buat query SQL untuk menghitung total subtask berdasarkan ID task
	query := "SELECT COUNT(*) FROM subtask WHERE id_task = ?"

	// Eksekusi query dan ambil hasilnya
	var totalSubtasks int
	err := config.DB.QueryRow(query, taskID).Scan(&totalSubtasks)
	if err != nil {
		// Jika terjadi kesalahan, Anda bisa menambahkan penanganan kesalahan di sini,
		// misalnya mencatat log atau melakukan tindakan lain yang sesuai.
		// Namun, dalam kasus ini, fungsi akan mengembalikan 0 jika terjadi kesalahan.
		return 0
	}

	// Mengembalikan total subtask yang ditemukan
	return totalSubtasks
}

/// GetCompletedSubTasksByTaskID mengembalikan jumlah subtask dengan atribut complete == true
// berdasarkan ID task.
func GetCompletedSubTasksByTaskID(taskID int) int {
	// Buat query SQL untuk menghitung jumlah subtask yang memiliki complete == true
	query := "SELECT COUNT(*) FROM subtask WHERE id_task = ? AND complete = true"

	// Eksekusi query dan ambil hasilnya
	var completedSubtasks int
	err := config.DB.QueryRow(query, taskID).Scan(&completedSubtasks)
	if err != nil {
		// Jika terjadi kesalahan, Anda bisa menambahkan penanganan kesalahan di sini,
		// misalnya mencatat log atau melakukan tindakan lain yang sesuai.
		// Namun, dalam kasus ini, fungsi akan mengembalikan 0 jika terjadi kesalahan.
		return 0
	}

	// Mengembalikan jumlah subtask yang memiliki complete == true
	return completedSubtasks
}
