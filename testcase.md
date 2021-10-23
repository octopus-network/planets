# Hermes relayer test cases.

## Introduction

In this repository we will be working with the [Hermes IBC Relayer](https://hermes.informal.systems) in order to transfer fungible tokens ([ics-020](https://github.com/cosmos/ics/tree/master/spec/ics-020-fungible-token-transfer)) between two [Starport](https://github.com/tendermint/starport) custom [Cosmos SDK](https://github.com/cosmos/cosmos-sdk) chains.



## Setup Environment
1. Install docker and docker-compose 
2. Install hermes.   
Please follow the [instructions to install](https://hermes.informal.systems/installation.html) it locally on your mahcine.
## Run the chains
Follow this instructions to run the two chains

### -------------------------------
### Start earth
### -------------------------------
Open a terminal prompt:

```
cd /your/path/to/planets

docker-compose up earth
```

### Restore key (Alice)

Open another terminal prompt in the same location (earth folder)

This command restores a key for a user in `earth` that we will be using during the workshop.

```
# into the container
docker exec -it earth bash

# restore the key for alice
earth keys add alice --recover --home .earth

# When prompted for the mnemonic please enter the following words:
picture switch picture soap flip dawn nerve easy rebuild company hawk stand menu rhythm unfold engine rug rally weapon raccoon glide mosquito lion dog

# query Alice's balance
earth --node tcp://localhost:26657 query bank balances $(earth --home .earth keys --keyring-backend="test" show alice -a)

# exit container
exit
```

### -------------------------------
### Start mars
### -------------------------------

Open another terminal prompt:

```
cd /your/path/to/planets

docker-compose up mars
```
### Restore key (Bob)

Open another terminal prompt in the same location (mars folder)

This command restores a key for a user in `mars` that we will be using during the workshop.

```
# into the container
docker exec -it mars bash

# restore the key for bob
mars keys add bob --recover --home .mars

# When prompted for the mnemonic please enter the following words:
gaze clay walk tail shove sphere follow twenty agent basket viable gun popular decide vanish coyote guilt carry toward exhaust hour six scout chest

# query Bob's balance
mars --node tcp://localhost:26658 query bank balances $(mars --home .mars keys --keyring-backend="test" show bob -a)

# exit container
exit
```


## IBC clients
### Configure Keys for relayer

#### Add keys for each chain

These keys will be used by Hermes to sign transactions sent to each chain. The keys need to have a balance on each respective chain in order to pay for the transactions.

```
hermes -c ./hermes.toml keys add earth -f ../chains/earth/alice_key.json 
```

```
hermes -c ./hermes.toml keys add mars -f ../chains/mars/bob_key.json 
```

### List keys

To ensure the keys were properly added let's list them.

```
hermes -c ./hermes.toml keys list earth
```

```
hermes -c ./hermes.toml keys list mars
```

### Create client

Create a mars client on earth:

```
hermes -c ./hermes.toml tx raw create-client earth mars 
```

Query the client state

```
hermes -c ./hermes.toml query client state earth 07-tendermint-0 
```

Create a earth client on mars

```
hermes -c ./hermes.toml tx raw create-client mars earth 
```

Query the client state

```
hermes -c ./hermes.toml query client state mars 07-tendermint-0
```

### Update client

Update the mars client on earth

```
hermes -c ./hermes.toml tx raw update-client earth 07-tendermint-0 
```

Query the client state

```
hermes -c ./hermes.toml query client state earth 07-tendermint-0 
```

## IBC Connection
In this section we will run all the steps required to establish a connection handshake 

## Connection Handshake (earth -> mars)

The steps below need to succeed in order to have a connection opened between earth and mars

### ConnOpenInit

```
hermes -c ./hermes.toml tx raw conn-init earth mars 07-tendermint-0 07-tendermint-0 
```

### ConnOpenTry

```
hermes -c ./hermes.toml tx raw conn-try mars earth 07-tendermint-0 07-tendermint-0 -s connection-0 
```

### ConnOpenAck

```
hermes -c ./hermes.toml tx raw conn-ack earth mars 07-tendermint-0 07-tendermint-0 -d connection-0 -s connection-0 
```

### ConnOpenConfirm

```
hermes -c ./hermes.toml tx raw conn-confirm mars earth 07-tendermint-0 07-tendermint-0 -d connection-0 -s connection-0
```

### Query connection

The commands below allow you to query the connection state on each chain.

```
hermes -c ./hermes.toml query connection end earth connection-0 
```

```
hermes -c ./hermes.toml query connection end mars connection-0 
```

## IBC Channel 
In this step, we will establish a channel between the chains

### Channel Handshake

The steps below need to succeed in order to have a channel opened between earth and mars

### ChanOpenInit

```
hermes -c ./hermes.toml tx raw chan-open-init earth mars connection-0 transfer transfer -o UNORDERED 
```

### ChanOpenTry

```
hermes -c ./hermes.toml tx raw chan-open-try mars earth connection-0 transfer transfer -s channel-0 
```

### ChanOpenAck

```
hermes -c ./hermes.toml tx raw chan-open-ack earth mars connection-0 transfer transfer -d channel-0 -s channel-0 
```

### ChanOpenConfirm

```
hermes -c ./hermes.toml tx raw chan-open-confirm mars earth connection-0 transfer transfer -d channel-0 -s channel-0 
```

### Query channel

Use these commands to query the channel end on each chain

```
hermes -c ./hermes.toml query channel end earth transfer channel-0 
```

```
hermes -c ./hermes.toml query channel end mars transfer channel-0 
```

## IBC Relay Packets (Transfer Tokens)  
First let's query the balance for Alice and Bob to view how much tokens each have.

In the terminal prompt for each chain run the commands below

### Query balance - Alice

```
docker exec -it earth bash
earth --node tcp://localhost:26657 query bank balances $(earth --home .earth keys --keyring-backend="test" show alice -a)
exit
```

> Please note the amound to `coina` tokens that Alice has. We will check it again later.

### Query balance - Bob

```
docker exec -it mars bash
mars --node tcp://localhost:26658 query bank balances $(mars --home .mars keys --keyring-backend="test" show bob -a)
exit
```

> Please note that Bob only has `coinb` tokens (also some `stake` tokens)

### Fungible token transfer

Now we will transfer some tokens

### Send packet

``` 
hermes -c ./hermes.toml tx raw ft-transfer mars earth transfer channel-0 999 -o 1000 -n 1 -d samoleans

```

### View response

The response has a data field that the information is encoded. To view it, you can use the command below. Just replace `[DATA]` with the value from the `data` field in the response

```
echo "[DATA]" | xxd -r -p | jq
```

### Query packet commitments

```
hermes -c ./hermes.toml query packet commitments earth transfer channel-0 
```

### Query unreceived packets on mars

```
hermes -c ./hermes.toml query packet unreceived-packets mars transfer channel-0 
```

### Send recv_packet to mars

```
hermes -c ./hermes.toml tx raw packet-recv mars earth transfer channel-0 
```

### Query unreceived ack on earth

```
hermes -c ./hermes.toml query packet unreceived-acks earth transfer channel-0 
```

### Send ack to earth

```
hermes -c ./hermes.toml tx raw packet-ack earth mars transfer channel-0 
```

### Query balance - Alice

To ensure the tokens were transferred, query Alice balance.

```
docker exec -it earth bash
earth --node tcp://localhost:26657 query bank balances $(earth --home .earth keys --keyring-backend="test" show alice -a)
exit

```

### Query balance - Bob

Now we will query Bob's balance to ensure the tokens were transferred.

```
docker exec -it mars bash
mars --node tcp://localhost:26658 query bank balances $(mars --home .mars keys --keyring-backend="test" show bob -a)
exit
```

### View denom trace

In the command above to show the balance, the denom is hashed. In order to view the denom trace you can call the API endpoint below. Just replace the `[denom]` value with the value from the `denom` field (the one that starts with `ibc/`)
```
curl http://localhost:1318/ibc/applications/transfer/v1beta1/denom_traces/[denom]
```

## Send tokens back to earth

### Tranfer tokens back

Just replace the `[HASH]` with the hash in the `denom` field

```
hermes -c ./hermes.toml tx raw ft-transfer earth mars transfer channel-0 999 -o 1000 -n 1 -d ibc/[HASH] 
```

### Send recv_packet to mars

```
hermes -c ./hermes.toml tx raw packet-recv earth mars transfer channel-0 
```

### Send ack to mars

```
hermes -c ./hermes.toml tx raw packet-ack mars earth transfer channel-0 
```

### Query balance - Alice

To ensure the tokens were transferred back, query Alice balance again:

```
docker exec -it earth bash
earth --node tcp://localhost:26657 query bank balances $(earth --home .earth keys --keyring-backend="test" show alice -a)
exit
```

### Query balance - Bob

Query Bob's balance to ensure his token balance for the `ibc/[hash]` denom shows `0` balance

```
docker exec -it mars bash
mars --node tcp://localhost:26658 query bank balances $(mars --home .mars keys --keyring-backend="test" show bob -a)
exit
```

## Stop and save docker container.
To avoid importing accounts repeatedly, you can just only stop docker
```
cd /path/to/your/planets
docker-compose stop
```

## Congratulations

If you successfully executed all the previous steps you should have a better understanding now how IBC prottocol works to connect two chains and transfer some tokens.

## References
* [Interchain Standards](https://github.com/cosmos/ibc)
* [IBC Protocol Website](https://ibcprotocol.org)
* [IBC Modules and Relayer in Rust](https://github.com/informalsystems/ibc-rs)
* [Hermes Documentation](https://hermes.informal.systems)
* [Starport](https://docs.starport.network/)