package main

import (
	"fmt"

	rl "github.com/gen2brain/raylib-go/raylib"
	hp "hydrogen-browser.com/html-parser"
)

const (
	W = 800
	H = 600
)

func main() {
	pageFilePath := "resources/TheProject.html"

	res, err := hp.ParseHTML(pageFilePath)
	if err != nil {
		fmt.Println("Error opening file:", err)
	}

	fmt.Print(res)

	rl.SetConfigFlags(rl.FlagWindowResizable)

	rl.InitWindow(W, H, "Hydrogen Browser")
	defer rl.CloseWindow()

	for !rl.WindowShouldClose() {
		rl.BeginDrawing()
		rl.ClearBackground(rl.White)

		// TODO: Search bar + Getting page from the web
		if key := rl.GetKeyPressed(); rl.IsWindowFocused() && key != 0 {
			fmt.Println(key)
		}

    rl.DrawText("Hello World", W/2, H/2, 20, rl.Pink)

		rl.EndDrawing()
	}
}
