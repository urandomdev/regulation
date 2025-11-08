package config

import (
	"strconv"

	validation "github.com/go-ozzo/ozzo-validation/v4"
	"github.com/redis/rueidis"
)

type Redis struct {
	Host     string `json:"host"`
	Port     int    `json:"port"`
	Database int    `json:"database"`
	Password string `json:"password"`
}

func (r *Redis) GetAddress() string {
	return r.Host + ":" + strconv.Itoa(r.Port)
}

func (r *Redis) GetURL() string {
	if r.Password != "" {
		return "redis://" + r.Password + "@" + r.GetAddress() + "/" + strconv.Itoa(r.Database)
	}
	return "redis://" + r.GetAddress() + "/" + strconv.Itoa(r.Database)
}

func (r *Redis) ClientOption() rueidis.ClientOption {
	opt := rueidis.ClientOption{
		InitAddress: []string{r.GetAddress()},
	}

	if r.Password != "" {
		opt.Password = r.Password
	}

	if r.Database != 0 {
		opt.SelectDB = r.Database
	}

	return opt
}

func (r *Redis) Client() rueidis.Client {
	client, err := rueidis.NewClient(r.ClientOption())
	if err != nil {
		panic(err)
	}

	return client
}

func (r *Redis) Validate() error {
	return validation.ValidateStruct(r,
		validation.Field(&r.Host, validation.Required),
		validation.Field(&r.Port, validation.Required, validation.Min(1), validation.Max(65535)),
		validation.Field(&r.Database, validation.Min(0), validation.Max(255)),
		validation.Field(&r.Password, validation.Length(0, 128)),
	)
}
