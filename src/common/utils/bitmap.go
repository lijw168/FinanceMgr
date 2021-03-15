package utils

type Bitmap struct {
	Mem     []byte `json:"mem"`
	BitSize uint64 `json:"bit_size"`
}

func NewBitmap(bitSize uint64, data []byte) *Bitmap {
	bmp := new(Bitmap)
	bmp.BitSize = bitSize

	if data == nil {
		byteSize := (bitSize + 7) >> 3
		// Golang zero it automaticaly.
		bmp.Mem = make([]byte, byteSize)
	} else {
		bmp.Mem = data
	}
	return bmp
}

func (b *Bitmap) Init(bitSize uint64, data []byte) {
	b.BitSize = bitSize
	b.Mem = data
}

// Set bit to 1.
func (b *Bitmap) Set(nr uint64) {
	Assert(nr < b.BitSize)

	bit := byte(1) << byte(nr&7)
	b.Mem[nr>>3] |= bit
}

// Set bit to 0.
func (b *Bitmap) Clear(nr uint64) {
	Assert(nr < b.BitSize)

	bit := byte(1) << byte(nr&7)
	b.Mem[nr>>3] &= ^bit
}

// Test whether bit is 1.
func (b *Bitmap) Test(nr uint64) (ret bool) {
	Assert(nr < b.BitSize)

	bit := byte(1) << byte(nr&7)
	return (b.Mem[nr>>3]&bit != 0)
}

// Get bitmap size.
func (b *Bitmap) Size() (size uint64) {
	return b.BitSize
}
