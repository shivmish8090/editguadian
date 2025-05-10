package database

import (
	"context"
	"fmt"
	"strings"

	"go.mongodb.org/mongo-driver/v2/bson"
	"go.mongodb.org/mongo-driver/v2/mongo/options"

	"main/config"
)

type EchoSettings struct {
	ChatID int64  `bson:"chat_id"`
	Mode   string `bson:"mode"`
	Limit  int    `bson:"limit"`
}

const (
	defaultEchoLimit = 800
	defaultEchoMode  = "MANUAL"
)

func SetEchoSettings(data *EchoSettings) error {
	key := fmt.Sprintf("echos:%d", data.ChatID)

	if val, ok := config.Cache.Load(key); ok {
		existing := val.(*EchoSettings)
		if data.Mode == "" {
			data.Mode = existing.Mode
		}
		if data.Limit == 0 {
			data.Limit = existing.Limit
		}
		if data.Mode == existing.Mode && data.Limit == existing.Limit {
			return nil
		}
	} else {
		existing, _ := GetEchoSettings(data.ChatID)
		if data.Mode == "" {
			data.Mode = existing.Mode
		}
		if data.Limit == 0 {
			data.Limit = existing.Limit
		}
	}

	ctx, cancel := context.WithTimeout(context.Background(), timeout)
	defer cancel()
	data.Mode = strings.ToUpper(data.Mode)

	update := bson.M{
		"$set": bson.M{
			"mode":  data.Mode,
			"limit": data.Limit,
		},
	}

	_, err := echoDB.UpdateOne(ctx, bson.M{"chat_id": data.ChatID}, update, options.UpdateOne().SetUpsert(true))
	if err != nil {
		return err
	}

	go config.Cache.Store(key, &EchoSettings{
		ChatID: data.ChatID,
		Mode:   data.Mode,
		Limit:  data.Limit,
	})

	return nil
}

func GetEchoSettings(chatID int64) (*EchoSettings, error) {
	key := fmt.Sprintf("echos:%d", chatID)
	if val, ok := config.Cache.Load(key); ok {
		return val.(*EchoSettings), nil
	}

	ctx, cancel := context.WithTimeout(context.Background(), timeout)
	defer cancel()

	var settings EchoSettings
	err := echoDB.FindOne(ctx, bson.M{"chat_id": chatID}).Decode(&settings)
	if err != nil {
		settings = EchoSettings{
			ChatID: chatID,
			Mode:   defaultEchoMode,
			Limit:  defaultEchoLimit,
		}
	}

	if settings.Mode == "" {
		settings.Mode = defaultEchoMode
	}
	if settings.Limit == 0 {
		settings.Limit = defaultEchoLimit
	}

	config.Cache.Store(key, &settings)

	return &settings, nil
}
