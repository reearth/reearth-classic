import { Meta, StoryObj } from "@storybook/react-vite";

import Component from ".";

const imageAsset = {
  url: `/sample.svg`,
  name: "hoge",
  id: "hoge",
  teamId: "hoge",
  size: 4300,
  contentType: "asset-image",
};

const fileAsset = {
  url: `somewhere/sample.kml`,
  name: "hoge.kml",
  id: "hoge",
  teamId: "hoge",
  size: 4300,
  contentType: "asset-image",
};

const meta: Meta<typeof Component> = {
  title: "classic/molecules/EarthEditor/AssetsModal/AssetListItem",
  component: Component,
};
export default meta;
type Story = StoryObj<typeof Component>;

export const Image: Story = {
  args: {
    checked: false,
    asset: imageAsset,
  },
};

export const File: Story = {
  args: {
    checked: false,
    asset: fileAsset,
  },
};

export const CheckedAndSelected: Story = {
  args: {
    checked: true,
    asset: imageAsset,
    selected: true,
  },
};
