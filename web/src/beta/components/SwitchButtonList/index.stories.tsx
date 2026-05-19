import { Meta, StoryObj } from "@storybook/react-vite";

import SwitchButtonList from ".";

export default {
  component: SwitchButtonList,
} as Meta;

type Story = StoryObj<typeof SwitchButtonList>;

export const Default: Story = {
  args: {
    list: [
      { id: "1", text: "PC", active: true },
      { id: "2", text: "Ipad", active: false },
      { id: "3", text: "Phone", active: false },
    ],
  },
};
