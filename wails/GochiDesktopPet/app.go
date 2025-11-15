package main
import "context"
type App struct {
    Tasks     []string
    Completed int
}

func (a *App) Startup(ctx context.Context) {
    // Called when the app starts. You can initialize tasks here.
    a.Tasks = []string{}
    a.Completed = 0
}


func NewApp() *App {
    return &App{
        Tasks: []string{},
    }
}

func (a *App) AddTask(task string) {
    a.Tasks = append(a.Tasks, task)
}

func (a *App) CompleteTask(index int) {
    if index >= 0 && index < len(a.Tasks) {
        a.Completed++
        a.Tasks = append(a.Tasks[:index], a.Tasks[index+1:]...)
    }
}

func (a *App) Mood() string {
    if a.Completed == 0 {
        return "Sad :("
    } else if a.Completed < len(a.Tasks)/2 {
        return "Neutral."
    } else {
        return "Happy!"
    }
}

