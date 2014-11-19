var Fortia = Fortia || {};

Fortia.World = function(settings){
	this.sight = 2;
	// ACtual cache size: sight + cacheSize 
	this.cacheSize = 1;

	this.api = Fortia.game.api;
	this.scene = Fortia.game.scene;

	this.chunks = {};
	this.fetchingChunks = {};
	this.chunkFetchQueue = {};
	
	this.settings = settings

	// Init the mesher
	var worker = new Worker("/web/js/game/backgroundworker.js");
	worker.onmessage = this.workerCallback.bind(this);
	this.workerRunning =  true;
	this.mesher = worker;
	this.mesher.postMessage({
		action: "init",
		layerSize: this.settings.LayerSize,
		height: this.settings.Height,
	});
}

Fortia.World.prototype.workerCallback = function(evt){
	var data = evt.data;
	switch (data.action){
		case "log":
			console.log(JSON.stringify(data.data));
			break;
		case "finChunk":
			var pos = new THREE.Vector2(data.position.x, data.position.y)
			var chunk = new Fortia.Chunk();
			chunk.vertices = data.vertices;
			chunk.colors = data.colors;
			chunk.uv = data.uv;
			chunk.normals = data.normals;
			if (data.blocks) {chunk.blocks = data.blocks};

			chunk.pos = pos;
			chunk.size = new THREE.Vector2(this.settings.LayerSize, this.settings.Height);
				
			this.addChunk(chunk);

			var index = pos.x + ":"  + pos.y;
			delete this.fetchingChunks[index];
			break;
	}

}

Fortia.World.prototype.addChunk = function(chunk){
	var index = chunk.pos.x + ":"  + chunk.pos.y;
	this.chunks[index] = chunk;
	//chunk.createWireMesh(0x00ff00);
	chunk.createSurfaceMesh(Fortia.game.blkMaterial);
	chunk.addToScene(this.scene);
	
	//Fortia.game.updateVerticeCount();
	Fortia.game.vertexCountDirty = true;
}

Fortia.World.prototype.removeChunk = function(x, y){
	var index = x + ":" + y;
	var chunk = this.chunks[index];
	if(chunk){
		chunk.removeFromScene(this.scene);
		chunk.dispose();
		delete chunk[index];
	}
	Fortia.game.vertexCountDirty = true;
}

Fortia.World.prototype.getBlock = function(x, y, z){

}

Fortia.World.prototype.updateChunks = function(newPos){
	var newLPos = this.worldToLayerPos(newPos)
	for (var x = -1*this.sight; x < this.sight; x++) {
		for (var y = -1*this.sight; y < this.sight; y++) {
			var pos = newLPos.clone().add(new THREE.Vector3(x, y, 0))
			var index = pos.x+":"+pos.y
			var chunk = this.chunks[index];
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
	this.cameraPos = newPos;
	this.clean();
}

Fortia.World.prototype.fetchChunks = function(){
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
	this.api.dType = "text"; // Set the data type to text temporarily, since were decoding it in the background worker, maybe use arraybuffer in future for that extra performance
	this.api.get("chunks?x="+xList+"&y="+yList, "",function(response){
		that.mesher.postMessage({
			action: "process",
			json: response,
		});
	}, function(e, r){
		//that.fetchingLayers[indexStr] = false
		console.log("Error fetching chunk", e)
	});

	this.api.dType = "json"; // Set it back to json
},

Fortia.World.prototype.clean = function(){
	var cameraLayerPos = this.worldToLayerPos(this.cameraPos);
	for( var key in this.chunks ){
		
		var chunk = this.chunks[key];
		var pos = chunk.pos.clone();
		var delta = cameraLayerPos.clone().sub(pos);

		// Remove entirely from cache if 
		if (delta.x > this.sight+this.cacheSize || delta.x < -1*this.sight-this.cacheSize || 
			delta.y > this.sight+this.cacheSize || delta.y < -1*this.sight-this.cacheSize){
			this.removeChunk(chunk.pos.x, chunk.pos.y);
		}else if (delta.x > this.sight || delta.x < -1*this.sight || 
				delta.y > this.sight || delta.y < -1*this.sight) {
			chunk.removeFromScene(this.scene)
			Fortia.game.vertexCountDirty = true;
		};
	}
}

// 
Fortia.World.prototype.blockCoordsToIndex = function(pos) {
	return this.settings.LayerSize*pos.x + pos.y
}

// Return a blocks x and y from the index in the layer slice
// x = index / size
// y = index - (x * size)
Fortia.World.prototype.blockIndexToCoords = function(index) {
	var x = index / this.settings.LayerSize
	var y = index - (x * this.settings.LayerSize)
	return new THREE.Vector2(x, y)
}

Fortia.World.prototype.worldToLayerPos = function(pos){
	var t = pos.clone();
	t.x = Math.floor(t.x / this.settings.LayerSize);
	t.y = Math.floor(t.y / this.settings.LayerSize);
	return t
}
Fortia.World.prototype.layerToWorldPos = function(pos){
	var t = pos.clone();
	t.multiplyScalar(this.settings.LayerSize);
	return t;
}