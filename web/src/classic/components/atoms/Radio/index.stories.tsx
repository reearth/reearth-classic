import { Meta } from "@storybook/react-vite";

import Radio from ".";

export default {
  title: "classic/atoms/Radio",
  component: Radio,
} as Meta;

export const Default = () => <Radio value="value">Radio</Radio>;
