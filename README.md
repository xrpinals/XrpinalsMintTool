# XrpinalsMintTool


![](https://github.com/xrpinals/XrpinalsMintTool/blob/main/assets/xrpinals.jpg)


### How To Build

* Install Golang Compiler

  [How To Install](https://go.dev/doc/install).
  It is strongly recommended to install the **Latest** version.

* Build XrpinalsMintTool
  ```
  git clone https://github.com/xrpinals/XrpinalsMintTool.git
  cd XrpinalsMintTool
  make XrpinalsMintTool_clean
  make XrpinalsMintTool

  cp conf_example.json conf.json   # and update this field "walletRpcUrl" to "http://api.xrpinals.com:1222"
  ```
  

### How To Use

* Import Private Key

```
./XrpinalsMintTool import_key -h
     
Usage of import_key:
  -key string
        your private key



./XrpinalsMintTool import_key -key 5JBvpSG46ipHruGae4oi7A4gE6sd8JxMXNKT4ZmsyF8Q5Cb6yjT

private key of address mnUbdaJcTiBUARHGMZpQ5dVkrcj1XUMame imported
 
```
![](https://github.com/xrpinals/XrpinalsMintTool/blob/main/assets/import-key.png)

* Check if the Address is in Storage

```
./XrpinalsMintTool check_address -h

Usage of check_address:
  -addr string
        address you want to check


./XrpinalsMintTool check_address -addr mnUbdaJcTiBUARHGMZpQ5dVkrcj1XUMame

address mnUbdaJcTiBUARHGMZpQ5dVkrcj1XUMame is in the storage

```
![](https://github.com/xrpinals/XrpinalsMintTool/blob/main/assets/check-address.png)

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
![](https://github.com/xrpinals/XrpinalsMintTool/blob/main/assets/get-balance.png)


* Get Cross Chain Deposit Address

```
./XrpinalsMintTool get_deposit_address -h 

Usage of get_deposit_address:
  -addr string
        your address


./XrpinalsMintTool get_deposit_address -addr mumtmaYKH3ttGpVaAJRCiiWZsn5zAB9hU

BTC deposit address:  2MvSkSdKtHC253TnDmmBj8uwWjcyeQ5sQkK

```
![](https://github.com/xrpinals/XrpinalsMintTool/blob/main/assets/get-deposit-address.png)


* Mint Brc20 Asset

```
./XrpinalsMintTool mint -h

Usage of mint:
  -addr string
        your address
  -asset string
        asset name you want to mint


./XrpinalsMintTool mint -addr mnUbdaJcTiBUARHGMZpQ5dVkrcj1XUMame -asset XX

mining success, txHash: 03d8216ce49753cbe6ae7a1a65e08b4fe841d5b6

```
![](https://github.com/xrpinals/XrpinalsMintTool/blob/main/assets/mint.png)



* Transfer Brc20 Asset

```
./XrpinalsMintTool transfer -h

Usage of mint:
  Usage of transfer:
  -from string
        sender address (must be imported before)
  -to string
        receiver address
  -amount string
        asset amount you want to transfer
  -asset string
        asset name you want to transfer


./XrpinalsMintTool transfer -from mfhGJnP5T7A5kYDJNxnHozxrVzC7WKHzKs -to mnUbdaJcTiBUARHGMZpQ5dVkrcj1XUMame -asset BTC -amount 0.02

transfer success, txHash: 6eb8918df731f29952cc00a4ae77c0a07e907742

```
![](https://github.com/xrpinals/XrpinalsMintTool/blob/main/assets/transfer.png)



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
last mint time: "2023-12-18T15:30:00"

```
![](https://github.com/xrpinals/XrpinalsMintTool/blob/main/assets/get-mint-info.png)


* Withdraw Btc Asset

```
./XrpinalsMintTool withdraw -h

Usage of mint:
  Usage of transfer:
  -fromAddr string
        sender address (must be imported before)
  -toAddr string
        the destination address to which you withdrew your funds
  -amount string
        your btc withdrawal amount. empty means withdraw all
  -memo string
        Remarks on Withdrawal Transactions


./XrpinalsMintTool withdraw -fromAddr mfhGJnP5T7A5kYDJNxnHozxrVzC7WKHzKs -toAddr mnUbdaJcTiBUARHGMZpQ5dVkrcj1XUMame -memo test -amount 0.02

transfer success, txHash: 6eb8918df731f29952cc00a4ae77c0a07e907742

```


* Query Withdraw Btc Info

```
./XrpinalsMintTool query_withdraw 
This command will list all the pending and in-process withdrawals on the network.


./XrpinalsMintTool query_withdraw

[Info]:  Waiting Withdrawal Info

----------------------------------------

[Info]:  Processing Withdrawal Info
[Info]:  Withdrawal account: mhk8YnXVEe6KdTNyoH4GRsoVYdHUaJLLaL Withdrawal amount: 0.002 Withdrawal to account: mjd8EBD4Q4mQePbTpysVjfYcAjyLzCxbTS

[Info]:  Withdrawal account: mhk8YnXVEe6KdTNyoH4GRsoVYdHUaJLLaL Withdrawal amount: 0.002 Withdrawal to account: mjd8EBD4Q4mQePbTpysVjfYcAjyLzCxbTS

[Info]:  Withdrawal account: mjd8EBD4Q4mQePbTpysVjfYcAjyLzCxbTS Withdrawal amount: 0.001 Withdrawal to account: mhk8YnXVEe6KdTNyoH4GRsoVYdHUaJLLaL

[Info]:  Withdrawal account: mjd8EBD4Q4mQePbTpysVjfYcAjyLzCxbTS Withdrawal amount: 0.0001 Withdrawal to account: mjd8EBD4Q4mQePbTpysVjfYcAjyLzCxbTS

[Info]:  Withdrawal account: mjd8EBD4Q4mQePbTpysVjfYcAjyLzCxbTS Withdrawal amount: 0.001 Withdrawal to account: mjd8EBD4Q4mQePbTpysVjfYcAjyLzCxbTS

[Info]:  Withdrawal account: mjd8EBD4Q4mQePbTpysVjfYcAjyLzCxbTS Withdrawal amount: 0.001 Withdrawal to account: mjd8EBD4Q4mQePbTpysVjfYcAjyLzCxbTS


```