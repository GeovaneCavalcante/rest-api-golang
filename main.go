package main

import (
	"fmt"
	"net/http"
	"sync/atomic"

	"github.com/gin-gonic/gin"
)
type Task struct {
    ID          int    `json:"id"`
    Description string `json:"description"`
}

var tasks = []Task{}
var atomicID int32

func main() {
    r := gin.Default()

    // Get all tasks
    r.GET("/tasks", func(c *gin.Context) {
        c.JSON(http.StatusOK, tasks)
    })

    // Get a specific task by ID
    r.GET("/tasks/:id", func(c *gin.Context) {
        id := c.Param("id")
        for _, task := range tasks {
            if id == fmt.Sprintf("%d", task.ID) {
                c.JSON(http.StatusOK, task)
                return
            }
        }
        c.JSON(http.StatusNotFound, gin.H{"error": "Task not found"})
    })

    // Create a new task
    r.POST("/tasks", func(c *gin.Context) {
        var newTask Task
        if err := c.ShouldBindJSON(&newTask); err != nil {
            c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
            return
        }

        newTask.ID = int(atomic.AddInt32(&atomicID, 1))  // Assign a new unique ID
        tasks = append(tasks, newTask)
        c.JSON(http.StatusOK, newTask)
    })

    // Update a task by ID
    r.PUT("/tasks/:id", func(c *gin.Context) {
        id := c.Param("id")
        var updatedTask Task
        if err := c.ShouldBindJSON(&updatedTask); err != nil {
            c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
            return
        }
        
        for i, task := range tasks {
            if id == fmt.Sprintf("%d", task.ID) {
                tasks[i].Description = updatedTask.Description
                c.JSON(http.StatusOK, tasks[i])
                return
            }
        }
        c.JSON(http.StatusNotFound, gin.H{"error": "Task not found"})
    })

    // Delete a task by ID
    r.DELETE("/tasks/:id", func(c *gin.Context) {
        id := c.Param("id")
        for i, task := range tasks {
            if id == fmt.Sprintf("%d", task.ID) {
                tasks = append(tasks[:i], tasks[i+1:]...)
                c.JSON(http.StatusOK, gin.H{"message": "Task deleted"})
                return
            }
        }
        c.JSON(http.StatusNotFound, gin.H{"error": "Task not found"})
    })

    r.Run(":8080")  // Start the server on port 8080
}
