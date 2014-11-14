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
	sight: 3,
	//sightHeight: 25,
	chunkFetchQueue: {},
	fetchingChunks: {},
	cachedChunks: {},
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
		$("#main-container").append(stats.domElement);
		this.stats = stats;

		this.initScene();
		$(window).resize(Fortia.game.resize)
		Fortia.game.resize();

		this.cpos = new THREE.Vector3(0,0,100)


		var boundWheelCB = this.onWheel.bind(this);
		window.addWheelListener(window, boundWheelCB);

		// Run the worker if it hasnt started yet
		if (!this.workerRunning) {
			this.initWorker();
		};
	},
	initScene: function(){
		this.blkMaterial = new THREE.MeshLambertMaterial({
			vertexColors: THREE.VertexColors,
		});
		
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
		var material = new THREE.MeshBasicMaterial( { color: 0x0ff0000, wireframe: true } );
		var cube = new THREE.Mesh( geometry, this.blkMaterial );
		scene.add( cube );
		this.helperCube = cube;

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
				case "finChunk":
					that.addChunk(data);
					break;
			}
		}
		this.workerRunning =  true;

		this.worker = worker;
	},
	addChunk: function(data){
		var pos = new THREE.Vector2(data.position.x, data.position.y)
		var chunk = new Fortia.Chunk();
		chunk.vertices = data.vertices;
		chunk.colors = data.colors;
		chunk.uv = data.uv;

		chunk.pos = pos;
		chunk.size = new THREE.Vector2(this.worldInfo.LayerSize, this.worldInfo.Height);
		
		var index = pos.x + ":"  + pos.y;
		delete this.fetchingChunks[index];
		this.cachedChunks[index] = chunk;

		//chunk.createWireMesh(0x00ff00);
		chunk.createSurfaceMesh(this.blkMaterial);
		chunk.addToScene(this.scene)
		this.updateVerticeCount();
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
				layerSize: that.worldInfo.LayerSize,
				height: that.worldInfo.Height,
			});
			that.moveCamera(new THREE.Vector3(0,0,0));
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

		if (Mouse.isDirty()) {
			Mouse.setDirty(false);
			this.updateHelperHelperBlock();
		};

		if (move) {
			this.moveCamera(moveBy);
		};
		// Update layers every 10th tick (140 ish millisecond)
		if (this.tickLayersUpdate >= 10) {
			this.fetchChunks();
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

		Fortia.game.canvasWidth = cwidth;
		Fortia.game.canvasHeight = cheight;

		Fortia.gameView.canvasWidth = cwidth;
		Fortia.gameView.canvasHeight = cheight;
		// $("#game-canvas")[0].width = cwidth;
		// $("#game-canvas")[0].height = cheight;

		Fortia.game.renderer.setSize(cwidth, cheight);

		var aspect = cwidth / cheight;
		Fortia.game.camera.aspect = aspect;
		Fortia.game.camera.updateProjectionMatrix();
	},
	updateHelperHelperBlock: function(chunk){
		var mv = new THREE.Vector3();
		mv.x = 2 * (Mouse.getX() / this.canvasWidth) - 1;
		mv.y = 1 - 2 * (Mouse.getY() / this.canvasHeight);

		mv = mv.unproject(this.camera);
		mv.sub(this.camera.position);
		mv.normalize();
		var raycaster = new THREE.Raycaster( this.camera.position, mv, 0, 100);
		var intersects = raycaster.intersectObjects(this.scene.children);
		// Change color if hit block
		if ( intersects.length > 0 ) {
			var inter = intersects[0]
			if (inter.object === this.helperCube) {
				if (intersects.length > 1) {
					inter = intersects[1];
				}else{
					inter = null
				}
			}

			if (inter) {
				pos = inter.point.clone();
				pos.floor();
				pos.add(new THREE.Vector3(0.5, 0.5, 0));
				this.helperCube.position.set(pos.x, pos.y, pos.z);
			};
		}
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

		if (newLPos.x !== oldLPos.x || newLPos.y !== oldLPos.y) {
			console.log("New chunk pos! ", oldLPos, newLPos);
			this.updateChunks();
		};
		
		this.camera.position.copy(newPos);
		this.camera.position.z += Fortia.game.cameraHeight;
		this.cpos = newPos

		var posString = "X: " + newPos.x 
		posString += ", Y: " + newPos.y;
		posString += ", Z: " + newPos.z; 

		$("#game-position").text(posString);
	},
	updateChunks: function(){
		if (!this.chunksDirty) {
			//return;
		};

		var newLPos = this.worldToLayerPos(this.cpos)
		for (var x = -1*this.sight; x < this.sight; x++) {
			for (var y = -1*this.sight; y < this.sight; y++) {
				var pos = newLPos.clone().add(new THREE.Vector3(x, y, 0))
				var index = pos.x+":"+pos.y
				var chunk = this.cachedChunks[index];
				if (chunk) {
					if (!chunk.addedToScene) {
						chunk.addToScene(this.scene);
					};
					continue;
				}else if(this.fetchingChunks[index]){
					continue;
				}else if(this.chunkFetchQueue[index]){
					continue;
				}
				this.chunkFetchQueue[index] = pos;
			};
		};
		this.cleanScene();
		this.chunksDirty = false
	},
	fetchChunks: function(){
		var num = 0;
		
		var xList = "";
		var yList = "";

		for(var key in this.chunkFetchQueue){
			var pos = this.chunkFetchQueue[key];
			var sep = "";
			if (num != 0) {
				sep = ",";
			};
			xList += sep + pos.x;
			yList += sep + pos.y;
			this.fetchingChunks[key] = true;			

			num++;
		}
		
		if (num < 1) {
			return;
		};

		this.chunkFetchQueue = {};

		var that = this;
		this.api.dType = "text"; // Set the data type to text temporarily, since were decoding it in the background worker
		this.api.get("chunks?x="+xList+"&y="+yList, "",function(response){
			that.processChunkResponse(response)
		}, function(e, r){
			//that.fetchingLayers[indexStr] = false
			console.log("Error fetching chunk", e)
		});

		this.api.dType = "json"; // Set it back to json
	},
	processChunkResponse: function(response){
		this.worker.postMessage({
			action: "process",
			json: response,
		});
	},
	cleanCache: function(){
		for( var key in this.cachedChunks ){
			var chunk = this.cachedChunks[key];
			var pos = chunk.pos.clone();

			var lpos = this.worldToLayerPos(this.cpos);

			var delta = lpos.sub(pos);

			if (delta.x > this.sight+1 || delta.x < -1*this.sight-1 || 
				delta.y > this.sight+1 || delta.y < -1*this.sight-1){
				chunk.removeFromScene(this.scene);
				chunk.dispose();
				delete this.cachedChunks[key];
				chunk = null;
				console.log("cleared a chunk from the cache");
			};
		}
		this.updateVerticeCount();
	},
	cleanScene: function(){
		for( var key in this.cachedChunks ){
			var chunk = this.cachedChunks[key];
			var pos = chunk.pos.clone();

			var lpos = this.worldToLayerPos(this.cpos);

			var delta = lpos.sub(pos);

			if (delta.x > this.sight || delta.x < -1*this.sight || 
				delta.y > this.sight || delta.y < -1*this.sight){
				chunk.removeFromScene(this.scene)
			};
		}
		this.updateVerticeCount();
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
		// This seems to be much higher on chrome than firefox
		var deltaY = details.deltaY / (navigator.userAgent.toLowerCase().indexOf('chrome') > -1 ? 5 : 1);
		this.cameraHeight += deltaY;
		this.camera.position.z += deltaY
	},
	updateVerticeCount: function(){
		function getVerticesCount(obj){
			var vertices = 0
			if (obj.children.length > 0) {
				for (var i = 0; i < obj.children.length; i++) {
					vertices += getVerticesCount(obj.children[i])
				};
			};
			if (obj.geometry) {
				var geom = obj.geometry;
				if (geom.vertices) {
					vertices += geom.vertices.length;
				}else{
					attrib = geom.getAttribute("position")
					if (attrib) {
						vertices += (attrib.length / 3);
					};
				}
			};
			return vertices;
		}
		var numVertices = getVerticesCount(this.scene);
		$("#game-vertices").text("Vertex count: " + numVertices);
	}
}

