package patternphone

import "service/pkg/repo"

//IPatternPhone ....
type IPatternPhone interface {
	Get() []string
}

//PatternPhone ...
var PatternPhone IPatternPhone = &repo.PatternPhone{}
