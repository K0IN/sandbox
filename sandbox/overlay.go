package sandbox

const (
	sandboxConfigFileName = "config.json"
	sandboxUpperDir       = "upper"
	sandboxWorkDir        = "workdir"
	sandboxMountPointDir  = "rootfs"
)

type SandboxInfo struct {
	StagedFiles  []string `json:"stagedFiles"`
	ChangedFiles []string `json:"changedFiles"`
}

type OverlayFs struct {
	SandboxDir string
}

func OpenOverlay(sandboxDir string) (*OverlayFs, error) {
	return nil, nil
}

func CreateOverlay(sandboxDir string) (*OverlayFs, error) {
	// create all paths and config file

	return nil, nil
}

func (s *OverlayFs) Destroy() error {
	return nil
}

func (s *OverlayFs) CommitToDisk() error {
	return nil
}

func (s *OverlayFs) StageFile(filePath string) error {
	return nil
}

func (s *OverlayFs) UnstageFile(filePath string) error {
	return nil
}

func (s *OverlayFs) GetStagedFiles() ([]string, error) {
	return nil, nil
}

func (s *OverlayFs) GetChangedFiles() ([]string, error) {
	return nil, nil
}
