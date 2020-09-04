package main

import (
    "github.com/faiface/pixel"
    "github.com/faiface/pixel/pixelgl"
    "golang.org/x/image/colornames"
    "io/ioutil"
    "time"
)

func main() {
    pixelgl.Run(run)
}

func run() {
    // p, _ := ioutil.ReadFile("/Users/erikvanbrakel/roms/display-test.c8")
   //  p, _ := ioutil.ReadFile("/Users/erikvanbrakel/roms/test_opcode.ch8")
     p, _ := ioutil.ReadFile("//Users/erikvanbrakel/repos/chip8/roms/programs/Chip8 emulator Logo [Garstyciuks].ch8")

   c := NewCPU(p)

    cfg := pixelgl.WindowConfig{
        Title:  "CHIP-8",
        Bounds: pixel.R(0, 0, 640, 320),
        VSync:  true,
    }
    win, err := pixelgl.NewWindow(cfg)
    if err != nil {
        panic(err)
    }

    for !win.Closed() {
        c.Cycle()
        time.Sleep(1/60 * time.Second)
        win.Clear(colornames.Aqua)
        screen := pixel.MakePictureData(pixel.R(0,0,64,32))
        for x := 0; x<64;x++ {
            for y := 0; y<32;y++ {
                screen.Pix[x + (31-y) * 64].R = c.DisplayBuffer[x][y] * 255
                screen.Pix[x + (31-y) * 64].G = c.DisplayBuffer[x][y] * 255
                screen.Pix[x + (31-y) * 64].B = c.DisplayBuffer[x][y] * 255
                screen.Pix[x + (31-y) * 64].A = 255
            }
        }

        sprite := pixel.NewSprite(screen, screen.Bounds())
        sprite.Draw(win, pixel.IM.Scaled(pixel.ZV, 10).Moved(win.Bounds().Center()))
        win.Update()
    }
}


