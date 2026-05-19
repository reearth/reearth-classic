import { Meta } from "@storybook/react-vite";

import Check from ".";

export default {
  title: "classic/atoms/Check",
  component: Check,
} as Meta;

export const Default = () => <Check value="value">Check</Check>;
