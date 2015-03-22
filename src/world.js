/*jslint plusplus: true */
/*jslint browser: true*/
/*global $, Phaser*/

var Config = {
    KEYS : [Phaser.Keyboard.UP, Phaser.Keyboard.DOWN, Phaser.Keyboard.LEFT, Phaser.Keyboard.RIGHT],
    COORMAP : [[0, -1],
               [0, 1],
               [-1, 0],
               [1, 0]],
    UP : 0,
    DOWN : 1,
    LEFT : 2,
    RIGHT : 3,
    SHOOT : Phaser.Keyboard.S,
    FIRE_RATE : 100,
    FIRE_SPEED : 500,
    FIRE_INIT_OFFSET : 3,
    SPEED : 500
};

var Globals = {
    backgroundImg : null,
    backgroundMusic : null,
    target : null,
    leftTop : null,
	zombieView : 200
};

// game variables
var game = new Phaser.Game(800, 600, Phaser.CANVAS, 'game-div');

var world = {
    zombies: null,
    humans: null,
    player: null,
    balls: null,
	npcs:null,
    nextFire:0,
    
    preload: function () {
        // sounds
        game.load.audio('background_music', 'res/snd/temp.mp3');

        // images
        game.load.image('zombie', 'res/img/zombie_simple.png');
        game.load.image('ball', 'res/img/brocoli.png');
		game.load.image('human', 'res/img/human_simple.png');
        game.load.spritesheet('medic_simple', 'res/img/medic_simple.png', 33, 58);
    },

    create: function () {
		game.physics.startSystem(Phaser.Physics.ARCADE);
	
        //Entities
        this.zombies = game.add.group();
		this.zombies.enableBody = true;
        
        this.humans = game.add.group();
        this.humans.enableBody = true;

        this.player = new Player(world);
		Globals.target = this.player;

        this.balls = game.add.group();
        this.balls.createMultiple(500, 'ball', 0, false);
        this.balls.enableBody = true;
        this.nextFire = game.time.now;
        
        // sound
        Globals.backgroundMusic = game.add.audio('background_music');
        Globals.backgroundMusic.play(null, 0, 1, true);

		// Humans 
		for(i = 0; i < 50; i++)
		{
			var human = this.humans.create(game.world.randomX, game.world.randomY, 'human');
			human.name = 'human' + i;
			human.body.collideWorldBounds = true;
			human.body.width = 30;
			human.body.height = 50;
			human.body.bounce.setTo(0.8, 0.8);
			human.body.velocity.setTo(10 + Math.random() * 40, 10 + Math.random() * 40);
		}
    },

    update: function () {
		game.physics.arcade.collide(this.player.sprite, this.zombies);
		game.physics.arcade.collide(this.player.sprite, this.humans);
		game.physics.arcade.collide(this.zombies, this.humans);
		game.physics.arcade.collide(this.zombies, this.zombies);
		game.physics.arcade.collide(this.humans, this.humans);
		
		this.player.update();
		this.zombies.forEach(this.gotToTarget, this, true);
        this.humans.update();
    },
	
	gotToTarget: function(zombie) {
		if (this.targetInRange(zombie)) {
			game.physics.arcade.moveToObject(zombie, Globals.target.sprite, 100);
		}
	},
	
	targetInRange: function(zombie) {
		dx = zombie.x - Globals.target.sprite.x;
		dy = zombie.y - Globals.target.sprite.y;
		return Math.sqrt(dx*dx + dy*dy) < Globals.zombieView;
	},
    
    shoot: function(position, direction) {
        if (game.time.now > this.nextFire){
            this.nextFire = game.time.now + Config.FIRE_RATE;
            var ball = this.balls.getFirstExists(false); // get the first created fireball that no exists atm
            console.log(position);
            console.log(ball);
            if (ball){
                ball.exists = true; 
                ball.lifespan = 2500;
                game.physics.enable(ball, Phaser.Physics.ARCADE);
                
                if(direction == Config.UP){  
                    ball.reset(position.x, position.y-Config.FIRE_INIT_OFFSET);
                    ball.body.velocity.x = 0;
                    ball.body.velocity.y = -Config.FIRE_SPEED;
                } else if(direction == Config.DOWN){
                    ball.reset(position.x, position.y+Config.FIRE_INIT_OFFSET);
                    ball.body.velocity.x = 0;
                    ball.body.velocity.y = Config.FIRE_SPEED;
                } else if(direction == Config.LEFT){
                    ball.reset(position.x-Config.FIRE_INIT_OFFSET, position.y);
                    ball.body.velocity.x = -Config.FIRE_SPEED;
                    ball.body.velocity.y = 0;
                } else if(direction == Config.RIGHT){
                    ball.reset(position.x+Config.FIRE_INIT_OFFSET, position.y);
                    ball.body.velocity.x = Config.FIRE_SPEED;
                    ball.body.velocity.y = 0;
                }
            }
        }
    }
}

game.state.add('world', world);
game.state.start('world');
