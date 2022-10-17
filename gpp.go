package main

import (
	"fmt"
	"image"
	"image/color"
	"log"

	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/ebitenutil"
	"github.com/hajimehoshi/ebiten/v2/vector"
)

var (
	appWidth      = 640
	appHeight     = 480
	emptyImage    = ebiten.NewImage(3, 3)
	emptySubImage = emptyImage.SubImage(image.Rect(1, 1, 2, 2)).(*ebiten.Image)
)

type app struct{}

func (g *app) Update() error {
	return nil
}

func drawPanel(s *ebiten.Image, x, y int) {
	var p vector.Path
	p.MoveTo(20, 20)
	p.LineTo(40, 40)
	//p.LineTo(50, 50)

	var v []ebiten.Vertex
	var i []uint16
	v, i = p.AppendVerticesAndIndicesForStroke(nil, nil, &vector.StrokeOptions{Width: 5, LineJoin: vector.LineJoinMiter})
	for i := range v {
		v[i].DstX = v[i].DstX + float32(x)
		v[i].DstY = v[i].DstY + float32(y)
		v[i].SrcX = 1
		v[i].SrcY = 1
		v[i].ColorR = 0xdb / float32(0xff)
		v[i].ColorG = 0x56 / float32(0xff)
		v[i].ColorB = 0x20 / float32(0xff)
	}
	s.DrawTriangles(v, i, emptySubImage, &ebiten.DrawTrianglesOptions{})
}

func (g *app) Draw(screen *ebiten.Image) {
	screen.Fill(color.White)
	ebitenutil.DebugPrint(screen, fmt.Sprintf("TPS: %0.2f\nFPS: %0.2f", ebiten.ActualTPS(), ebiten.ActualFPS()))
	drawPanel(screen, 40, 40)
}

func (g *app) Layout(outsideWidth, outsideHeight int) (int, int) {
	return appWidth, appHeight
}

func main() {
	emptyImage.Fill(color.White)
	ebiten.SetWindowSize(appWidth, appHeight)
	ebiten.SetWindowTitle("GPP")
	if err := ebiten.RunGame(&app{}); err != nil {
		log.Fatal(err)
	}
}
