import { Meta, StoryObj } from "@storybook/react-vite";

import SpacingInput from "./index";

const meta: Meta<typeof SpacingInput> = {
  component: SpacingInput,
};
export default meta;
type Story = StoryObj<typeof SpacingInput>;

export const Default: Story = {
  args: {
    name: "Padding",
    description: "Adjust the padding values",
  },
};

export const WithValues: Story = {
  args: {
    name: "Padding",
    description: "Adjust the padding values",
    value: {
      top: 10,
      left: 20,
      right: 30,
      bottom: 40,
    },
  },
};
