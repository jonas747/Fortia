package vec

type Vec3I struct {
	X, Y, Z int
}

func (v *Vec3I) Add(other Vec3I) {
	v.X += other.X
	v.Y += other.Y
	v.Z += other.Z
}

func (v *Vec3I) AddScalar(scalar int) {
	v.X += scalar
	v.Y += scalar
	v.Z += scalar
}

func (v *Vec3I) Sub(other Vec3I) {
	v.X -= other.X
	v.Y -= other.Y
	v.Z -= other.Z
}

func (v *Vec3I) MltiplyScalar(scalar float64) {
	v.X = int(float64(v.X) * scalar)
	v.Y = int(float64(v.Y) * scalar)
	v.Z = int(float64(v.Z) * scalar)
}

func (v *Vec3I) Multiply(other Vec3I) {
	v.X *= other.X
	v.Y *= other.Y
	v.Z *= other.Z
}

func (v *Vec3I) Dot() float64 {
	return float64(v.X)*float64(v.X) + float64(v.Y)*float64(v.Y) + float64(v.Z)*float64(v.Z)
}

func (v *Vec3I) Copy(other Vec3I) {
	v.X = other.X
	v.Y = other.Y
	v.Z = other.Z
}

func (v *Vec3I) Clone() Vec3I {
	return Vec3I{
		X: v.X,
		Y: v.Y,
		Z: v.Z,
	}
}
