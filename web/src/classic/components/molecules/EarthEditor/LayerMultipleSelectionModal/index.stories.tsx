import { Meta, StoryObj } from "@storybook/react-vite";

import Component, { Layer } from ".";

const meta: Meta<typeof Component> = {
  title: "classic/molecules/EarthEditor/LayerMultipleSelectionModal",
  component: Component,
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
      { id: "d", title: "xxxxxxxxxxxxxxx" },
      { id: "e", title: "ああああああああああああああああ" },
      { id: "f", title: "F", group: true, children: [{ id: "g", title: "G" }] },
    ],
  },
  { id: "c", title: "C" },
];

export const Basic: Story = {
  render: args => <Component {...args} layers={layers} selected={["d", "e"]} />,
  args: {
    active: true,
  },
};
