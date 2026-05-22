import { Meta, StoryObj } from "@storybook/react-vite";

import Component from ".";

const meta: Meta<typeof Component> = {
  title: "classic/atoms/TextBox",
  component: Component,
  argTypes: {
    color: { control: "color" },
    backgroundColor: { control: "color" },
    borderColor: { control: "color" },
    floatedTextColor: { control: "color" },
  },
  parameters: { actions: { argTypesRegex: "^on.*" } },
};
export default meta;
type Story = StoryObj<typeof Component>;

export const Basic: Story = {
  args: {
    color: "#fff",
    backgroundColor: "#000",
    borderColor: "#fff",
    floatedTextColor: "#ccc",
    disabled: false,
    placeholder: "",
    prefix: "",
    suffix: "",
    multiline: false,
    throttle: false,
    throttleTimeout: 1000,
    value: "",
  },
};
