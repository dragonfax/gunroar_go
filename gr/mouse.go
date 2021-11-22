package main

const MOUSE_SCREEN_MAPPING_RATIO_X = 26.0
const MOUSE_SCREEN_MAPPING_RATIO_Y = 19.5

/**
 * Mouse input.
 */
type RecordableMouse struct {
  mouse.RecordableMouse
}

  public this(SizableScreen screen) {
    super();
    this.screen = screen;
  }

  protected override void adjustPos(MouseState ms) {
    ms.x =  (ms.x - screen.width  / 2) * MOUSE_SCREEN_MAPPING_RATIO_X / screen.width;
    ms.y = -(ms.y - screen.height / 2) * MOUSE_SCREEN_MAPPING_RATIO_Y / screen.height;
  }
}
