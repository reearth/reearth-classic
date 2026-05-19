import { Meta, StoryObj } from "@storybook/react-vite";

import Component from ".";

const meta: Meta<typeof Component> = {
  title: "classic/atoms/HelpButton",
  component: Component,
};
export default meta;
type Story = StoryObj<typeof Component>;

const descriptionTitle = "Title";
const description = "Description";
const img = {
  imagePath: "/sample.svg",
  alt: "sample image",
};

export const Default: Story = {
  args: {
    description: description,
    descriptionTitle: descriptionTitle,
    img: img,
  },
};
