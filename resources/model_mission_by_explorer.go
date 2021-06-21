/*
 * GENERATED. Do not modify. Your changes might be overwritten!
 */

package resources

type MissionByExplorer struct {
	Key
	Attributes MissionByExplorerAttributes `json:"attributes"`
}
type MissionByExplorerResponse struct {
	Data     MissionByExplorer `json:"data"`
	Included Included          `json:"included"`
}

type MissionByExplorerListResponse struct {
	Data     []MissionByExplorer `json:"data"`
	Included Included            `json:"included"`
	Links    *Links              `json:"links"`
}

// MustMissionByExplorer - returns MissionByExplorer from include collection.
// if entry with specified key does not exist - returns nil
// if entry with specified key exists but type or ID mismatches - panics
func (c *Included) MustMissionByExplorer(key Key) *MissionByExplorer {
	var missionByExplorer MissionByExplorer
	if c.tryFindEntry(key, &missionByExplorer) {
		return &missionByExplorer
	}
	return nil
}
