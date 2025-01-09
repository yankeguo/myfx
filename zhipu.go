package myfx

import (
	"github.com/go-resty/resty/v2"
	"github.com/yankeguo/zhipu"
)

func NewZhipuClient(client *resty.Client) (*zhipu.Client, error) {
	return zhipu.NewClient(zhipu.WithRestyClient(client))
}
