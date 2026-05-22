import { Meta, StoryObj } from "@storybook/react-vite";
import { Cartesian3, Color } from "cesium";
import { Entity } from "resium";

import Component from ".";

const meta: Meta<typeof Component> = {
  title: "classic/molecules/Visualizer/Engine/Cesium",
  component: Component,
  argTypes: {
    onCameraChange: { action: "onCameraChange" },
    onLayerSelect: { action: "onLayerSelect" },
  },
};
export default meta;
type Story = StoryObj<typeof Component>;

const defaultArgs = {
  isBuilt: false,
  isEditable: false,
  small: false,
  ready: true,
  selectedLayerId: undefined,
  property: {
    default: {
      terrain: true,
      terrainExaggeration: 10,
      bgcolor: "#ff0",
      skybox: false,
    },
    tiles: [{ id: "default", tile_type: "default" }],
    atmosphere: {
      enable_lighting: true,
      enable_sun: true,
      sky_atmosphere: true,
      ground_atmosphere: true,
    },
  },
};

export const Default: Story = {
  args: defaultArgs,
};

export const Selected: Story = {
  args: {
    ...defaultArgs,
    children: (
      <Entity
        id="a"
        point={{ color: Color.WHITE, pixelSize: 10 }}
        position={Cartesian3.fromDegrees(0, 0, 0)}
        selected
      />
    ),
    selectedLayerId: "a",
  },
};

export const DefaultCamera: Story = {
  args: {
    ...defaultArgs,
    camera: {
      fov: 1.0471975511965976,
      heading: 0.23410230091957818,
      height: 21075.98272847632,
      lat: 36.382785984750846,
      lng: 139.59084714252558,
      pitch: 0,
      roll: 0,
    },
  },
};

export const InitialCamera: Story = {
  args: {
    ...defaultArgs,
    property: {
      ...defaultArgs.property,
      default: {
        ...defaultArgs.property?.default,
        camera: {
          fov: 1.0471975511965976,
          heading: 0.23410230091957818,
          height: 21075.98272847632,
          lat: 36.382785984750846,
          lng: 139.59084714252558,
          pitch: 0,
          roll: 0,
        },
      },
    },
  },
};
