package utilities

import "math"

type Vec2[T Number] struct {
	X T
	Y T
}

type Vec3[T Number] struct {
	X T
	Y T
	Z T
}

type Vec2i Vec2[int]

func (v Vec2[T]) Dot(other Vec2[T]) T {
	return (v.X * other.X) + (v.Y * other.Y)
}

func (v Vec2[T]) Len() T {
	return T(math.Sqrt(float64(v.LenSquared())))
}

func (v Vec2[T]) LenSquared() T {
	return (v.X * v.X) + (v.Y * v.Y)
}

func (v Vec2[T]) To(other Vec2[T]) Vec2[T] {
	return Vec2[T]{
		X: v.X - other.X,
		Y: v.Y - other.Y,
	}
}

func (v Vec2[T]) AngleBetween(other Vec2[T]) float64 {
	rad := math.Atan2(float64(other.Y-v.Y), float64(other.X-v.X))
	return rad * 180 / math.Pi
}

func (v Vec2[T]) Equals(other Vec2[T]) bool {
	return v.X == other.X &&
		v.Y == other.Y
}

func (v Vec2[T]) ManhattanDistance(other Vec2[T]) T {
	return T(math.Abs(float64(v.X-other.X)) + math.Abs(float64(v.Y-other.Y)))
}

func VecBetween[T Number](a, b Vec2[T]) Vec2[T] {
	return a.To(b)
}

func ManhattanDistance[T Number](a, b Vec2[T]) T {
	return a.ManhattanDistance(b)
}

func (v Vec3[T]) Dot(other Vec3[T]) T {
	return (v.X * other.X) + (v.Y * other.Y) + (v.Z * other.Z)
}

func (v Vec3[T]) Len() T {
	return T(math.Sqrt(float64(v.LenSquared())))
}

func (v Vec3[T]) LenSquared() T {
	return (v.X * v.X) + (v.Y * v.Y) + (v.Z * v.Z)
}

func (v *Vec3[T]) Add(other Vec3[T]) {
	v.X += other.X
	v.Y += other.Y
	v.Z += other.Z
}

func (v Vec3[T]) Equals(other Vec3[T]) bool {
	return v.X == other.X &&
		v.Y == other.Y &&
		v.Z == other.Z
}
