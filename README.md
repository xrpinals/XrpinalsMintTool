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