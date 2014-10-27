var Fortia = Fortia || {};
Fortia.initGameView = function(){
	var GameView = Backbone.View.extend({
		el: "#main-container",
		template: templates.game,
		templateHeader: templates.lobbyHeader,

		initialize: function() {
		},

		render: function(){
			//var header = this.templateHeader({username: localStorage.getItem("username")})
			var body = this.template(Fortia)
			this.$el.html(body)
		},
		update: function(){

		},
		switchTo: function(){
			this.render();
			$(window).resize(Fortia.game.resize)
			Fortia.game.resize();
		},
	});
	gameView = new GameView()

	return gameView;
}

Fortia.game = {
	chunks: [],
	init: function(view, world){
		this.api = new Fortia.REST({urlRoot:"http://localhost:8081/"});
		if (Fortia.Production) {
			// if this is a pruduction build
			this.api.urlRoot = "http://"+world+".jonas747.com/"
		};

		this.view = view;
		this.worldName = world;
		this.running = false;
		this.initKeybinds();
		this.camera = Fortia.camera;
		Keyboard.storeEvents = true;

		var stats = new Stats();
		stats.setMode(0);
		document.body.appendChild(stats.domElement);
		this.stats = stats;

	},
	initKeybinds: function(){
	},
	start: function(){
		console.log("Starting the game")
		var that = this;
		/*
		this.api.get("worldinfo", "", function(response){
			that.worldInfo = response;
			running = true;
		});
		this.api.get("visiblechunks", "", function(response){
			that.visibleCunks = response 
		})
*/
	
		this.camera.move({x:0, y:0, z: 50});
	
		this.running = true;
		this.loop();
	},
	stop: function(){
		this.running = false;
	},

	render2: function(){ // 2d renderer, df like
		this.camera.render2();
	},
	render3: function(){}, // 3d renderer
	update: function(delta){
		var moveMap = {
			"Up": {x: 0, y: -1, z: 0},
			"Down": {x: 0, y: 1, z: 0},
			"Left": {x: -1, y: 0, z: 0},
			"Right": {x: 1, y: 0, z: 0},
		}
		var moveBy = {x: 0, y: 0, z: 0}; 
		var move = false;

		for(var key in moveMap){
			if(Keyboard.isKeyDown(key)){
				move = true;
				var madd = moveMap[key]
				moveBy.x += madd.x
				moveBy.y += madd.y
			}
		}

		for (var i = 0; i < Keyboard.events.length; i++) {
			var evt = Keyboard.events[i];
			if (evt.down) {
				switch(evt.keyStr){
					case ",":
						// go up
						move = true;
						moveBy.z -= 1;
						break;
					case ".":
						move = true;
						moveBy.z += 1;
						// go down
						break;
				}
			};
		};
		Keyboard.events = [];

		if (move) {
			this.camera.move(moveBy);
		};
	},
	loop: function(){
		var that = Fortia.game;
		that.stats.end();
		that.stats.begin();

		if (that.running) {
			window.requestAnimationFrame(that.loop);
		};
		if (!that.lastUpdate) {
			that.lastUpdate = Date.now();
		};

		var now = Date.now();
		var delta = now - that.lastUpdate;
		lastUpdate = now;

		that.update(delta);
		that.render2();
	},
	resize: function(){
		var cwidth = window.innerWidth;
		var cheight = window.innerHeight - 50;

		Fortia.gameView.canvasWidth = cwidth;
		Fortia.gameView.canvasHeight = cheight;
		$("#game-canvas")[0].width = cwidth;
		$("#game-canvas")[0].height = cheight;
	},
	setChunk: function(chunk){
		console.log(chunk)
	},
}