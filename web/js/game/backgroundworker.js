importScripts("/web/js/game/layer.js", "/web/js/libs/voxelmesh.js", "/web/js/libs/greedymesher.js", "/web/js/libs/three.js");

var settings;
	
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
			break;
		case "process":
			process(data);
			break;
	}
}

function process(data){
	log("Processing response now!")
	var jsonStr = data.json;
	var decoded = JSON.parse(jsonStr);

	for(var i = 0; i < decoded.length; i++){
		var rawLayer = decoded[i];
		if (rawLayer == null) {
			continue;
		};

		if (rawLayer.IsAir) {
			continue
		}

		var layer = new Fortia.Layer({x:0,y:0,z:0}, settings.size);
	 	layer.fromJson(rawLayer);
	 	layer.generateMesh();

	 	var out = {
	 		action: "finlayer",
	 		vertices: layer.vMesh.result.vertices.buffer,
	 		colors: layer.vMesh.result.colors.buffer,
	 		uv: layer.vMesh.uv.buffer,
	 		position: {x: layer.pos.x, y: layer.pos.y, z: layer.pos.z},
	 		//blocks: layer.blocks
	 	}
	 	self.postMessage(out, [layer.vMesh.result.vertices.buffer, layer.vMesh.result.colors.buffer, layer.vMesh.uv.buffer]);
	 	// todo
	}
}
