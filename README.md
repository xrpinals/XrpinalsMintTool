# XrpinalsMintTool


### How to use

* Import Private Key

```
./XrpinalsMintTool import_key -h
     
Usage of import_key:
  -key string
        your private key



./XrpinalsMintTool import_key -key 5JBvpSG46ipHruGae4oi7A4gE6sd8JxMXNKT4ZmsyF8Q5Cb6yjT

private key of address mnUbdaJcTiBUARHGMZpQ5dVkrcj1XUMame imported
 
```


* Check if the Address is in Storage

```
./XrpinalsMintTool check_address -h

Usage of check_address:
  -addr string
        your address


./XrpinalsMintTool check_address -addr mnUbdaJcTiBUARHGMZpQ5dVkrcj1XUMame

address mnUbdaJcTiBUARHGMZpQ5dVkrcj1XUMame is in the storage

```

* Get balance of Address

```
./XrpinalsMintTool get_balance -h

Usage of get_balance:
  -addr string
        your address
  -asset string
        asset name you want to query


./XrpinalsMintTool get_balance -addr mfhGJnP5T7A5kYDJNxnHozxrVzC7WKHzKs -asset BTC

balance:  0

```


* Get Cross Chain Deposit Address

```

./XrpinalsMintTool get_deposit_address 

BTC deposit address:  2MvSkSdKtHC253TnDmmBj8uwWjcyeQ5sQkK


```


* Mint Brc20 Asset

```
./XrpinalsMintTool mint -h

Usage of mint:
  -addr string
        your address
  -asset string
        asset name you want to mint


./XrpinalsMintTool mint -addr mnUbdaJcTiBUARHGMZpQ5dVkrcj1XUMame -asset XX

mining success, txHash:03d8216ce49753cbe6ae7a1a65e08b4fe841d5b6

```


* Transfer Brc20 Asset

```
./XrpinalsMintTool mint -h

Usage of mint:
  Usage of transfer:
  -from string
        your address
  -to string
        receiver address
  -amount string
        asset amount you want to transfer
  -asset string
        asset name you want to transfer


./XrpinalsMintTool transfer -from mfhGJnP5T7A5kYDJNxnHozxrVzC7WKHzKs -to n2frLX2z972TxP7vvqGEMUyi4CFxmmTnk8 -asset 1.3.0 -amount 2

transfer success,txHash:6eb8918df731f29952cc00a4ae77c0a07e907742

```


* Query Mint Brc20 Info

```
./XrpinalsMintTool get_mint_info -h

Usage of get_mint_info:
  -addr string
        your address
  -asset string
        asset name you want to query


./XrpinalsMintTool get_mint_info -addr mnUbdaJcTiBUARHGMZpQ5dVkrcj1XUMame -asset XX

mint info:
 mint amount: 0
 mint count: 0


```
