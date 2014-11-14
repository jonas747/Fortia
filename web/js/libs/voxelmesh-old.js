/*
data is a voxel containing voxels and dims

data.voxels is an array containig colors for all the blocks, 0 being no block
data.dims is an array of size 3, containing the x,y and z size of the voxel
*/

// Edit by jonas747:
// Returns uv, position and color buffers
function VoxelMesh() {
  this.THREE = THREE || three;
}

VoxelMesh.prototype.mesh = function(data, mesher, scaleFactor, mesherExtraData){
  this.data = data;
  this.scale = scaleFactor || new this.THREE.Vector3(10, 10, 10)
  var result = mesher( data.voxels, data.dims, mesherExtraData )
  this.result = result;

  var numFaces = result.vertices.length / 18
  var uv = new Float32Array(numFaces * 12)

  for (var i = 0; i < numFaces; i++) {
    var fi = i * 18
    var vs = [
      [result.vertices[fi], result.vertices[fi+1], result.vertices[fi+2]],
      [result.vertices[fi+3], result.vertices[fi+4], result.vertices[fi+5]],
      [result.vertices[fi+6], result.vertices[fi+7], result.vertices[fi+8]],
      // Skip to last vertex
      [result.vertices[i+15], result.vertices[i+16], result.vertices[i+17]]
    ];
    
    var rawUv = this.faceVertexUv(vs);
    //var processedUv = new Float32Array(rawUv.length * 3) //maybe right or nto?>
    var uva = rawUv[0]
    var uvb = rawUv[1]
    var uvc = rawUv[2]
    var uvd = rawUv[3]

    var ui = i * 12;
    uv[ui]   = uva.x
    uv[ui+1] = uva.y

    uv[ui+2]   = uvb.x
    uv[ui+3]   = uvb.y

    uv[ui+4]   = uvd.x
    uv[ui+5]   = uvd.y
  
    // second
    uv[ui+6]   = uvb.x
    uv[ui+7]   = uvb.y

    uv[ui+8]   = uvc.x
    uv[ui+9]   = uvc.y

    uv[ui+10]   = uvd.x
    uv[ui+11]   = uvd.y

  };
  this.result.uv = uv;
}
// Generates the buffergeomtry object with attributes and everything
VoxelMesh.prototype.bufferGeometry = function(){
  var geometry = this.geometry = new THREE.BufferGeometry();

  geometry.addAttribute("position", new THREE.BufferAttribute(this.result.vertices, 3));
  geometry.addAttribute("color", new THREE.BufferAttribute(this.result.colors, 3));
}
VoxelMesh.prototype.createWireMesh = function(hexColor) {    
  var wireMaterial = new this.THREE.MeshBasicMaterial({
    color : hexColor || 0xffffff,
    wireframe : true
  })
  wireMesh = new this.THREE.Mesh(this.geometry, wireMaterial)
  wireMesh.scale = this.scale
  wireMesh.doubleSided = true
  this.wireMesh = wireMesh
  return wireMesh
}

VoxelMesh.prototype.createSurfaceMesh = function(material) {
  material = material || new this.THREE.MeshNormalMaterial()
  var surfaceMesh  = new this.THREE.Mesh( this.geometry, material )
  surfaceMesh.scale = this.scale
  surfaceMesh.doubleSided = false
  this.surfaceMesh = surfaceMesh
  return surfaceMesh
}

VoxelMesh.prototype.addToScene = function(scene) {
  if (this.wireMesh) scene.add( this.wireMesh )
  if (this.surfaceMesh) scene.add( this.surfaceMesh )
}

VoxelMesh.prototype.setPosition = function(x, y, z) {
  if (this.wireMesh) this.wireMesh.position = new this.THREE.Vector3(x, y, z)
  if (this.surfaceMesh) this.surfaceMesh.position = new this.THREE.Vector3(x, y, z)
}

VoxelMesh.prototype.faceVertexUv = function(vs) {
 
  var spans = {
    x0: vs[0][0] - vs[1][0],
    x1: vs[1][0] - vs[2][0],
    y0: vs[0][1] - vs[1][1],
    y1: vs[1][1] - vs[2][1],
    z0: vs[0][2] - vs[1][2],
    z1: vs[1][2] - vs[2][2]
  }

  var size = {
    x: Math.max(Math.abs(spans.x0), Math.abs(spans.x1)),
    y: Math.max(Math.abs(spans.y0), Math.abs(spans.y1)),
    z: Math.max(Math.abs(spans.z0), Math.abs(spans.z1))
  }
  if (size.x === 0) {
    if (spans.y0 > spans.y1) {
      var width = size.y
      var height = size.z
    }
    else {
      var width = size.z
      var height = size.y
    }
  }
  if (size.y === 0) {
    if (spans.x0 > spans.x1) {
      var width = size.x
      var height = size.z
    }
    else {
      var width = size.z
      var height = size.x
    }
  }
  if (size.z === 0) {
    if (spans.x0 > spans.x1) {
      var width = size.x
      var height = size.y
    }
    else {
      var width = size.y
      var height = size.x
    }
  }
  if ((size.z === 0 && spans.x0 < spans.x1) || (size.x === 0 && spans.y0 > spans.y1)) {
    return [
      new this.THREE.Vector2(height, 0),
      new this.THREE.Vector2(0, 0),
      new this.THREE.Vector2(0, width),
      new this.THREE.Vector2(height, width)
    ]
  } else {
    return [
      new this.THREE.Vector2(0, 0),
      new this.THREE.Vector2(0, height),
      new this.THREE.Vector2(width, height),
      new this.THREE.Vector2(width, 0)
    ]
  }
};