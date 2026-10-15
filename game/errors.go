package game

import "errors"

var (
	ErrWrongPhase      = errors.New("wrong game phase")
	ErrNoRolls         = errors.New("no rolls remaining")
	ErrInvalidDiceRoll = errors.New("invalid dice roll")
	ErrInvalidOption   = errors.New("invalid score slot")
	ErrNoDiceRolled    = errors.New("no dice have been rolled")
)
