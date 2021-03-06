package model

import (
	rpc "abbysoft/gardarike-online/rpc/generated"
	"bytes"
	"encoding/gob"
	"fmt"
)

type EventWrapper struct {
	Topic string
	Event *rpc.Event
}

type Account struct {
	ID            int64  `db:"id"`
	Login         string `db:"login"`
	Password      string `db:"password"`
	Salt          string `db:"salt"`
	IsOnline      bool   `db:"is_online"`
	LastSessionID string `db:"last_session_id"`
}

type ChatMessage struct {
	ID       int64
	Sender   string
	Text     string
	IsSystem bool `db:"is_system"`
}

func (c ChatMessage) ToRPC() *rpc.ChatMessage {
	var messageType rpc.ChatMessage_Type
	if c.IsSystem {
		messageType = rpc.ChatMessage_SYSTEM
	} else {
		messageType = rpc.ChatMessage_NORMAL
	}

	return &rpc.ChatMessage{
		Id:     c.ID,
		Sender: c.Sender,
		Text:   c.Text,
		Type:   messageType,
	}
}

type Vector2D struct {
	X float32
	Y float32
}

func ToModelVector(vec *rpc.Vector2D) Vector2D {
	if vec == nil {
		return Vector2D{}
	}

	return Vector2D{
		X: vec.X,
		Y: vec.Y,
	}
}

func (v Vector2D) ToRPC() *rpc.Vector2D {
	return &rpc.Vector2D{
		X: v.X,
		Y: v.Y,
	}
}

type Town struct {
	ID         int64
	X          int64
	Y          int64
	OwnerName  string `db:"owner_name"`
	Population uint64
	Name       string
	Buildings  []Building
	Rotation   float32
}

func (t Town) ToRPC() *rpc.Town {
	return &rpc.Town{
		Id:         t.ID,
		X:          t.X,
		Y:          t.Y,
		Name:       t.Name,
		OwnerName:  t.OwnerName,
		Population: t.Population,
		Rotation:   t.Rotation,
	}
}

type ChunkRange struct {
	MinX int `db:"min_x"`
	MaxX int `db:"max_x"`
	MinY int `db:"min_y"`
	MaxY int `db:"max_y"`
}

type ChunkResources struct {
	Trees   uint64
	Stones  uint64
	Animals uint64
	Plants  uint64
}

type WorldMapChunk struct {
	Number int32
	X      int64
	Y      int64
	Width  int32
	Height int32
	Data   []byte
	Towns  []Town
	ChunkResources
}

func NewWorldMapChunkFromRPC(rpcChunk rpc.WorldMapChunk) (WorldMapChunk, error) {
	var terrain []byte
	result := WorldMapChunk{
		Number: 0,
		X:      int64(rpcChunk.X),
		Y:      int64(rpcChunk.Y),
		Width:  rpcChunk.Width,
		Height: rpcChunk.Height,
		Data:   nil,
		Towns:  nil,
		ChunkResources: ChunkResources{
			Trees:   rpcChunk.Trees,
			Stones:  rpcChunk.Stones,
			Animals: rpcChunk.Animals,
			Plants:  rpcChunk.Plants,
		},
	}

	buffer := bytes.NewBuffer(terrain)
	encoder := gob.NewEncoder(buffer)

	if err := encoder.Encode(&rpcChunk.Data); err != nil {
		return WorldMapChunk{}, fmt.Errorf("failed to encode map chunk: %w", err)
	}

	result.Data = buffer.Bytes()

	return result, nil
}

func (w WorldMapChunk) ToRPC() (*rpc.WorldMapChunk, error) {
	var terrain []float32
	decoder := gob.NewDecoder(bytes.NewBuffer(w.Data))
	if err := decoder.Decode(&terrain); err != nil {
		return nil, fmt.Errorf("failed to decode terrain data: %w", err)
	}

	mapChunk := &rpc.WorldMapChunk{
		X:       int32(w.X),
		Y:       int32(w.Y),
		Width:   w.Width,
		Height:  w.Height,
		Data:    terrain,
		Towns:   nil,
		Trees:   w.Trees,
		Stones:  w.Stones,
		Animals: w.Animals,
		Plants:  w.Plants,
	}

	for _, town := range w.Towns {
		mapChunk.Towns = append(mapChunk.Towns, town.ToRPC())
	}

	return mapChunk, nil
}

func (w WorldMapChunk) ToLocalRPC() (*rpc.LocalMapChunk, error) {
	worldChunk, err := w.ToRPC()
	if err != nil {
		return nil, err
	}

	return &rpc.LocalMapChunk{
		GlobalChunkCoords: &rpc.IntVector2D{X: worldChunk.X, Y: worldChunk.Y},
		LocalChunkNumber:  w.Number,
		Data:              worldChunk.Data,
		Width:             worldChunk.Height,
		Height:            worldChunk.Width,
	}, nil
}

type Character struct {
	ID                int64
	AccountID         int64 `db:"account_id"`
	Name              string
	MaxPopulation     uint64 `db:"max_population"`
	CurrentPopulation uint64 `db:"current_population"`
	Towns             []Town
	Resources         Resources
	ProductionRate    Resources
}

func (c Character) HasTown(townID int64) bool {
	for _, town := range c.Towns {
		if town.ID == townID {
			return true
		}
	}

	return false
}

func (c Character) ToRPC() *rpc.Character {
	return &rpc.Character{
		Id:                c.ID,
		Name:              c.Name,
		MaxPopulation:     c.MaxPopulation,
		CurrentPopulation: c.CurrentPopulation,
	}
}

type Resources struct {
	CharacterID int64 `db:"character_id"`
	Wood        uint64
	Food        uint64
	Stone       uint64
	Leather     uint64
}

func (r Resources) ToRPC() *rpc.Resources {
	return &rpc.Resources{
		Wood:    r.Wood,
		Stone:   r.Stone,
		Food:    r.Food,
		Leather: r.Leather,
	}
}

// Subtract - decrement resources by the provided values if there is enough
// resources or do nothing
// return true if the resources were subtracted
func (r *Resources) Subtract(resources Resources) bool {
	if !r.IsEnough(resources) {
		return false
	}

	r.Food -= resources.Food
	r.Wood -= resources.Wood
	r.Stone -= resources.Stone
	r.Leather -= resources.Leather

	return true
}

func (r *Resources) Add(resources Resources) {
	r.Food += resources.Food
	r.Wood += resources.Wood
	r.Stone += resources.Stone
	r.Leather += resources.Leather

	*r = minResources(*r, ResourcesLimit)
}

func minResources(a, b Resources) (r Resources) {
	r = a
	if a.Food > b.Food {
		r.Food = b.Food
	}
	if a.Wood > b.Wood {
		r.Wood = b.Wood
	}
	if a.Stone > b.Stone {
		r.Stone = b.Stone
	}
	if a.Leather > b.Leather {
		r.Leather = b.Leather
	}

	return
}

func (r Resources) IsEnough(requested Resources) bool {
	return r.Food >= requested.Food &&
		r.Stone >= requested.Stone &&
		r.Wood >= requested.Wood &&
		r.Leather >= requested.Leather
}

func (r Resources) IsLimitReached() bool {
	return r.IsEnough(ResourcesLimit)
}
