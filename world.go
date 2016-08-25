package main

import (
	"fmt"
	"math"
	"time"

	"github.com/gopherjs/gopherjs/js"
	"github.com/thoratou/go-phaser/generated/phaser"
)

type Config struct {
	KEYS                 []int
	COORMAP              [][]int
	UP                   int
	DOWN                 int
	LEFT                 int
	RIGHT                int
	SPEED                int
	SHOOT                int
	FIRE_RATE            int
	FIRE_SPEED           int
	FIRE_INIT_OFFSET     int
	ZOMBIE_VIEW_RADIUS   int
	MIN_MUTATE_TIME      int
	MAX_MUTATE_TIME      int
	MIN_MOVE_TIME        int
	MAX_MOVE_TIME        int
	PLAYER_LIFE          int
	INVULNERABILITY_TIME int
}

var (
	config = Config{
		KEYS: []int{phaser.GetKeyboardKey("UP"), phaser.GetKeyboardKey("DOWN"), phaser.GetKeyboardKey("LEFT"), phaser.GetKeyboardKey("RIGHT")},
		COORMAP: [][]int{
			[]int{0, -1},
			[]int{0, 1},
			[]int{-1, 0},
			[]int{1, 0},
		},
		UP:                   0,
		DOWN:                 1,
		LEFT:                 2,
		RIGHT:                3,
		SPEED:                200,
		SHOOT:                phaser.GetKeyboardKey("S"),
		FIRE_RATE:            300,
		FIRE_SPEED:           500,
		FIRE_INIT_OFFSET:     3,
		ZOMBIE_VIEW_RADIUS:   400,
		MIN_MUTATE_TIME:      3000,
		MAX_MUTATE_TIME:      40000,
		MIN_MOVE_TIME:        500,
		MAX_MOVE_TIME:        2000,
		PLAYER_LIFE:          5,
		INVULNERABILITY_TIME: 500,
	}
)

var (
	game *phaser.Game
)

func Game() *phaser.Game {
	if game == nil {
		game = phaser.NewGame4O(800, 600, phaser.GetGlobalConst("CANVAS"), "game-div")
	}
	return game
}

type World struct {
	backgroundImage *phaser.Sprite
	backgroundMusic *phaser.Sound
	zombies         *phaser.Group
	humans          *phaser.Group
	player          *Player
	target          *phaser.Sprite
	balls           *phaser.Group
	nextFire        int
	mobSpeed        int
	mobsNb          int
	lastZombieNb    int
	rng             *phaser.RandomDataGenerator
}

func NewWorld() *World {
	return &World{
		nextFire:     0,
		mobSpeed:     40,
		mobsNb:       35,
		lastZombieNb: 0,
		rng:          phaser.NewRandomDataGenerator1O(time.Now().String()),
	}
}

func (w *World) Preload() {
	fmt.Println("Start Preload")
	// sounds
	Game().Load().Audio("background_music", "res/snd/temp.mp3")

	// images
	Game().Load().Spritesheet("zombie", "res/img/zombie_sprite.png", 48, 63)
	Game().Load().Image1O("ball", "res/img/brocoli.png")
	Game().Load().Spritesheet("human", "res/img/human_sprite.png", 55, 65)
	Game().Load().Spritesheet("medic_simple", "res/img/medic_sprite.png", 35, 62)
	Game().Load().Image1O("background", "res/img/background.png")

	fmt.Println("End Preload")
}

func (w *World) Create() {
	fmt.Println("Start Create")

	Game().Physics().StartSystem(Game().Physics().ARCADE())

	fmt.Println("Physics started")

	// background
	w.backgroundImage = Game().Add().Sprite3O(0, 0, "background")

	// sound
	w.backgroundMusic = Game().Add().Audio("background_music")
	w.backgroundMusic.SetPositionA(0)
	w.backgroundMusic.SetVolumeA(1)
	w.backgroundMusic.SetLoopA(true)
	w.backgroundMusic.Play()
	js.Global.Get("console").Call("log", w.backgroundMusic.Object)

	fmt.Println("Music started")

	//Entities
	w.zombies = Game().Add().Group()
	w.zombies.SetEnableBodyA(true)
	fmt.Println("Zombie group created")

	w.humans = Game().Add().Group()
	w.humans.SetEnableBodyA(true)
	fmt.Println("Human group created")

	w.player = NewPlayer(w)
	w.target = w.player.Sprite
	fmt.Println("Player created")

	w.balls = Game().Add().Group()
	w.balls.CreateMultiple2O(500, "ball", 0, false)
	w.balls.SetEnableBodyA(true)
	w.nextFire = Game().Time().Now()
	fmt.Println("Balls created")

	// Humans
	for i := 0; i < w.mobsNb; i++ {
		valid := false
		x := 0
		y := 0
		for !valid {
			x = Game().World().RandomX()
			y = Game().World().RandomY()
			valid = math.Sqrt(float64(x*x+y*y)) > 150
			//fmt.Println("x %d, y %d, valid %t", x, y, valid)
		}

		human := phaser.Sprite{Object: w.humans.Create1O(x, y, "human").Object}
		human.SetNameA(fmt.Sprintf("human%d", i))
		body := phaser.ToPhysicsArcadeBody(human.Body())
		body.SetCollideWorldBoundsA(true)
		body.SetWidthA(30)
		body.SetHeightA(50)
		body.Bounce().SetToI(0.1, 0.1)
		human.SetLifespanA(w.rng.Between(config.MIN_MUTATE_TIME, config.MAX_MUTATE_TIME))
		human.Events().OnKilled().Add1O(w.HumanMutate, w)
		human.Set("lastmove", 0)

		//animations
		human.Animations().Add3O("walk-up", []interface{}{9, 10, 9, 11}, 7, true)
		human.Animations().Add3O("walk-down", []interface{}{0, 1, 0, 2}, 7, true)
		human.Animations().Add3O("walk-left", []interface{}{3, 4, 3, 5}, 7, true)
		human.Animations().Add3O("walk-right", []interface{}{6, 7, 6, 8}, 7, true)
	}
	fmt.Println("End Create")
}

func (w *World) Update() {
	//fmt.Println("Start Update")
	Game().Physics().Arcade().Collide3O(w.player.Sprite,
		w.zombies,
		func(_ *phaser.Sprite, _ *phaser.Sprite) {
			w.player.Damage()
		}, nil, w)
	Game().Physics().Arcade().Collide(w.player.Sprite, w.humans)
	Game().Physics().Arcade().Collide(w.zombies, w.humans)
	Game().Physics().Arcade().Collide(w.zombies, w.zombies)
	Game().Physics().Arcade().Collide(w.humans, w.humans)
	Game().Physics().Arcade().Collide1O(w.zombies,
		w.balls,
		func(zombie *phaser.Sprite, ball *phaser.Sprite) {
			ball.Kill()
			zombie.Kill()
		})
	Game().Physics().Arcade().Collide1O(w.humans,
		w.balls,
		func(human *phaser.Sprite, ball *phaser.Sprite) {
			ball.Kill()
			human.Kill()
		})

	w.player.Update()

	// zombies
	w.zombies.ForEach1O(w.GotToTarget, w, true)
	w.zombies.ForEach1O(DrawEntity, w, true)

	// humans
	w.humans.ForEach1O(w.RandomMove, w, true)
	w.humans.ForEach1O(DrawEntity, w, true)

	if w.player.Life == 0 {
		w.player.Sprite.Kill()
		w.target = nil
		w.mobSpeed = w.mobSpeed * 5

		style := map[string]interface{}{
			"font":  "65px Arial",
			"fill":  "#ffffff",
			"align": "center",
		}
		text := Game().Add().Text4O(Game().World().CenterX(), Game().World().CenterY(), "this.player undefined...\nTHE GAMEOVER\nF5 to reload", style)
		text.Anchor().SetI(0.5)
		text.AddColor("#ff0000", 24)
		text.AddColor("#ffffff", 20)
		text.AddColor("#ffff00", 36)
	}
	//fmt.Println("End Update")
}

func (w *World) HumanMutate(human *phaser.Sprite) {
	zombie := phaser.ToSprite(w.zombies.Create1O(human.X(), human.Y(), "zombie").Object)
	zombie.SetNameA(fmt.Sprintf("zombie%d", w.lastZombieNb))
	w.lastZombieNb++
	body := phaser.ToPhysicsArcadeBody(zombie.Body())
	body.SetCollideWorldBoundsA(true)
	body.SetWidthA(30)
	body.SetHeightA(50)
	body.Bounce().SetToI(0.1, 0.1)
	zombie.Set("lastmove", 0)

	//animations
	zombie.Animations().Add3O("walk-up", []interface{}{9, 10, 9, 11}, 7, true)
	zombie.Animations().Add3O("walk-down", []interface{}{0, 1, 0, 2}, 7, true)
	zombie.Animations().Add3O("walk-left", []interface{}{3, 4, 3, 5}, 7, true)
	zombie.Animations().Add3O("walk-right", []interface{}{6, 7, 6, 8}, 7, true)
}

func (w *World) GotToTarget(zombie *phaser.Sprite) {
	if w.target != nil && Game().Physics().Arcade().DistanceBetween(zombie, w.target) < config.ZOMBIE_VIEW_RADIUS {
		Game().Physics().Arcade().MoveToObject1O(zombie, w.target, 50)
	} else {
		w.RandomMove(zombie)
	}
}

func (w *World) RandomMove(entity *phaser.Sprite) {
	if Game().Time().Now()-entity.Get("lastmove").Int() > w.rng.Between(config.MIN_MOVE_TIME, config.MAX_MOVE_TIME) {
		body := phaser.ToPhysicsArcadeBody(entity.Body())
		body.Velocity().SetTo1O(w.rng.Between(-1, 1)*w.mobSpeed, w.rng.Between(-1, 1)*w.mobSpeed)
		entity.Set("lastmove", Game().Time().Now())
	}
}

func (w *World) Shoot(position *phaser.Point, direction int) {
	if Game().Time().Now() > w.nextFire {
		w.nextFire = Game().Time().Now() + config.FIRE_RATE
		ball := phaser.ToSprite(w.balls.GetFirstExists1O(false).Object) // get the first created fireball that no exists atm
		if ball != nil {
			ball.SetExistsA(true)
			ball.SetLifespanA(2500)
			Game().Physics().Enable1O(ball, Game().Physics().ARCADE())

			body := phaser.ToPhysicsArcadeBody(ball.Body())
			if direction == config.UP {
				ball.Reset(position.X(), position.Y()-config.FIRE_INIT_OFFSET)
				body.Velocity().SetXA(0)
				body.Velocity().SetYA(-config.FIRE_SPEED)
			} else if direction == config.DOWN {
				ball.Reset(position.X(), position.Y()+config.FIRE_INIT_OFFSET)
				body.Velocity().SetXA(0)
				body.Velocity().SetYA(config.FIRE_SPEED)
			} else if direction == config.LEFT {
				ball.Reset(position.X()-config.FIRE_INIT_OFFSET, position.Y())
				body.Velocity().SetXA(-config.FIRE_SPEED)
				body.Velocity().SetYA(0)
			} else if direction == config.RIGHT {
				ball.Reset(position.X()+config.FIRE_INIT_OFFSET, position.Y())
				body.Velocity().SetXA(config.FIRE_SPEED)
				body.Velocity().SetYA(0)
			}
		}
	}
}

func DrawEntity(entity *phaser.Sprite) {
	body := phaser.ToPhysicsArcadeBody(entity.Body())
	if math.Abs(float64(body.Velocity().X())) > math.Abs(float64(body.Velocity().Y())) {
		if body.Velocity().X() > 0 {
			entity.Animations().Play("walk-right")
		} else {
			entity.Animations().Play("walk-left")
		}
	} else {
		if body.Velocity().Y() > 0 {
			entity.Animations().Play("walk-down")
		} else {
			entity.Animations().Play("walk-up")
		}
	}
}
