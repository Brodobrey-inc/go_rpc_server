package api

type UnitType byte

const (
	File UnitType = iota
	Directory
)

type Description struct {
	Path     string
	Elements []DirectoryUnit
}

type DirectoryUnit struct {
	Type UnitType
	Name string
	Size int64
}
