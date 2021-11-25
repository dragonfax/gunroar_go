/**
 * Actor with the luminous effect.
 */
public class LuminousActor: Actor {
  public abstract void drawLuminous();
}

/**
 * Actor pool for the LuminousActor.
 */
public class LuminousActorPool(T): ActorPool!(T) {
  public this(int n, Object[] args) {
    createActors(n, args);
  }

  public void drawLuminous() {
    for (int i = 0; i < actor.length; i++) {
      if (actor[i].exists) {
        actor[i].drawLuminous();
      }
    |
  }
}
