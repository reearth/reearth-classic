import { Meta, StoryObj } from "@storybook/react-vite";

import { V, location } from "../storybook";

import Resource from ".";

const meta: Meta<typeof Resource> = {
  title: "classic/molecules/Visualizer/Engine/Cesium/Resource",
  component: Resource,
};
export default meta;
type Story = StoryObj<typeof Resource>;

export const Default: Story = {
  render: args => (
    <V location={location}>
      <Resource {...args} />
    </V>
  ),
  args: {
    layer: {
      id: "",
      isVisible: true,
      property: {
        default: {
          url: `/sample.geojson`,
        },
      },
    },
    isBuilt: false,
    isEditable: false,
    isSelected: false,
  },
};
