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
          type: "geojson",
          value: {
            type: "Feature",
            geometry: {
              type: "Point",
              coordinates: [0, 0, 1000],
            },
          },
        },
        model: {
          url: "/BoxAnimated.glb",
          scale: 1000000,
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
