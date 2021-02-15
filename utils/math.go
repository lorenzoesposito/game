package utils

import "math"

// Triangle 2D float64
type Tri2f struct {
	a, b, c Vec2f
}

// Triangle 2D uint16
type Tri2ui struct {
	a, b, c Vec2ui
}

// Triangle 3D float64
type Tri3f struct {
	a, b, c Vec3f
}

// Vector 3D float64
type Vec3f struct {
	X, Y, Z float64
}

// Vector 2D float64
type Vec2f struct {
	x, y float64
}

// Vector 2D uint16
type Vec2ui struct {
	x, y uint16
}

func ValVec(n float64) Vec3f {
	return Vec3f{n, n, n}
}

func Plus(a, b Vec3f) Vec3f {
	return Vec3f{a.X + b.X, a.Y + b.Y, a.Z + b.Z}
}

func Minus(a, b Vec3f) Vec3f {
	return Vec3f{a.X - b.X, a.Y - b.Y, a.Z - b.Z}
}

func Times(a, b Vec3f) Vec3f {
	return Vec3f{a.X * b.X, a.Y * b.Y, a.Z * b.Z}
}

func Div(a, b Vec3f) Vec3f {
	return Vec3f{a.X / b.X, a.Y / b.Y, a.Z / b.Z}
}

func Len(vec Vec3f) float64 {
	return math.Sqrt(vec.X*vec.X + vec.Y*vec.Y + vec.Z*vec.Z)
}

func Dist(vec1, vec2 Vec3f) (dist float64) {
	return Len(Minus(vec1, vec2))
}

func Normalize(vec Vec3f) (q Vec3f) {
	return Div(vec, ValVec(Len(vec)))
}

func MathAbs(n float64) float64 {
	if n < 0 {
		return -n
	}
	return n
}

func DegToRad(deg float64) float64 {
	return deg * 0.01745
}

func MathClamp(n, min, max float64) float64 {
	return math.Max(math.Min(n, max), min)
}

func Abs(vec Vec3f) Vec3f {
	return Vec3f{MathAbs(vec.X), MathAbs(vec.Y), MathAbs(vec.Z)}
}

func Max(vec1, vec2 Vec3f) Vec3f {
	return Vec3f{math.Max(vec1.X, vec2.X), math.Max(vec1.Y, vec2.Y), math.Max(vec1.Z, vec2.Z)}
}

func Min(vec1, vec2 Vec3f) Vec3f {
	return Vec3f{math.Min(vec1.X, vec2.X), math.Min(vec1.Y, vec2.Y), math.Min(vec1.Z, vec2.Z)}
}

func Dot(vec1, vec2 Vec3f) (dist float64) {
	return vec1.X*vec2.X + vec1.Y*vec2.Y + vec1.Z*vec2.Z
}

func Cross(vec1, vec2 Vec3f) Vec3f {
	return Vec3f{vec1.Y*vec2.Z - vec1.Z*vec2.Y,
		vec1.Z*vec2.X - vec1.X*vec2.Z,
		vec1.X*vec2.Y - vec1.Y*vec2.X}
}

func ProjectOnVector(point, normal Vec3f) Vec3f {
	return Times(ValVec(Dot(point, normal)), normal)
}

func ProjectOnPlane(point, normal Vec3f) Vec3f {
	return Minus(point, Times(ValVec(Dot(normal, point)), normal))
}

func Clamp(vec Vec3f, min, max float64) Vec3f {
	return Vec3f{MathClamp(vec.X, min, max), MathClamp(vec.Y, min, max), MathClamp(vec.Z, min, max)}
}

func VecToColor(vec Vec3f) (r, g, b uint8) {
	vec = Times(vec, ValVec(255))
	return uint8(vec.X), uint8(vec.Y), uint8(vec.Z)
}

// Matrix 4x4 float64

func IdentityMatrix() [4][4]float64 {
	return [4][4]float64{
		[4]float64{1, 0, 0, 0},
		[4]float64{0, 1, 0, 0},
		[4]float64{0, 0, 1, 0},
		[4]float64{0, 0, 0, 1}}
}

func Translate(mat [4][4]float64, vec Vec3f) [4][4]float64 {
	if vec.X != 0 {
		mat[3][0] += vec.X
	}
	if vec.Y != 0 {
		mat[3][1] += vec.Y
	}
	if vec.Z != 0 {
		mat[3][2] += vec.Z
	}
	return mat
}

func SetPosition(mat [4][4]float64, vec Vec3f) [4][4]float64 {
	mat[3][0] = vec.X
	mat[3][1] = vec.Y
	mat[3][2] = vec.Z
	return mat
}

func Rotate(mat [4][4]float64, vec Vec3f) [4][4]float64 {
	if vec.X != 0 {
		mat[1][1] = math.Cos(vec.X)
		mat[2][1] = -math.Sin(vec.X)
		mat[1][1] = math.Cos(vec.X)
		mat[2][2] = math.Sin(vec.X)
	}
	if vec.Y != 0 {
		mat[0][0] = math.Cos(vec.Y)
		mat[2][0] = math.Sin(vec.Y)
		mat[0][2] = -math.Sin(vec.Y)
		mat[2][2] = math.Cos(vec.Y)
	}
	if vec.Z != 0 {
		mat[0][0] = math.Cos(vec.Z)
		mat[1][0] = -math.Sin(vec.Z)
		mat[0][1] = math.Sin(vec.Z)
		mat[1][1] = math.Cos(vec.Z)
	}
	return mat
}

func SetRotation(mat [4][4]float64, vec Vec3f) [4][4]float64 {
	if vec.X != 0 {
		mat[1][1] = math.Cos(vec.X)
		mat[2][1] = -math.Sin(vec.X)
		mat[1][1] = math.Cos(vec.X)
		mat[2][2] = math.Sin(vec.X)
	}
	if vec.Y != 0 {
		mat[0][0] = math.Cos(vec.Y)
		mat[2][0] = math.Sin(vec.Y)
		mat[0][2] = -math.Sin(vec.Y)
		mat[2][2] = math.Cos(vec.Y)
	}
	if vec.Z != 0 {
		mat[0][0] = math.Cos(vec.Z)
		mat[1][0] = -math.Sin(vec.Z)
		mat[0][1] = math.Sin(vec.Z)
		mat[1][1] = math.Cos(vec.Z)
	}
	return mat
}

func Inverse(mat [4][4]float64) [4][4]float64 {
	if Determinant(mat) == 0 {
		return mat
	}

	mat1 := IdentityMatrix()

	for i := 0; i < 4; i++ {
		pivot := mat[i][i]

		if pivot != 1 && pivot != 0 {
			for t := i; t < 4; t++ {
				mat[i][t] = mat[i][t] / pivot
				mat1[i][t] = mat1[i][t] / pivot
			}
		}

		//Update to the new pivot which must be 1.0
		pivot = mat[i][i]

		for j := 0; j < 4; j++ {
			if j == i {
				continue
			} else {
				l := mat[j][i] / pivot
				for m := 0; m < 4; m++ {
					mat[j][m] = mat[j][m] - l*mat[i][m]
					mat1[j][m] = mat1[j][m] - (l * mat1[i][m])
				}
			}
		}
	}
	return mat1
}

func Determinant(mat [4][4]float64) float64 {
	return mat[0][0]*mat[1][1]*mat[2][2]*mat[3][3] + mat[1][0]*mat[2][1]*mat[3][2]*mat[0][3] +
		mat[2][0]*mat[3][1]*mat[0][2]*mat[1][3] + mat[3][0]*mat[0][1]*mat[1][2]*mat[2][3] -
		mat[3][0]*mat[2][1]*mat[1][2]*mat[0][3] - mat[2][0]*mat[1][1]*mat[0][2]*mat[3][3] -
		mat[1][0]*mat[0][1]*mat[3][2]*mat[2][3] - mat[0][0]*mat[3][1]*mat[2][2]*mat[1][3]
}

func MultPointMatrix(p Vec3f, mat [4][4]float64) (p1 Vec3f) {
	p1.X = p.X*mat[0][0] + p.Y*mat[1][0] + p.Z*mat[2][0] + mat[3][0]
	p1.Y = p.X*mat[0][1] + p.Y*mat[1][1] + p.Z*mat[2][1] + mat[3][1]
	p1.Z = p.X*mat[0][2] + p.Y*mat[1][2] + p.Z*mat[2][2] + mat[3][2]
	w := p.X*mat[0][3] + p.Y*mat[1][3] + p.Z*mat[2][3] + mat[3][3]
	if w != 1 && w != 0 {
		p.X, p.Y, p.Z = p.X/w, p.Y/w, p.Z/w
	}
	return
}
