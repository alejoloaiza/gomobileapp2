package main

import (
	"time"

	"golang.org/x/mobile/app"
	"golang.org/x/mobile/event/lifecycle"
	"golang.org/x/mobile/event/paint"
	"golang.org/x/mobile/event/size"
	"golang.org/x/mobile/event/touch"
	"golang.org/x/mobile/gl"
)

func main() {
	app.Main(func(a app.App) {
		var glctx gl.Context

		sz := size.Event{}
		for {
			select {
			case e := <-a.Events():
				switch e := a.Filter(e).(type) {
				case lifecycle.Event:
					glctx, _ = e.DrawContext.(gl.Context)
				case size.Event:
					sz = e
				case paint.Event:
					if glctx == nil {
						continue
					}
					onDraw(glctx, sz)
					a.Publish()
				case touch.Event:
					if time.Now().After(t2) {
						t2 = time.Now().Add(time.Second * 5)
						ok = !ok
						a.Send(paint.Event{})
					}
					if counter == 0 {
						_ = GetConfig("config.json")
						counter++
						go StartIRCprocess(InChan)
					}

				}
			}
		}
	})
}

var (
	ok      = false
	counter = 0
	t2      = time.Now().Add(time.Second * 2)
	InChan  chan string
)

func onDraw(glctx gl.Context, sz size.Event) {
	select {
	case msg := <-InChan:
		ok = !ok
	default:

	}

	if ok {
		glctx.ClearColor(1, 1, 1, 1)
	} else {
		glctx.ClearColor(0, 0, 0, 1)
	}

	glctx.Clear(gl.COLOR_BUFFER_BIT)
}
