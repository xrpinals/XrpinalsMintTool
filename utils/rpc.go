package utils

import (
	"encoding/json"
)

type BaseRsp struct {
	Id     int         `json:"id"`
	Result interface{} `json:"result"`
}

type AssetInfoResult struct {
	Id         string `json:"id"`
	Symbol     string `json:"symbol"`
	Precision  int    `json:"precision"`
	Issuer     string `json:"issuer"`
	CreateAddr string `json:"create_addr"`
	Options    struct {
		MaxSupply         string `json:"max_supply"`
		MaxPerMint        int    `json:"max_per_mint"`
		CreateTime        string `json:"create_time"`
		BosToken          bool   `json:"bos_token"`
		MintInterval      int    `json:"mint_interval"`
		MaxMintCountLimit int    `json:"max_mint_count_limit"`
		MarketFeePercent  int    `json:"market_fee_percent"`
		MaxMarketFee      string `json:"max_market_fee"`
		IssuerPermissions int    `json:"issuer_permissions"`
		Flags             int    `json:"flags"`
	} `json:"options"`
}

type AssetInfoRsp struct {
	BaseRsp
	AssetInfoResult
}

func GetAssetInfo(url string, assetName string) (rsp *AssetInfoRsp, err error) {
	assetInfoReq := RpcReq{
		Id:     1,
		Method: "get_asset_imp",
		Params: []interface{}{assetName},
	}

	body, err := HttpClient{
		Timeout: 30,
	}.HttpPost(url, assetInfoReq)

	if err != nil {
		return nil, err
	}

	var response AssetInfoRsp

	err = json.Unmarshal(body, &response)

	return &response, nil
}

type AddressBalanceRsp struct {
	BaseRsp
	Result []AddressBalanceResult `json:"result"`
}

type AddressBalanceResult struct {
	Amount  interface{} `json:"amount"`
	AssetId string      `json:"asset_id"`
}

func GetAddressBalance(url, address string) (rsp *AddressBalanceRsp, err error) {
	assetInfoReq := RpcReq{
		Id:     1,
		Method: "get_addr_balances",
		Params: []interface{}{address},
	}

	body, err := HttpClient{
		Timeout: 30,
	}.HttpPost(url, assetInfoReq)

	if err != nil {
		return nil, err
	}

	var response AddressBalanceRsp

	err = json.Unmarshal(body, &response)

	return &response, nil
}
