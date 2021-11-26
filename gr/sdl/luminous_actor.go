package sdl

import "github.com/dragonfax/gunroar/gr/actor"

/**
 * Actor with the luminous effect.
 */

type LuminousActor interface {
	actor.Actor
	DrawLuminous()
}

/**
 * Actor pool for the LuminousActor.
 */

type LuminousActorPool struct {
	actor.ActorPool
}

func NewLuminousActorPool(f actor.CreateActorFunc, n int, args []interface{}) *LuminousActorPool {
	this := &LuminousActorPool{
		ActorPool: actor.NewActorPoolInternal(f, n, args),
	}
	return this
}

func (this *LuminousActorPool) DrawLuminous() {
	for _, a := range this.Actor {
		if a.Exists() {
			a.(LuminousActor).DrawLuminous()
		}
	}
}
