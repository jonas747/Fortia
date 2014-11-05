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
		Keyboard.storeEvents = true;

		var stats = new Stats();
		stats.setMode(0);
		document.body.appendChild(stats.domElement);
		this.stats = stats;

		this.initScene();
		console.log("init!!!!!")
		$(window).resize(Fortia.game.resize)
		Fortia.game.resize();

		this.cameraHeight = 50;

		this.sight = 3;
		this.sightHeight = 3;

		this.cpos = new THREE.Vector3(0,0,100)

		this.blkMaterial = new THREE.MeshPhongMaterial({
			vertexColors: THREE.FaceColors,
			blending: THREE.AdditiveBlending,
			shininess: 80,
		});

		var boundWheelCB = this.onWheel.bind(this);
		window.addWheelListener(this.canvas, boundWheelCB);
	},
	initScene: function(){
		this.canvas = $("#game-canvas")[0];
		var width = Fortia.gameView.canvasWidth;
		var height = Fortia.gameView.canvasHeigh;

		var renderer = new THREE.WebGLRenderer({canvas: this.canvas});
		renderer.setSize( width, height );

		var scene = new THREE.Scene();

		var camera = new THREE.PerspectiveCamera(
		                                70,             	// Field of view
		                                width / height,     // Aspect ratio
		                                1,            		// Near plane
		                                10000				// Far plane
		                            );
		camera.position.set( 0, -2, 10 );
		camera.lookAt(new THREE.Vector3(0,0,0));
		//camera.rotation.z = -0.78;
		scene.add( camera );

		var light = new THREE.HemisphereLight( 0xffffff, 0xff0000, 1 );
		light.position.set( 0, 0, 1000 );
		scene.add( light );

		var geometry = new THREE.BoxGeometry( 1, 1, 1 );
		var material = new THREE.MeshBasicMaterial( { color: 0x00ff00 } );
		var cube = new THREE.Mesh( geometry, material );
		scene.add( cube );

		this.scene = scene;
		this.camera = camera;
		this.renderer = renderer;
	},
	initKeybinds: function(){
	},
	start: function(){
		console.log("Starting the game")
		var that = this;
		this.api.get("info", "", function(response){
			that.worldInfo = response;
			that.worldInfo.layerSize = parseInt(that.worldInfo.layerSize);
			console.log("Response", response)
			that.running = true;
			that.loop();
		});
		
		var boundUptLayer = this.updateLayers.bind(this);
		setInterval(boundUptLayer, 1000)
		this.loop();
	},
	stop: function(){
		this.running = false;
	},
	render3: function(){
		this.renderer.render(this.scene, this.camera);
	}, // 3d renderer
	update: function(delta){
		var moveMap = {
			"Up": new THREE.Vector3(0, 1, 0),
			"Down": new THREE.Vector3(0, -1, 0),
			"Left": new THREE.Vector3(-1, 0, 0),
			"Right": new THREE.Vector3(1, 0, 0),
		}
		var moveBy = new THREE.Vector3(0,0,0);
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
			this.moveCamera(moveBy);
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
		that.render3();
	},
	resize: function(){
		var cwidth = window.innerWidth;
		var cheight = window.innerHeight - 50;

		Fortia.gameView.canvasWidth = cwidth;
		Fortia.gameView.canvasHeight = cheight;
		// $("#game-canvas")[0].width = cwidth;
		// $("#game-canvas")[0].height = cheight;

		Fortia.game.renderer.setSize(cwidth, cheight);

		var aspect = cwidth / cheight;
		Fortia.game.camera.aspect = aspect;
		Fortia.game.camera.updateProjectionMatrix();
	},
	setChunk: function(chunk){
		console.log(chunk)
	},

	// index = size * x + y
	coordsToIndex: function(pos) {
		return this.worldInfo.layerSize*pos.x + pos.y
	},

	// Return a blocks x and y from the index in the layer slice
	// x = index / size
	// y = index - (x * size)
	indexToCoords: function(index) {
		var x = index / this.worldInfo.layerSize
		var y = index - (x * this.worldInfo.layerSize)
		return new THREE.Vector2(x, y)
	},	
	moveCamera: function(by){
		if (!this.cpos){
			this.cpos = new THREE.Vector3(0,0,0); // Default camera position
		}

		var newPos = this.cpos.clone();
		newPos.add(by);

		// if we moved to a new layer/chunk we fetch/update the cache
		var oldLPos = this.worldToLayerPos(this.cpos)
		var newLPos = this.worldToLayerPos(newPos)

		if (!newLPos.equals(oldLPos)) {
			console.log("Moved to a new Layer!")
			this.layersDirty = true;
		};
		
		this.camera.position.copy(newPos);
		this.camera.position.z += Fortia.game.cameraHeight;
		this.cpos = newPos

	},
	updateLayers: function(){
		if (!this.layersDirty) {
			return;
		};

		this.cachedLayers = this.cachedLayers || {};
		this.fetchingLayers = this.fetchingLayers || {};

		var newLPos = this.worldToLayerPos(this.cpos)

		for (var x = -1*this.sight; x < this.sight; x++) {
			for (var y = -1*this.sight; y < this.sight; y++) {
				for (var z = -1*this.sightHeight; z < this.sightHeight; z++) {
					var pos = newLPos.clone().add(new THREE.Vector3(x, y, 0))
					var index = pos.x+":"+pos.y+":"+pos.z;
					var l = this.cachedLayers[index];
					if (l) {
						if (!l.mesh.surfaceMesh.parent) {
							this.scene.add(l.mesh.surfaceMesh);
						};
						continue;
					}else if(this.fetchingLayers[pos.x+":"+pos.y+":"+pos.z]){
						continue;
					}
					this.fetchLayer(pos)
				};
			};
		};
		this.cleanScene();
		this.layersDirty = false
	},
	cleanCache: function(){

	},
	cleanScene: function(){
		for( var key in this.cachedLayers ){
			var l = this.cachedLayers[key];
			var pos = l.pos.clone();

			var lpos = this.worldToLayerPos(this.cpos);

			var delta = lpos.sub(pos);

			if (delta.x > this.sight || delta.x < -1*this.sight || 
				delta.y > this.sight || delta.y < -1*this.sight ||
				delta.z > this.sightHeight || delta.z < 0) {
				this.scene.remove(l.mesh.surfaceMesh);
			};
		}
	},
	fetchLayer: function(pos){
		this.fetchingLayers = this.fetchingLayers || {};

		var indexStr = pos.x+":"+pos.y+":"+pos.z
		this.fetchingLayers[indexStr] = true;

		var that = this;
		this.api.get("layers?x="+pos.x+"&y="+pos.y+"&layers="+pos.z, "",function(response){
			for(var i = 0; i < response.length; i++){
				var rawLayer = response[i];
				var layer = new Fortia.Layer();
				layer.fromJson(rawLayer);
				layer.generateMesh();
				that.scene.add(layer.mesh.surfaceMesh);
				console.log("Added layer to scene")

				that.cachedLayers[indexStr] = layer;
				that.fetchingLayers[indexStr] = false;
			}
		}, function(e, r){
			that.fetchingLayers[indexStr] = false
			console.log("Error fetching layer", e)
		})
	},
	worldToLayerPos: function(pos){
		var t = pos.clone();
		t.x = Math.floor(t.x / this.worldInfo.layerSize);
		t.y = Math.floor(t.y / this.worldInfo.layerSize);
		return t
	},
	onWheel: function(details){
		var deltaY = details.deltaY;
		this.cameraHeight += deltaY;
		this.camera.position = this.cpos.clone().add(new THREE.Vector3(0,0,this.cameraHeight));
	}
}

