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
	move()
	draw()
	remove()
}

/* NOTE: you must always insert a POINTER to the struct implementing Actor
 *	golang and the map below will accept a struct value just the same.
 *	but the pointer is needed to make every new struct act as a unique key.
 *
 * the value of the map replaces `exists` and is used to identify if an actor
 * is still in play. or has been decommissions, but not yet garbage collected.
 *
 * TODO use weak references throughout the system to improve this.
 */
var actors = make(map[Actor]bool)
