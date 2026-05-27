package errors

import (
	"github.com/google/uuid"
	sf "github.com/wissance/stringFormatter"
)

type AttackerStatDataNotFoundError struct {
	attackerSign string
	statDataId   uuid.UUID
}

func NewAttackerStatDataNotFoundError(attackerSign string, statDataId uuid.UUID) AttackerStatDataNotFoundError {
	return AttackerStatDataNotFoundError{
		attackerSign: attackerSign,
		statDataId:   statDataId,
	}
}

func (e AttackerStatDataNotFoundError) Error() string {
	return sf.Format("Attacker stat data for ID \"{0}\" was not found, however sign :\"{1}\" (device id or IP address exists)",
		e.attackerSign, e.statDataId)
}

type AttackerNotFoundError struct {
	attackerSign string
}

func NewAttackerNotFoundError(attackerSign string) AttackerNotFoundError {
	return AttackerNotFoundError{
		attackerSign: attackerSign,
	}
}

func (e AttackerNotFoundError) Error() string {
	return sf.Format("Attacker with sign \"{0}\" was not found", e.attackerSign)
}
