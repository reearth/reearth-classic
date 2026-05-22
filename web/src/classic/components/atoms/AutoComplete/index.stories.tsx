import { Meta, StoryObj } from "@storybook/react-vite";

import AutoComplete from ".";

const meta: Meta<typeof AutoComplete> = {
  title: "classic/atoms/AutoComplete",
  component: AutoComplete,
};
export default meta;
type Story = StoryObj<typeof AutoComplete>;

const sampleItems: { value: string; label: string }[] = [
  {
    value: "hoge",
    label: "hoge",
  },
  {
    value: "fuga",
    label: "fuga",
  },
];

const addItem = (value: string) => {
  sampleItems.push({ value, label: value });
};

const handleSelect = (value: string) => {
  console.log("select ", value);
};

const handleCreate = (value: string) => {
  console.log("create ", value);
  addItem(value);
};

export const Default: Story = {
  render: () => (
    <AutoComplete items={sampleItems} onCreate={handleCreate} onSelect={handleSelect} />
  ),
};

export const Creatable: Story = {
  render: () => (
    <AutoComplete items={sampleItems} onCreate={handleCreate} onSelect={handleSelect} creatable />
  ),
};
