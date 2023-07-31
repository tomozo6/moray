package ui

import (
	"fmt"

	"github.com/AlecAivazis/survey/v2"
)

func AskBastion(EC2InstancesName []string) string{

	var qs = []*survey.Question{
		{
			Name: "color",
			Prompt: &survey.Select{
				Message: "Choose Bastion EC2:",
				Options: EC2InstancesName,
			},
		},
	}

	answers := struct {
		FavoriteColor string `survey:"color"` // or you can tag fields to match a specific name
	}{}

	err := survey.Ask(qs, &answers)
	if err != nil {
		fmt.Println(err.Error())
	}

	return answers.FavoriteColor
}
