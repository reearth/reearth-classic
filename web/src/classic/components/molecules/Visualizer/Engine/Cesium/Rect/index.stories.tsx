import { Meta, StoryObj } from "@storybook/react-vite";

import { V, location } from "../storybook";

import Rect from ".";

const meta: Meta<typeof Rect> = {
  title: "classic/molecules/Visualizer/Engine/Cesium/Rect",
  component: Rect,
};
export default meta;
type Story = StoryObj<typeof Rect>;

export const Default: Story = {
  render: args => (
    <V location={location}>
      <Rect {...args} />
    </V>
  ),
  args: {
    layer: {
      id: "",
      isVisible: true,
      property: {
        default: {
          rect: { west: 139, east: 140, north: 36, south: 35 },
          fillColor: "#f00a",
          extrudedHeight: 10000,
          outlineColor: "yellow",
          outlineWidth: 10,
        },
      },
    },
    isBuilt: false,
    isEditable: false,
    isSelected: false,
  },
};

export const Image: Story = {
  render: args => (
    <V location={location}>
      <Rect {...args} />
    </V>
  ),
  args: {
    layer: {
      id: "",
      isVisible: true,
      property: {
        default: {
          rect: { west: 139, east: 140, north: 36, south: 35 },
          style: "image",
          image: `/sample.png`,
        },
      },
    },
    isBuilt: false,
    isEditable: false,
    isSelected: false,
  },
};
