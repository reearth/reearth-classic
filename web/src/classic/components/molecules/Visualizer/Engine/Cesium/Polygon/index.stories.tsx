import { Meta, StoryObj } from "@storybook/react-vite";

import { V, location } from "../storybook";

import Polygon from ".";

const meta: Meta<typeof Polygon> = {
  title: "classic/molecules/Visualizer/Engine/Cesium/Polygon",
  component: Polygon,
};
export default meta;
type Story = StoryObj<typeof Polygon>;

export const Default: Story = {
  render: args => (
    <V location={location}>
      <Polygon {...args} />
    </V>
  ),
  args: {
    layer: {
      id: "",
      property: {
        default: {
          fill: true,
          fillColor: "#ffffffaa",
          stroke: true,
          strokeColor: "red",
          strokeWidth: 10,
          polygon: [
            [
              { lat: 35.652832, lng: 139.839478, height: 100 },
              { lat: 36.652832, lng: 140.039478, height: 100 },
              { lat: 34.652832, lng: 141.839478, height: 100 },
            ],
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
