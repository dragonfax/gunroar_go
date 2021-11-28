package actor

type Actor interface {
	Init([]interface{})
	Move()
	Draw()
	Exists() bool
	SetExists(bool) bool
}

/**
 * Actor in the game that has the interface to move and draw.
 */
type ExistsImpl struct {
	_exists bool
}

func (this ExistsImpl) Exists() bool {
	return this._exists
}

func (this *ExistsImpl) SetExists(value bool) bool {
	this._exists = value
	return value
}

/**
 * Object pooling for actors.
 */

type CreateActorFunc func() Actor

type ActorPool struct {
	Actor       []Actor
	actorIdx    int
	createActor CreateActorFunc
}

func NewActorPool(f CreateActorFunc, n int, args []interface{}) ActorPool {
	this := NewActorPoolInternal(f, n, args)
	return this
}

func NewActorPoolInternal(f CreateActorFunc, n int, args []interface{}) ActorPool {
	this := ActorPool{
		Actor:       nil,
		createActor: f,
	}
	this.createActors(n, args)
	return this
}

func (this *ActorPool) createActors(n int, args []interface{} /* = null */) {
	this.Actor = make([]Actor, n, n)
	for i := range this.Actor {
		a := this.createActor()
		a.SetExists(false)
		a.Init(args)
		this.Actor[i] = a
	}
	this.actorIdx = 0
}

func (this *ActorPool) GetInstance() Actor {
	for i := 0; i < len(this.Actor); i++ {
		this.actorIdx--
		if this.actorIdx < 0 {
			this.actorIdx = len(this.Actor) - 1
		}
		if !this.Actor[this.actorIdx].Exists() {
			return this.Actor[this.actorIdx]
		}
	}
	return nil
}

func (this *ActorPool) GetInstanceForced() Actor {
	this.actorIdx--
	if this.actorIdx < 0 {
		this.actorIdx = len(this.Actor) - 1
	}
	return this.Actor[this.actorIdx]
}

func (this *ActorPool) GetMultipleInstances(n int) []Actor {
	rsl := make([]Actor, n, n)
	for i := 0; i < n; i++ {
		inst := this.GetInstance()
		if inst == nil {
			for _, r := range rsl {
				r.SetExists(false)
			}
			return nil
		}
		inst.SetExists(true)
		rsl = append(rsl, inst)
	}
	for _, r := range rsl {
		r.SetExists(false)
	}
	return rsl
}

func (this *ActorPool) Move() {
	for _, ac := range this.Actor {
		if ac.Exists() {
			ac.Move()
		}
	}
}

func (this *ActorPool) Draw() {
	for _, ac := range this.Actor {
		if ac.Exists() {
			ac.Draw()
		}
	}
}

func (this *ActorPool) Clear() {
	for _, ac := range this.Actor {
		ac.SetExists(false)
	}
	this.actorIdx = 0
}
