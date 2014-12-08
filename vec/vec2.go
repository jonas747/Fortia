package vec

type Vec2 interface {
}

type Vec2I struct {
	X, Y int
}

func (v *Vec2I) Add(other Vec2I) {
	v.X += other.X
	v.Y += other.Y
}

func (v *Vec2I) AddScalar(scalar int) {
	v.X += scalar
	v.Y += scalar
}

func (v *Vec2I) Sub(other Vec2I) {
	v.X -= other.X
	v.Y -= other.Y
}

func (v *Vec2I) MultiplyScalar(scalar float64) {
	v.X = int(float64(v.X) * scalar)
	v.Y = int(float64(v.Y) * scalar)
}

func (v *Vec2I) Multiply(other Vec2I) {
	v.X *= other.X
	v.Y *= other.Y
}

func (v *Vec2I) Copy(other Vec2I) {
	v.X = other.X
	v.Y = other.Y
}

func (v *Vec2I) Clone() Vec2I {
	return Vec2I{
		X: v.X,
		Y: v.Y,
	}
}
