import { Meta, StoryObj } from "@storybook/react-vite";
import { action } from "storybook/actions";

import DateTimeField from ".";

const meta: Meta<typeof DateTimeField> = {
  component: DateTimeField,
};

export default meta;

type Story = StoryObj<typeof DateTimeField>;

export const DateTimeFieldInput: Story = {
  render: () => <DateTimeField name="Date Time Field" onChange={action("onchange")} />,
};
