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
	
		Fortia.camera2.move({x:0, y:0, z: 50});
	
		this.running = true;
		this.loop();
	},
	stop: function(){
		this.running = false;
	},

	render2: function(){ // 2d renderer, df like
		Fortia.camera2.render();
	},
	render3: function(){}, // 3d renderer
	update: function(delta){
		var keys = KeyboardJS.activeKeys();
		var moveMap = {
			"up": {x: 0, y: -1, z: 0},
			"down": {x: 0, y: 1, z: 0},
			"left": {x: -1, y: 0, z: 0},
			"right": {x: 1, y: 0, z: 0},
			"period": {x: 0, y: 0, z: 1},
			"comma": {x: 0, y: 0, z: -1},
		}
		var moveBy = {x: 0, y: 0, z: 0}; 
		var move = false;
		for (var i = 0; i < keys.length; i++) {
			var key =keys[i]
			switch(key){
				case "up":
				case "down":
				case "left":
				case "right":
				case "comma":
				case "period":
					move = true;
					moveBy.x += moveMap[key].x;
					moveBy.y += moveMap[key].y;
					moveBy.z += moveMap[key].z;
					break;
			}	
		};
		if (move) {
			Fortia.camera2.move(moveBy);
		};
	},
	loop: function(){
		var that = Fortia.game;
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