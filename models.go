package discord_chatbot

import (
	"context"
	"errors"
	"fmt"
	"os"

	yaml "github.com/goccy/go-yaml"
	"github.com/rs/zerolog/log"

	"github.com/taythebot/discord_chatbot/provider"
)

// Model yaml struct
type Model struct {
	Name  string `yaml:"name"`
	Value string `yaml:"value"`
}

// ParseModels parses and validates models from the YAML file
func ParseModels(ctx context.Context, file string, provider *provider.Client) ([]Model, error) {
	data, err := os.ReadFile(file)
	if err != nil {
		return nil, fmt.Errorf("failed to read file: %w", err)
	}

	var models []Model
	if err := yaml.Unmarshal(data, &models); err != nil {
		return nil, fmt.Errorf("failed to unmarshal yaml: %w", err)
	}

	if len(models) == 0 {
		return nil, errors.New("no models found")
	} else if len(models) > 25 {
		return nil, errors.New("cannot load more than 25 models")
	}

	// Validate all models with provider
	apiModels, err := provider.GetModels(ctx, log.Logger)
	if err != nil {
		return nil, fmt.Errorf("failed to get models from provider: %w", err)
	}

	for _, model := range models {
		var found bool
		for _, apiModel := range apiModels.Data {
			if apiModel.ID == model.Value {
				found = true
				break
			}
		}

		if !found {
			return nil, fmt.Errorf("failed to find model '%s' at provider", model.Name)
		}
	}

	return models, nil
}

var DefaultModelPrompt = "You are a very loyal and helpful assistant. You will always answer each question and request to the best of your abilities! Each message from this point on will have the user's name and message in the format `$user: $message`."
