var Fortia = Fortia || {};
Fortia.Layer = function(pos, size){
	this.pos = pos || new THREE.Vector3(0, 0);

	this.blocks = [];
	this.voxels = [];

	this.size = size || (Fortia.game !== undefined ? Fortia.game.worldInfo.LayerSize : 1);
	this.dims = [this.size, this.size, 1];
	
	this.mesh;
	this.geometry;

	this.vertices;
	this.colors;
}

/*
	World    *World `json:"-"`
	Position vec.Vec3I
	Blocks   []*Block
	Flags    int
	IsAir    bool // True if this layer is just air

*/
Fortia.Layer.prototype.fromJson = function(json){
	this.pos = new THREE.Vector3(json.Position.X, json.Position.Y, json.Position.Z);
	this.blocks = json.Blocks;
}

Fortia.Layer.prototype.generateMesh = function(){
	var num = 0;
	for (var k = 0; k < this.size; k++) {
		for (var j = 0; j < this.size; j++, num++) {
			var curBlock = this.blocks[this.coordsToIndex(new THREE.Vector2(j, k))];
			if (!curBlock) {
				this.voxels[num] = 0;
				continue;
			};
			if (curBlock.Id > 0) {
				var color = 0x0000ff
				switch (curBlock.Id){
					case 1: // stone 
						//color = 0x555555;
						var color1c = (Math.random() * 0x05)+0x20 
						color = (color1c) | (color1c << 8) | (color1c << 16);
						break;
					case 2:
						var color1c = ((Math.random() * 0x08) + 0x70) << 8
						color = 0x220022 | color1c
						break;
				}
				// if (curBlock.Flags && curBlock.Flags & 1) { // fully covered
				// 	color = 0xbbbbbb;
				// };
				this.voxels[num] = color;
			}else{
				this.voxels[num] = 0;
			}
		};
	};
	this.vMesh = new VoxelMesh()
	this.vMesh.mesh(this, GreedyMesh, new THREE.Vector3(1,1,1));
}

Fortia.Layer.prototype.createGeometry = function(){
	var geometry = this.geometry = new THREE.BufferGeometry();

	geometry.addAttribute("position", new THREE.BufferAttribute(new Float32Array(this.vertices), 3));
	geometry.addAttribute("color", new THREE.BufferAttribute(new Float32Array(this.colors), 3));
	geometry.addAttribute("uv", new THREE.BufferAttribute(new Float32Array(this.uv), 2));
	geometry.computeFaceNormals();
	geometry.computeVertexNormals();
	return geometry;
}

Fortia.Layer.prototype.createSurfaceMesh = function(material) {
	if (!this.geometry) {this.createGeometry()};
	material = material || new THREE.MeshNormalMaterial()
	var surfaceMesh  = new THREE.Mesh( this.geometry, material )
	surfaceMesh.scale = new THREE.Vector3(1,1,1);
	surfaceMesh.doubleSided = false;
	var pos = new THREE.Vector3(this.pos.x * this.size, this.pos.y * this.size, this.pos.z);
	surfaceMesh.position.x = pos.x;
	surfaceMesh.position.y = pos.y;
	surfaceMesh.position.z = pos.z;
	this.surfaceMesh = surfaceMesh;
	return surfaceMesh;
}

Fortia.Layer.prototype.createWireMesh = function(hexColor) {    
	if (!this.geometry) {this.createGeometry()};
	var wireMaterial = new THREE.MeshBasicMaterial({
		color : hexColor || 0xffffff,
		wireframe : true
	})
	wireMesh = new THREE.Mesh(this.geometry, wireMaterial)
	wireMesh.scale = new THREE.Vector3(1,1,1);
	wireMesh.doubleSided = false
	var pos = new THREE.Vector3(this.pos.x * this.size, this.pos.y * this.size, this.pos.z);
	wireMesh.position.x = pos.x;
	wireMesh.position.y = pos.y;
	wireMesh.position.z = pos.z;
	this.wireMesh = wireMesh;
	return wireMesh;
}

Fortia.Layer.prototype.coordsToIndex =  function(pos) {
	return this.size*pos.x + pos.y;
}



