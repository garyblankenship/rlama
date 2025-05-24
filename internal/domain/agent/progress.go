package agent

import (
	"fmt"
	"sync"
	"time"
)

// ProgressDisplay manages simple progress display for agent tasks
type ProgressDisplay struct {
	mu      sync.Mutex
	tasks   []*TaskProgress // Keep order
	taskMap map[string]*TaskProgress
	verbose bool
	started bool
}

// TaskProgress represents the progress state of a single task
type TaskProgress struct {
	ID          string
	Description string
	Status      TaskProgressStatus
	StartTime   time.Time
	EndTime     time.Time
	Result      string
	Error       error
}

// TaskProgressStatus represents the progress status
type TaskProgressStatus string

const (
	TaskProgressPending   TaskProgressStatus = "pending"
	TaskProgressRunning   TaskProgressStatus = "running"
	TaskProgressCompleted TaskProgressStatus = "completed"
	TaskProgressFailed    TaskProgressStatus = "failed"
)

// NewProgressDisplay creates a new progress display
func NewProgressDisplay(verbose bool) *ProgressDisplay {
	return &ProgressDisplay{
		tasks:   make([]*TaskProgress, 0),
		taskMap: make(map[string]*TaskProgress),
		verbose: verbose,
	}
}

// StartTask starts tracking a new task
func (p *ProgressDisplay) StartTask(id, description string) {
	p.mu.Lock()
	defer p.mu.Unlock()

	if p.verbose {
		return // Don't show progress in verbose mode
	}

	if !p.started {
		fmt.Printf("\n🤖 Orchestration des tâches:\n")
		p.started = true
	}

	task := &TaskProgress{
		ID:          id,
		Description: description,
		Status:      TaskProgressPending,
		StartTime:   time.Now(),
	}

	p.tasks = append(p.tasks, task)
	p.taskMap[id] = task

	// Truncate description to fit nicely
	desc := description
	if len(desc) > 50 {
		desc = desc[:47] + "..."
	}

	fmt.Printf("⏳ %s\n", desc)
}

// UpdateTaskStatus updates the status of a task
func (p *ProgressDisplay) UpdateTaskStatus(id string, status TaskProgressStatus) {
	p.mu.Lock()
	defer p.mu.Unlock()

	if p.verbose {
		return
	}

	task, exists := p.taskMap[id]
	if !exists {
		return
	}

	task.Status = status
	if status == TaskProgressRunning {
		task.StartTime = time.Now()

		// Find task position and update display
		for i, t := range p.tasks {
			if t.ID == id {
				desc := task.Description
				if len(desc) > 50 {
					desc = desc[:47] + "..."
				}

				// Move cursor up to the task line and update
				fmt.Printf("\033[%dA", len(p.tasks)-i)
				fmt.Printf("\r⚡ %s", desc)
				fmt.Printf("\033[K")                   // Clear rest of line
				fmt.Printf("\033[%dB", len(p.tasks)-i) // Move back down
				break
			}
		}

		// Small pause to let user see the change
		time.Sleep(500 * time.Millisecond)
	}
}

// CompleteTask marks a task as completed with a result
func (p *ProgressDisplay) CompleteTask(id string, result string, err error) {
	p.mu.Lock()
	defer p.mu.Unlock()

	if p.verbose {
		return
	}

	task, exists := p.taskMap[id]
	if !exists {
		return
	}

	task.EndTime = time.Now()
	task.Result = result
	task.Error = err

	var icon string
	if err != nil {
		task.Status = TaskProgressFailed
		icon = "❌"
	} else {
		task.Status = TaskProgressCompleted
		icon = "✅"
	}

	// Find task position and update display
	for i, t := range p.tasks {
		if t.ID == id {
			desc := task.Description
			if len(desc) > 50 {
				desc = desc[:47] + "..."
			}

			duration := task.EndTime.Sub(task.StartTime)

			// Move cursor up to the task line and update
			fmt.Printf("\033[%dA", len(p.tasks)-i)
			fmt.Printf("\r%s %s (%.1fs)", icon, desc, duration.Seconds())
			fmt.Printf("\033[K")                   // Clear rest of line
			fmt.Printf("\033[%dB", len(p.tasks)-i) // Move back down
			break
		}
	}

	// Check if all tasks are completed
	allCompleted := true
	for _, t := range p.tasks {
		if t.Status != TaskProgressCompleted && t.Status != TaskProgressFailed {
			allCompleted = false
			break
		}
	}

	if allCompleted {
		p.renderFinal()
	}

	// Pause to let user see the completion
	time.Sleep(800 * time.Millisecond)
}

// renderFinal displays a simple final summary
func (p *ProgressDisplay) renderFinal() {
	if p.verbose {
		return
	}

	completed := 0
	failed := 0
	totalDuration := time.Duration(0)

	for _, task := range p.tasks {
		if task.Status == TaskProgressCompleted {
			completed++
			totalDuration += task.EndTime.Sub(task.StartTime)
		} else if task.Status == TaskProgressFailed {
			failed++
		}
	}

	fmt.Printf("\n🎯 %d tâches terminées", completed)
	if failed > 0 {
		fmt.Printf(" (%d échecs)", failed)
	}
	fmt.Printf(" en %.1fs\n", totalDuration.Seconds())
}

// DebugPrint prints debug information conditionally
func (p *ProgressDisplay) DebugPrint(format string, args ...interface{}) {
	if p.verbose {
		fmt.Printf("\n🔍 DEBUG: "+format+"\n", args...)
	}
}

// Message prints a regular message
func (p *ProgressDisplay) Message(format string, args ...interface{}) {
	fmt.Printf(format+"\n", args...)
}
