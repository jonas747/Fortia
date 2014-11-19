/* SimplexNoise1234, Simplex noise with true analytic
 * derivative in 1D to 4D.
 *
/*
 * This implementation is "Simplex Noise" as presented by
 * Ken Perlin at a relatively obscure and not often cited course
 * session "Real-Time Shading" at Siggraph 2001 (before real
 * time shading actually took on), under the title "hardware noise".
 * The 3D function is numerically equivalent to his Java reference
 * code available in the PDF course notes, although I re-implemented
 * it from scratch to get more readable code. The 1D, 2D and 4D cases
 * were implemented from scratch by me from Ken Perlin's text.
*/

package simplex

import (
	"math"
	"math/rand"
)

type Noise struct {
	Perm [512]uint8
}

func NewNoise(rng *rand.Rand) *Noise {
	if rng == nil {
		rng = rand.New(rand.NewSource(rand.Int63()))
	}

	perm := [512]uint8{}
	for i := 0; i < 512; i++ {
		perm[i] = uint8(rng.Intn(256))
		if i >= 256 {
			perm[i] = perm[i&255]
		}
	}
	return &Noise{
		Perm: perm,
	}
}

//---------------------------------------------------------------------

/*
 * Helper functions to compute gradients-dot-residualvectors (1D to 4D)
 * Note that these generate gradients of more than unit length. To make
 * a close match with the value range of classic Perlin noise, the final
 * noise values need to be rescaled to fit nicely within [-1,1].
 * (The simplex noise functions as such also have different scaling.)
 * Note also that these noise functions are the most practical and useful
 * signed version of Perlin noise. To return values according to the
 * RenderMan specification from the SL noise() and pnoise() functions,
 * the noise values need to be scaled and offset to [0,1], like this:
 * float SLnoise = (noise(x,y,z) + 1.0) * 0.5;
 */

func Q(cond bool, v1 float64, v2 float64) float64 {
	if cond {
		return v1
	}
	return v2
}

func grad1(hash uint8, x float64) float64 {
	h := hash & 15
	grad := float64(1 + h&7) // Gradient value 1.0, 2.0, ..., 8.0
	if h&8 != 0 {
		grad = -grad // Set a random sign for the gradient
	}
	return grad * x // Multiply the gradient with the distance
}

func grad2(hash uint8, x float64, y float64) float64 {
	h := hash & 7       // Convert low 3 bits of hash code
	u := Q(h < 4, x, y) // into 8 simple gradient directions,
	v := Q(h < 4, y, x) // and compute the dot product with (x,y).
	return Q(h&1 != 0, -u, u) + Q(h&2 != 0, -2*v, 2*v)
}

func grad3(hash uint8, x, y, z float64) float64 {
	h := hash & 15                                // Convert low 4 bits of hash code into 12 simple
	u := Q(h < 8, x, y)                           // gradient directions, and compute dot product.
	v := Q(h < 4, y, Q(h == 12 || h == 14, x, z)) // Fix repeats at h = 12 to 15
	return Q(h&1 != 0, -u, u) + Q(h&2 != 0, -v, v)
}

// 1D simplex noise
func (n *Noise) Noise1(x float64) float64 {
	i0 := int(math.Floor(x))
	i1 := i0 + 1
	x0 := x - float64(i0)
	x1 := x0 - 1

	t0 := 1 - x0*x0
	t0 *= t0
	n0 := t0 * t0 * grad1(n.Perm[i0&0xff], x0)

	t1 := 1 - x1*x1
	t1 *= t1
	n1 := t1 * t1 * grad1(n.Perm[i1&0xff], x1)
	// The maximum value of this noise is 8*(3/4)^4 = 2.53125
	// A factor of 0.395 would scale to fit exactly within [-1,1].
	// fmt.Printf("Noise1 x %.4f, i0 %v, i1 %v, x0 %.4f, x1 %.4f, perm0 %d, perm1 %d: %.4f,%.4f\n", x, i0, i1, x0, x1, perm[i0&0xff], perm[i1&0xff], n0, n1)
	// The algorithm isn't perfect, as it is assymetric. The correction will normalize the result to the interval [-1,1], but the average will be off by 3%.
	return (n0 + n1 + 0.076368899) / 2.45488110001
}

// 2D simplex noise
func (n *Noise) Noise2(x, y float64) float64 {

	const F2 = 0.366025403 // F2 = 0.5*(sqrt(3.0)-1.0)
	const G2 = 0.211324865 // G2 = (3.0-Math.sqrt(3.0))/6.0

	var n0, n1, n2 float64 // Noise contributions from the three corners

	// Skew the input space to determine which simplex cell we're in
	s := (x + y) * F2 // Hairy factor for 2D
	xs := x + s
	ys := y + s
	i := int(math.Floor(xs))
	j := int(math.Floor(ys))

	t := float64(i+j) * G2
	X0 := float64(i) - t // Unskew the cell origin back to (x,y) space
	Y0 := float64(j) - t
	x0 := x - X0 // The x,y distances from the cell origin
	y0 := y - Y0

	// For the 2D case, the simplex shape is an equilateral triangle.
	// Determine which simplex we are in.
	var i1, j1 int // Offsets for second (middle) corner of simplex in (i,j) coords
	if x0 > y0 {
		i1 = 1
		j1 = 0 // lower triangle, XY order: (0,0)->(1,0)->(1,1)
	} else {
		i1 = 0
		j1 = 1
	} // upper triangle, YX order: (0,0)->(0,1)->(1,1)

	// A step of (1,0) in (i,j) means a step of (1-c,-c) in (x,y), and
	// a step of (0,1) in (i,j) means a step of (-c,1-c) in (x,y), where
	// c = (3-sqrt(3))/6

	x1 := x0 - float64(i1) + G2 // Offsets for middle corner in (x,y) unskewed coords
	y1 := y0 - float64(j1) + G2
	x2 := x0 - 1 + 2*G2 // Offsets for last corner in (x,y) unskewed coords
	y2 := y0 - 1 + 2*G2

	// Wrap the integer indices at 256, to avoid indexing perm[] out of bounds
	ii := i & 0xff
	jj := j & 0xff

	// Calculate the contribution from the three corners
	t0 := 0.5 - x0*x0 - y0*y0
	if t0 < 0 {
		n0 = 0
	} else {
		t0 *= t0
		n0 = t0 * t0 * grad2(n.Perm[ii+int(n.Perm[jj])], x0, y0)
	}

	t1 := 0.5 - x1*x1 - y1*y1
	if t1 < 0 {
		n1 = 0
	} else {
		t1 *= t1
		n1 = t1 * t1 * grad2(n.Perm[ii+i1+int(n.Perm[jj+j1])], x1, y1)
	}

	t2 := 0.5 - x2*x2 - y2*y2
	if t2 < 0 {
		n2 = 0
	} else {
		t2 *= t2
		n2 = t2 * t2 * grad2(n.Perm[ii+1+int(n.Perm[jj+1])], x2, y2)
	}

	// Add contributions from each corner to get the final noise value.
	// The result is scaled to return values in the interval [-1,1].
	return (n0 + n1 + n2) / 0.022108854818853867
}

// 3D simplex noise
func (n *Noise) Noise3(x, y, z float64) float64 {

	// Simple skewing factors for the 3D case
	const F3 = 0.333333333
	const G3 = 0.166666667

	var n0, n1, n2, n3 float64 // Noise contributions from the four corners

	// Skew the input space to determine which simplex cell we're in
	s := (x + y + z) * F3 // Very nice and simple skew factor for 3D
	xs := x + s
	ys := y + s
	zs := z + s
	i := int(math.Floor(xs))
	j := int(math.Floor(ys))
	k := int(math.Floor(zs))

	t := float64(i+j+k) * G3
	X0 := float64(i) - t // Unskew the cell origin back to (x,y,z) space
	Y0 := float64(j) - t
	Z0 := float64(k) - t
	x0 := float64(x) - X0 // The x,y,z distances from the cell origin
	y0 := float64(y) - Y0
	z0 := float64(z) - Z0

	// For the 3D case, the simplex shape is a slightly irregular tetrahedron.
	// Determine which simplex we are in.
	var i1, j1, k1 int // Offsets for second corner of simplex in (i,j,k) coords
	var i2, j2, k2 int // Offsets for third corner of simplex in (i,j,k) coords

	/* This code would benefit from a backport from the GLSL version! */
	if x0 >= y0 {
		if y0 >= z0 {
			i1 = 1
			j1 = 0
			k1 = 0
			i2 = 1
			j2 = 1
			k2 = 0 // X Y Z order
		} else if x0 >= z0 {
			i1 = 1
			j1 = 0
			k1 = 0
			i2 = 1
			j2 = 0
			k2 = 1 // X Z Y order
		} else {
			i1 = 0
			j1 = 0
			k1 = 1
			i2 = 1
			j2 = 0
			k2 = 1 // Z X Y order
		}
	} else { // x0<y0
		if y0 < z0 {
			i1 = 0
			j1 = 0
			k1 = 1
			i2 = 0
			j2 = 1
			k2 = 1 // Z Y X order
		} else if x0 < z0 {
			i1 = 0
			j1 = 1
			k1 = 0
			i2 = 0
			j2 = 1
			k2 = 1 // Y Z X order
		} else {
			i1 = 0
			j1 = 1
			k1 = 0
			i2 = 1
			j2 = 1
			k2 = 0 // Y X Z order
		}
	}

	// A step of (1,0,0) in (i,j,k) means a step of (1-c,-c,-c) in (x,y,z),
	// a step of (0,1,0) in (i,j,k) means a step of (-c,1-c,-c) in (x,y,z), and
	// a step of (0,0,1) in (i,j,k) means a step of (-c,-c,1-c) in (x,y,z), where
	// c = 1/6.

	x1 := x0 - float64(i1) + G3 // Offsets for second corner in (x,y,z) coords
	y1 := y0 - float64(j1) + G3
	z1 := z0 - float64(k1) + G3
	x2 := x0 - float64(i2) + 2*G3 // Offsets for third corner in (x,y,z) coords
	y2 := y0 - float64(j2) + 2*G3
	z2 := z0 - float64(k2) + 2*G3
	x3 := x0 - 1 + 3*G3 // Offsets for last corner in (x,y,z) coords
	y3 := y0 - 1 + 3*G3
	z3 := z0 - 1 + 3*G3

	// Wrap the integer indices at 256, to avoid indexing perm[] out of bounds
	ii := i & 0xff
	jj := j & 0xff
	kk := k & 0xff

	// Calculate the contribution from the four corners
	t0 := 0.6 - x0*x0 - y0*y0 - z0*z0
	if t0 < 0 {
		n0 = 0
	} else {
		t0 *= t0
		n0 = t0 * t0 * grad3(n.Perm[ii+int(n.Perm[jj+int(n.Perm[kk])])], x0, y0, z0)
	}

	t1 := 0.6 - x1*x1 - y1*y1 - z1*z1
	if t1 < 0 {
		n1 = 0
	} else {
		t1 *= t1
		n1 = t1 * t1 * grad3(n.Perm[ii+i1+int(n.Perm[jj+j1+int(n.Perm[kk+k1])])], x1, y1, z1)
	}

	t2 := 0.6 - x2*x2 - y2*y2 - z2*z2
	if t2 < 0 {
		n2 = 0
	} else {
		t2 *= t2
		n2 = t2 * t2 * grad3(n.Perm[ii+i2+int(n.Perm[jj+j2+int(n.Perm[kk+k2])])], x2, y2, z2)
	}

	t3 := 0.6 - x3*x3 - y3*y3 - z3*z3
	if t3 < 0 {
		n3 = 0
	} else {
		t3 *= t3
		n3 = t3 * t3 * grad3(n.Perm[ii+1+int(n.Perm[jj+1+int(n.Perm[kk+1])])], x3, y3, z3)
	}

	// Add contributions from each corner to get the final noise value.
	// The result is scaled to stay just inside [-1,1]
	return (n0 + n1 + n2 + n3) / 0.030555466710745972
}
