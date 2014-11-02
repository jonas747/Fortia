var Fortia = Fortia || {};
Fortia.Layer = function(pos){
	this.pos = pos || new THREE.Vector3(0, 0);

	this.blocks = [];
	this.voxels = [];

	this.size = Fortia.game.worldInfo.layerSize;
	this.dims = [this.size, this.size, 1];
	
	this.mesh;
	this.geometry;
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
			var curBlock = this.blocks[Fortia.game.coordsToIndex(new THREE.Vector2(k, j))];
			if (!curBlock) {
				this.voxels[num] = 0;
				continue;
			};
			if (curBlock.Id > 0) {
				this.voxels[num] = 0xCC9900;
			}else{
				this.voxels[num] = 0;
			}
		};
	};
	this.mesh = new Mesh(this, GreedyMesh, new THREE.Vector3(1,1,1));
	this.mesh.createSurfaceMesh();
	this.mesh.setPosition(this.pos.x*this.size, this.pos.y*this.size, this.pos.z);
	console.log(this)
}



