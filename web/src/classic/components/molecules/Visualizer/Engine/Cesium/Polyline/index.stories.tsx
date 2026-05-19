import { Meta, StoryObj } from "@storybook/react-vite";

import { V, location } from "../storybook";

import Polyline from ".";

const meta: Meta<typeof Polyline> = {
  title: "classic/molecules/Visualizer/Engine/Cesium/Polyline",
  component: Polyline,
};
export default meta;
type Story = StoryObj<typeof Polyline>;

export const Default: Story = {
  render: args => (
    <V location={location}>
      <Polyline {...args} />
    </V>
  ),
  args: {
    layer: {
      id: "",
      property: {
        default: {
          strokeColor: "#ccaa",
          strokeWidth: 10,
          coordinates: [
            { lat: 35.652832, lng: 139.839478, height: 100 },
            { lat: 36.652832, lng: 140.039478, height: 100 },
            { lat: 34.652832, lng: 141.839478, height: 100 },
          ],
        },
      },
      isVisible: true,
    },
    isBuilt: false,
    isEditable: false,
    isSelected: false,
  },
};
