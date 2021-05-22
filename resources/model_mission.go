/*
 * GENERATED. Do not modify. Your changes might be overwritten!
 */

package resources

type Mission struct {
	Key
	Attributes MissionAttributes `json:"attributes"`
}
type MissionResponse struct {
	Data     Mission  `json:"data"`
	Included Included `json:"included"`
}

type MissionListResponse struct {
	Data     []Mission `json:"data"`
	Included Included  `json:"included"`
	Links    *Links    `json:"links"`
}

// MustMission - returns Mission from include collection.
// if entry with specified key does not exist - returns nil
// if entry with specified key exists but type or ID mismatches - panics
func (c *Included) MustMission(key Key) *Mission {
	var mission Mission
	if c.tryFindEntry(key, &mission) {
		return &mission
	}
	return nil
}
