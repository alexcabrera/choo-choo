package state

type Phase string

const (
	PhaseInit       Phase = "init"
	PhaseDesign     Phase = "design"
	PhasePlan       Phase = "plan"
	PhaseValidate   Phase = "validate"
	PhaseExecution  Phase = "execution"
	PhaseVerify     Phase = "verification"
	PhaseCloseGaps  Phase = "gap-closure"
	PhaseDone       Phase = "done"
)

func (p Phase) IsValid() bool {
	switch p {
	case PhaseInit, PhaseDesign, PhasePlan, PhaseValidate,
		PhaseExecution, PhaseVerify, PhaseCloseGaps, PhaseDone:
		return true
	default:
		return false
	}
}

func (p Phase) String() string {
	return string(p)
}
