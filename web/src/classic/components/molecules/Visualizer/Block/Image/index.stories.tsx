import { Meta, StoryObj } from "@storybook/react-vite";

import Image from ".";

const meta: Meta<typeof Image> = {
  title: "classic/molecules/Visualizer/Block/Image",
  component: Image,
  argTypes: {
    onClick: { action: "onClick" },
    onChange: { action: "onChange" },
  },
};
export default meta;
type Story = StoryObj<typeof Image>;

export const Default: Story = {
  render: args => (
    <div style={{ background: "#fff" }}>
      <Image {...args} />
    </div>
  ),
  args: {
    block: { id: "", property: { default: { image: `/sample.svg` } } },
    isSelected: false,
    isBuilt: false,
    isEditable: false,
  },
};

export const NoImage: Story = {
  render: args => (
    <div style={{ background: "#fff" }}>
      <Image {...args} />
    </div>
  ),
  args: {
    isSelected: false,
    isBuilt: false,
    isEditable: true,
  },
};

export const Title: Story = {
  render: args => (
    <div style={{ background: "#fff" }}>
      <Image {...args} />
    </div>
  ),
  args: {
    block: {
      id: "",
      property: { default: { image: `/sample.svg`, title: "Title" } },
    },
    isSelected: false,
    isBuilt: false,
    isEditable: false,
  },
};

export const FullSize: Story = {
  render: args => (
    <div style={{ background: "#fff" }}>
      <Image {...args} />
    </div>
  ),
  args: {
    block: {
      id: "",
      property: { default: { image: `/sample.svg`, fullSize: true } },
    },
    isSelected: false,
    isBuilt: false,
    isEditable: false,
  },
};

export const Cover: Story = {
  render: args => (
    <div style={{ background: "#fff" }}>
      <Image {...args} />
    </div>
  ),
  args: {
    block: {
      id: "",
      property: {
        default: {
          image: `/sample.svg`,
          imageSize: "cover",
          title: "Title",
        },
      },
    },
    isSelected: false,
    isBuilt: false,
    isEditable: false,
  },
};

export const Contain: Story = {
  render: args => (
    <div style={{ background: "#fff" }}>
      <Image {...args} />
    </div>
  ),
  args: {
    block: {
      id: "",
      property: {
        default: {
          image: `/sample.svg`,
          imageSize: "contain",
          title: "Title",
        },
      },
    },
    isSelected: false,
    isBuilt: false,
    isEditable: false,
  },
};

export const Position: Story = {
  render: args => (
    <div style={{ background: "#fff" }}>
      <Image {...args} />
    </div>
  ),
  args: {
    block: {
      id: "",
      property: {
        default: {
          image: `/sample.svg`,
          imageSize: "cover",
          title: "Title",
          imagePositionX: "left",
          imagePositionY: "top",
        },
      },
    },
    isSelected: false,
    isBuilt: false,
    isEditable: false,
  },
};
