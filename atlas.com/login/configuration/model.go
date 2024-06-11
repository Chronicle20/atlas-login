package configuration

import "errors"

type Model struct {
	Data Data `json:"data"`
}

func (d *Model) FindTask(task string) (Task, error) {
	for _, v := range d.Data.Attributes.Tasks {
		if v.Type == task {
			return v, nil
		}
	}
	return Task{}, errors.New("task not found")
}

// Data contains the main data configuration.
type Data struct {
	Type       string     `json:"type"`
	Attributes Attributes `json:"attributes"`
}

// Attributes contain all settings under attributes key.
type Attributes struct {
	Tasks   []Task   `json:"tasks"`
	Servers []Server `json:"servers"`
}

// Task represents a task in the configuration.
type Task struct {
	Type       string         `json:"type"`
	Attributes TaskAttributes `json:"attributes"`
}

// TaskAttributes contains settings specific to a task.
type TaskAttributes struct {
	Interval int64 `json:"interval"`
	Duration int64 `json:"duration"`
}

// Server represents a server in the configuration.
type Server struct {
	Tenant   string    `json:"tenant"`
	Region   string    `json:"region"`
	Port     string    `json:"port"`
	Version  Version   `json:"version"`
	Handlers []Handler `json:"handlers"`
	Writers  []Writer  `json:"writers"`
}

// Version represents a server version.
type Version struct {
	Major string `json:"major"`
	Minor string `json:"minor"`
}

// Handler represents a server handler.
type Handler struct {
	OpCode    string `json:"opCode"`
	Validator string `json:"validator"`
	Handler   string `json:"handler"`
}

// Writer represents a server writer.
type Writer struct {
	OpCode string `json:"opCode"`
	Writer string `json:"writer"`
}
