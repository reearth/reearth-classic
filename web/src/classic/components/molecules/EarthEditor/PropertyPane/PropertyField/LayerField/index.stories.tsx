import { Meta, StoryObj } from "@storybook/react-vite";

import Component, { Layer } from ".";

const meta: Meta<typeof Component> = {
  title: "classic/molecules/EarthEditor/PropertyPane/PropertyField/LayerField",
  component: Component,
  argTypes: { onChange: { action: "onChange" } },
  parameters: { actions: { argTypesRegex: "^on.*" } },
};
export default meta;
type Story = StoryObj<typeof Component>;

const layers: Layer[] = [
  { id: "a", title: "A" },
  {
    id: "b",
    title: "B",
    group: true,
    children: [
      { id: "d", title: "xxx" },
      { id: "e", title: "aaa" },
      { id: "f", title: "F", group: true, children: [{ id: "g", title: "G" }] },
    ],
  },
  { id: "c", title: "C" },
];

export const Default: Story = {
  render: args => <Component {...args} layers={layers} />,
  args: {
    value: "xxxx",
    linked: false,
    overridden: false,
    disabled: false,
  },
};
