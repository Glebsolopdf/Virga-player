package audio

const (
	sampleRate      = 11025
	readChunkBytes  = 256
	analysisWinSize = 512
	minNormEnvelope = 0.003
	envSensitivity  = 28.0
)

type Bands struct {
	Low      float64
	Mid      float64
	High     float64
	Envelope float64
}
