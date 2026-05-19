import { Meta, StoryObj } from "@storybook/react-vite";

import Text from ".";

const meta: Meta<typeof Text> = {
  component: Text,
  parameters: { actions: { argTypesRegex: "^on.*" } },
};
export default meta;
type Story = StoryObj<typeof Text>;

const markdownText = `
# Hello
This is **markdown**.
## H2
### H3
#### H4
##### H5
`;

export const Default: Story = {
  args: {
    block: { id: "", property: { text: "aaaaaa" } },
    isSelected: false,
    isBuilt: false,
    isEditable: false,
  },
};

export const Title: Story = {
  args: {
    block: { id: "", property: { text: "aaaaaa", title: "Title" } },
    isSelected: false,
    isBuilt: false,
    isEditable: false,
  },
};

export const Markdown: Story = {
  args: {
    block: { id: "", property: { text: markdownText, markdown: true } },
    isSelected: false,
    isBuilt: false,
    isEditable: false,
  },
};

export const Typography: Story = {
  args: {
    block: {
      id: "",
      property: {
        text: markdownText,
        markdown: true,
        typography: { color: "red", fontSize: 16 },
      },
    },
    isSelected: false,
    isBuilt: false,
    isEditable: false,
  },
};

export const NoText: Story = {
  args: {
    isSelected: false,
    isBuilt: false,
    isEditable: true,
  },
};
