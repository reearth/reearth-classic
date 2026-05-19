import { action } from "storybook/actions";
import { Meta } from "@storybook/react-vite";

import ColorField from ".";

export default {
  title: "classic/molecules/EarthEditor/PropertyPane/PropertyField/ColorField",
  component: ColorField,
} as Meta;

export const Default = () => <ColorField onChange={action("onchange")} />;
