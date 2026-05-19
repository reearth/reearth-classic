import { Meta } from "@storybook/react-vite";

import { Option } from ".";

export default {
  title: "classic/atoms/SelectOption",
  component: Option,
} as Meta;

export const Default = () => (
  <Option value="value" label="label">
    Option
  </Option>
);
