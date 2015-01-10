/*
 * $Id: actor.d,v 1.1.1.1 2005/06/18 00:46:00 kenta Exp $
 *
 * Copyright 2004 Kenta Cho. Some rights reserved.
 */
package gr

/**
 * Actor in the game that has the interface to move and draw.
 */
type Actor struct {
  Exists bool
}

var actors = make(map[Actor]bool)

func NewActor() *Actor {
  a = new(Actor)
  a.Init()
}

func (a *Actor) Init() {
  a.Exists = true
  actors[a] = false
}

func (a *Actor) Done() {
  a.Exists = false
  delete(actors,a)
}

