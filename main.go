package main

import "github.com/gopherjs/gopherjs/js"

type JsWorld struct {
	*js.Object
	World *World
}

func main() {
	world := NewWorld()
	jsWorld := &JsWorld{
		Object: js.MakeWrapper(world),
		World:  world,
	}

	jsWorld.Set("preload", jsWorld.World.Preload)
	jsWorld.Set("create", jsWorld.World.Create)
	jsWorld.Set("update", jsWorld.World.Update)

	Game().State().Add("world", jsWorld)
	Game().State().StartI("world")
}
