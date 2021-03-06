/*
 * $Id: actor.d,v 1.1.1.1 2005/06/18 00:46:00 kenta Exp $
 *
 * Copyright 2004 Kenta Cho. Some rights reserved.
 */
package main

/**
 * Actor in the game that has the interface to move and draw.
 */
type Actor interface {
	move()
	draw()
	close()
}

/* Each Actor must insert its pointer into actors when created
 *
 * NOTE: You must always insert a POINTER to the Actor.
 *	Golang will accept a struct value just the same as a pointer.
 *	But the pointer is needed to make every struct act as a unique key in the map.
 */
var actors = make(map[Actor]bool)

func clearActors() {
	for a, _ := range actors {
		a.close()
	}
	actors = make(map[Actor]bool)
}
