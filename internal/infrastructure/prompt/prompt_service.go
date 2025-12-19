package prompt

import (
	"github.com/AlecAivazis/survey/v2"
	"github.com/ahmed-abdelgawad92/lockify/internal/domain/service"
)

type PromptService struct{}

func NewPromptService() service.PromptService {
	return &PromptService{}
}

func (p *PromptService) GetUserInputForKeyAndValue(isSecret bool) (key, value string) {
	prompt := &survey.Input{Message: "Enter key:"}
	survey.AskOne(prompt, &key)

	if isSecret {
		prompt := &survey.Password{Message: "Enter secret:"}
		survey.AskOne(prompt, &value)
	} else {
		prompt = &survey.Input{Message: "Enter value:"}
		survey.AskOne(prompt, &value)
	}

	return key, value
}
