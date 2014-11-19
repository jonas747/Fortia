importScripts("/web/js/game/layer.js", "/web/js/game/chunk.js", "/web/js/libs/voxelmesh.js", "/web/js/libs/greedymesher.js", "/web/js/libs/three.js");

var settings;

var chunkSize;

function log(){
	self.postMessage({
		"action": "log",
		"data": JSON.stringify(arguments),
	});
}

self.onmessage = function(evt){
	var data = evt.data;

	switch (data.action){
		case "init":
			settings = data;
			chunkSize = new THREE.Vector2(data.layerSize, data.height)
			break;
		case "process":
			process(data);
			break;
	}
}

function process(data){
	var jsonStr = data.json;
	var decoded = JSON.parse(jsonStr);
	for (var i = 0; i < decoded.length; i++) {
		processChunk(decoded[i])
	};	
}

function processChunk(rawChunk){
	if (rawChunk == null) {
		return;
	};

	if (rawChunk.IsAir) {
		return
	}
	var chunk = Fortia.chunkFromJson(rawChunk, chunkSize);
 	chunk.generateMesh();

 	var transferOwnership = [];
 	transferOwnership.push(chunk.vertices.buffer);
 	transferOwnership.push(chunk.colors.buffer);
 	transferOwnership.push(chunk.uv.buffer);
 	transferOwnership.push(chunk.normals.buffer);

 	for (var i = 0; i < chunk.layers.length; i++) {
 		if (!chunk.bufferLayers) {chunk.bufferLayers = []};
 		var l = chunk.layers[i]; // Convert to a typed array instead of a object array
 		if (!l) {
 			continue
 		};

 		var voxelBuffer = new Int32Array(l.blocks.length);
 		chunk.bufferLayers[i] = voxelBuffer.buffer;
 		transferOwnership.push(voxelBuffer.buffer);

 		for (var j = 0; j < l.blocks.length; j++) {
 			if(l.blocks[j]){
 				voxelBuffer[j] = l.blocks[j].Id;
 			}
 		};
 	};

 	var out = {
 		action: "finChunk",
 		vertices: chunk.vertices.buffer,
 		colors: chunk.colors.buffer,
 		uv: chunk.uv.buffer,
 		position: {x: chunk.pos.x, y: chunk.pos.y},
 		normals: chunk.normals.buffer,
 		layers: chunk.bufferLayers, 
 	}
 	self.postMessage(out, transferOwnership);
}
