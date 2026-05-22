import { Meta, StoryObj } from "@storybook/react-vite";

import { V, location } from "../storybook";

import Ellipsoid from ".";

const meta: Meta<typeof Ellipsoid> = {
  title: "classic/molecules/Visualizer/Engine/Cesium/Ellipsoid",
  component: Ellipsoid,
};
export default meta;
type Story = StoryObj<typeof Ellipsoid>;

export const Default: Story = {
  render: args => (
    <V location={location}>
      <Ellipsoid {...args} />
    </V>
  ),
  args: {
    layer: {
      id: "",
      property: {
        default: {
          radius: 1000,
          fillColor: "#f00a",
          position: location,
          height: location.height,
        },
      },
      isVisible: true,
    },
  },
};
