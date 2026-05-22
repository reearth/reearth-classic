import { Meta, StoryObj } from "@storybook/react-vite";

import HTML from ".";

const meta: Meta<typeof HTML> = {
  component: HTML,
  parameters: { actions: { argTypesRegex: "^on.*" } },
};
export default meta;
type Story = StoryObj<typeof HTML>;

export const Default: Story = {
  args: {
    block: { id: "", property: { html: "<p style='color: red;'>aaaaaa</p>" } },
    isSelected: false,
    isBuilt: false,
    isEditable: false,
  },
};

export const Title: Story = {
  args: {
    block: { id: "", property: { html: "<p style='color: red'>aaaaaa</p>", title: "Title" } },
    isSelected: false,
    isBuilt: false,
    isEditable: false,
  },
};

export const NoText: Story = {
  args: {
    isSelected: false,
    isBuilt: false,
    isEditable: true,
  },
};
