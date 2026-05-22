import { Meta, StoryObj } from "@storybook/react-vite";
import { Math as CesiumMath } from "cesium";

import { V, location } from "../storybook";

import PhotoOverlay from ".";

const meta: Meta<typeof PhotoOverlay> = {
  title: "classic/molecules/Visualizer/Engine/Cesium/PhotoOverlay",
  component: PhotoOverlay,
};
export default meta;
type Story = StoryObj<typeof PhotoOverlay>;

export const Default: Story = {
  render: args => (
    <V location={location}>
      <PhotoOverlay {...args} />
    </V>
  ),
  args: {
    layer: {
      id: "",
      isVisible: true,
      property: {
        default: {
          location,
          photoOverlayImage: `/sample.svg`,
          camera: {
            ...location,
            fov: CesiumMath.toRadians(30),
            heading: 0,
            pitch: 0,
            roll: 0,
          },
        },
      },
    },
    isBuilt: false,
    isEditable: false,
    isSelected: false,
  },
};

export const Selected: Story = {
  render: args => (
    <V location={location}>
      <PhotoOverlay {...args} />
    </V>
  ),
  args: {
    layer: {
      id: "",
      isVisible: true,
      property: {
        default: {
          location,
          photoOverlayImage: `/sample.svg`,
          camera: {
            ...location,
            fov: CesiumMath.toRadians(30),
            heading: 0,
            pitch: 0,
            roll: 0,
          },
        },
      },
    },
    isBuilt: false,
    isEditable: false,
    isSelected: true,
  },
};
