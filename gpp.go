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
	defWidth      = 200
	defHeight     = 300
	panelMargin   = 5
	panelHeight   = 45
	emptyImage    = ebiten.NewImage(3, 3)
	emptySubImage = emptyImage.SubImage(image.Rect(1, 1, 2, 2)).(*ebiten.Image)
)

type app struct{}

func (g *app) Update() error {
	return nil
}

func drawRect(s *ebiten.Image, x, y, w, h float32, r, g, b float32) {
	var p vector.Path
	p.MoveTo(x, y)
	p.LineTo(w, y)
	p.LineTo(w, h)
	p.LineTo(x, h)
	p.LineTo(x, y)

	var v []ebiten.Vertex
	var i []uint16
	v, i = p.AppendVerticesAndIndicesForStroke(nil, nil, &vector.StrokeOptions{Width: 1, LineJoin: vector.LineJoinMiter})
	for i := range v {
		v[i].ColorR = r
		v[i].ColorG = g
		v[i].ColorB = b
	}
	s.DrawTriangles(v, i, emptySubImage, &ebiten.DrawTrianglesOptions{})
}

func drawPanel(s *ebiten.Image, no int) {
	w, _ := s.Size()
	drawRect(s, float32(panelMargin), float32(no*panelMargin), float32(w-panelMargin), float32(panelHeight), 0xdb/float32(0xff), 0x56/float32(0xff), 0x20/float32(0xff))
}

func (g *app) Draw(screen *ebiten.Image) {
	screen.Fill(color.Gray16{0x8000})
	ebitenutil.DebugPrint(screen, fmt.Sprintf("TPS: %0.2f\nFPS: %0.2f", ebiten.ActualTPS(), ebiten.ActualFPS()))
	drawPanel(screen, 1)
}

func (g *app) Layout(outsideWidth, outsideHeight int) (int, int) {
	s := ebiten.DeviceScaleFactor()
	return int(float64(outsideWidth) * s), int(float64(outsideHeight) * s)
}

func main() {
	emptyImage.Fill(color.White)
	ebiten.DeviceScaleFactor()
	ebiten.SetWindowSize(defWidth, defHeight)
	ebiten.SetWindowResizingMode(ebiten.WindowResizingModeEnabled)
	ebiten.SetWindowTitle("GPP")
	if err := ebiten.RunGame(&app{}); err != nil {
		log.Fatal(err)
	}
}
