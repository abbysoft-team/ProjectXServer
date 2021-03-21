package db

import (
	"abbysoft/gardarike-online/model"
	rpc "abbysoft/gardarike-online/rpc/generated"
)

type CharacterDatabaseTransaction interface {
	GetCharacter(id int64) (model.Character, error)
	AddCharacter(name string) (id int, err error)
	AddAccountCharacter(characterID, accountID int) error
	DeleteCharacter(id int64) error
	GetCharacters(accountID int64) ([]model.Character, error)
	UpdateCharacter(character model.Character) error
	GetResources(characterID int64) (model.Resources, error)
	GetProductionRates(characterID int64) (model.Resources, error)
	GetEmpiresByCriteria(
		characterName string, offset, limit uint32, criteria rpc.EmpiresRatingCriteria) ([]*rpc.RatingEntry, *rpc.RatingEntry, error)
}

type AccountDatabaseTransaction interface {
	GetAccount(login string) (model.Account, error)
	AddAccount(login string, password string, salt string) (int, error)
}

type WorldDatabaseTransaction interface {
	AddChatMessage(message model.ChatMessage) (int64, error)
	GetChatMessages(offset int, count int) ([]model.ChatMessage, error)
	GetMapChunk(x, y int64) (model.WorldMapChunk, error)
	GetChunkRange() (model.ChunkRange, error)
	IncrementMapResources(resources model.ChunkResources, limit model.ChunkResources) error
	SaveMapChunkOrUpdate(chunk model.WorldMapChunk) error
	GetTowns(ownerName string) ([]model.Town, error)
	GetAllTowns() ([]model.Town, error)
	GetTownsForRect(xStart, xEnd, yStart, yEnd int) ([]model.Town, error)
	AddOrUpdateResources(resources model.Resources) error
	AddOrUpdateProductionRates(rates model.Resources) error
	AddTown(town model.Town) error
	AddTownBuilding(townID int64, building model.Building) error
	GetAllBuildings() (map[int64]model.CharacterBuildings, error)
	RenameTown(townID int64, newName string) error
}

type DatabaseTransaction interface {
	CharacterDatabaseTransaction
	AccountDatabaseTransaction
	WorldDatabaseTransaction

	EndTransaction() error
	IsCompleted() bool
	IsFailed() bool
	IsSucceed() bool
	SetAutoCommit(value bool)
	SetAutoRollBack(value bool)
}

type Database interface {
	BeginTransaction(autoCommit bool, autoRollBack bool) (DatabaseTransaction, error)
}
