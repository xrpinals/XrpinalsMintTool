package utils

import (
	"encoding/json"
	"fmt"
	"math/big"
	"reflect"
	"strconv"
	"strings"
)

type AssetInfoRsp struct {
	Id     int             `json:"id"`
	Result AssetInfoResult `json:"result"`
	Error  interface{}     `json:"error"`
}

type AssetInfoResult struct {
	Id         string `json:"id"`
	Symbol     string `json:"symbol"`
	Precision  int    `json:"precision"`
	Issuer     string `json:"issuer"`
	CreateAddr string `json:"create_addr"`
	Options    struct {
		MaxSupply         interface{} `json:"max_supply"`
		MaxPerMint        interface{} `json:"max_per_mint"`
		CreateTime        string      `json:"create_time"`
		Brc20Token        bool        `json:"brc20_token"`
		MintInterval      int         `json:"mint_interval"`
		MaxMintCountLimit interface{} `json:"max_mint_count_limit"`
		MarketFeePercent  int         `json:"market_fee_percent"`
		MaxMarketFee      interface{} `json:"max_market_fee"`
		IssuerPermissions int         `json:"issuer_permissions"`
		Flags             int         `json:"flags"`
	} `json:"options"`
	DynamicData struct {
		CurrentSupply interface{} `json:"current_supply"`
		CurrentNBits  uint32      `json:"current_nBits"`
	} `json:"dynamic_data"`
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
	if err != nil {
		return nil, err
	}

	if response.Error != nil {
		errStr, _ := json.Marshal(response.Error)
		return nil, fmt.Errorf(string(errStr))
	}

	return &response, nil
}

type AddressBalanceRsp struct {
	Id     int                    `json:"id"`
	Result []AddressBalanceResult `json:"result"`
	Error  interface{}            `json:"error"`
}

type AddressBalanceResult struct {
	Amount  interface{} `json:"amount"`
	AssetId string      `json:"asset_id"`
}

func GetAddressBalance(url, address, assetId string) (*big.Int, error) {
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
	if err != nil {
		return nil, err
	}

	if response.Error != nil {
		errStr, _ := json.Marshal(response.Error)
		return nil, fmt.Errorf(string(errStr))
	}

	for _, k := range response.Result {
		if k.AssetId == assetId {
			typeStr := reflect.TypeOf(k.Amount).String()
			if typeStr == "string" {
				balance, err := strconv.ParseUint(k.Amount.(string), 10, 64)
				return big.NewInt(int64(balance)), err
			} else {
				return big.NewInt(int64(k.Amount.(float64))), nil
			}
		}
	}
	return big.NewInt(0), nil
}

type InfoRsp struct {
	Id     int         `json:"id"`
	Result InfoResult  `json:"result"`
	Error  interface{} `json:"error"`
}

type InfoResult struct {
	ChainId string `json:"chain_id"`
}

func GetChainId(url string) (string, error) {
	assetInfoReq := RpcReq{
		Id:     1,
		Method: "info",
		Params: []interface{}{},
	}

	body, err := HttpClient{
		Timeout: 30,
	}.HttpPost(url, assetInfoReq)

	if err != nil {
		return "", err
	}

	var response InfoRsp

	err = json.Unmarshal(body, &response)
	if err != nil {
		return "", err
	}

	if response.Error != nil {
		errStr, _ := json.Marshal(response.Error)
		return "", fmt.Errorf(string(errStr))
	}

	return response.Result.ChainId, err
}

type RefBlockInfoRsp struct {
	Id     int         `json:"id"`
	Result string      `json:"result"`
	Error  interface{} `json:"error"`
}

func GetRefBlockInfo(url string) (uint16, uint32, error) {
	assetInfoReq := RpcReq{
		Id:     1,
		Method: "lightwallet_get_refblock_info",
		Params: []interface{}{},
	}

	body, err := HttpClient{
		Timeout: 30,
	}.HttpPost(url, assetInfoReq)

	if err != nil {
		return 0, 0, err
	}

	var response RefBlockInfoRsp

	err = json.Unmarshal(body, &response)
	if err != nil {
		return 0, 0, err
	}

	if response.Error != nil {
		errStr, _ := json.Marshal(response.Error)
		return 0, 0, fmt.Errorf(string(errStr))
	}

	l := strings.Split(response.Result, ",")
	refBlockNum, _ := strconv.ParseInt(l[0], 10, 64)
	refBlockPrefix, _ := strconv.ParseInt(l[1], 10, 64)

	return uint16(refBlockNum), uint32(refBlockPrefix), nil
}

type BroadcastTxRsp struct {
	Id     int         `json:"id"`
	Result string      `json:"result"`
	Error  interface{} `json:"error"`
}

func BroadcastTx(url string, signTx interface{}) (string, error) {
	BroadcastTxReq := RpcReq{
		Id:     1,
		Method: "lightwallet_broadcast",
		Params: []interface{}{signTx},
	}

	body, err := HttpClient{
		Timeout: 30,
	}.HttpPost(url, BroadcastTxReq)

	if err != nil {
		return "", err
	}

	var response BroadcastTxRsp

	err = json.Unmarshal(body, &response)
	if err != nil {
		return "", err
	}

	if response.Error != nil {
		errStr, _ := json.Marshal(response.Error)
		return "", fmt.Errorf(string(errStr))
	}

	return response.Result, err
}
