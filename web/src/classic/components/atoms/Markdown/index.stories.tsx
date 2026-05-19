import { Meta, StoryObj } from "@storybook/react-vite";

import Component from ".";

const markdown = `
> A block quote with ~strikethrough~ and a URL: https://reactjs.org.

* Lists
* [ ] todo
* [x] done

A table:

| a | b |
| - | - |
`;

const meta: Meta<typeof Component> = {
  title: "classic/atoms/Markdown",
  component: Component,
  parameters: { actions: { argTypesRegex: "^on.*" } },
};
export default meta;
type Story = StoryObj<typeof Component>;

export const Default: Story = {
  args: {
    children: markdown,
    backgroundColor: "#fff",
  },
};
