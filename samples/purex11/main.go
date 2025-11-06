package main

import (
	"fmt"
	"image"
	"image/color"
	"image/draw"
	"log"
	"os"

	"github.com/BurntSushi/freetype-go/freetype"
	"github.com/BurntSushi/freetype-go/freetype/truetype"
	"github.com/BurntSushi/xgb/xproto"
	"github.com/BurntSushi/xgbutil"
	"github.com/BurntSushi/xgbutil/icccm"
	"github.com/BurntSushi/xgbutil/xevent"
	"github.com/BurntSushi/xgbutil/xgraphics"
	"github.com/BurntSushi/xgbutil/xwindow"
)

func main() {
	// Connect to the X server.
	xu, err := xgbutil.NewConn()
	if err != nil {
		log.Fatal(err)
	}

	// Create a new window.
	win, err := xwindow.Generate(xu)
	if err != nil {
		log.Fatal(err)
	}
	win.Create(xu.RootWin(), 0, 0, 400, 200, xproto.CwBackPixel, xu.Screen().WhitePixel)

	// Create a pixmap.
	pix, err := xproto.NewPixmapId(xu.Conn())
	if err != nil {
		log.Fatal(err)
	}
	err = xproto.CreatePixmapChecked(xu.Conn(), xu.Screen().RootDepth, pix, xproto.Drawable(win.Id), 400, 200).Check()
	if err != nil {
		log.Fatal(err)
	}

	// Create an xgraphics.Image from the pixmap.
	img := xgraphics.New(xu, image.Rect(0, 0, 400, 200))
	img.Pixmap = pix

		// Fill the image with white color
	draw.Draw(img, img.Bounds(), image.White, image.Point{}, draw.Src)

	// Make the window visible and set a title.
	icccm.WmNameSet(xu, win.Id, fmt.Sprintf("Text Example (id: %d)", win.Id))

	// Load the font file. You may need to change the path.
	fontBytes, err := os.ReadFile("/usr/share/fonts/truetype/dejavu/DejaVuSans-Bold.ttf")
	if err != nil {
		log.Fatalf("Could not read font file: %v", err)
	}
	font, err := truetype.Parse(fontBytes)
	if err != nil {
		log.Fatalf("Could not parse font: %v", err)
	}

	// Create a freetype context.
	ft := freetype.NewContext()
	ft.SetDPI(72)
	ft.SetFont(font)
	ft.SetFontSize(24)
	ft.SetClip(img.Bounds())
	ft.SetDst(img)
	ft.SetSrc(image.NewUniform(color.Black))

	// Draw the text "hello world x11".
	pt := freetype.Pt(10, 24)
	_, err = ft.DrawString("hello world x11", pt)
	if err != nil {
		log.Fatalf("Error drawing string: %v", err)
	}

	// Copy the pixmap to the window.
	img.XDraw()
	img.XPaint(win.Id)

	// Map the window.
	win.Map()

	// Handle events. The program will exit when the window is closed.
	xevent.Main(xu)
}
