import { Meta, StoryObj } from "@storybook/react-vite";

import { V } from "../storybook";

import Component from ".";

const meta: Meta<typeof Component> = {
  title: "classic/molecules/Visualizer/Engine/Cesium/Tileset",
  component: Component,
};
export default meta;
type Story = StoryObj<typeof Component>;

export const Default: Story = {
  render: args => (
    <V lookAt={{ lng: -75.61209430779367, lat: 40.05083633101078, height: 0, range: 1200 }}>
      <Component {...args} />
    </V>
  ),
  args: {
    layer: {
      id: "",
      isVisible: true,
      property: {
        default: {
          tileset: `/tileset/tileset.json`,
        },
      },
    },
  },
};
