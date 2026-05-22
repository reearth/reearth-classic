import { fn } from "storybook/test";

import type { Context } from ".";

export const contextEvents: Context = {
  onCameraOrbit: fn(),
  onCameraRotateRight: fn(),
  onFlyTo: fn(),
  onLayerSelect: fn(),
  onLookAt: fn(),
  onPause: fn(),
  onPlay: fn(),
  onSpeedChange: fn(),
  onTick: fn(),
  onTimeChange: fn(),
  onZoomIn: fn(),
  onZoomOut: fn(),
};
