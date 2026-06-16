package events

import "github.com/gdamore/tcell/v2"

func Start(screen tcell.Screen) <-chan tcell.Event {
	ch := make(chan tcell.Event, 16)
	go func() {
		for {
			event := screen.PollEvent()
			if event == nil {
				close(ch)
				return
			}
			ch <- event
		}
	}()
	return ch
}
