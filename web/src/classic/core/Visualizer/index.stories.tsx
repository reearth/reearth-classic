import { Meta, StoryObj } from "@storybook/react-vite";

import Component from ".";

const meta: Meta<typeof Component> = {
  component: Component,
  parameters: { actions: { argTypesRegex: "^on.*" } },
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
