package filesystem

import "go.uber.org/zap"

type PlatformNamer struct {
	logger                    *zap.Logger
	atcoCodeToPlatformNameMap map[string]*string
}

func NewPlatformNamer(logger *zap.Logger) *PlatformNamer {
	atcoCodeToPlatformNameMap := make(map[string]*string)
	atcoCodeToPlatformNameMap["9400ZZMASTP1"] = strToStrPtr("D")
	atcoCodeToPlatformNameMap["9400ZZMASTP2"] = strToStrPtr("C")
	atcoCodeToPlatformNameMap["9400ZZMASTP3"] = strToStrPtr("B")
	atcoCodeToPlatformNameMap["9400ZZMASTP4"] = strToStrPtr("A")
	atcoCodeToPlatformNameMap["9400ZZMAVIC1"] = strToStrPtr("D")
	atcoCodeToPlatformNameMap["9400ZZMAVIC2"] = strToStrPtr("C")
	atcoCodeToPlatformNameMap["9400ZZMAVIC3"] = strToStrPtr("B")
	atcoCodeToPlatformNameMap["9400ZZMAVIC4"] = strToStrPtr("A")

	return &PlatformNamer{
		logger:                    logger,
		atcoCodeToPlatformNameMap: atcoCodeToPlatformNameMap,
	}
}

func (p *PlatformNamer) GetPlatformNameForAtcoCode(atcoCode string) (*string, error) {
	return p.atcoCodeToPlatformNameMap[atcoCode], nil
}

func strToStrPtr(s string) *string {
	return &s
}
