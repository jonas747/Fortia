var Fortia = Fortia || {};
Fortia.Layer = function(pos, size){
	this.pos = pos || new THREE.Vector3(0, 0);

	this.blocks = [];
	this.voxels = [];

	this.size = size || (Fortia.game !== undefined ? Fortia.game.worldInfo.LayerSize : 1);
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

Fortia.Layer.prototype.coordsToIndex =  function(pos) {
	return this.size*pos.x + pos.y;
}



