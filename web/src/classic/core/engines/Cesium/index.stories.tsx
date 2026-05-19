import { Meta, StoryObj } from "@storybook/react-vite";

import Component from "../../Map";

import { engine } from ".";

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
          type: "geojson",
          value: { type: "Feature", geometry: { type: "Point", coordinates: [0, 0] } },
        },
        marker: {
          imageColor: "#fff",
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
