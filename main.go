package main


import engine "github.com/vorpalgame/raylib-engine"

func main() {


	c := engine.NewEngine()
	go NewGame()

	c.Start()
}
