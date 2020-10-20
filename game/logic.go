package game

import (
	"errors"
	"fmt"
	"github.com/sirupsen/logrus"
	"projectx-server/model"
	"projectx-server/model/postgres"
	rpc "projectx-server/rpc/generated"
)

const (
	mapChunkSize = 100
)

var ErrNotAuthorized = errors.New("user not authorized")

type Logic interface {
	GetMap(request *rpc.GetMapRequest) (*rpc.GetMapResponse, error)
	Login(request *rpc.LoginRequest) (*rpc.LoginResponse, error)
	SelectCharacter(request *rpc.SelectCharacterRequest) (*rpc.SelectCharacterResponse, error)
}

type SimpleLogic struct {
	gameMap  rpc.Map
	db       model.Database
	log      *logrus.Entry
	sessions map[string]*PlayerSession
}

func NewLogic(generator TerrainGenerator) (*SimpleLogic, error) {
	width := mapChunkSize
	height := mapChunkSize

	terrain := generator.GenerateTerrain(mapChunkSize, mapChunkSize)
	database, err := postgres.NewDatabase(postgres.Config{
		Port:      5432,
		Host:      "localhost",
		User:      "admin",
		Password:  "admin",
		DBName:    "game",
		EnableSSL: false,
	})
	if err != nil {
		return nil, fmt.Errorf("failed to init db: %w", err)
	}

	return &SimpleLogic{
		gameMap: rpc.Map{
			Width:  int32(width),
			Height: int32(height),
			Points: terrain,
		},
		db:       database,
		log:      logrus.WithField("module", "logic"),
		sessions: make(map[string]*PlayerSession),
	}, nil
}

func (s *SimpleLogic) GetMap(request *rpc.GetMapRequest) (*rpc.GetMapResponse, error) {
	s.log.WithField("location", request.GetLocation()).
		WithField("sessionID", request.GetSessionID()).
		Debugf("GetMap request")

	_, authorized := s.sessions[request.GetSessionID()]
	if !authorized {
		return nil, ErrNotAuthorized
	}

	return &rpc.GetMapResponse{Map: &s.gameMap}, nil
}

func (s *SimpleLogic) SelectCharacter(request *rpc.SelectCharacterRequest) (*rpc.SelectCharacterResponse, error) {
	s.log.WithField("characterID", request.GetCharacterID()).
		WithField("sessionID", request.GetSessionID()).
		Debugf("SelectCharacter request")

	session, authorized := s.sessions[request.GetSessionID()]
	if !authorized {
		return nil, ErrNotAuthorized
	}

	char, err := s.db.GetCharacter(int(request.GetCharacterID()))
	if err != nil {
		return nil, fmt.Errorf("character not found")
	}

	session.SelectedCharacter = &char
	s.log.WithFields(logrus.Fields{
		"accountID": request.GetSessionID(),
		"character": char,
	}).Info("User selected character")

	return &rpc.SelectCharacterResponse{}, nil
}
