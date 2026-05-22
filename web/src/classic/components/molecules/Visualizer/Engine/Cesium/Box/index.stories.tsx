import { Meta, StoryObj } from "@storybook/react-vite";

import { V, location } from "../storybook";

import Box from ".";

const meta: Meta<typeof Box> = {
  title: "classic/molecules/Visualizer/Engine/Cesium/Box",
  component: Box,
};
export default meta;
type Story = StoryObj<typeof Box>;

export const Default: Story = {
  render: args => (
    <V location={location}>
      <Box {...args} />
    </V>
  ),
  args: {
    layer: {
      id: "",
      property: {
        default: {
          fillColor: "rgba(0, 0, 0, 0.5)",
          location,
          width: 1000,
          length: 1000,
          height: 1000,
        },
      },
      isVisible: true,
    },
  },
};
