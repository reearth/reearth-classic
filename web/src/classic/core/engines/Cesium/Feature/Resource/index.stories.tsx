import { Meta, StoryObj } from "@storybook/react-vite";

import { engine } from "../..";
import Component from "../../../../Map";

const meta: Meta<typeof Component> = {
  component: Component,
  parameters: { actions: { argTypesRegex: "^on.*" } },
};
export default meta;
type Story = StoryObj<typeof Component>;

export const Default: Story = {
  args: {
    engine: "cesium",
    engines: {
      cesium: engine,
    },
    ready: true,
    layers: [
      {
        id: "l",
        type: "simple",
        data: {
          type: "czml",
          url: "/testdata/sample.czml",
        },
        resource: {},
        marker: {
          pointColor: {
            expression: {
              conditions: [["${id} === '1'", "color('red')"]],
            },
          },
        },
      },
    ],
    property: {
      tiles: [
        {
          id: "default",
          tile_type: "default",
        },
      ],
    },
  },
};
