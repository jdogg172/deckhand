package modes

type Mode string

const (
	OpsMode      Mode = "ops"
	PipelineMode Mode = "pipeline"
)

func Normalize(mode string) Mode {
	if mode == "pipeline" {
		return PipelineMode
	}
	return OpsMode
}
