# Go Pixel

Pixel is a hand-crafted 2D game library in Go.
This is an intro tutorial that will cover the basics and make a simple version of [Conway's Game of Life](https://en.wikipedia.org/wiki/Conway%27s_Game_of_Life)
###### Credit to the [Pixel Wiki](https://github.com/faiface/pixel/wiki) for some explanations, code, and quotes. Please explore the wiki before or after this guide, there's lots covered there that isn't in here.

## Initialization
First we need to get some imports out of the way, such as [Pixel](https://pkg.go.dev/github.com/faiface/pixel?utm_source=godoc) and [PixelGL](https://pkg.go.dev/github.com/faiface/pixel/pixelgl?utm_source=godoc)
```go
package main

import (
    "github.com/faiface/pixel"
    "github.com/faiface/pixel/pixelgl"
)
```

In order to allow PixelGL to use the main thread for all the windowing and graphics code, we need to make this call from our main function:

```go
func main() {
    pixelgl.Run(run)
}
```

`run` will be seperate function in which you will run all your code, basically a new `main` function.

## Creating a Window
Before we create a window we will need to need to configure it, which we can do using a [pixelgl.WindowConfig](https://godoc.org/github.com/faiface/pixel/pixelgl#WindowConfig) struct, which lets us set up the parameters in a convinient fashion.
```go
func run() {
    cfg := pixelgl.WindowConfig{
        Title: "Game of Life",
        Bounds: pixel.R(0, 0, 960, 540),
        Vsync: true,
    }
}
```

Let's break this down. We create a pixelgl.WindowConfig struct value and assign it to the cfg variable for later use. We only need to change three fields in the pixelgl.WindowConfig struct, the rest uses sensible defaults. The first field is the window title. The second field is the bounds of the window, set to a rectangle in Pixel. The `pixel.R` function creates a new rectangle, with the first two arguments being the coordinates of the lower-left corner of the rectangle, and the last two being the coordinates of the upper-right corner. In this we are creating a window with the size 960x540 pixels. The third field turns on VSync, which matches the window framerate to the monitor framerate. We do this in order to not use 100% of the CPU by updating the window as fast as we can. 

Speaking of the window, we still need to create that. We can use the function [pixelgl.NewWindow](https://godoc.org/github.com/faiface/pixel/pixelgl#NewWindow) to create a new window.
```go
func run() {
    cfg := pixelgl.WindowConfig{
        Title: "Game of Life",
        Bounds: pixel.R(0, 0, 960, 540),
        Vsync: true,
    }
    win, err := pixelgl.NewWindow(cfg)
    if err != nil {
        panic(err)
    }
}
```
We can see here that the [pixelgl.NewWindow](https://godoc.org/github.com/faiface/pixel/pixelgl#NewWindow) takes the parameters we just created as input and returns a window. It also returns a potential error, such as no graphics driver existing, which is handeled by panicking if there is an error.

Next, we need to create a main loop, to keep the window up and running until a user clicks the close button.
```go
func run() {
    cfg := pixelgl.WindowConfig{
        Title: "Game of Life",
        Bounds: pixel.R(0, 0, 960, 540),
        Vsync: true,
    }
    win, err := pixelgl.NewWindow(cfg)
    if err != nil {
        panic(err)
    }

    for !win.Closed() {
        win.Update()
    }
}
```

When our run function finishes, the whole program exists, so we need to make sure, that run is running until we want to actually exit our program.

Here we run a loop that finishes when a user closes our window. We need to call win.Update periodically. The function win.Update fetches new events (key presses, mouse moves and clicks, etc.) and redraws the window.

The final code for creating a window looks like this:
```go
package main

import (
    "github.com/faiface/pixel"
    "github.com/faiface/pixel/pixelgl"
)

func run() {
    cfg := pixelgl.WindowConfig{
        Title: "Game of Life",
        Bounds: pixel.R(0, 0, 960, 540),
        VSync: true,
    }
    win, err := pixelgl.NewWindow(cfg)
    if err != nil {
        panic(err)
    }

    for !win.Closed() {
        win.Update()
    }
}

func main() {
    pixelgl.Run(run)
}
```

This should just create a blank window like so:
![empty](https://raw.githubusercontent.com/commonkestrel/Replit-Pixel/main/screenshots/blank.png)

## Drawing
A blank window is only exciting for so long. What if we want a different color background?
There's a package `"golang.org/x/image/colornames"` which provides helpful color bindings, which will help us color our background.
We can set the background when we clear the window, using the `win.Clear` function, like so:
```go
for !win.Closed() {
    win.Clear(colornames.Skyblue)
    win.Update()
}
```

Now you should see something like this:
![empty](https://raw.githubusercontent.com/commonkestrel/Replit-Pixel/main/screenshots/skyblue.png)

But wait! I thought we were going to draw things, not just change the background color. Don't worry, drawing primitive shapes is made easy using the [IMDraw](https://pkg.go.dev/github.com/faiface/pixel/imdraw?utm_source=godoc) package. First we have to add it to our import statement:
```go
import (
	"github.com/faiface/pixel"
	"github.com/faiface/pixel/imdraw"
	"github.com/faiface/pixel/pixelgl"
	"golang.org/x/image/colornames"
)
```

The `imdraw` package exports a type called [IMDraw](https://godoc.org/github.com/faiface/pixel/imdraw#IMDraw), which can be created like so:
```
    win, err := pixelgl.NewWindow(cfg)
	if err != nil {
		panic(err)
	}

	imd := imdraw.New(nil)
```

Now that we have our `IMDraw`, we can use it to make some basic shapes, for example a triangle.
```go
    imd := imdraw.New(nil)
    for !win.Closed() {
        imd.Clear()

        imd.Color = pixel.RGB(1, 0, 0)
        imd.Push(pixel.V(200, 100))
        imd.Color = pixel.RGB(0, 1, 0)
        imd.Push(pixel.V(800, 100))
        imd.Color = pixel.RGB(0, 0, 1)
        imd.Push(pixel.V(500, 700))
        imd.Polygon(0)

        win.Clear(colornames.Black)
        imd.Draw(win)
        win.Update()
    }
```

Here we get this:
![empty](https://raw.githubusercontent.com/commonkestrel/Replit-Pixel/main/screenshots/triangle.png)

It works! But hold on, what exactly is going on here?

`IMDraw` is basically a pretty convenient state machine. There are three kinds things we can do with it.

Fields such as `imd.Color`, `imd.EndShape` or `imd.Precision` are properties. They are easily settable using `=`assignment. All of these properties affect points before they are pushed.

The second kind is the [imd.Push](https://godoc.org/github.com/faiface/pixel/imdraw#IMDraw.Push) method, which takes variable number of arguments: vectors representing the positions of points. This method pushes points to the IMDraw. The points take all of the currently set properties with themselves (remembers them). Changing the properties later does not affect any previously pushed points.

The last kind of methods is shape finalizers. These methods include [imd.Line](https://godoc.org/github.com/faiface/pixel/imdraw#IMDraw.Line), [imd.Polygon](https://godoc.org/github.com/faiface/pixel/imdraw#IMDraw.Polygon), [imd.Rectangle](https://godoc.org/github.com/faiface/pixel/imdraw#IMDraw.Rectangle) or [imd.Circle](https://godoc.org/github.com/faiface/pixel/imdraw#IMDraw.Circle). Each of these methods collects all of the pushed points and draws a shape according to them. For example, the `imd.Line` method draws a line between the pushed points and `imd.Circle` draws a circle around each of the pushed points. These methods take additional arguments further describing the specific shape.

The drawn shapes are then remembered inside the `IMDraw`, so when we call `imd.Draw(win)`, the `IMDraw` draws all of the drawn shapes to the window.

So, looking back at the triangle code, it's quite easy. We set a color and push a point for each of the three points in the triangle. Finally, we draw a polygon with the 0 thickness, which means a filled polygon.

## pixel.RGBA

The only part we don't really understand about the triangle code yet is the [pixel.RGB](https://godoc.org/github.com/faiface/pixel#RGB) function. As you probably already know, the standard ["image/color"](https://godoc.org/image/color) package defines the [color.Color](https://godoc.org/image/color#Color) interface. It's possible to create our own color formats just by implementing this interface.

Pixel does that and implements it's own color format (which you may or might not use, it's up to you): [pixel.RGBA](https://godoc.org/github.com/faiface/pixel#RGBA). It's an alpha-premultiplied RGBA color with `float64` components in range [0, 1] and additional useful methods (e.g. multiplying two colors).

There are two constructors. One is [pixel.RGB](https://godoc.org/github.com/faiface/pixel#RGB), which creates a fully opaque RGB color. The other one is [pixel.Alpha](https://godoc.org/github.com/faiface/pixel#Alpha) constructor which creates a transparent white color. Creating a transparent RGBA color is achieved by creating a opaque RGB color and multiplying it by a transparent white.

## Icon
Let's say we want to add an icon to spruce up our little window, make it less dull. How would we go about that? Well, Pixel has it's own system for sprites and icons, which is implemented through the `pixel.Picture` struct (Sprites work a little bit differently, but we won't get into that here. For info on sprites see the wiki [here](https://github.com/faiface/pixel/wiki/Drawing-a-Sprite)). Pictures can be loaded using this handy function shamelessly stolen from the [Pixel Wiki](https://github.com/faiface/pixel/wiki). 
First, we need to import the "image" package and we also need to "underscore import" the "image/png" package to load PNG files. We also need to import the "os" package to load files from the filesystem.
```go
import (
	"image"
	"os"

	_ "image/png"

	"github.com/faiface/pixel"
	"github.com/faiface/pixel/pixelgl"
	"golang.org/x/image/colornames"
)
```
After this we can insert our helper function
```go
func loadPicture(path string) (pixel.Picture, error) {
	file, err := os.Open(path)
	if err != nil {
		return nil, err
	}
	defer file.Close()
	img, _, err := image.Decode(file)
	if err != nil {
		return nil, err
	}
	return pixel.PictureDataFromImage(img), nil
}
```
This simply loads an image from a path, decodes it into a standard `image.Image`, then converts it to a `pixel.Picture`

Now we can insert this snippet into our window config, adding an icon to the window.
```go
    icon, err := LoadPicture("icon.png")         //
    if err != nil {                              //
        panic(err)                               //
    }                                            //

    cfg := pixelgl.WindowConfig{
        Title:  "Go Pixel",
        Bounds: pixel.R(0, 0, SCREENX, SCREENY),
        Icon:  []pixel.Picture{icon},            //
        VSync: true,
    }
```
Now, you might have noticed that when you add the icon to the config, you use a slice. This is because you can add multiple differently sized icons to the window, and the OS can choose which one to display given the context. For example, Windows uses a 32x32 icon for the Taskbar and Start menu, a 48x48 icon for Desktop shorcuts, and 16x16 icons for Task Manager. You can stylize theses smaller icons differently vs. the default scaling. 

## Game Creation
IMDraw can create many more shapes, and there are countless ways to utilize it, but for now all we need is the `imd.Rectangle` function.

If you're not familiar with [Conway's Game of Life](https://en.wikipedia.org/wiki/Conway%27s_Game_of_Life), don't worry, it's a pretty famous cellular automation. It's a zero-player game, meaning that it requires no input past the beginning. It starts with a theoretically infinite grid and checks three simple rules each generation. Based off of each cell's eight neighbours, or adjacent cells. The rules are as follows.
1. Any live cell with two or three live neighbours survives.
2. Any dead cell with three live neighbours becomes a live cell.
3. All other live cells die in the next generation. Similarly, all other dead cells stay dead.

So how do we implement this? One way we can do this is creating a struct to store the game data and methods:

```go
const BOARDX, BOARDY = 48, 27

type Board struct {
    cells [BOARDX][BOARDY]
}

func (b *Board) Set(x, y int, state bool) {
    b.cells[x][y] = state
}

func (b *Board) Get(x, y int) bool {
    return b.cells[x][y]
}

func (b *Board) Invert(x, y int) {
    b.cells[x][y] = !b.cells[x][y]
}
```

This struct provides a way to set, get, and invert indices in a board of a set size. But wait, why a set size? I thought this was supposed to be an infinite grid! Unfortunatly an infinite grid is much more complicated than a set size, and most implementations of this simply loop the board if you reach the end, which is what we'll do. Now that the simple stuff is out of the way, how let's program the logic for advancing to the next generation:
```go
func (b *Board) Neighbours(x, y int) int {
    // loop min and max if they are outside the bounds
    min_y := y - 1
    max_y := y + 1
    if min_y < 0 {
        min_y = BOARDY - 1
    } else if max_y >= BOARDY {
        max_y = 0
    }

    min_x := x - 1
    max_x := x + 1
    if min_x < 0 {
        min_x = BOARDX - 1
    } else if max_x >= BOARDX {
        max_x = 0
    }

    // check each neighboring cell, incrementing count if alive
    var neighbours int

    if b.Get(min_x, min_y) {neighbours++}
    if b.Get(x, min_y)     {neighbours++}
    if b.Get(max_x, min_y) {neighbours++}

    if b.Get(min_x, y)     {neighbours++}
    if b.Get(max_x, y)     {neighbours++}

    if b.Get(min_x, max_y) {neighbours++}
    if b.Get(x, max_y)     {neighbours++}
    if b.Get(max_x, max_y) {neighbours++}

    return neighbours
}

func (b *Board) Update() {
    temp := b.cells
    for x := 0; x < BOARDX; x++ {
        for y := 0; y < BOARDY; y++ {
            neighbours := b.Neighbours(x, y)
            if b.Get(x, y) { // check if the cell is alive
                if neighbours != 2 && neighbours != 3 {
                    temp[x][y] = false
                }
            } else if neighbours == 3 {
                temp[x][y] = true
            }
        }
    }
    b.cells = temp
}
```
Woah! That's quite a bit, what's even happening here? In the `Board.Update` function we create a temporary array so we don't overwrite cells before we read it. After this we itterate through each cell, checking the neighbors against the rules we specified earlier. But how do we get the neighbours of a cell? That's where the `Board.Neighbours` function comes in. In the first half of the function it calculates the minimum and maximum coordinate, and loops it to the other side of the board if it is outside the bounds. The last half of the function checks each adjacent tile and increments the value if it is alive.

Now that we can advance to the next generation, we need a way to show the board on the screen. We can do this by using the `IMDraw.Rectangle` function, which we can use to color the living cells.
```go
func (b *Board) Draw() {
    imd.Color = colornames.White
    for x := 0; x < BOARDX; x++ {
        for y := 0; y < BOARDY; y++ {
            if b.Get(x, y) {
                imd.Push(pixel.V(float64(x*20), float64(y*20)))
                imd.Push(pixel.V(float64(x*20)+20, float64(y*20)+20))
                imd.Rectangle(0)
            }
        }
    }
}
```
This function itterates through each cell, drawing a white rectangle in the position if it is alive. But what's this `pixel.V` function. This function returns a `pixel.Vec`, which is a vector with X and Y coordinates that Pixel uses to position things around the screen. We are also multiplying the index by 20 in order to scale the cells to the screen, since the screen size is 20 times the board size.

Now we need a way to run the simulation and change cells on the screen, which we can do by checking for mouse clicks and inverting the cell that was pressed. This code will go inside our main game loop, like so:
```go
    const UPDATETIME = time.Second/8
    update := time.Now()
    running := false

    for !win.Closed() {
        imd.Clear()

        if win.JustPressed(pixelgl.KeySpace) {
            running = !running
        }

        if running && time.Since(update) >= UPDATETIME {
            board.Update()
            update = time.Now()
        }

        // Check if the left mouse button was just pressed and the sim is not running
        if !running && win.JustPressed(pixelgl.MouseButtonLeft) {
            // get mouse position and convert to cell index
            pos := win.MousePosition()
            cell_x, cell_y := int(pos.X/20), int(pos.Y/20)

            board.Invert(cell_x, cell_y)
        }

        board.Draw()

        win.Clear(colornames.Black)
        imd.Draw(win)
        win.Update()
    }
```
What we're doing here is checking if certain keys were just pressed with the `Window.JustPressed` method. You can check a specific button with the coorisponding constant in `PixelGL`. At the beginning we check if the space bar was pressed, and if so play/pause the simulation. After we do that we check if the simulation is running and also if the time since the last update has exceeded the minimum time between updates. We add a minimum between time because without the simulation would be updating at the same speed as your FPS, which is a little nausiating. After that check, we also check if the left mouse button was pressed, and if the program isn't running we get the mouse coordinates, which we can get with the `Window.MousePosition` function. The mouse coordinates are in screen coordinates, which doesn't work for our smaller grid, so we have to scale the coordinates down before inverting the clicked cell. After all this logic we can simply draw the board to the screen! Here's what all of this looks like with a simple glider:

![GIF of game](https://raw.githubusercontent.com/commonkestrel/Replit-Pixel/main/screenshots/game.gif)

Here's the final code:
```go
package main

import (
    "os"
    "time"
    "image"

    _ "image/png"

    "github.com/faiface/pixel"
    "github.com/faiface/pixel/imdraw"
    "github.com/faiface/pixel/pixelgl"
    "golang.org/x/image/colornames"
)

const (
    SCREENX, SCREENY = 960, 540
    BOARDX, BOARDY   = 48, 27 // resolution/20
    UPDATETIME = time.Second/4
)

var (
    win     *pixelgl.Window
    imd     *imdraw.IMDraw
    running bool
)

type Board struct {
    cells [BOARDX][BOARDY]bool
}

func NewBoard() *Board {
    return &Board{}
}

func (b *Board) Set(x, y int, state bool) {
    b.cells[x][y] = state
}

func (b *Board) Get(x, y int) bool {
    return b.cells[x][y]
}

func (b *Board) Invert(x, y int) {
    b.cells[x][y] = !b.cells[x][y]
}

func (b *Board) Neighbours(x, y int) int {
    // loop min and max if they are outside the bounds
    min_y := y - 1
    max_y := y + 1
    if min_y < 0 {
        min_y = BOARDY - 1
    } else if max_y >= BOARDY {
        max_y = 0
    }

    min_x := x - 1
    max_x := x + 1
    if min_x < 0 {
        min_x = BOARDX - 1
    } else if max_x >= BOARDX {
        max_x = 0
    }

    // check each neighboring cell, incrementing count if alive
    var neighbours int

    if b.Get(min_x, min_y) {neighbours++}
    if b.Get(x, min_y) {neighbours++}
    if b.Get(max_x, min_y) {neighbours++}

    if b.Get(min_x, y) {neighbours++}
    if b.Get(max_x, y) {neighbours++}

    if b.Get(min_x, max_y) {neighbours++}
    if b.Get(x, max_y) {neighbours++}
    if b.Get(max_x, max_y) {neighbours++}

    return neighbours
}

func (b *Board) Update() {
    temp := b.cells
    for x := 0; x < BOARDX; x++ {
        for y := 0; y < BOARDY; y++ {
            neighbours := b.Neighbours(x, y)
            if b.Get(x, y) { // check if the cell is alive
                if neighbours != 2 && neighbours != 3 {
                    temp[x][y] = false
                }
            } else if neighbours == 3 {
                temp[x][y] = true
            }
        }
    }
    b.cells = temp
}

func (b *Board) Draw() {
    imd.Color = colornames.White
    for x := 0; x < BOARDX; x++ {
        for y := 0; y < BOARDY; y++ {
            if b.Get(x, y) {
                imd.Push(pixel.V(float64(x*20), float64(y*20)))
                imd.Push(pixel.V(float64(x*20)+20, float64(y*20)+20))
                imd.Rectangle(0)
            }
        }
    }
}

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

    board := NewBoard()

    update := time.Now()

    for !win.Closed() {
        imd.Clear()

        if win.JustPressed(pixelgl.KeyEscape) {
            win.SetClosed(true)
        }

        if win.JustPressed(pixelgl.KeySpace) {
            running = !running
        }

        if running && time.Since(update) >= UPDATETIME {
            board.Update()
            update = time.Now()
        }

        // Check if the left mouse button was just pressed and the sim is not running
        if !running && win.JustPressed(pixelgl.MouseButtonLeft) {
            // get mouse position and convert to cell index
            pos := win.MousePosition()
            cell_x, cell_y := int(pos.X/20), int(pos.Y/20)
            
            board.Invert(cell_x, cell_y)
        }

        board.Draw()

        win.Clear(colornames.Black)
        imd.Draw(win)
        win.Update()
    }
}

func main() {
    pixelgl.Run(run)
}

```
