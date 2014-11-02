var Fortia = Fortia || {}

Fortia.camera = {
	pos: {x: 0, y: 0, z: 0},
	layersVisible: 5,

	tileSize: 32,

	cachedLayers: [],
	layerSize: 50,

	canvas: null,
	ctx: null,

	colors: ["#fff", "#af5", "#2e9"],

	render3: function(){
		
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
			x: Math.floor(oldPos.x / this.layerSize),
			y: Math.floor(oldPos.y / this.layerSize),
		}

		var newSection = {
			x: Math.floor(this.pos.x / this.layerSize),
			y: Math.floor(this.pos.y / this.layerSize),
		}

		// If we moved to a new section or we moved up/down
		if (newSection.x !== oldSection.x || newSection.y !== oldSection.y || oldPos.z !== this.pos.z) {
			// Fetch the new layer(s)
			var layers = this.pos.z + "";
			var that = this;
			Fortia.game.api.get("layers?x="+newSection.x+"&y="+newSection.y+"&layers="+layers,"", function(response){
				for(var i = 0; i < response.length; i++){
					// Check if its allready in the cachhe
					var setLayer = false;
					var newLayer = response[i];
					for(var j = 0; j < that.cachedLayers.length; j++){
						var oldLayer = that.cachedLayers[j];
						//console.log(newLayer.Position.X === oldLayer.Position.X)
						if (newLayer.Position.X === oldLayer.Position.X && newLayer.Position.Y === oldLayer.Position.Y && newLayer.Position.Z === oldLayer.Position.Z) {
							that.cachedLayers[j] = newLayer;
							setLayer = true;
							break;
						};
					}
					if (!setLayer) {
						that.cachedLayers.push(newLayer);
					};
				};
			})
		};
		console.log("New position: ", this.pos, "LPos: ", newSection)
	},
}