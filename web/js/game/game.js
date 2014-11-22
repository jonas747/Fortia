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
	sight: 4,
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
		Fortia.game.resize();

		this.cpos = new THREE.Vector3(0,0,100)

		this.addListeners();
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

		var geometry = new THREE.BoxGeometry( 1.1, 1.1, 1.1 );
		var material = new THREE.MeshBasicMaterial( { color: 0x0ff0000, wireframe: true } );
		var cube = new THREE.Mesh( geometry, material );
		scene.add( cube );
		this.helperCube = cube;

		this.scene = scene;
		this.camera = camera;
		this.renderer = renderer;
	},
	initKeybinds: function(){
	},
	addListeners: function(){
		this.boundWheelCB = this.onWheel.bind(this);
		window.addWheelListener(window, this.boundWheelCB);

		this.boundResize = this.resize.bind(this);
		$(window).on("resize", this.boundResize)
	},
	removeListeners: function(){
		window.removeWheelListener(this.boundWheelCB);
		$(window).off("resize", this.boundResize)
	},
	start: function(){
		console.log("Starting the game")
		var that = this;
		this.api.get("info", "", function(response){
			that.worldInfo = response;
			console.log("Response", response)
			that.running = true;


			that.world = new Fortia.World(response);

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
	tickChunksFetch: 0,
	//tickCleanCache: 0,
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

		// Fetch chunks every 10th tick (140 ish millisecond)
		if (this.tickChunksFetch >= 10) {
			this.world.fetchChunks();
			this.tickChunksFetch = 0
		};
		this.tickChunksFetch++

		if (this.vertexCountDirty) {
			this.updateVertexCount();
			this.vertexCountDirty = false;
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

		this.canvasWidth = cwidth;
		this.canvasHeight = cheight;

		this.view.canvasWidth = cwidth;
		this.view.canvasHeight = cheight;
		// $("#game-canvas")[0].width = cwidth;
		// $("#game-canvas")[0].height = cheight;

		this.renderer.setSize(cwidth, cheight);

		var aspect = cwidth / cheight;
		this.camera.aspect = aspect;
		this.camera.updateProjectionMatrix();
	},
	updateHelperHelperBlock: function(chunk){
		var mv = new THREE.Vector3();
		mv.x = 2 * (Mouse.getX() / this.canvasWidth) - 1;
		mv.y = 1 - 2 * (Mouse.getY() / this.canvasHeight);

		mv = mv.unproject(this.camera);
		mv.sub(this.camera.position);
		mv.normalize();
		
		var cpos = this.camera.position
		var origin = [cpos.x, cpos.y, cpos.z];
		var direction = [mv.x, mv.y, mv.z];

		var hit_pos = [];
		var hit_norm = [];
		var hit = traceRay(this.world, origin, direction, 100, hit_pos, hit_norm);
		if (hit) {
			//console.log(hit_pos, hit_norm);
			var pos = new THREE.Vector3(Math.floor(hit_pos[0])+0.5, Math.floor(hit_pos[1])+0.5, Math.floor(hit_pos[2])+0.5);
	 		this.helperCube.position.set(pos.x, pos.y, pos.z);
	 		this.onHelperMove(pos, hit);
		};
	},
	onHelperMove: function(pos, id){
		var cloned = pos.clone();
		var lPos = this.worldToLayerPos(pos);
		$("#game-cursor").text("Cursor WPos: "+cloned.x + ", "+ cloned.y + ", " + cloned.z + " LPos: " + lPos.x + ", " + lPos.y + " TypeId: "+id);
	},
	moveCamera: function(by){

		var newPos = this.camera.position.clone();
		newPos.add(by);

		// if we moved to a new layer/chunk we fetch/update the cache
		var oldLPos = this.worldToLayerPos(this.camera.position)
		var newLPos = this.worldToLayerPos(newPos)

		if (newLPos.x !== oldLPos.x || newLPos.y !== oldLPos.y) {
			this.world.updateChunks(newPos);
		};
		
		this.camera.position.copy(newPos);
		this.updateStatusCoords();		
	},
	updateStatusCoords: function(){
		var posString = "X: " + this.camera.position.x 
		posString += ", Y: " + this.camera.position.y;
		posString += ", Z: " + this.camera.position.z; 

		$("#game-position").text(posString);
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
		this.updateStatusCoords();		
	},
	updateVertexCount: function(){
		function getVertexCount(obj){
			var vertices = 0
			if (obj.children.length > 0) {
				for (var i = 0; i < obj.children.length; i++) {
					vertices += getVertexCount(obj.children[i])
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
		var numVertices = getVertexCount(this.scene);
		$("#game-vertices").text("Vertex count: " + numVertices);
	}
}

