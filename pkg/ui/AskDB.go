package ui

import (
	"github.com/AlecAivazis/survey/v2"
)

func AskDB(DBNames []string) (*string, error) {

	var qs = []*survey.Question{
		{
			Name: "color",
			Prompt: &survey.Select{
				Message: "Choose Target DB:",
				Options: DBNames,
			},
		},
	}

	answers := struct {
		DBName string `survey:"color"` // or you can tag fields to match a specific name
	}{}

	err := survey.Ask(qs, &answers)
	if err != nil {
		return nil, err
	}

	return &answers.DBName, nil
}
