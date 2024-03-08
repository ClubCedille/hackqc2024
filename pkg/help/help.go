package help

type Help struct {
	MapObjectId  int    `clover:"map_object_id"`
	ContactInfos string `clover:"contact_infos"`
	NeedHelp     bool   `clover:"need_help"`
	HowToHelp    string `clover:"how_to_help"`
	HowToUseHelp string `clover:"how_to_use_help"`
}
