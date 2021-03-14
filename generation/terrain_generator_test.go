package generation

import (
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"testing"
	"time"
)

func TestSimplexTerrainGenerator_GenerateTerrain(t *testing.T) {
	testGenerator := NewSimplexTerrainGenerator(TerrainGeneratorConfig{
		Octaves:     7,
		Persistence: 0.8,
		ScaleFactor: 1,
		Normalize:   true,
		Debug:       true,
	}, time.Now().UnixNano())

	terrain := testGenerator.GenerateTerrain(10, 10, 0, 0)
	for _, point := range terrain {
		assert.True(t, point >= 0.0, "point bellow zero")
		assert.True(t, point <= 1.0, "point greater then zero")
	}

	zeroTerrain := make([]float32, 100)
	require.NotEqual(t, zeroTerrain, terrain)
}

func benchmarkGenerateTerrain(b *testing.B, octaves int, size int) {
	testGenerator := NewSimplexTerrainGenerator(TerrainGeneratorConfig{
		Octaves:     octaves,
		Persistence: 1,
		ScaleFactor: 1,
		Normalize:   true,
		Debug:       false,
	}, time.Now().UnixNano())

	for i := 0; i < b.N; i++ {
		testGenerator.GenerateTerrain(size, size, float64(size)*float64(i), float64(size)*float64(i))
	}
}

func BenchmarkSimplexTerrainGenerator_GenerateTerrain10(b *testing.B) {
	benchmarkGenerateTerrain(b, 7, 10)
}

func BenchmarkSimplexTerrainGenerator_GenerateTerrain100(b *testing.B) {
	benchmarkGenerateTerrain(b, 7, 100)
}

func BenchmarkSimplexTerrainGenerator_GenerateTerrain1000(b *testing.B) {
	benchmarkGenerateTerrain(b, 7, 1000)
}

func BenchmarkSimplexTerrainGenerator_GenerateTerrain5000(b *testing.B) {
	benchmarkGenerateTerrain(b, 7, 5000)
}
