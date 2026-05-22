import { Meta, StoryObj } from "@storybook/react-vite";

import Field from ".";

const meta: Meta<typeof Field> = {
  component: Field,
  parameters: { actions: { argTypesRegex: "^on.*" } },
};
export default meta;
type Story = StoryObj<typeof Field>;

export const Default: Story = {
  render: args => (
    <Field {...args}>
      <h1>HogeHoge</h1>
    </Field>
  ),
  args: {
    id: "aaa",
  },
};

export const Selected: Story = {
  render: args => (
    <Field {...args}>
      <h1>HogeHoge</h1>
    </Field>
  ),
  args: {
    id: "aaa",
    isEditable: true,
    isSelected: true,
  },
};
