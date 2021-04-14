package generation

import (
	"math"

	simplex "github.com/ojrac/opensimplex-go"
	log "github.com/sirupsen/logrus"
	"time"
)

type TerrainGenerator interface {
	GenerateTerrain(width int, height int, offsetX, offsetY float64) []float32
	Seed() int64
	SetSeed(seed int64)
}

type SimplexTerrainGenerator struct {
	config    TerrainGeneratorConfig
	generator simplex.Noise
}

type TerrainGeneratorConfig struct {
	Octaves     int
	Persistence float64
	ScaleFactor float64
	Normalize   bool
	Seed        int64
	Debug       bool
}

func NewSimplexTerrainGenerator(config TerrainGeneratorConfig) SimplexTerrainGenerator {
	logger := log.
		WithField("module", "terrain_generator").
		WithField("config", config).
		WithField("seed", config.Seed)

	logger.Info("Simplex terrain generator initialized")

	if config.Seed == 0 {
		config.Seed = time.Now().Unix()
		log.WithField("seed", config.Seed).Debug("Generated seed")
	}

	return SimplexTerrainGenerator{
		config:    config,
		generator: simplex.New(config.Seed),
	}
}

func (s SimplexTerrainGenerator) Seed() int64 {
	return s.config.Seed
}

func (s SimplexTerrainGenerator) SetSeed(seed int64) {
	s.config.Seed = seed
	s.generator = simplex.New(seed)
}

func (s SimplexTerrainGenerator) GenerateTerrain(width, height int, offsetX, offsetY float64) (result []float32) {
	pixels := make([][]float64, width)
	maxNoise := 0.0
	minNoise := 0.0

	finishChan := make(chan interface{}, width)
	for x := 0; x < width; x++ {
		pixels[x] = make([]float64, height)

		x := x

		go func() {
			for y := 0; y < height; y++ {
				noise := 0.0
				freq := 2.0

				for octave := 0; octave < s.config.Octaves; octave++ {
					// Freq is always growing
					freq = math.Pow(2, float64(octave))
					amplitude := math.Pow(s.config.Persistence, float64(octave))

					nx := (float64(x) + offsetX) / float64(width)
					ny := (float64(y) + offsetY) / float64(height)

					noiseVal := s.generator.Eval2(freq*nx, freq*ny)
					if s.config.Debug {
						log.WithFields(log.Fields{
							"module":    "terrain_generator",
							"noise":     noiseVal,
							"x":         x,
							"y":         y,
							"amplitude": amplitude,
							"freq":      freq,
							"nx":        nx,
							"ny":        ny,
							"octave":    octave,
						}).Debugf("Nose generated")
					}

					// Map noise from [-1;1] to [0;1)
					noiseVal = (noiseVal + 1.0) / 2.0
					noise += amplitude * noiseVal
				}

				pixels[x][y] = noise
			}

			finishChan <- true
		}()
	}

	// Wait all routines to complete
	for i := 0; i < width; i++ {
		<-finishChan
	}

	//Calculate min/max noise
	if s.config.Normalize {
		for x := 0; x < width; x++ {
			for y := 0; y < height; y++ {
				noise := pixels[x][y]
				maxNoise = math.Max(noise, maxNoise)
				minNoise = math.Min(noise, minNoise)
			}
		}
	}

	for x := 0; x < width; x++ {
		for y := 0; y < height; y++ {
			var normalized float64
			if s.config.Normalize && maxNoise != minNoise {
				normalized = (pixels[x][y] - minNoise) / (maxNoise - minNoise)
			} else {
				normalized = pixels[x][y]
			}
			//normalized = math.Pow(normalized, s.config.Persistence)

			result = append(result, float32(normalized*s.config.ScaleFactor))
		}
	}

	return result
}
