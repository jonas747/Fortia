Fortia = Fortia || {};
Fortia.Chunk = function(){
	this.layers = [];
	this.size = 0;
	this.pos;

	this.vertices;
	this.colors;
	this.uv;

	this.geometry;
	this.surfaceMesh;
	this.wireMesh;
}

Fortia.chunkFromJson = function(json, size){
	var layers = [];
	var pos = new THREE.Vector2(json.Position.X, json.Position.Y);
	for (var i = 0; i < json.Layers.length; i++) {
		var l = json.Layers[i];
		if (l === null) {continue};
		var layer = new Fortia.Layer(new THREE.Vector3(pos.x, pos.y, l.Position.Z), size.x);

		layer.fromJson(l);
		layers[l.Position.Z] = layer
	};
	var chunk = new Fortia.Chunk();
	chunk.layers = layers;
	chunk.pos = pos;
	chunk.size = size
	return chunk;
}

Fortia.Chunk.prototype.createVoxelData = function(){
	var voxels = new Uint32Array(this.size.x * this.size.x * this.size.y);
	var num = 0;
	for (var z = 0; z < this.size.y; z++) {
		for (var x = 0; x <this.size.x; x++) {
			for (var y = 0; y < this.size.x; y++, num++) {
				var l = this.layers[z]
				if (l == undefined) {
					voxels[num] = 0;
					continue
				};
				var curBlock = l.blocks[l.coordsToIndex(new THREE.Vector2(y, x))];
				if (!curBlock) {
					voxels[num] = 0;
					continue;
				};
				if (curBlock.Flags && curBlock.Flags & 8) { // fully covered
					voxels[num] = 0;
					continue
				};
				if (curBlock.Id > 0) {
					var color = 0x0000ff
					switch (curBlock.Id){
						case 1: // stone 
							var color1c = (Math.random() * 0x05)+0x20 
							color = (color1c) | (color1c << 8) | (color1c << 16);
							break;
						case 2:
							var color1c = ((Math.random() * 0x08) + 0x70) << 8
							color = 0x220022 | color1c
							break;
					}
					voxels[num] = color;
				}else{
					voxels[num] = 0;
				}
			};
		};
	};
	this.voxels = voxels;
	return voxels;
}

Fortia.Chunk.prototype.generateMeshOld = function(){
	// Convert it to an array of colors
	var dims = [this.size.x, this.size.x, this.size.y];
	var voxels = this.createVoxelData();

	var data = {dims: dims, voxels: voxels};
	var vmesh = new VoxelMesh();
	vmesh.mesh(data, GreedyMesh, new THREE.Vector3(1,1,1));
	this.vertices = vmesh.result.vertices;
	this.colors = vmesh.result.colors;
	this.uv = vmesh.result.uv;
}

Fortia.Chunk.prototype.generateMesh = function(){
	var dims = [this.size.x, this.size.x, this.size.y];
	var voxels = this.createVoxelData();

	var data = {dims: dims, voxels: voxels};
	var vmesh = new VoxelMesh(data, GreedyMesh, new THREE.Vector3(1,1,1));
		
	var bufferGeometry = new THREE.BufferGeometry();
	bufferGeometry.fromGeometry(vmesh.geometry, {"vertexColors": THREE.FaceColors});
	//log(bufferGeometry.attributesKeys);
	this.vertices = bufferGeometry.getAttribute("position").array;
	this.colors = bufferGeometry.getAttribute("color").array;
	this.uv = bufferGeometry.getAttribute("uv").array;
}

Fortia.Chunk.prototype.createGeometry = function(){
	var geometry = new THREE.BufferGeometry();
	geometry.addAttribute("position", new THREE.BufferAttribute(new Float32Array(this.vertices), 3));
	geometry.addAttribute("color", new THREE.BufferAttribute(new Float32Array(this.colors), 3));
	geometry.addAttribute("uv", new THREE.BufferAttribute(new Float32Array(this.uv), 2));
	//geometry.computeFaceNormals();
	geometry.computeVertexNormals();
	this.geometry = geometry
	return geometry;
}
Fortia.Chunk.prototype.createSurfaceMesh = function(material){
	if (!this.geometry) {this.createGeometry()};
	material = material || new THREE.MeshNormalMaterial()
	var surfaceMesh  = new THREE.Mesh( this.geometry, material )
	surfaceMesh.scale = new THREE.Vector3(1,1,1);
	surfaceMesh.doubleSided = false;
	surfaceMesh.position.set(this.pos.x * this.size.x, this.pos.y * this.size.x, 0)
	this.surfaceMesh = surfaceMesh;
	return surfaceMesh;
}

Fortia.Chunk.prototype.createWireMesh = function(hexColor){
	if (!this.geometry) {this.createGeometry()};
	var wireMaterial = new THREE.MeshBasicMaterial({
		color : hexColor || 0xffffff,
		wireframe : true
	});
	var wireMesh  = new THREE.Mesh( this.geometry, wireMaterial )
	wireMesh.scale = new THREE.Vector3(1,1,1);
	wireMesh.doubleSided = true;
	wireMesh.position.set(this.pos.x * this.size.x, this.pos.y * this.size.x, 0)
	this.wireMesh = wireMesh;
	return wireMesh;
}

Fortia.Chunk.prototype.addToScene = function(scene){
	this.addedToScene = true;
	if (this.surfaceMesh) {
		scene.add(this.surfaceMesh);
	};
	if (this.wireMesh) {
		scene.add(this.wireMesh);
	};
}

Fortia.Chunk.prototype.removeFromScene = function(scene){
	this.addedToScene = false
	console.log("removing from scene")
	if (this.surfaceMesh) {
		scene.remove(this.surfaceMesh);
	};
	if (this.wireMesh) {
		scene.remove(this.wireMesh);
	};
}

// Clear the buffers on the gpu
Fortia.Chunk.prototype.dispose = function(){
	this.geometry.dispose();
}