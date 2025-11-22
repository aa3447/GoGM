package player

type GM struct{
	Name string
	Description string
}

func NewGM(name string, description string) *GM{
	return &GM{
		Name: name,
		Description: description,
	}
}