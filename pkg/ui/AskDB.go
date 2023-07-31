package ui

import (
	"fmt"

	"github.com/AlecAivazis/survey/v2"
)

func AskDB(DBNames []string) string{

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
		fmt.Println(err.Error())
	}

	return answers.DBName
}
