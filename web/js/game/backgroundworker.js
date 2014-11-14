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

 	var out = {
 		action: "finChunk",
 		vertices: chunk.vertices.buffer,
 		colors: chunk.colors.buffer,
 		uv: chunk.uv.buffer,
 		position: {x: chunk.pos.x, y: chunk.pos.y},
 	}
 	self.postMessage(out, [chunk.vertices.buffer, chunk.colors.buffer, chunk.uv.buffer]);
 	log("Done generating chunk")
}
