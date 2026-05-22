import { Meta, StoryObj } from "@storybook/react-vite";
import { ReactNode } from "react";

import PropertyList, { Layer } from ".";

const Wrapper: React.FC<{ children?: ReactNode }> = ({ children }) => (
  <div style={{ width: 300 }}>{children}</div>
);

const meta: Meta<typeof PropertyList> = {
  title: "classic/molecules/EarthEditor/PropertyPane/PropertyList",
  component: PropertyList,
  parameters: { actions: { argTypesRegex: "^on.*" } },
};
export default meta;
type Story = StoryObj<typeof PropertyList>;

const items = [
  { id: "a", title: "Tokyo", layerId: "a" },
  { id: "b", title: "Osaka", layerId: "c" },
  { id: "c", title: "Kyoto", layerId: "i" },
  { id: "d", title: "Tokyo", layerId: "e" },
  { id: "e", title: "Osaka", layerId: "g" },
  { id: "f", title: "Kyoto", layerId: "h" },
];

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
  { id: "h", title: "H" },
  { id: "i", title: "I" },
];

export const Default: Story = {
  render: args => (
    <Wrapper>
      <PropertyList {...args} items={items} />
    </Wrapper>
  ),
  args: {
    name: "Items",
  },
};

export const Selected: Story = {
  render: args => (
    <Wrapper>
      <PropertyList {...args} items={items} />
    </Wrapper>
  ),
  args: {
    name: "Items",
    selectedIndex: 1,
  },
};

export const LayerMode: Story = {
  render: args => (
    <Wrapper>
      <PropertyList {...args} layers={layers} items={items} />
    </Wrapper>
  ),
  args: {
    name: "Items",
    selectedIndex: 1,
    layerMode: true,
  },
};
