package bot_father

type Father struct {
	BotCreator
}

type BotCreator interface {
	AddBot(token string, name string) error
}

func (f *Father) AddBot(token string, name string) error {
	if err := f.AddBot(token, name); err != nil {
		return err
	}

	return nil
}
