/*
 * GENERATED. Do not modify. Your changes might be overwritten!
 */

package resources

type MissionByExplorerMissions struct {
	Key
	Attributes MissionByExplorerMissionsAttributes `json:"attributes"`
}
type MissionByExplorerMissionsResponse struct {
	Data     MissionByExplorerMissions `json:"data"`
	Included Included                  `json:"included"`
}

type MissionByExplorerMissionsListResponse struct {
	Data     []MissionByExplorerMissions `json:"data"`
	Included Included                    `json:"included"`
	Links    *Links                      `json:"links"`
}

// MustMissionByExplorerMissions - returns MissionByExplorerMissions from include collection.
// if entry with specified key does not exist - returns nil
// if entry with specified key exists but type or ID mismatches - panics
func (c *Included) MustMissionByExplorerMissions(key Key) *MissionByExplorerMissions {
	var missionByExplorerMissions MissionByExplorerMissions
	if c.tryFindEntry(key, &missionByExplorerMissions) {
		return &missionByExplorerMissions
	}
	return nil
}
