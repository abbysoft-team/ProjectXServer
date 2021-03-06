package model

var (
	ResourcesPlaceTown = Resources{
		Wood:    1000,
		Food:    1000,
		Stone:   1000,
		Leather: 0,
	}

	ResourcesLimit = Resources{
		Wood:    2000,
		Food:    2000,
		Stone:   2000,
		Leather: 2000,
	}

	ChunkResourcesLimit = ChunkResources{
		Trees:   200,
		Stones:  200,
		Animals: 200,
		Plants:  200,
	}
)
