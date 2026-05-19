import { Meta, StoryObj } from "@storybook/react-vite";

import Component from ".";

const meta: Meta<typeof Component> = {
  title: "classic/atoms/Scroll",
  component: Component,
  parameters: { actions: { argTypesRegex: "^on.*" } },
};
export default meta;
type Story = StoryObj<typeof Component>;

export const Default: Story = {
  render: args => (
    <div
      style={{
        width: "300px",
        height: "300px",
        border: "1px solid #fff",
        background: "#000",
        color: "#fff",
      }}>
      <Component {...args}>
        {new Array(100).fill("hogehoge").map((t, i) => (
          <div style={{ padding: "10px" }} key={i}>
            {t}
          </div>
        ))}
      </Component>
    </div>
  ),
  args: {},
};
