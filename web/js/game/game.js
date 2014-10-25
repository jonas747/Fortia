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
		this.api.get("chunk?x=1&y=1", "", function(response){
			that.setChunk(response);
		})

		this.running = true;
		this.loop();
	},
	stop: function(){
		this.running = false;
	},

	render2: function(){ // 2d renderer, df like

	},
	render3: function(){}, // 3d renderer
	update: function(delta){},
	loop: function(){
		if (this.running) {
			window.requestAnimationFrame(this.loop);
		};
		if (!this.lastUpdate) {
			this.lastUpdate = Date.now();
		};

		var now = Date.now();
		var delta = now - this.lastUpdate;
		lastUpdate = now;

		this.render2();
		this.update(delta);
	},
	resize: function(){
		var cwidth = window.innerWidth;
		var cheight = window.innerHeight - 50;

		Fortia.gameView.canvasWidth = cwidth;
		Fortia.gameView.canvasHeight = cheight;
		$("#game-canvas").width(cwidth);
		$("#game-canvas").height(cheight);
	},
	setChunk: function(chunk){
		console.log(chunk)
	},
}