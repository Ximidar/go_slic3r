package slice

// ScalingFactor for a conversion from coord_t to coordf_t: 10e-6
// This scaling generates a following fixed point representation with for a 32bit integer:
// 0..4294mm with 1nm resolution
const ScalingFactor float64 = 0.000001

func Scale(val float64) float64 {
	return val / ScalingFactor
}

func UnScale(val float64) float64 {
	return val * ScalingFactor
}

const Epsilon float64 = 1e-4

var ScaledEpsilon float64 = Scale(Epsilon)

const Resolution float64 = 0.0125

var ScaledResolution float64 = Scale(Resolution)
