package main

import (
	"errors"
	"fmt"
	"image/color"
	"log"

	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/ebitenutil"
	"github.com/hajimehoshi/ebiten/v2/examples/resources/fonts"
	"github.com/hajimehoshi/ebiten/v2/text"
	"golang.org/x/image/font"
	"golang.org/x/image/font/opentype"
)

const (
	title  string = "Worm1"
	width  int    = 640
	height int    = 480
	wField int    = 40
	hField int    = 30
	maxLen rune   = 300
	dpi           = 72
)

var (
	//clGreen = color.RGBA{0x75, 0xf9, 0x4d, 0xff}
	pressStartFont font.Face
	fontSize       = 16
	clGreen        = color.RGBA{0x60, 0xa6, 0x65, 0xff}
	clPink         = color.RGBA{0xea, 0x36, 0x80, 0xff}
	clBrown        = color.RGBA{0x78, 0x43, 0x15, 0xff}
	errorExit      = errors.New("regular termination")
)

type Pos struct {
	x rune
	y rune
}
type Game struct {
	wormField   [wField][hField]rune
	wormBody    [maxLen]Pos
	wormHead    rune
	wormTail    rune
	wormEat     rune
	delay       rune
	count       rune
	lenght      rune
	direct      rune
	change      bool
	pause       bool
	debug       bool
	printText   string
	returnError error
}

func init() {
	tt1, err := opentype.Parse(fonts.PressStart2P_ttf)
	if err != nil {
		log.Fatal(err)
	}

	pressStartFont, err = opentype.NewFace(tt1, &opentype.FaceOptions{
		Size:    float64(fontSize),
		DPI:     dpi,
		Hinting: font.HintingFull,
	})
	if err != nil {
		log.Fatal(err)
	}

}

func (g *Game) Update() error {

	g.count--
	if g.count < 0 {
		g.count = g.delay
		return g.stepWorm()
	}

	if ebiten.IsKeyPressed(ebiten.KeyEscape) {
		if !g.change {
			if g.pause {
				return errorExit
			}
			g.printText = fmt.Sprintf("Выход, длинна червяка %d", g.lenght)
			g.pause = true
			g.returnError = errorExit
			g.change = true
			return nil
		}
		return nil
	}
	if ebiten.IsKeyPressed(ebiten.KeyArrowRight) {
		if !g.change {
			g.directRight()
			g.change = true
			return nil
		}
		return nil
	}
	if ebiten.IsKeyPressed(ebiten.KeyArrowLeft) {
		if !g.change {
			g.directLeft()
			g.change = true
			return nil
		}
		return nil
	}
	if ebiten.IsKeyPressed(ebiten.KeySpace) {
		if !g.change {
			if g.returnError == nil {
				g.pause = !g.pause
				if g.pause {
					g.printText = fmt.Sprintf("Пауза, длинна %d  delay %d", g.lenght, g.delay)
				}
			}
			g.change = true
			return nil
		}
		return nil
	}
	if ebiten.IsKeyPressed(ebiten.KeyD) {
		if !g.change {
			g.debug = !g.debug
			g.change = true
			return nil
		}
		return nil
	}
	g.change = false

	if !g.pause {
		return g.returnError
	} else {
		return nil
	}
}

func (g *Game) Draw(screen *ebiten.Image) {
	for j := 0; j < hField; j++ {
		for i := 0; i < wField; i++ {
			text.Draw(screen, string(g.wormField[i][j]), pressStartFont, i*fontSize, j*fontSize+fontSize, clGreen)
		}
	}

	op := &ebiten.DrawImageOptions{}
	op.GeoM.Translate(0, float64(hField*fontSize))
	op.ColorM.ScaleWithColor(clBrown)
	op.CompositeMode = ebiten.CompositeModeCopy
	text.DrawWithOptions(screen, fmt.Sprintf(g.printText), pressStartFont, op)

	if g.debug {
		ebitenutil.DebugPrint(screen, fmt.Sprintf("TPS: %0.2f\nFPS: %0.2f\nCount: %d\ndirect: %d\nchange: %v\nhead: %d\ntail: %d",
			ebiten.CurrentTPS(), ebiten.CurrentFPS(), g.count, g.direct, g.change, g.wormHead, g.wormTail))
	}

}

func (g *Game) Layout(outsideWidth, outsideHeight int) (int, int) {
	return width, height
}

func main() {
	ebiten.SetWindowSize(width, height)
	ebiten.SetWindowTitle(title)
	ggg := &Game{}
	ggg.initWorm()
	if err := ebiten.RunGame(ggg); err != nil {
		log.Fatal(err)
	}

}
