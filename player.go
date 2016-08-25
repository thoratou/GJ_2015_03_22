package main

import "github.com/thoratou/go-phaser/generated/phaser"

type Player struct {
	World                  *World
	Sprite                 *phaser.Sprite
	Life                   int
	InvulnerabilityEndTime int
	Direction              int
}

func NewPlayer(world *World) *Player {
	player := &Player{
		World:  world,
		Sprite: Game().Add().Sprite3O(0, 0, "medic_simple"),
		Life:   config.PLAYER_LIFE,
		InvulnerabilityEndTime: 0,
		Direction:              config.DOWN,
	}

	Game().Physics().Enable1O(player.Sprite, Game().Physics().ARCADE())
	body := phaser.ToPhysicsArcadeBody(player.Sprite.Body())
	body.Bounce().SetToI(0.1, 0.1)
	body.SetCollideWorldBoundsA(true)

	player.Sprite.Animations().Add3O("walk-up", []interface{}{12, 13, 15, 14}, 7, true)
	player.Sprite.Animations().Add3O("walk-down", []interface{}{0, 1, 2, 3}, 7, true)
	player.Sprite.Animations().Add3O("walk-left", []interface{}{4, 5, 6, 7}, 7, true)
	player.Sprite.Animations().Add3O("walk-right", []interface{}{8, 9, 11, 10}, 7, true)
	player.Sprite.Animations().Add3O("idle-up", []interface{}{12}, 0, false)
	player.Sprite.Animations().Add3O("idle-down", []interface{}{0}, 0, false)
	player.Sprite.Animations().Add3O("idle-left", []interface{}{4}, 0, false)
	player.Sprite.Animations().Add3O("idle-right", []interface{}{9}, 0, false)

	return player
}

func (p *Player) Damage() {
	if p.InvulnerabilityEndTime < Game().Time().Now() {
		p.Life = p.Life - 1
		p.InvulnerabilityEndTime = Game().Time().Now() + config.INVULNERABILITY_TIME
	}
}

func (p *Player) Update() {
	body := phaser.ToPhysicsArcadeBody(p.Sprite.Body())
	body.Velocity().SetXA(0)
	body.Velocity().SetYA(0)
	for direction := 0; direction < 4; direction++ {
		if Game().Input().Keyboard().IsDown(config.KEYS[direction]) {
			//compute velocity
			body.Velocity().SetXA(config.COORMAP[direction][0] * config.SPEED)
			body.Velocity().SetYA(config.COORMAP[direction][1] * config.SPEED)
		}
	}

	//compute tile
	//up and down sprites get the priority on left and right
	if body.Velocity().Y() != 0 {
		if body.Velocity().Y() > 0 {
			p.Sprite.Animations().Play("walk-down")
			p.Direction = config.DOWN
		} else {
			p.Sprite.Animations().Play("walk-up")
			p.Direction = config.UP
		}
	} else if body.Velocity().X() != 0 {
		if body.Velocity().X() > 0 {
			p.Sprite.Animations().Play("walk-right")
			p.Direction = config.RIGHT
		} else {
			p.Sprite.Animations().Play("walk-left")
			p.Direction = config.LEFT
		}
	} else {
		//idle
		if p.Direction == config.UP {
			p.Sprite.Animations().Play("idle-up")
		} else if p.Direction == config.DOWN {
			p.Sprite.Animations().Play("idle-down")
		} else if p.Direction == config.LEFT {
			p.Sprite.Animations().Play("idle-left")
		} else if p.Direction == config.RIGHT {
			p.Sprite.Animations().Play("idle-right")
		}
	}

	if Game().Input().Keyboard().IsDown(config.SHOOT) {
		p.World.Shoot(phaser.NewPoint2O(p.Sprite.X(), p.Sprite.Y()), p.Direction)
	}
}
