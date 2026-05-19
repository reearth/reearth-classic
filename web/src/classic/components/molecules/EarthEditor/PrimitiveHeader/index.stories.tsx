import { Meta, StoryObj } from "@storybook/react-vite";

import PrimitiveHeader from ".";

const primitives = [
  {
    id: "hoge",
    name: "Elipsoid",
    description: "This is an Elipsoid. This is an Elipsoid.",
    icon: "ellipsoid",
  },
  {
    id: "hoge",
    name: "marker",
    description: "This is a marker Elipsoid. This is a marker.",
    icon: "marker",
  },
  {
    id: "hoge",
    name: "resource",
    description: "This is an resource. This is an resource.",
    icon: "resource",
  },
  {
    id: "hoge",
    name: "polyline",
    description: "This is an polyline. This is an polyline.",
    icon: "polyline",
  },
];
const defaultProps = {
  primitives: primitives,
};

const meta: Meta<typeof PrimitiveHeader> = {
  title: "classic/molecules/EarthEditor/PrimitiveHeader",
  component: PrimitiveHeader,
};
export default meta;
type Story = StoryObj<typeof PrimitiveHeader>;

export const Default: Story = {
  args: defaultProps,
};
