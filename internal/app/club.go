package app

type App struct {
	handler Handler
}

type Handler interface {
	ProcessEvents() error
	EndDay() error
}

func NewApp(handler Handler) *App {
	return &App{
		handler: handler,
	}
}

func (a *App) Run() error {
	if err := a.handler.ProcessEvents(); err != nil {
		return err
	}
	err := a.handler.EndDay()
	if err != nil {
		return err
	}
	return nil
}
