package data


type Mission struct{
	missionId     uint64
	description   string
	name          string
	boardingTime  uint64
	launchTime    uint64
	endTime       uint64
	duration      uint64
	missionType   uint64
	reward        uint64
	spaceshipCost uint64
	missionPower  uint64
	totalShips    uint64
	nftContract   []byte
	nftTokenURI   string
}