package main

import (
	"fmt"
	"math/rand"
	"time"
)

var (
	xWrm, yWrm, fld rune
)

func (g *Game) initWorm() {
	rand.Seed(time.Now().UnixNano()) // рандомизация
	g.wormHead = 0
	g.wormTail = 0
	g.wormBody[g.wormHead] = Pos{19, 14}
	g.delay = 30
	g.lenght = 1
	g.wormEat = 2
	g.direct = 1
	g.count = g.delay
	g.returnError = nil

	g.initField()
}

func (g *Game) appleSet() {
	for {
		x := 1 + rand.Intn(wField-1)
		y := 1 + rand.Intn(hField-1)
		if (g.wormField[x][y]) == 32 { // space
			g.wormField[x][y] = rune(48 + rand.Intn(10))
			if g.wormField[x][y] == 48 { //0
				g.wormField[x][y] = 64 //@
			}
			break
		}
	}

}
func (g *Game) initField() {
	for j := 1; j < hField-1; j++ {
		for i := 1; i < wField-1; i++ {
			g.wormField[i][j] = 32 //пробелы
		}
	}
	for i := 0; i < wField; i++ {
		g.wormField[i][0] = 88 // рамка XXXXXX
		g.wormField[i][hField-1] = 88
	}
	for j := 0; j < hField-1; j++ {
		g.wormField[0][j] = 88 // рамка XXXXXX
		g.wormField[wField-1][j] = 88
	}

	g.wormField[19][14] = 88

	for i := 0; i < 10; i++ {
		g.appleSet()
	}

}

func (g *Game) directRight() {
	g.direct++
	if g.direct > 3 {
		g.direct = 0
	}
}

func (g *Game) directLeft() {
	g.direct--
	if g.direct < 0 {
		g.direct = 3
	}
}

func (g *Game) getField(xx, yy rune) rune {
	return g.wormField[xx][yy]
}

func (g *Game) headTailMove(n rune) (rune, rune, rune) {
	x := g.wormBody[n].x
	y := g.wormBody[n].y
	n++
	if n >= maxLen {
		n = 0
	}
	return x, y, n
}
func (g *Game) stepWorm() error {

	if g.pause {

		return nil
	}

	g.printText = fmt.Sprintf("Длинна червяка %d delay %d", g.lenght, g.delay)

	xWrm, yWrm, g.wormHead = g.headTailMove(g.wormHead)

	switch g.direct {
	case 0:
		yWrm--
	case 1:
		xWrm++
	case 2:
		yWrm++
	case 3:
		xWrm--
	}
	fld = g.getField(xWrm, yWrm)
	if (fld == 88) || (fld == 64) {
		g.pause = true
		g.printText = fmt.Sprintf("!!! Game Over !!! длинна %d  delay %d", g.lenght, g.delay)
		g.returnError = errorExit
		return nil
	}
	if (fld > 48) && (fld < 58) {

		g.delay--
		if g.delay < 3 {
			g.pause = true
			g.printText = fmt.Sprintf("!!! Вы выиграли !!! длинна червяка %d", g.lenght)
			g.returnError = errorExit
			return nil
		}
		g.appleSet()
		g.wormEat = g.wormEat + fld - 48
	}

	g.wormBody[g.wormHead] = Pos{xWrm, yWrm}
	g.wormField[xWrm][yWrm] = 88

	if g.wormEat > 0 {
		g.wormEat--
		g.lenght++
		return nil
	}

	xWrm, yWrm, g.wormTail = g.headTailMove(g.wormTail)
	g.wormField[xWrm][yWrm] = 32

	return g.returnError
}
