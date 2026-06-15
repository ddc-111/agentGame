package agent

import "context"

type Agent struct {
	ID       string
	Name     string
	Persona  string
	Location string
}

type Message struct {
	Role    string `json:"role"`
	Content string `json:"content"`
}

func New(id, name, persona string) *Agent {
	return &Agent{
		ID:      id,
		Name:    name,
		Persona: persona,
	}
}

func (a *Agent) Chat(ctx context.Context, messages []Message) (string, error) {
	return "Response from " + a.Name, nil
}
