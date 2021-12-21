package models

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"time"

	"github.com/brndedhero/finance/config"
	"github.com/brndedhero/finance/helpers"
	"github.com/opensearch-project/opensearch-go/opensearchapi"
	"github.com/sirupsen/logrus"
)

type Account struct {
	Id        uint64     `json:"id" gorm:"primary_key"`
	CreatedAt time.Time  `json:"createdAt"`
	UpdatedAt time.Time  `json:"updatedAt"`
	DeletedAt *time.Time `json:"-" sql:"index"`
	Name      string     `json:"name"`
	Balance   float32    `json:"balance"`
}

func GetAllAccounts() (string, error) {
	var accounts []Account
	data, err := config.Redis.Get(config.Redis.Context(), config.RedisKeyAll).Result()
	if err != nil {
		config.Log.WithFields(logrus.Fields{
			"app":      "redis",
			"func":     "getAllaccount",
			"redisKey": config.RedisKeyAll,
		}).Warn("key not found")
		res := config.DB.Find(&accounts)
		if res.Error != nil {
			config.Log.WithFields(logrus.Fields{
				"app":  "mysql",
				"func": "getAllaccount",
			}).Error(res.Error)
			message, _ := helpers.PrepareString(404, nil)
			return message, res.Error
		}
		json, err := json.Marshal(accounts)
		if err != nil {
			config.Log.WithFields(logrus.Fields{
				"app":  config.AppName,
				"func": "getAllaccount",
			}).Error(err)
			message, _ := helpers.PrepareString(500, nil)
			return message, err
		}
		redisErr := config.Redis.Set(config.Redis.Context(), config.RedisKeyAll, json, 0).Err()
		if redisErr != nil {
			config.Log.WithFields(logrus.Fields{
				"app":      "redis",
				"func":     "getAllaccount",
				"redisKey": config.RedisKeyAll,
			}).Warn(redisErr)
		}
		message, _ := helpers.PrepareString(200, json)
		return message, nil
	}
	message, _ := helpers.PrepareString(200, []byte(data))
	return message, nil
}

func CreateAccount(name string, balance float32) (string, error) {
	account := Account{Name: name, Balance: balance}
	res := config.DB.Create(&account)
	if res.Error != nil {
		config.Log.WithFields(logrus.Fields{
			"app":  "mysql",
			"func": "createAccount",
		}).Error(res.Error)
		message := helpers.PrepareErrorString(500, res.Error)
		return message, res.Error
	}
	redisKey := fmt.Sprintf("%s:%d", config.RedisKey, account.Id)
	json, err := json.Marshal(account)
	if err != nil {
		config.Log.WithFields(logrus.Fields{
			"app":  config.AppName,
			"func": "createAccount",
		}).Error(err)
		message := helpers.PrepareErrorString(500, err)
		return message, err
	}
	redisErr := config.Redis.Set(config.Redis.Context(), redisKey, json, 0).Err()
	if redisErr != nil {
		config.Log.WithFields(logrus.Fields{
			"app":      "redis",
			"func":     "createAccount",
			"redisKey": redisKey,
		}).Warn(redisErr)
	}
	redisDelErr := config.Redis.Del(config.Redis.Context(), config.RedisKeyAll).Err()
	if redisDelErr != nil {
		config.Log.WithFields(logrus.Fields{
			"app":      "redis",
			"func":     "createAccount",
			"redisKey": config.RedisKeyAll,
		}).Warn(redisDelErr)
	}
	opensearchReq := opensearchapi.IndexRequest{
		Index:      config.OpensearchIndex,
		DocumentID: fmt.Sprint(account.Id),
		Body:       bytes.NewReader(json),
	}
	_, opensearchErr := opensearchReq.Do(context.Background(), config.Opensearch)
	if opensearchErr != nil {
		config.Log.WithFields(logrus.Fields{
			"app":  "opensearch",
			"func": "createAccount",
		}).Warn(opensearchErr)
	}
	data, _ := helpers.PrepareCreateJson(int(account.Id), int(res.RowsAffected))
	message, _ := helpers.PrepareString(201, data)
	return message, nil
}

func GetAccount(id uint64) (string, error) {
	var account Account
	redisKey := fmt.Sprintf("%s:%d", config.RedisKey, id)
	data, err := config.Redis.Get(config.Redis.Context(), redisKey).Result()
	if err != nil {
		config.Log.WithFields(logrus.Fields{
			"app":      "redis",
			"func":     "getAccount",
			"redisKey": redisKey,
		}).Warn("key not found")
		res := config.DB.First(&account, id)
		if res.Error != nil {
			config.Log.WithFields(logrus.Fields{
				"app":  "mysql",
				"func": "getAccount",
			}).Error(res.Error)
			message, _ := helpers.PrepareString(404, nil)
			return message, res.Error
		}
		json, err := json.Marshal(account)
		if err != nil {
			config.Log.WithFields(logrus.Fields{
				"app":  config.AppName,
				"func": "getAccount",
			}).Error(err)
			message := helpers.PrepareErrorString(500, err)
			return message, err
		}
		redisErr := config.Redis.Set(config.Redis.Context(), redisKey, json, 0).Err()
		if redisErr != nil {
			config.Log.WithFields(logrus.Fields{
				"app":      "redis",
				"func":     "getAccount",
				"redisKey": redisKey,
			}).Warn(redisErr)
			redisDelErr := config.Redis.Del(config.Redis.Context(), redisKey).Err()
			if redisDelErr != nil {
				config.Log.WithFields(logrus.Fields{
					"app":      "redis",
					"func":     "getAccount",
					"redisKey": redisKey,
				}).Warn(redisDelErr)
			}
		}
		message, _ := helpers.PrepareString(200, json)
		return message, nil
	}
	message, _ := helpers.PrepareString(200, []byte(data))
	return message, nil
}

func UpdateAccount(id uint64, name string, balance float32) (string, error) {
	var account Account
	res := config.DB.First(&account, id)
	if res.Error != nil {
		config.Log.WithFields(logrus.Fields{
			"app":  "mysql",
			"func": "updateAccount",
		}).Error(res.Error)
		message, _ := helpers.PrepareString(404, nil)
		return message, res.Error
	}
	account.Name = name
	account.Balance = balance
	res = config.DB.Save(&account)
	if res.Error != nil {
		config.Log.WithFields(logrus.Fields{
			"app":  "mysql",
			"func": "updateAccount",
		}).Error(res.Error)
		message, _ := helpers.PrepareString(500, nil)
		return message, res.Error
	}

	jsonData, err := json.Marshal(account)
	if err != nil {
		config.Log.WithFields(logrus.Fields{
			"app":  config.AppName,
			"func": "updateAccount",
		}).Error(err)
		message := helpers.PrepareErrorString(500, err)
		return message, err
	}
	redisKey := fmt.Sprintf("%s:%d", config.RedisKey, account.Id)
	redisErr := config.Redis.Set(config.Redis.Context(), redisKey, jsonData, 0).Err()
	if redisErr != nil {
		config.Log.WithFields(logrus.Fields{
			"app":      "redis",
			"func":     "updateAccount",
			"redisKey": redisKey,
		}).Warn(redisErr)
		redisDelErr := config.Redis.Del(config.Redis.Context(), redisKey)
		if redisDelErr != nil {
			config.Log.WithFields(logrus.Fields{
				"app":      "redis",
				"func":     "updateAccount",
				"redisKey": redisKey,
			}).Warn(redisDelErr)
		}
	}
	redisDelErr := config.Redis.Del(config.Redis.Context(), config.RedisKeyAll)
	if redisDelErr != nil {
		config.Log.WithFields(logrus.Fields{
			"app":      "redis",
			"func":     "updateAccount",
			"redisKey": config.RedisKeyAll,
		}).Warn(redisDelErr)
	}
	opensearchReq := opensearchapi.IndexRequest{
		Index:      config.OpensearchIndex,
		DocumentID: fmt.Sprint(account.Id),
		Body:       bytes.NewReader(jsonData),
	}
	_, opensearchErr := opensearchReq.Do(context.Background(), config.Opensearch)
	if opensearchErr != nil {
		config.Log.WithFields(logrus.Fields{
			"app":  "opensearch",
			"func": "updateAccount",
		}).Warn(opensearchErr)
	}
	data, _ := helpers.PrepareCreateJson(int(account.Id), int(res.RowsAffected))
	message, _ := helpers.PrepareString(200, data)
	return message, nil
}

func DeleteAccount(id uint64) (string, error) {
	var account Account
	res := config.DB.First(&account, id)
	if res.Error != nil {
		config.Log.WithFields(logrus.Fields{
			"app":  "mysql",
			"func": "deleteAccount",
		}).Error(res.Error)
		message, _ := helpers.PrepareString(404, nil)
		return message, res.Error
	}
	res = config.DB.Delete(&account)
	if res.Error != nil {
		config.Log.WithFields(logrus.Fields{
			"app":  "mysql",
			"func": "deleteAccount",
		}).Error(res.Error)
		message, _ := helpers.PrepareString(500, nil)
		return message, res.Error
	}
	redisKey := fmt.Sprintf("%s:%d", config.RedisKey, account.Id)
	redisErr := config.Redis.Del(config.Redis.Context(), redisKey).Err()
	if redisErr != nil {
		config.Log.WithFields(logrus.Fields{
			"app":      "redis",
			"func":     "deleteAccount",
			"redisKey": redisKey,
		}).Warn(redisErr)
	}
	redisDelErr := config.Redis.Del(config.Redis.Context(), config.RedisKeyAll).Err()
	if redisDelErr != nil {
		config.Log.WithFields(logrus.Fields{
			"app":      "redis",
			"func":     "deleteAccount",
			"redisKey": config.RedisKeyAll,
		}).Warn(redisDelErr)
	}
	opensearchReq := opensearchapi.DeleteRequest{
		Index:      config.OpensearchIndex,
		DocumentID: fmt.Sprint(account.Id),
	}
	_, opensearchErr := opensearchReq.Do(context.Background(), config.Opensearch)
	if opensearchErr != nil {
		config.Log.WithFields(logrus.Fields{
			"app":  "opensearch",
			"func": "deleteAccount",
		}).Warn(opensearchErr)
	}
	data, _ := helpers.PrepareCreateJson(int(account.Id), int(res.RowsAffected))
	message, _ := helpers.PrepareString(200, data)
	return message, nil
}
