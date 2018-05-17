// Copyright 2015 The Go Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

// +build darwin linux

package main

import (
	"image"
	"log"

	_ "image/png"

	"golang.org/x/mobile/asset"
	"golang.org/x/mobile/exp/f32"
	"golang.org/x/mobile/exp/sprite"
	"golang.org/x/mobile/exp/sprite/clock"
)

const (
	tileWidth, tileHeight = 16, 16 // width and height of each tile
	tilesX, tilesY        = 16, 16 // number of horizontal tiles

	gopherTile = 4 // which tile the gopher is standing on (0-indexed)

	initScrollV = 1     // initial scroll velocity
	scrollA     = 0.001 // scroll accelleration
	gravity     = 0.1   // gravity
	jumpV       = -5    // jump velocity
	flapV       = -1.5  // flap velocity

	deadScrollA         = -0.01 // scroll deceleration after the gopher dies
	deadTimeBeforeReset = 240   // how long to wait before restarting the game

	groundChangeProb = 5 // 1/probability of ground height change
	groundWobbleProb = 3 // 1/probability of minor ground height change
	groundMin        = tileHeight * (tilesY - 2*tilesY/5)
	groundMax        = tileHeight * tilesY
	initGroundY      = tileHeight * (tilesY - 1)

	climbGrace = tileHeight / 3 // gopher won't die if it hits a cliff this high
)

type Game struct {
	gopher struct {
		y        float32    // y-offset
		v        float32    // velocity
		atRest   bool       // is the gopher on the ground?
		flapped  bool       // has the gopher flapped since it became airborne?
		dead     bool       // is the gopher dead?
		deadTime clock.Time // when the gopher died
		size     float32
	}
	scroll struct {
		x float32 // x-offset
		v float32 // velocity
	}
	groundY   [tilesX + 3]float32 // ground y-offsets
	groundTex [tilesX + 3]int     // ground texture
	lastCalc  clock.Time          // when we last calculated a frame

}

func NewGame() *Game {
	var g Game
	g.reset()
	return &g
}

func (g *Game) reset() {
	g.gopher.y = 120
	g.gopher.v = 0
	g.scroll.x = 0
	g.scroll.v = initScrollV
	for i := range g.groundY {
		g.groundY[i] = initGroundY
		g.groundTex[i] = 8
	}
	g.gopher.atRest = false
	g.gopher.flapped = false
	g.gopher.dead = false
	g.gopher.deadTime = 0
	g.gopher.size = 20
}

func (g *Game) Scene(eng sprite.Engine) *sprite.Node {
	texs := loadTextures(eng)

	scene := &sprite.Node{}
	eng.Register(scene)
	eng.SetTransform(scene, f32.Affine{
		{1, 0, 0},
		{0, 1, 0},
	})

	newNode := func(fn arrangerFunc) {
		n := &sprite.Node{Arranger: arrangerFunc(fn)}
		eng.Register(n)
		scene.AppendChild(n)
	}

	// The gopher.
	newNode(func(eng sprite.Engine, n *sprite.Node, t clock.Time) {
		a := f32.Affine{
			{tileWidth * 2, 0, tileWidth*(gopherTile-1) + tileWidth/8},
			{0, tileHeight * 2, g.gopher.y - tileHeight + tileHeight/4},
		}
		var x int

		x = frame(t, 8, texGopherFlap1, texGopherFlap2)
		a.Scale(&a, 1+g.gopher.size/20, 1+g.gopher.size/20)
		eng.SetSubTex(n, texs[x])
		eng.SetTransform(n, a)
	})

	return scene
}

type arrangerFunc func(e sprite.Engine, n *sprite.Node, t clock.Time)

func (a arrangerFunc) Arrange(e sprite.Engine, n *sprite.Node, t clock.Time) { a(e, n, t) }

// frame returns the frame for the given time t
// when each frame is displayed for duration d.
func frame(t, d clock.Time, frames ...int) int {
	total := int(d) * len(frames)
	return frames[(int(t)%total)/int(d)]
}

const (
	texGopherRun1 = iota
	texGopherRun2
	texGopherFlap1
	texGopherFlap2
	texGopherDead1
	texGopherDead2
)

func loadTextures(eng sprite.Engine) []sprite.SubTex {
	a, err := asset.Open("sprite.png")
	if err != nil {
		log.Fatal(err)
	}
	defer a.Close()

	m, _, err := image.Decode(a)
	if err != nil {
		log.Fatal(err)
	}
	t, err := eng.LoadTexture(m)
	if err != nil {
		log.Fatal(err)
	}

	const n = 128
	// The +1's and -1's in the rectangles below are to prevent colors from
	// adjacent textures leaking into a given texture.
	// See: http://stackoverflow.com/questions/19611745/opengl-black-lines-in-between-tiles
	return []sprite.SubTex{
		texGopherRun1:  sprite.SubTex{t, image.Rect(n*0+1, 0, n*1-1, n)},
		texGopherRun2:  sprite.SubTex{t, image.Rect(n*1+1, 0, n*2-1, n)},
		texGopherFlap1: sprite.SubTex{t, image.Rect(n*2+1, 0, n*3-1, n)},
		texGopherFlap2: sprite.SubTex{t, image.Rect(n*3+1, 0, n*4-1, n)},
		texGopherDead1: sprite.SubTex{t, image.Rect(n*4+1, 0, n*5-1, n)},
		texGopherDead2: sprite.SubTex{t, image.Rect(n*5+1, 0, n*6-1, n)},
	}
}

func (g *Game) Press(down bool, y float32, x float32) {

	if down {
		//	g.gopher.y = y
		g.gopher.size = g.gopher.size + 1

	}
}
