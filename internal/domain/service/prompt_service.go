package service

// PromptService defines the interface for prompting user input.
type PromptService interface {
	GetUserInputForKeyAndValue(isSecret bool) (key, value string, err error)
	GetPassphraseInput(message string) (string, error)
	GetConfirmation(message string, defaultValue bool) (bool, error)
}
