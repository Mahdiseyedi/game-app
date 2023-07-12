package accesscontrol

type AccessControl struct {
	ID           uint
	ActorID      uint
	ActorType    ActorType
	PermissionID uint
}

type ActorType uint

const (
	RoleActorType = "role"
	UserActorType = "user"
)
