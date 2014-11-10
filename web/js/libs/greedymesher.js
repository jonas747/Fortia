// Edit by jonas747:
// Returns a typed array of vertices only, to be used with buffer geometry

var GreedyMesh = (function() {
//Cache buffer internally
var mask = new Int32Array(4096);

return function(volume, dims) {
  function f(i,j,k) {
    return volume[i + dims[0] * (j + dims[1] * k)];
  }
  //Sweep over 3-axes
  var vertices = new Float32Array(100);
  var colors = new Float32Array(100);
  var vaCounter = 0;

  for(var d=0; d<3; ++d) {
    var i, j, k, l, w, h
      , u = (d+1)%3
      , v = (d+2)%3
      , x = [0,0,0]
      , q = [0,0,0];
    if(mask.length < dims[u] * dims[v]) {
      mask = new Float32Array(dims[u] * dims[v]);
    }
    q[d] = 1;
    for(x[d]=-1; x[d]<dims[d]; ) {
      //Compute mask
      var n = 0;
      for(x[v]=0; x[v]<dims[v]; ++x[v])
      for(x[u]=0; x[u]<dims[u]; ++x[u], ++n) {
        var a = (0    <= x[d]      ? f(x[0],      x[1],      x[2])      : 0)
          , b = (x[d] <  dims[d]-1 ? f(x[0]+q[0], x[1]+q[1], x[2]+q[2]) : 0);
        if((!!a) === (!!b) ) {
          mask[n] = 0;
        } else if(!!a) {
          mask[n] = a;
        } else {
          mask[n] = -b;
        }
      }
      //Increment x[d]
      ++x[d];
      //Generate mesh for mask using lexicographic ordering
      n = 0;
      for(j=0; j<dims[v]; ++j)
      for(i=0; i<dims[u]; ) {
        var c = mask[n];
        if(!!c) {
          //Compute width
          for(w=1; c === mask[n+w] && i+w<dims[u]; ++w) {
          }
          //Compute height (this is slightly awkward
          var done = false;
          for(h=1; j+h<dims[v]; ++h) {
            for(k=0; k<w; ++k) {
              if(c !== mask[n+k+h*dims[u]]) {
                done = true;
                break;
              }
            }
            if(done) {
              break;
            }
          }
          //Add quad
          x[u] = i;  x[v] = j;
          var du = [0,0,0]
            , dv = [0,0,0]; 
          if(c > 0) {
            dv[v] = h;
            du[u] = w;
          } else {
            c = -c;
            du[v] = h;
            dv[u] = w;
          }
          if (vaCounter+18 > vertices.length) {
            // Expand it
            var newVA = new Float32Array(vertices.length + 1000);
            newVA.set(vertices);
            vertices = newVA;

            // Expand color buffer aswell
            var newColors = new Float32Array(colors.length + 1000);
            newColors.set(colors);
            colors = newColors;

          };

          // Set colors
          for (var ic = 0; ic < 18; ic += 3) {
             var cr = c >> 16;
             var cg = (c >> 8) & 0x0000ff;
             var cb = c & 0x0000ff;
             cr /= 255
             cg /= 255
             cb /= 255
            colors[vaCounter + ic] = cr;
            colors[vaCounter + ic+1] = cg;
            colors[vaCounter + ic+2] = cb;
          };

          
          vertices[vaCounter++] = x[0];
          vertices[vaCounter++] = x[1];
          vertices[vaCounter++] = x[2];
          
          vertices[vaCounter++] = x[0]+du[0];
          vertices[vaCounter++] = x[1]+du[1];
          vertices[vaCounter++] = x[2]+du[2];
          
          vertices[vaCounter++] = x[0] +dv[0];
          vertices[vaCounter++] = x[1] +dv[1];
          vertices[vaCounter++] = x[2] +dv[2];

          vertices[vaCounter++] = x[0]+du[0];
          vertices[vaCounter++] = x[1]+du[1];
          vertices[vaCounter++] = x[2]+du[2];
          
          vertices[vaCounter++] = x[0]+du[0]+dv[0];
          vertices[vaCounter++] = x[1]+du[1]+dv[1];
          vertices[vaCounter++] = x[2]+du[2]+dv[2];
          
          vertices[vaCounter++] = x[0] +dv[0];
          vertices[vaCounter++] = x[1] +dv[1];
          vertices[vaCounter++] = x[2] +dv[2];

          //faces.push([vertex_count, vertex_count+1, vertex_count+2, vertex_count+3, c]);
          
          //Zero-out mask
          for(l=0; l<h; ++l)
          for(k=0; k<w; ++k) {
            mask[n+k+l*dims[u]] = 0;
          }
          //Increment counters and continue
          i += w; n += w;
        } else {
          ++i;    ++n;
        }
      }
    }
  }
  // Trim the vertices array
  if (vertices.length > vaCounter) { 
    var newVA = new Float32Array(vaCounter);
    var subbedVA = vertices.subarray(0, vaCounter);
    newVA.set(subbedVA);
    vertices = newVA;
    
    // Also trim the color buffer
    var newColors = new Float32Array(vaCounter);
    var subbedColors = colors.subarray(0, vaCounter);
    newColors.set(subbedColors);
    colors = newColors;
  };
  return { vertices:vertices, colors: colors};
}
})();