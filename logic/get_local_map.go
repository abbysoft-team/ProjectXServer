package logic

import (
	"abbysoft/gardarike-online/model"
	rpc "abbysoft/gardarike-online/rpc/generated"
	"database/sql"
	"errors"
	log "github.com/sirupsen/logrus"
)

func getGlobalChunkCoordsForPosition(location rpc.Vector2D, mapChunkSize int) (int64, int64, int64) {
	xOffset := int64(int(location.X) % mapChunkSize)
	yOffset := int64(int(location.Y) % mapChunkSize)

	if xOffset < 0.0 {
		xOffset = int64(mapChunkSize) + xOffset
	}
	if yOffset < 0.0 {
		yOffset = int64(mapChunkSize) + yOffset
	}

	globalX := int64(int(location.X) / mapChunkSize)
	globalY := int64(int(location.Y) / mapChunkSize)

	if globalX == 0 && location.X < 0 {
		globalX = -1
	}
	if globalY == 0 && location.Y < 0 {
		globalY = -1
	}

	//chunkStartX := globalX * int64(s.MapChunkSize())
	//chunkStartY := globalY * int64(s.MapChunkSize())

	halfChunk := int64(mapChunkSize / 2)

	isLeftPartOfChunk := func() bool {
		return xOffset < halfChunk
	}

	isUpperPartOfChunk := func() bool {
		return yOffset >= halfChunk
	}

	var localChunkNumber int64 = 0

	// [x][.]
	// [x][.]
	// 1 or 3 square
	if isLeftPartOfChunk() {
		if isUpperPartOfChunk() {
			localChunkNumber = 1
		} else {
			localChunkNumber = 3
		}
	} else { // 2 or 4 square
		if isUpperPartOfChunk() {
			localChunkNumber = 2
		} else {
			localChunkNumber = 4
		}
	}

	return globalX, globalY, localChunkNumber
}

func (s *SimpleLogic) generateNewLocalChunk(session *PlayerSession, x, y int64, number int32) (*rpc.GetLocalMapResponse, model.Error) {
	var offsetX, offsetY float64
	if number == 2 || number == 4 {
		offsetX = float64(s.MapChunkSize())/2 - 1
	}

	if number == 1 || number == 2 {
		offsetY = float64(s.MapChunkSize())/2 - 1
	}

	data := s.localGenerator.GenerateTerrain(s.MapChunkSize(), s.MapChunkSize(), offsetX, offsetY)

	chunk := rpc.WorldMapChunk{
		X:          int32(x),
		Y:          int32(y),
		Width:      int32(s.MapChunkSize()),
		Height:     int32(s.MapChunkSize()),
		Data:       data,
		Towns:      nil,
		Trees:      0,
		Stones:     0,
		Animals:    0,
		Plants:     0,
		WaterLevel: s.config.WaterLevel,
	}

	modelChunk, err := model.NewWorldMapChunkFromRPC(chunk)
	if err != nil {
		s.log.WithError(err).Error("Failed to convert rpc chunk to model chunk")
		return nil, model.ErrInternalServerError
	}

	modelChunk.Number = number
	if err := session.Tx.SaveMapChunkOrUpdate(modelChunk); err != nil {
		s.log.WithError(err).Error("Failed to save local map chunk")
		return nil, model.ErrInternalServerError
	}

	return &rpc.GetLocalMapResponse{
		Map: &rpc.LocalMapChunk{
			GlobalChunkCoords: &rpc.IntVector2D{X: chunk.X, Y: chunk.Y},
			LocalChunkNumber:  number,
			Data:              chunk.Data,
			Width:             chunk.Width,
			Height:            chunk.Height,
		},
	}, nil
}

func (s *SimpleLogic) GetLocalMap(session *PlayerSession, request *rpc.GetLocalMapRequest) (*rpc.GetLocalMapResponse, model.Error) {
	s.log.WithFields(log.Fields{
		"sessionID": session.SessionID,
		"location":  *request.Location,
	}).Info("GetLocalMap")

	globalChunkX, globalChunkY, number := getGlobalChunkCoordsForPosition(*request.Location, s.MapChunkSize())
	chunk, err := session.Tx.GetMapChunk(globalChunkX, globalChunkY, number)
	if err != nil && errors.Is(err, sql.ErrNoRows) {
		return s.generateNewLocalChunk(session, globalChunkX, globalChunkY, int32(number))
	} else if err != nil {
		s.log.WithError(err).Error("Failed to get map chunk")
		return nil, model.ErrInternalServerError
	}

	local, err := chunk.ToLocalRPC()
	if err != nil {
		s.log.WithError(err).Error("Failed to convert chunk to local")
		return nil, model.ErrInternalServerError
	}

	return &rpc.GetLocalMapResponse{Map: local}, nil
}
