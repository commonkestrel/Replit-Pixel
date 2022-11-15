package main

import (
    "os"
    "image"

    _ "image/png"

    "github.com/faiface/pixel"
    "github.com/faiface/pixel/imdraw"
    "github.com/faiface/pixel/pixelgl"
    "golang.org/x/image/colornames"
)

const SCREENX, SCREENY = 960, 540

var (
    win     *pixelgl.Window
    imd     *imdraw.IMDraw
)

// used for loading icons and sprites
func LoadPicture(path string) (pixel.Picture, error) {
    // loads and decodes PNG
    file, err := os.Open(path)
    if err != nil {
        panic(err)
    }
    defer file.Close()
    img, _, err := image.Decode(file)
    if err != nil {
        panic(err)
    }
    // converts to Pixel picture
    return pixel.PictureDataFromImage(img), nil
}

func run() {
    icon, err := LoadPicture("icon.png")
    if err != nil {
        panic(err)
    }

    cfg := pixelgl.WindowConfig{
        Title:  "Go Pixel",
        Bounds: pixel.R(0, 0, SCREENX, SCREENY),
        Icon:  []pixel.Picture{icon},
        VSync: true,
    }
    win, err = pixelgl.NewWindow(cfg)
    if err != nil {
        panic(err)
    }

    imd = imdraw.New(nil)

    for !win.Closed() {
        imd.Clear()

        // game loop here

        win.Clear(colornames.Black)
        imd.Draw(win)
        win.Update()
    }
}

func main() {
    pixelgl.Run(run)
}