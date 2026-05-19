import { Meta, StoryObj } from "@storybook/react-vite";

import { engine } from "../..";
import Component from "../../../../Map";

const meta: Meta<typeof Component> = {
  component: Component,
  parameters: { actions: { argTypesRegex: "^on.*" } },
};
export default meta;
type Story = StoryObj<typeof Component>;

export const WMS: Story = {
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
          type: "wms",
          url: "https://gibs.earthdata.nasa.gov/wms/epsg3857/best/wms.cgi",
          layers: "IMERG_Precipitation_Rate",
        },
        raster: {
          maximumLevel: 100,
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

export const MVT: Story = {
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
          type: "mvt",
          url: "https://example.com/exmaple.mvt", // You need to set MVT URL.
          layers: "road",
        },
        polygon: {
          fillColor: "white",
          strokeColor: "white",
          strokeWidth: 1,
          lineJoin: "round",
        },
        raster: {},
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
