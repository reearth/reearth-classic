import { Meta, StoryObj } from "@storybook/react-vite";

import Component from ".";

const meta: Meta<typeof Component> = {
  title: "classic/molecules/EarthEditor/AssetsModal/AssetCard",
  component: Component,
};
export default meta;
type Story = StoryObj<typeof Component>;

export const DefaultMedium: Story = {
  args: {
    checked: false,
    cardSize: "medium",
    url: `/sample.svg`,
    name: "hoge",
  },
};

export const MediumChecked: Story = {
  args: {
    checked: true,
    cardSize: "medium",
    url: `/sample.svg`,
    name: "hoge",
  },
};

export const SmallChecked: Story = {
  args: {
    checked: true,
    cardSize: "small",
    url: `/sample.svg`,
    name: "hoge",
  },
};

export const LargeCheckedAndSelected: Story = {
  args: {
    checked: true,
    cardSize: "large",
    url: `/sample.svg`,
    name: "hoge",
    selected: true,
  },
};
