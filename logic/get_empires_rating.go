package logic

import (
	"abbysoft/gardarike-online/model"
	rpc "abbysoft/gardarike-online/rpc/generated"
	log "github.com/sirupsen/logrus"
)

func (s *SimpleLogic) GetEmpiresRating(session *PlayerSession, request *rpc.GetEmpiresRatingRequest) (*rpc.GetEmpiresRatingResponse, model.Error) {
	s.log.WithFields(log.Fields{
		"sessionID": request.SessionID,
		"offset":    request.Offset,
		"limit":     request.Limit,
		"criteria":  request.Criteria,
	}).Info("GetEmpiresRating")

	entries, playerEntry, err := session.Tx.GetEmpiresByCriteria(
		session.SelectedCharacter.Name, request.Offset, request.Limit, request.Criteria)

	if err != nil {
		s.log.WithError(err).Error("Failed to get empires rating")
		return nil, model.ErrInternalServerError
	}

	return &rpc.GetEmpiresRatingResponse{
		Entries:      entries,
		PlayerRating: playerEntry,
	}, nil
}
