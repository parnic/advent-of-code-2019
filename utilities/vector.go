package utilities

import "math"

type Vec2[T Number] struct {
	X T
	Y T
}

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

func VecBetween[T Number](a, b Vec2[T]) Vec2[T] {
	return Vec2[T]{
		X: a.X - b.X,
		Y: a.Y - b.Y,
	}
}
