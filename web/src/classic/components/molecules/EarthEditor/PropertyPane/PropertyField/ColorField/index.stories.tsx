import { Meta } from "@storybook/react-vite";
import { action } from "storybook/actions";

import ColorField from ".";

export default {
  title: "classic/molecules/EarthEditor/PropertyPane/PropertyField/ColorField",
  component: ColorField,
} as Meta;

export const Default = () => <ColorField onChange={action("onchange")} />;
