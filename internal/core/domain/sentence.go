package domain

type SentenceLevelType string

const (
	SentenceLevelEasy   SentenceLevelType = "EASY"
	SentenceLevelNormal SentenceLevelType = "NORMAL"
	SentenceLevelHard   SentenceLevelType = "HARD"
)

type Sentence struct {
	Base
	Modifier

	Text    string
	Grammar Grammar
	Level   SentenceLevelType
	Status  StatusType
}
