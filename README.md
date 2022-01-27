# missions-api
An event listener for contracts deployed on binance smart chain to store the events information in the database

[![Visual Source](https://img.shields.io/badge/visual-source-orange)](https://www.visualsource.net/repo/github.com/Alien-Worlds/missions-api)


## Requirements

* [Docker 20.10.6+](https://www.docker.com/get-started)
* [Compose 3.3+](https://docs.docker.com/compose/install/)
* [Go 1.15+](https://golang.org/) 
* [Postgresql 12.6](https://www.postgresql.org/)

## Running the service
#### For development purposes
1. Modify *config.yaml* file with your needs:

   * Provide a database url where the information will be stored in the form as provided [here](https://www.postgresql.org/docs/current/libpq-connect.html#LIBPQ-CONNSTRING)

		```sh
		db:
		   url: "postgresql://[userspec@][hostspec][/dbname][?paramspec]"
		```
	 
   * Provide contract address instead of 0x0..0

		```sh
		contract:
		   address: "0x0000000000000000000000000000000000000000"
		```
   * Provide block number to start with
		```sh
		blockchain_info:
                   from_block_num: 8325019
		```
   * If you are using local blockchain (e.g Ganache), provide other endpoint

		```sh
		rpc:
		   endpoint: "http://127.0.0.1:8545" # or localhost
		```
  
2. Open project using IDE (e.g GoLand from JetBrains), open IDE terminal, run
   
	```sh
	go mod init
	go mod vendor
	```

3. At *bsc-checker-events/assets/main.go* run two migration related scripts:
   ```sh
   //go:generate packr2 clean
   //go:generate packr2
	```
4. Entry point is *bsc-checker-events/main.go*
5. Modify run configuration as follows:
	* KV_VIPER_FILE=config.yaml *(environment variable)*
6. Run service twice with the following command arguments:
   
	```sh
	migrate up
	run service
	```

#### For deployment purposes
1. Navigate to the cloned repository
2. Do the step 1 from development build, except modify config at *configs/spaceship-staking.yaml*, changing contract address and database url (for the contract deployed on the Binance Smart Chain leave the endpoint as it is 
3. Build container image:
   
   ```sh
   docker build -t spaceship-staking .
	```
4. Run using docker-compose
   ```sh
   docker-compose down -v
   docker-compose up -d
	```

### API
To change port, configure 
```sh
listener:
  addr: :8888
```
where *8888* is a port to listen on.

#### Endpoints
```sh
/missions # get all missions
/missions/{mission-id} # get mission by it's id
/explorers/{explorer-address} # get missions joined by explorer address
```
