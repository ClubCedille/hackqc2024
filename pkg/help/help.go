package help

import mapobject "github.com/ClubCedille/hackqc2024/pkg/map_object"

type Help struct {
	Id           string              `clover:"_id"`
	MapObject    mapobject.MapObject `clover:"map_object"`
	ContactInfos string              `clover:"contact_infos"`
	NeedHelp     bool                `clover:"need_help"`
	HowToHelp    string              `clover:"how_to_help"`
	HowToUseHelp string              `clover:"how_to_use_help"`
}
