package ui

import (
	"github.com/AlecAivazis/survey/v2"
)

func AskEC2(EC2InstancesName []string) (*string, error) {

	var qs = []*survey.Question{
		{
			Name: "color",
			Prompt: &survey.Select{
				Message: "Choose EC2 Instance:",
				Options: EC2InstancesName,
			},
		},
	}

	answers := struct {
		FavoriteColor string `survey:"color"` // or you can tag fields to match a specific name
	}{}

	err := survey.Ask(qs, &answers)
	if err != nil {
		return nil, err
	}

	return &answers.FavoriteColor, nil
}
