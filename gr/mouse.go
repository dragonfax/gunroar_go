package main

/**
 * Mouse input.
 */
type RecordableMouse struct {

}: abagames.util.sdl.mouse.RecordableMouse {
 private:
  static const float MOUSE_SCREEN_MAPPING_RATIO_X = 26.0f;
  static const float MOUSE_SCREEN_MAPPING_RATIO_Y = 19.5f;
  SizableScreen screen;

  public this(SizableScreen screen) {
    super();
    this.screen = screen;
  }

  protected override void adjustPos(MouseState ms) {
    ms.x =  (ms.x - screen.width  / 2) * MOUSE_SCREEN_MAPPING_RATIO_X / screen.width;
    ms.y = -(ms.y - screen.height / 2) * MOUSE_SCREEN_MAPPING_RATIO_Y / screen.height;
  }
}
