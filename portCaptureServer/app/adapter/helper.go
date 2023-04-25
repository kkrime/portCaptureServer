package adapter

import (
	"portCaptureServer/app/api/pb"
	"portCaptureServer/app/entity"
)

// convertPBPortToEntityPort converts the pb (protobuf) format for ports
// to the entity format
func convertPBPortToEntityPort(port *pb.Port) *entity.Port {
	// alias
	alias := make([]entity.Alias, 0, len(port.Alias))
	for _, a := range port.Alias {
		alias = append(alias, entity.Alias{
			Name: a,
		})
	}

	// regions
	regions := make([]entity.Region, 0, len(port.Regions))
	for _, r := range port.Regions {
		regions = append(regions, entity.Region{
			Name: r,
		})
	}
	// unlocs
	unlocs := make([]entity.Unloc, 0, len(port.Unlocs))
	for _, u := range port.Unlocs {
		unlocs = append(unlocs, entity.Unloc{
			Name: u,
		})
	}

	return &entity.Port{
		Name:         port.Name,
		PrimaryUnloc: port.PrimaryUnloc,
		Code:         port.Code,
		City:         port.City,
		Country:      port.Country,
		Alias:        &alias,
		Regions:      &regions,
		Coordinantes: func() [2]float32 {
			if len(port.Coordinates) == 2 {
				return [2]float32{port.Coordinates[0], port.Coordinates[1]}
			}
			// return default -1, -1
			return [2]float32{-1, -1}
		}(),
		Province: port.Province,
		Timezone: port.Timezone,
		Unlocs:   &unlocs,
	}
}
