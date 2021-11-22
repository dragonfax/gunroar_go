public class RecordableMouse: Mouse {
  mixin RecordableInput!(MouseState);
 private:

  public MouseState getState(bool doRecord = true) {
    MouseState s = super.getState();
    if (doRecord)
      record(s);
    return s;
  }
}
