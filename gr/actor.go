/*
 * $Id: actor.d,v 1.1.1.1 2005/06/18 00:46:00 kenta Exp $
 *
 * Copyright 2004 Kenta Cho. Some rights reserved.
 */
package gr

/**
 * Actor in the game that has the interface to move and draw.
 */
type Actor interface {
	Exists() bool
	SetExists(b bool)
}

var actors = make(map[Actor]bool)

type ActorImpl struct {
	exists bool
}

func (a *ActorImpl) Exists() bool {
	return a.exists
}

func (a *ActorImpl) SetExists(b bool) {
	a.exists = b
}

func (a *ActorImpl) Init() {
	a.Exists = true
	actors[a] = false
}

func (a *ActorImpl) Done() {
	a.Exists = false
	delete(actors, a)
}
