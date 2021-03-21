package logic

import (
	"abbysoft/gardarike-online/model"
	rpc "abbysoft/gardarike-online/rpc/generated"
	log "github.com/sirupsen/logrus"
)

func haveTown(towns []model.Town, id int64) bool {
	for _, town := range towns {
		if town.ID == id {
			return true
		}
	}

	return false
}

func (s *SimpleLogic) RenameTown(session *PlayerSession, request *rpc.RenameTownRequest) (*rpc.RenameTownResponse, model.Error) {
	s.log.WithFields(log.Fields{
		"sessionID": session.SessionID,
		"townID":    request.TownID,
		"newName":   request.NewName,
	}).Info("RenameTown request")

	towns, err := session.Tx.GetTowns(session.SelectedCharacter.Name)
	if err != nil {
		s.log.WithError(err).Error("Failed to get characters towns")
		return nil, model.ErrInternalServerError
	}

	if !haveTown(towns, request.TownID) {
		return nil, model.ErrNotAuthorized
	}

	if err := session.Tx.RenameTown(request.TownID, request.NewName); err != nil {
		s.log.WithError(err).Error("Failed to rename town")
		return nil, model.ErrInternalServerError
	}

	return &rpc.RenameTownResponse{}, nil
}
