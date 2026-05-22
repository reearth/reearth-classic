import { Meta, StoryObj } from "@storybook/react-vite";
import { fn } from "storybook/test";

import Component from ".";

const meta: Meta<typeof Component> = {
  component: Component,
  args: {
    onCameraChange: fn(),
    onLayerDrop: fn(),
    onLayerSelect: fn(),
    onWidgetLayoutUpdate: fn(),
    onWidgetAlignmentUpdate: fn(),
    onWidgetAreaSelect: fn(),
    onInfoboxMaskClick: fn(),
    onBlockSelect: fn(),
    onBlockChange: fn(),
    onBlockMove: fn(),
    onBlockDelete: fn(),
    onBlockInsert: fn(),
    onZoomToLayer: fn(),
  },
};
export default meta;
type Story = StoryObj<typeof Component>;

export const Cesium: Story = {
  args: {
    ready: true,
    engine: "cesium",
  },
};

export const Plugin: Story = {
  args: {
    ready: true,
    engine: "cesium",
    widgetAlignSystem: {
      outer: {
        left: {
          top: {
            align: "start",
            widgets: [
              {
                id: "plugin",
                pluginId: "plugin",
                extensionId: "test",
                ...({
                  __REEARTH_SOURCECODE: `reearth.layers.add(${JSON.stringify({
                    id: "l",
                    type: "simple",
                    data: {
                      type: "geojson",
                      value: { type: "Feature", geometry: { type: "Point", coordinates: [0, 0] } },
                    },
                    marker: {
                      imageColor: "#fff",
                    },
                  })});`,
                } as any),
              },
            ],
          },
        },
      },
    },
  },
};
