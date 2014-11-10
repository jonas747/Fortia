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
	cameraHeight: 50,
	sight: 5,
	sightHeight: 25,
	layerFetchQueue: {},
	fetchingLayers: {},
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

		this.cpos = new THREE.Vector3(0,0,100)

		this.blkMaterial = new THREE.MeshLambertMaterial({
			 vertexColors: THREE.VertexColors,
		});

		var boundWheelCB = this.onWheel.bind(this);
		window.addWheelListener($("#main-container")[0], boundWheelCB);

		// Run the worker if it hasnt started yet
		if (!this.workerRunning) {
			this.initWorker();
		};
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

		var light = new THREE.HemisphereLight( 0xffffff, 0xff00ff, 1 );
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
	initWorker: function(){
		var worker = new Worker("/web/js/game/backgroundworker.js");
		var that = this;
		worker.onmessage = function(evt){
			var data = evt.data;
			switch (data.action){
				case "log":
					//var list = data.data;
					console.log(JSON.stringify(data.data));
					break;
				case "finlayer":
					that.addLayer(data);
					break;
			}
		}
		this.workerRunning =  true;

		this.worker = worker;
	},
	addLayer: function(data){
		//console.log(data)
		var pos = new THREE.Vector3(data.position.x, data.position.y, data.position.z)
		var layer = new Fortia.Layer(pos);
		//layer.blocks = data.blocks;
		layer.vertices = data.vertices;
		layer.colors = data.colors;
		layer.uv = data.uv;
		var surfaceMesh = layer.createSurfaceMesh(this.blkMaterial);
		this.scene.add(surfaceMesh);

		//var wiremesh = layer.createWireMesh(0x00ff00);
		//this.scene.add(wiremesh);

		var index = pos.x + ":"  + pos.y + ":" +pos.z;
		delete this.fetchingLayers[index];
		this.cachedLayers[index] = layer;
	},
	start: function(){
		console.log("Starting the game")
		var that = this;
		this.api.get("info", "", function(response){
			that.worldInfo = response;
			console.log("Response", response)
			that.running = true;

			that.worker.postMessage({
				action: "init",
				size: that.worldInfo.LayerSize,
			});

			that.loop();
		});
		
		this.loop();
	},
	stop: function(){
		this.running = false;
	},
	render3: function(){
		this.renderer.render(this.scene, this.camera);
	}, // 3d renderer
	tickLayersUpdate: 0,
	tickCleanCache: 0,
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
		// Update layers every 10th tick (140 ish millisecond)
		if (this.tickLayersUpdate >= 10) {
			this.fetchLayers();
			this.tickLayersUpdate = 0
		};
		this.tickLayersUpdate++

		if (this.tickCleanCache >= 13) {
			this.cleanCache();
			this.tickCleanCache = 0
		};
		this.tickCleanCache++
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
		return this.worldInfo.LayerSize*pos.x + pos.y
	},

	// Return a blocks x and y from the index in the layer slice
	// x = index / size
	// y = index - (x * size)
	indexToCoords: function(index) {
		var x = index / this.worldInfo.LayerSize
		var y = index - (x * this.worldInfo.LayerSize)
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
			console.log("New layer pos! ", oldLPos, newLPos);
			this.updateLayers();
		};
		
		this.camera.position.copy(newPos);
		this.camera.position.z += Fortia.game.cameraHeight;
		this.cpos = newPos

		var posString = "X: " + newPos.x 
		posString += ", Y: " + newPos.y;
		posString += ", Z: " + newPos.z; 

		$("#game-position").text(posString);
	},
	updateLayers: function(){
		if (!this.layersDirty) {
			//return;
		};

		this.cachedLayers = this.cachedLayers || {};
		this.fetchingLayers = this.fetchingLayers || {};

		var newLPos = this.worldToLayerPos(this.cpos)

		for (var x = -1*this.sight; x < this.sight; x++) {
			for (var y = -1*this.sight; y < this.sight; y++) {
				for (var z = -1 * this.sightHeight; z < 0; z++) {
					var pos = newLPos.clone().add(new THREE.Vector3(x, y, z))
					var index = pos.x+":"+pos.y+":"+pos.z;
					var l = this.cachedLayers[index];
					if (l) {
						if (!l.surfaceMesh.parent) {
							this.scene.add(l.surfaceMesh);
						};
						continue;
					}else if(this.fetchingLayers[index]){
						continue;
					}else if(this.layerFetchQueue[index]){
						continue;
					}
					this.layerFetchQueue[index] = pos;
				};
			};
		};
		this.cleanScene();
		this.layersDirty = false
	},
	fetchLayers: function(){
		var num = 0;
		
		var xList = "";
		var yList = "";
		var zList = "";

		for(var key in this.layerFetchQueue){
			var pos = this.layerFetchQueue[key];
			var sep = "";
			if (num != 0) {
				sep = ",";
			};
			xList += sep + pos.x;
			yList += sep + pos.y;
			zList += sep + pos.z;
			this.fetchingLayers[key] = true;			

			num++;
		}
		
		if (num < 1) {
			return;
		};

		this.layerFetchQueue = {};

		var that = this;
		this.api.dType = "text"; // Set the data type to text temporarily
		this.api.get("layers?x="+xList+"&y="+yList+"&z="+zList, "",function(response){
			that.processLayerReply(response)
		}, function(e, r){
			//that.fetchingLayers[indexStr] = false
			console.log("Error fetching layer", e)
		});

		this.api.dType = "json"; // Set it back to json
	},
	processLayerReply: function(response){
		this.worker.postMessage({
			action: "process",
			json: response,
		});
	},
	cleanCache: function(){
		for( var key in this.cachedLayers ){
			var l = this.cachedLayers[key];
			var pos = l.pos.clone();

			var lpos = this.worldToLayerPos(this.cpos);

			var delta = lpos.sub(pos);

			if (delta.x > this.sight+1 || delta.x < -1*this.sight-1 || 
				delta.y > this.sight+1 || delta.y < -1*this.sight-1 ||
				delta.z > this.sightHeight+1 || delta.z < -1) {
				this.scene.remove(l.surfaceMesh);
				l.geometry.dispose();
				delete this.cachedLayers[key];
				l = null;
				console.log("a layer from the cache");
			};
		}
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
				this.scene.remove(l.surfaceMesh);
			};
		}
	},
	worldToLayerPos: function(pos){
		var t = pos.clone();
		t.x = Math.floor(t.x / this.worldInfo.LayerSize);
		t.y = Math.floor(t.y / this.worldInfo.LayerSize);
		return t
	},
	layerToWorldPos: function(pos){
		var t = pos.clone();
		t.multiplyScalar(this.worldInfo.LayerSize);
		return t;
	},
	onWheel: function(details){
		var deltaY = details.deltaY;
		this.cameraHeight += deltaY;
		this.camera.position.z += deltaY
	}
}

