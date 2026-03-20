package command

type Command struct {
	Name string
	Args []string
}

func NewCommand(name string) Command {
	return Command{Name: name}
}
