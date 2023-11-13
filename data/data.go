package data

//Global data stored for the lifetime of the program
type Data struct {
	AllPeople       []string
	GroupNames []string
	Groups map[string][]string
	PastAssignments map[string][]string
	NumAssignments  int
	PossibleMatches []Possibility
	TriesAllowed    int
	Depth           int
	RawHistoryData  []map[string]string
}

//A single person's possible people they can give a gift to.  
type Possibility struct {
	Name          string
	Possibilities []string
}

//Strange variable that is designed to overcome an idiosyncrasy with scanf() in which it doesn't accept "" very well, so instead we have this recognizable, but outlandish string.
var BSpaceHold string = "-$#asdfas@#$aswefadsfagasdgasf@#$#@!$!radfasdfasdfas@#$!@#%!#$^tudfaf"