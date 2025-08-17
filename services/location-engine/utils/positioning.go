package utils

import (
	"errors"
	"math"
	"sort"
)

type Obs struct {
	X, Y, Z float64
	D       float64
	W       float64
}

// --- Helper functions ---

// Convert RSSI to distance using log-distance path loss model.
func RssiToDistance(rssi, txPower, n float64) float64 {
	if n <= 0 {
		n = 2.0
	}
	return math.Pow(10, (txPower-rssi)/(10.0*n))
}

// Normalize weights so that sum = 1.
func NormalizeWeights(data []Obs) {
	var sum float64
	for _, o := range data {
		sum += o.W
	}
	if sum == 0 {
		return
	}
	for i := range data {
		data[i].W /= sum
	}
}

// Check if all Z are close enough (within tolerance).
// If yes, we treat it as a 2D problem and return median Z.
func CommonZ(data []Obs, tol float64) (float64, bool) {
	zs := make([]float64, 0, len(data))
	for _, o := range data {
		zs = append(zs, o.Z)
	}
	sort.Float64s(zs)
	median := zs[len(zs)/2]
	var maxDev float64
	for _, v := range zs {
		if math.Abs(v-median) > maxDev {
			maxDev = math.Abs(v - median)
		}
	}
	return median, maxDev <= tol
}

// Initial guess for solver: weighted centroid of beacons.
func initialGuess(data []Obs, fixZ bool, zFixed float64) (x, y, z float64) {
	var sx, sy, sz, sw float64
	for _, o := range data {
		sx += o.X * o.W
		sy += o.Y * o.W
		if !fixZ {
			sz += o.Z * o.W
		}
		sw += o.W
	}
	if sw == 0 {
		sw = 1
	}
	if fixZ {
		return sx / sw, sy / sw, zFixed
	}
	return sx / sw, sy / sw, sz / sw
}

// gaussNewton runs Levenberg-Marquardt style Gauss-Newton optimization
// to minimize residuals between measured distances and estimated distances.
func GaussNewton(
	data []Obs,
	fixZ bool,
	zFixed float64,
) (float64, float64, float64, error) {

	x, y, z := initialGuess(data, fixZ, zFixed)

	lambda := 1e-2
	const (
		maxIter = 50
		tolStep = 1e-4
	)

	for iter := 0; iter < maxIter; iter++ {
		var jtj [3][3]float64
		var jtr [3]float64
		var rss float64

		for _, o := range data {
			dx := x - o.X
			dy := y - o.Y
			dz := z - o.Z
			if fixZ {
				dz = 0
			}
			rng := math.Sqrt(dx*dx + dy*dy + dz*dz)
			if rng < 1e-6 {
				continue
			}
			ri := rng - o.D
			rss += o.W * ri * ri

			jx := o.W * (dx / rng)
			jy := o.W * (dy / rng)
			jz := 0.0
			if !fixZ {
				jz = o.W * (dz / rng)
			}

			jtj[0][0] += jx * jx
			jtj[0][1] += jx * jy
			jtj[0][2] += jx * jz
			jtj[1][0] += jy * jx
			jtj[1][1] += jy * jy
			jtj[1][2] += jy * jz
			jtj[2][0] += jz * jx
			jtj[2][1] += jz * jy
			jtj[2][2] += jz * jz

			jtr[0] += jx * ri
			jtr[1] += jy * ri
			jtr[2] += jz * ri
		}

		for i := 0; i < 3; i++ {
			jtj[i][i] += lambda
		}
		step, ok := solve3x3(jtj, [3]float64{-jtr[0], -jtr[1], -jtr[2]}, fixZ)
		if !ok {
			lambda *= 10
			continue
		}

		normStep := math.Sqrt(step[0]*step[0] + step[1]*step[1] + step[2]*step[2])
		nx, ny, nz := x+step[0], y+step[1], z
		if !fixZ {
			nz += step[2]
		}

		newRss := residualRSS(data, nx, ny, nz, fixZ)
		if newRss < rss {
			x, y, z = nx, ny, nz
			lambda = math.Max(lambda*0.3, 1e-6)
			if normStep < tolStep {
				break
			}
		} else {
			lambda *= 5
		}
	}

	if !IsFinite(x) || !IsFinite(y) || !IsFinite(z) {
		return 0, 0, 0, errors.New("failed to converge")
	}
	return x, y, z, nil
}

// Compute total weighted residual error for given coordinates.
func residualRSS(data []Obs, x, y, z float64, fixZ bool) float64 {
	var s float64
	for _, o := range data {
		dx := x - o.X
		dy := y - o.Y
		dz := z - o.Z
		if fixZ {
			dz = 0
		}
		r := math.Sqrt(dx*dx+dy*dy+dz*dz) - o.D
		s += o.W * r * r
	}
	return s
}

// Solve small 3x3 linear system (with optional Z fixed).
func solve3x3(A [3][3]float64, b [3]float64, fixZ bool) ([3]float64, bool) {
	if fixZ {
		B := [2][2]float64{{A[0][0], A[0][1]}, {A[1][0], A[1][1]}}
		c := [2]float64{b[0], b[1]}
		s2, ok := solve2x2(B, c)
		if !ok {
			return [3]float64{}, false
		}
		return [3]float64{s2[0], s2[1], 0}, true
	}
	// Gaussian elimination for 3x3
	M := [3][4]float64{
		{A[0][0], A[0][1], A[0][2], b[0]},
		{A[1][0], A[1][1], A[1][2], b[1]},
		{A[2][0], A[2][1], A[2][2], b[2]},
	}
	for i := 0; i < 3; i++ {
		// partial pivoting
		p := i
		for r := i + 1; r < 3; r++ {
			if math.Abs(M[r][i]) > math.Abs(M[p][i]) {
				p = r
			}
		}
		if math.Abs(M[p][i]) < 1e-12 {
			return [3]float64{}, false
		}
		if p != i {
			M[i], M[p] = M[p], M[i]
		}
		piv := M[i][i]
		for c := i; c < 4; c++ {
			M[i][c] /= piv
		}
		for r := i + 1; r < 3; r++ {
			f := M[r][i]
			for c := i; c < 4; c++ {
				M[r][c] -= f * M[i][c]
			}
		}
	}
	for i := 2; i >= 0; i-- {
		for r := 0; r < i; r++ {
			f := M[r][i]
			M[r][i] = 0
			M[r][3] -= f * M[i][3]
		}
	}
	return [3]float64{M[0][3], M[1][3], M[2][3]}, true
}

// Solve 2x2 linear system.
func solve2x2(A [2][2]float64, b [2]float64) ([2]float64, bool) {
	d := A[0][0]*A[1][1] - A[0][1]*A[1][0]
	if math.Abs(d) < 1e-12 {
		return [2]float64{}, false
	}
	return [2]float64{
		(b[0]*A[1][1] - b[1]*A[0][1]) / d,
		(-b[0]*A[1][0] + b[1]*A[0][0]) / d,
	}, true
}

// Check if value is finite (not NaN or Inf).
func IsFinite(v float64) bool {
	return !math.IsNaN(v) && !math.IsInf(v, 0)
}
