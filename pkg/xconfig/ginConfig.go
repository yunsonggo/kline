package xconfig

type GinConf struct {
	Mode         string   `json:"mode"`
	UseCors      bool     `json:"useCors"`
	UsePProf     bool     `json:"usePProf"`
	AllowHeaders []string `json:"allowHeaders"`
	StaticPath   string   `json:"staticPath"`
	StaticPrefix string   `json:"staticPrefix"`
}
