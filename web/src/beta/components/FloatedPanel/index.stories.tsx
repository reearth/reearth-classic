import { css } from "@emotion/react";
import { Meta, StoryObj } from "@storybook/react-vite";

import FloatedPanel from ".";

const meta: Meta<typeof FloatedPanel> = {
  component: FloatedPanel,
};
export default meta;
type Story = StoryObj<typeof FloatedPanel>;

export const Default: Story = {
  args: {
    visible: true,
    children: "This is the content of the floated panel",
  },
};

export const Hidden: Story = {
  args: {
    visible: false,
    children: "This is the content of the floated panel",
  },
};

export const CustomStyles: Story = {
  args: {
    visible: true,
    children: "This is the content of the floated panel",
    styles: css`
      background-color: lightblue;
      color: white;
      padding: 10px;
    `,
  },
};
