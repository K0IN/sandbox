package container

type Config struct {
	AllowNetwork bool
	AllowProc    bool
	AllowEnv     bool
	Rootfs       string
	Entrypoint   string
	Env          []string
	Volume       []string
	Workdir      string
	User         string
	Group        string
	Hostname     string
	Mount        []string
}

type Container struct {
	AllowNetwork bool
	AllowProc    bool
	AllowEnv     bool
	Rootfs       string
	Entrypoint   string
	Env          []string
	Volume       []string
	Workdir      string
	User         string
	Group        string
	Hostname     string
	Mount        []string
}
