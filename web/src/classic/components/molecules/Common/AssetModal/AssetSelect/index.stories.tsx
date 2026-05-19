import { Meta, StoryObj } from "@storybook/react-vite";

import Component from ".";

const filterOptions: { key: string; label: string }[] = [
  { key: "time", label: "Time" },
  { key: "size", label: "File size" },
  { key: "name", label: "Alphabetical" },
];

const meta: Meta<typeof Component> = {
  title: "classic/molecules/EarthEditor/AssetsModal/AssetSelect",
  component: Component,
};
export default meta;
type Story = StoryObj<typeof Component>;

export const Default: Story = {
  args: {
    items: filterOptions,
    value: "time",
  },
};
