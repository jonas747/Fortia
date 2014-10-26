var Fortia = Fortia || {}

Fortia.camera2 = {
	pos: {x: 0, y: 0, z: 0},
	layersVisible: 5,

	cachedLayers: [],

	render: function(){

	},

	// Move the camera by amount(vector3)
	// If ignoreCache is true it will fetch the layers even if they are in the cache
	move: function(amount, ignoreCache){
		// Check if we moved to a new section
		var oldPos = {x: this.pos.x, y: this.pos.y, z: this.pos.z};

		this.pos.x += amount.x
		this.pos.y += amount.y
		this.pos.z += amount.z

		var oldSection = {
			x: Math.floor(oldPos.x / 10),
			y: Math.floor(oldPos.y / 10),
		}

		var newSection = {
			x: Math.floor(this.pos.x / 10),
			y: Math.floor(this.pos.y / 10),
		}

		// If we moved to a new section or we moved up/down
		if (newSection !== oldSection || oldPos.z !== this.pos.z) {
			// Fetch the new layer(s)
			var layers = this.pos.z + "";
			var that = this;
			Fortia.game.api.get("layers?x="+newSection.x+"&y="+newSection.y+"&layers="+layers,"", function(response){
				// TODO put in cache
			})
		};
	},
}