package main

import (
	"embed"
	"fmt"
	"github.com/blizzy78/ebitenui"
	"github.com/blizzy78/ebitenui/image"
	"github.com/blizzy78/ebitenui/widget"
	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/inpututil"
	"golang.org/x/image/font/basicfont"
	stdImage "image"
	"image/color"
	"image/png"
	"log"
	"math/rand"
	"os"
	"time"
)

//go:embed assets/*
var EmbeddedAssets embed.FS
var simpleGame Game
var textWidget *widget.Text
var button *widget.Button
var button2 *widget.Button

var counter = 0
var fillTime = false

const (
	GameWidth   = 700
	GameHeight  = 700
	PlayerSpeed = 10
)

type Sprite struct {
	pict    *ebiten.Image
	xloc    int
	yloc    int
	dX      int
	dY      int
	drawOps ebiten.DrawImageOptions
}

type Game struct {
	player     Sprite
	enemy      Sprite
	score      int
	drawOps    ebiten.DrawImageOptions
	enemyList  []Sprite
	enemyCount int
	AppUI      *ebitenui.UI
}

func (g *Game) Update() error {
	processPlayerInput(g)
	g.AppUI.Update()

	if fillTime == true {
		fillList(g)
		fillTime = false
	}
	return nil
}

func (g *Game) Draw(screen *ebiten.Image) {
	g.AppUI.Draw(screen)
	g.drawOps.GeoM.Reset()
	g.drawOps.GeoM.Scale(.5, .5)
	g.drawOps.GeoM.Translate(float64(g.player.xloc), float64(g.player.yloc))
	screen.DrawImage(g.player.pict, &g.drawOps)
	textWidget.SetLocation(stdImage.Rectangle{
		Min: stdImage.Point{
			X: 0,
			Y: 600},
		Max: stdImage.Point{
			X: 700,
			Y: 700,
		},
	})
	textWidget.Label = fmt.Sprintf("Score:%d", counter)
	button.SetLocation(stdImage.Rectangle{
		Min: stdImage.Point{
			X: 0,
			Y: -100},
		Max: stdImage.Point{
			X: 1000,
			Y: -100,
		},
	})
	button.Text().Label = ""

	button2.SetLocation(stdImage.Rectangle{
		Min: stdImage.Point{
			X: 0,
			Y: -100},
		Max: stdImage.Point{
			X: 1000,
			Y: -100,
		},
	})
	button2.Text().Label = ""
	// This collision detection is from jsantore firstGameDemo
	for num, enemy := range g.enemyList {
		enemyWidth, enemyHeight := enemy.pict.Size()
		playerWidth, playerHeight := g.player.pict.Size()
		if counter == 10 {
			button.SetLocation(stdImage.Rectangle{
				Min: stdImage.Point{
					X: 0,
					Y: 0},
				Max: stdImage.Point{
					X: 700,
					Y: 100,
				},
			})
			button.Text().Label = "I want to leave!"
			button.Text().SetLocation(stdImage.Rectangle{
				Min: stdImage.Point{
					X: 0,
					Y: 0},
				Max: stdImage.Point{
					X: 700,
					Y: 100,
				},
			})
			button2.SetLocation(stdImage.Rectangle{
				Min: stdImage.Point{
					X: 0,
					Y: 200},
				Max: stdImage.Point{
					X: 700,
					Y: 300,
				},
			})
			button2.Text().SetLocation(stdImage.Rectangle{
				Min: stdImage.Point{
					X: 0,
					Y: 200},
				Max: stdImage.Point{
					X: 700,
					Y: 300,
				},
			})
			button2.Text().Label = "I want to Play again!"
		}
		if g.player.xloc < enemy.xloc+enemyWidth && g.player.xloc+playerWidth > enemy.xloc-enemyWidth &&
			g.player.yloc < enemy.yloc+enemyHeight && g.player.yloc+playerHeight > enemy.yloc-enemyHeight {
			remove(g.enemyList, num)
			counter += 1
			//message := fmt.Sprintf("score: %d", counter)
			//textWidget.Label = message
		} else {
			g.drawOps.GeoM.Reset()
			g.drawOps.GeoM.Translate(float64(enemy.xloc), float64(enemy.yloc))
			screen.DrawImage(enemy.pict, &g.drawOps)

		}

	}

}

func (g Game) Layout(outsideWidth, outsideHeight int) (screenWidth, screenHeight int) {
	return GameWidth, GameHeight
}
func remove(s []Sprite, index int) []Sprite {
	return append(s[:index], s[index+1:]...)
}
func fillList(g *Game) {
	r := rand.New(rand.NewSource(time.Now().UnixNano()))
	for i := 0; i < 10; i++ {
		g.enemy = Sprite{
			pict: loadPNGImageFromEmbedded("smallhammer.png"),
			xloc: r.Intn(650),
			yloc: r.Intn(650),
			dX:   0,
			dY:   0,
		}
		g.enemyList[i] = g.enemy

	}
}
func main() {
	ebiten.SetWindowSize(GameWidth, GameHeight)
	ebiten.SetWindowTitle("Minimal Game")

	simpleGame := Game{AppUI: MakeUIWindow()}

	simpleGame.player = Sprite{
		pict: loadPNGImageFromEmbedded("f1-ship1-3.png"),
		xloc: 200,
		yloc: 300,
		dX:   0,
		dY:   0,
	}

	simpleGame.enemyList = make([]Sprite, 11)
	fillList(&simpleGame)
	simpleGame.enemy = Sprite{
		pict: loadPNGImageFromEmbedded("smallhammer.png"),
		xloc: 10000,
		yloc: 10000,
		dX:   0,
		dY:   0,
	}
	simpleGame.enemyList[10] = simpleGame.enemy
	if err := ebiten.RunGame(&simpleGame); err != nil {
		log.Fatal("Oh no! something terrible happened and the game crashed", err)
	}
	textInfo := widget.TextOptions{}.Text("score: 0", basicfont.Face7x13, color.White)

	textWidget = widget.NewText(textInfo)
}

func loadPNGImageFromEmbedded(name string) *ebiten.Image {
	pictNames, err := EmbeddedAssets.ReadDir("assets")
	if err != nil {
		log.Fatal("failed to read embedded dir ", pictNames, " ", err)
	}
	embeddedFile, err := EmbeddedAssets.Open("assets/" + name)
	if err != nil {
		log.Fatal("failed to load embedded image ", embeddedFile, err)
	}
	rawImage, err := png.Decode(embeddedFile)
	if err != nil {
		log.Fatal("failed to load embedded image ", name, err)
	}
	gameImage := ebiten.NewImageFromImage(rawImage)
	return gameImage
}

func processPlayerInput(theGame *Game) {
	if inpututil.IsKeyJustPressed(ebiten.KeyUp) {
		theGame.player.dY = -PlayerSpeed
	} else if inpututil.IsKeyJustPressed(ebiten.KeyDown) {
		theGame.player.dY = PlayerSpeed
	} else if inpututil.IsKeyJustReleased(ebiten.KeyUp) || inpututil.IsKeyJustReleased(ebiten.KeyDown) {
		theGame.player.dY = 0
	} else if inpututil.IsKeyJustPressed(ebiten.KeyLeft) {
		theGame.player.dX = -PlayerSpeed
	} else if inpututil.IsKeyJustPressed(ebiten.KeyRight) {
		theGame.player.dX = PlayerSpeed
	} else if inpututil.IsKeyJustReleased(ebiten.KeyLeft) || inpututil.IsKeyJustReleased(ebiten.KeyRight) {
		theGame.player.dX = 0
	}
	theGame.player.yloc += theGame.player.dY
	theGame.player.xloc += theGame.player.dX
	if theGame.player.yloc <= 0 {
		theGame.player.dY = 0
		theGame.player.yloc = 0
	} else if theGame.player.yloc > 675 {
		theGame.player.dY = 0
		theGame.player.yloc = 675
	}
	if theGame.player.xloc <= 0 {
		theGame.player.dX = 0
		theGame.player.xloc = 0
	} else if theGame.player.xloc > GameWidth-30 {
		theGame.player.dX = 0
		theGame.player.xloc = GameWidth - 30
	}
}

func MakeUIWindow() (GUIhandler *ebitenui.UI) {
	background := image.NewNineSliceColor(color.Gray16{})
	rootContainer := widget.NewContainer(
		widget.ContainerOpts.Layout(widget.NewGridLayout(
			widget.GridLayoutOpts.Columns(1),
			widget.GridLayoutOpts.Stretch([]bool{true}, []bool{false, true, false}),

			widget.GridLayoutOpts.Spacing(0, 20))),
		widget.ContainerOpts.BackgroundImage(background))
	textInfo := widget.TextOptions{}.Text("score: 0", basicfont.Face7x13, color.White)

	idle, err := loadImageNineSlice("button-idle.png", 20, 0)
	if err != nil {
		log.Fatalln(err)
	}
	hover, err := loadImageNineSlice("button-hover.png", 20, 0)
	if err != nil {
		log.Fatalln(err)
	}
	pressed, err := loadImageNineSlice("button-pressed.png", 20, 0)
	if err != nil {
		log.Fatalln(err)
	}
	disabled, err := loadImageNineSlice("button-disabled.png", 20, 0)
	if err != nil {
		log.Fatalln(err)
	}
	buttonImage := &widget.ButtonImage{
		Idle:     idle,
		Hover:    hover,
		Pressed:  pressed,
		Disabled: disabled,
	}

	button = widget.NewButton(
		// specify the images to use
		widget.ButtonOpts.Image(buttonImage),
		// specify the button's text, the font face, and the color
		widget.ButtonOpts.Text("I want to go :(", basicfont.Face7x13, &widget.ButtonTextColor{
			Idle: color.RGBA{0xdf, 0xf4, 0xff, 0xff},
		}),
		// specify that the button's text needs some padding for correct display
		widget.ButtonOpts.TextPadding(widget.Insets{
			Left:  30,
			Right: 30,
		}),
		// ... click handler, etc. ...
		widget.ButtonOpts.ClickedHandler(quit),
	)
	button2 = widget.NewButton(
		// specify the images to use
		widget.ButtonOpts.Image(buttonImage),
		// specify the button's text, the font face, and the color
		widget.ButtonOpts.Text("Let's Play again!", basicfont.Face7x13, &widget.ButtonTextColor{
			Idle: color.RGBA{0xdf, 0xf4, 0xff, 0xff},
		}),
		// specify that the button's text needs some padding for correct display

		// ... click handler, etc. ...
		widget.ButtonOpts.ClickedHandler(playAgain),
	)

	rootContainer.AddChild(button)
	rootContainer.AddChild(button2)
	textWidget = widget.NewText(textInfo)
	rootContainer.AddChild(textWidget)
	GUIhandler = &ebitenui.UI{Container: rootContainer}

	return GUIhandler
}

func loadImageNineSlice(path string, centerWidth int, centerHeight int) (*image.NineSlice, error) {
	i := loadPNGImageFromEmbedded(path)

	w, h := i.Size()
	return image.NewNineSlice(i,
			[3]int{(w - centerWidth) / 2, centerWidth, w - (w-centerWidth)/2 - centerWidth},
			[3]int{(h - centerHeight) / 2, centerHeight, h - (h-centerHeight)/2 - centerHeight}),
		nil
}

func playAgain(args *widget.ButtonClickedEventArgs) {

	fillTime = true
	counter = 0
}
func quit(args *widget.ButtonClickedEventArgs) {
	os.Exit(3)
}
