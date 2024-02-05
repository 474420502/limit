package limit

// Float64
type Float64 float64

func (x Float64) CountedFloat() float64 {
	return float64(x)
}

// Float
type Float Float64

// Int64
type Int64 int64

func (x Int64) CountedFloat() float64 {
	return float64(x)
}

// Int
type Int32 int32

func (x Int32) CountedFloat() float64 {
	return float64(x)
}

// Int
type Int int

func (x Int) CountedFloat() float64 {
	return float64(x)
}

// Uint64
type Uint64 uint64

func (x Uint64) CountedFloat() float64 {
	return float64(x)
}

// Uint32
type Uint32 uint32

func (x Uint32) CountedFloat() float64 {
	return float64(x)
}

// Uint
type Uint uint

func (x Uint) CountedFloat() float64 {
	return float64(x)
}
