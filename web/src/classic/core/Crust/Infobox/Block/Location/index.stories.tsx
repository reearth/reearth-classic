import { Meta, StoryObj } from "@storybook/react-vite";

import Location from ".";

const meta: Meta<typeof Location> = {
  component: Location,
  parameters: { actions: { argTypesRegex: "^on.*" } },
};
export default meta;
type Story = StoryObj<typeof Location>;

export const Default: Story = {
  args: {
    block: { id: "", property: { location: { lat: 0, lng: 0 } } },
    isSelected: false,
    isBuilt: false,
    isEditable: false,
  },
};

export const Title: Story = {
  args: {
    block: { id: "", property: { location: { lat: 0, lng: 0 }, title: "Location" } },
    isSelected: false,
    isBuilt: false,
    isEditable: false,
  },
};

export const NoLocation: Story = {
  args: {
    isSelected: false,
    isBuilt: false,
    isEditable: true,
  },
};
