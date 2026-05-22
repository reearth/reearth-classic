import { Meta, StoryObj } from "@storybook/react-vite";

import DataList, { Item } from ".";

const meta: Meta<typeof DataList> = {
  title: "classic/molecules/Visualizer/Block/DataList",
  component: DataList,
  argTypes: {
    onClick: { action: "onClick" },
    onChange: { action: "onChange" },
  },
};
export default meta;
type Story = StoryObj<typeof DataList>;

const items: Item[] = [
  { id: "a", item_title: "Name", item_datastr: "Foo bar", item_datatype: "string" },
  { id: "b", item_title: "Age", item_datanum: 20, item_datatype: "number" },
  { id: "c", item_title: "Address", item_datastr: "New York", item_datatype: "string" },
];

export const Default: Story = {
  args: {
    block: { id: "", property: { items } },
  },
};

export const Title: Story = {
  args: {
    block: { id: "", property: { default: { title: "Title" }, items } },
  },
};

export const Typography: Story = {
  args: {
    block: {
      id: "",
      property: {
        default: { typography: { color: "red", fontSize: 16 } },
        items,
      },
    },
  },
};

export const NoItems: Story = {
  args: { isEditable: true },
};

export const Selected: Story = {
  args: { block: { id: "", property: { items } }, isSelected: true },
};

export const Editable: Story = {
  args: {
    block: { id: "", property: { items } },
    isSelected: true,
    isEditable: true,
  },
};

export const Built: Story = {
  args: {
    block: { id: "", property: { items } },
    isBuilt: true,
  },
};
