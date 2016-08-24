package main

import "github.com/thoratou/go-phaser/generated/phaser"

func main() {
	Game().State().Add("world", phaser.WrapState(NewWorld()))
	Game().State().StartI("world")
}
