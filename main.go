package main

import (
	"fmt"
	"time"

	"golang.org/x/mobile/app"
	"golang.org/x/mobile/event/lifecycle"
	"golang.org/x/mobile/event/paint"
	"golang.org/x/mobile/event/size"
	"golang.org/x/mobile/event/touch"
	"golang.org/x/mobile/exp/app/debug"
	"golang.org/x/mobile/exp/gl/glutil"
	"golang.org/x/mobile/exp/sprite"
	"golang.org/x/mobile/exp/sprite/clock"
	"golang.org/x/mobile/exp/sprite/glsprite"
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
					switch e.Crosses(lifecycle.StageVisible) {
					case lifecycle.CrossOn:
						glctx, _ = e.DrawContext.(gl.Context)
						onStart(glctx)
						a.Send(paint.Event{})
					case lifecycle.CrossOff:
						onStop()
						glctx = nil
					}
				case size.Event:
					sz = e
				case paint.Event:
					if glctx == nil {
						continue
					}
					onDraw(glctx, sz)
					a.Publish()
					a.Send(paint.Event{}) // keep animating
				case touch.Event:
					if time.Now().After(t2) {
						t2 = time.Now().Add(time.Second * 2)
						ok = !ok
						//a.Send(paint.Event{})
						game.Press(true, e.Y/float32(sz.HeightPx), e.X)
					}

					if counter == 10 {
						_ = GetConfig("config.json")
						counter++
						myApp = a
						InChan = make(chan string)
						go StartIRCprocess(InChan)
						go RoutineWriter()
					}

				}
			}
		}
	})
}

var (
	myApp     app.App
	ok        = false
	counter   = 0
	t2        = time.Now().Add(time.Second * 2)
	InChan    chan string
	startTime = time.Now()
	images    *glutil.Images
	eng       sprite.Engine
	scene     *sprite.Node
	game      *Game
	fps       *debug.FPS
)

func onStart(glctx gl.Context) {
	images = glutil.NewImages(glctx)
	eng = glsprite.Engine(images)
	game = NewGame()
	scene = game.Scene(eng)
	fps = debug.NewFPS(images)
	debug.NewFPS(images)
}
func onStop() {
	eng.Release()
	images.Release()
	game = nil
	fps.Release()
}
func onDraw(glctx gl.Context, sz size.Event) {

	if ok {
		glctx.ClearColor(1, 1, 1, 1)
	} else {
		glctx.ClearColor(0, 0, 0, 1)
	}
	glctx.Clear(gl.COLOR_BUFFER_BIT)
	now := clock.Time(time.Since(startTime) * 60 / time.Second)
	eng.Render(scene, now, sz)
	fps.Draw(sz)

}
func RoutineWriter() {
	for {
		select {
		case msg := <-InChan:
			ok = !ok
			fmt.Println(msg)
			myApp.Send(paint.Event{})
		}
	}
}
