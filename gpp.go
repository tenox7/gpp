/*
 * Copyright 2019 Google LLC
 *
 * Licensed under the Apache License, Version 2.0 (the "License");
 * you may not use this file except in compliance with the License.
 * You may obtain a copy of the License at
 *
 *     https://www.apache.org/licenses/LICENSE-2.0
 *
 * Unless required by applicable law or agreed to in writing, software
 * distributed under the License is distributed on an "AS IS" BASIS,
 * WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
 * See the License for the specific language governing permissions and
 * limitations under the License.
 */
package main

import (
    "container/ring"
    "flag"
    "fmt"
    "io/ioutil"
    "math"
    "net"
    "os"
    "os/exec"
    "time"

    ping "github.com/digineo/go-ping"
    "github.com/digineo/go-ping/monitor"
    "github.com/veandco/go-sdl2/sdl"
    "github.com/veandco/go-sdl2/ttf"

    _ "./statik"
    "github.com/rakyll/statik/fs"
)

var (
    pingInterval       = 1 * time.Second
    pingTimeout        = 4 * time.Second
    size         uint  = 56
    title              = "GoPingPlot"
    windowWidth  int32 = 180
    windowHeight int32 = 100
    windowMargin int32 = 5
    panelHeight  int32 = 45
    targetSize   int32 = 90
    panelWidth         = windowWidth - (windowMargin * 2)
    fontName           = "noto.ttf"
    fontSize     int32 = 11
)

func errbox(format string, args ...interface{}) {
    sdl.Init(sdl.INIT_VIDEO)
    window, _ := sdl.CreateWindow(title, sdl.WINDOWPOS_UNDEFINED, sdl.WINDOWPOS_UNDEFINED, 100, 100, sdl.WINDOW_HIDDEN)
    sdl.ShowSimpleMessageBox(10, "Error", fmt.Sprintf(format, args...), window)
    window.Destroy()
    os.Exit(1)
}

func color(v float64, max float64) (uint8, uint8, uint8, uint8) {
    var red, green int = 0, 0xFF

    red = int(v / max * float64(0xFF))

    if red > 0xFF {
        green = 0xFF - (red - 0xFF)
        red = 0xFF
    } else if red < 0 {
        red = 0
    }

    if green > 0xFF {
        green = 0xFF
    } else if green < 0x80 {
        green = 0x80
    }
    return uint8(red), uint8(green), 0x00, 0xFF
}

func drawText(renderer *sdl.Renderer, font *ttf.Font, color sdl.Color, x int32, y int32, text string) (int32, int32) {
    var surface *sdl.Surface
    var texture *sdl.Texture

    surface, _ = font.RenderUTF8Blended(text, color)
    texture, _ = renderer.CreateTextureFromSurface(surface)
    _, _, w, h, _ := texture.Query()
    renderer.Copy(texture, nil, &sdl.Rect{x, y, w, h})
    surface.Free()
    texture.Destroy()
    return w, h
}

func plotRing(r *ring.Ring, host string, tgtnum int32, badping float64, renderer *sdl.Renderer, font *ttf.Font) {
    var min, max, avg, lst, tot float64 = 10000.0, 0, 0, 0, 0
    var i, h int32 = 0, 0
    var txt string
    var vs int32 = (tgtnum * targetSize) + 1
    var v float64

    renderer.SetDrawColor(0xFF, 0xFF, 0xFF, 0xFF)
    renderer.DrawRect(&sdl.Rect{windowMargin, vs + windowMargin + 15, windowWidth - (windowMargin * 2), panelHeight})
    drawText(renderer, font, sdl.Color{0xFF, 0xFF, 0xFF, 0xFF}, windowMargin, vs+windowMargin-3, host)

    r.Do(func(x interface{}) {
        v, _ = x.(float64)
        lst = v
        if v > max {
            max = v
        }
        if v < min {
            min = v
        }
        if v > 0 {
            tot += v
            i++
        }
    })

    avg = tot / float64(i)
    i = 0

    r.Do(func(x interface{}) {
        v, _ = x.(float64)
        if math.IsNaN(v) {
            h = panelHeight - 1
            renderer.SetDrawColor(0xFF, 0x00, 0x00, 0x00)
            renderer.DrawLine(windowMargin+i, vs+windowMargin+panelHeight+15-h,
                windowMargin+i, vs+windowMargin+panelHeight+15-2)
        } else if v > 0 {
            h = int32((v/max)*float64(panelHeight-2)) - 1
            renderer.SetDrawColor(color(v, badping))
            renderer.DrawLine(windowMargin+1+i, vs+windowMargin+panelHeight+13-h,
                windowMargin+1+i, vs+windowMargin+panelHeight+13)
        } else {
            h = 2
            renderer.SetDrawColor(0x00, 0xFF, 0x00, 0x00)
            renderer.DrawLine(windowMargin+i, vs+windowMargin+panelHeight+15-h,
                windowMargin+i, vs+windowMargin+panelHeight+15-2)
        }
        i++
    })
    txt = fmt.Sprintf("L=%.1f M=%.1f A=%.1f", lst, max, avg)
    drawText(renderer, font, sdl.Color{0xFF, 0xFF, 0xFF, 0xFF}, windowMargin, vs+windowMargin+panelHeight+15, txt)
}

func main() {
    var foreground bool
    var badPing float64

    flag.Usage = func() {
        errbox("Usage: \n%s host [host [...]]", os.Args[0])
    }

    flag.Float64Var(&badPing, "t", 100, "Pings longer than this will be be more redish color")
    flag.BoolVar(&foreground, "f", false, "Run in foreground, do not detach from terminal")
    flag.Parse()

    if !foreground {
        cwd, err := os.Getwd()
        if err != nil {
            errbox("Getcwd error: %s\n", err)
        }
        args := []string{"-f"}
        args = append(args, os.Args[1:]...)
        cmd := exec.Command(os.Args[0], args...)
        cmd.Dir = cwd
        if err := cmd.Start(); err != nil {
            errbox("Startup error: %s\n", err)
        }
        cmd.Process.Release()
        os.Exit(0)
    }

    if na := flag.NArg(); na == 0 {
        flag.Usage()
    } else if na > int(^byte(0)) {
        errbox("Too many targets")
    }

    // Ping Init
    pinger, err := ping.New("0.0.0.0", "::")
    if err != nil {
        errbox("Unable to bind:\n%s\nAre you running as root?\n", err)
    }
    pinger.SetPayloadSize(uint16(size))
    defer pinger.Close()

    monitor := monitor.New(pinger, pingInterval, pingTimeout)
    defer monitor.Stop()

    rings := make(map[string]*ring.Ring)
    targets := flag.Args()
    for i, target := range targets {
        ipAddr, err := net.ResolveIPAddr("", target)
        if err != nil {
            errbox("invalid target '%s':\n %s", target, err)
        }
        monitor.AddTargetDelayed(string([]byte{byte(i)}), *ipAddr, 10*time.Millisecond*time.Duration(i))
        rings[target] = ring.New(int(panelWidth - 2))
    }

    // SDL Init
    if err := sdl.Init(sdl.INIT_VIDEO); err != nil {
        errbox("Failed to initialize SDL:\n %s\n", err)
    }

    if err := ttf.Init(); err != nil {
        errbox("Failed to initialize TTF:\n %s\n", err)
    }

    windowHeight = int32(len(rings)) * targetSize
    window, err := sdl.CreateWindow(title, sdl.WINDOWPOS_UNDEFINED, sdl.WINDOWPOS_UNDEFINED, windowWidth, windowHeight, sdl.WINDOW_OPENGL)
    if err != nil {
        errbox("Failed to create window:\n %s\n", err)
    }
    defer window.Destroy()

    renderer, err := sdl.CreateRenderer(window, -1, sdl.RENDERER_ACCELERATED)
    if err != nil {
        errbox("Failed to create renderer:\n %s\n", err)
    }
    defer renderer.Destroy()

    // Font stuff
    sfs, err := fs.New()
    if err != nil {
        errbox("Unable to initialize Statik:\n%s\n", err)
    }

    fh, err := sfs.Open("/" + fontName)
    if err != nil {
        errbox("Unable to open font from Statik:\n%s\n", err)
    }

    fontData, err := ioutil.ReadAll(fh)
    if err != nil {
        errbox("Unable to read font from Statik:\n%s\n", err)
    }

    tmp, err := ioutil.TempFile(os.TempDir(), fontName)
    if err != nil {
        errbox("Unable to create temp file:\n%s\n", err)
    }
    defer os.Remove(tmp.Name())

    if _, err := tmp.Write(fontData); err != nil {
        errbox("Unable to write temp font:\n%s\n", err)
    }
    tmp.Close()

    font, err := ttf.OpenFont(tmp.Name(), int(fontSize))
    if err != nil {
        errbox("Failed to open font:\n %s\n", err)
    }

    // Main Loop
    for {
        for event := sdl.PollEvent(); event != nil; event = sdl.PollEvent() {
            switch event.(type) {
            case *sdl.QuitEvent:
                os.Exit(0)
            }
        }
        renderer.Clear()
        renderer.SetDrawColor(100, 100, 100, 0x20)
        renderer.FillRect(&sdl.Rect{0, 0, windowWidth, windowHeight})
        for i, metrics := range monitor.ExportAndClear() {
            n := int32([]byte(i)[0])
            t := targets[n]
            rings[t].Value = float64(metrics.Median)
            rings[t] = rings[t].Next()
            plotRing(rings[t], t, n, badPing, renderer, font)
        }
        renderer.Present()
        sdl.Delay(1000)
    }
}
