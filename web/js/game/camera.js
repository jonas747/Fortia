var Fortia = Fortia || {}

Fortia.camera2 = {
	pos: {x: 0, y: 0, z: 0},
	layersVisible: 5,

	tileSize: 32,

	cachedLayers: [],
	layerSize: 50,

	canvas: null,
	ctx: null,

	colors: ["#fff", "#af5", "#2e9"],

	render: function(){


		if (!this.ctx) {
			this.canvas = $("#game-canvas");
			this.ctx = this.canvas[0].getContext("2d");
		};

		visibleXLayers = (this.tileSize*this.layerSize) / this.canvas.width();
		visibleYLayers = (this.tileSize*this.layerSize) / this.canvas.height();
		Fortia.visibleXLayers = visibleXLayers
		// Find the current layer in the cache
		var curLayerPos = {
			x: Math.floor(this.pos.x / this.layerSize),
			y: Math.floor(this.pos.y / this.layerSize),
			z: this.pos.z,
		};

		// Clear the canvas
		this.ctx.fillStyle = "#fff";
		this.ctx.fillRect(0, 0, this.canvas.width(), this.canvas.height());

		for (var i = 0; i < this.cachedLayers.length; i++) {
			var layer = this.cachedLayers[i]
			if (layer.Position.X < curLayerPos.x + visibleXLayers && layer.Position.X >= curLayerPos.x &&
				layer.Position.Y < curLayerPos.y + visibleYLayers && layer.Position.Y >= curLayerPos.y && 
				layer.Position.Z === curLayerPos.z) {
				this.renderLayer(layer);
			};
		};
	},

	renderLayer: function(layer){
		var layerWorldPos = {
			x: layer.Position.X * this.layerSize,
			y: layer.Position.Y * this.layerSize,
		};

		var lScreenPos = {
			x: (layerWorldPos.x - this.pos.x) * this.tileSize,
			y: (layerWorldPos.y - this.pos.y) * this.tileSize,
		}	

		for(var i = 0; i < layer.Blocks.length; i++) {
			// Find out the world position of each block then the screen position
			var block = layer.Blocks[i];
			var lX = Math.round(i / this.layerSize);
			var lY = Math.round(i - (lX * this.layerSize));

			// Screen pos
			var screenPos = {
				x: lScreenPos.x + (lX * this.tileSize),
				y: lScreenPos.y + (lY * this.tileSize),
			};

			// Finally draw it
			this.ctx.fillStyle = this.colors[block.Id];
			this.ctx.fillRect(screenPos.x, screenPos.y, this.tileSize, this.tileSize);
		};
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
		if (newSection !== oldSection || oldPos.z !== this.pos.z) {
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