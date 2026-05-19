import { Meta } from "@storybook/react-vite";
import { action } from "storybook/actions";

import CheckBox from ".";

export default {
  title: "classic/atoms/CheckBox",
  component: CheckBox,
} as Meta;

export const Default = () => <CheckBox checked onChange={action("onchange")} />;
