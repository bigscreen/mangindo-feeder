package adapter

type Args map[string]interface{}

type Job struct {
	Queue   string
	Args    Args
	Handler string
}
