import { action } from "storybook/actions";
import { Meta } from "@storybook/react-vite";

import FontSizeField from ".";

export default {
  title: "classic/molecules/EarthEditor/PropertyPane/PropertyField/TypographyField/FontSizeField",
  component: FontSizeField,
} as Meta;

export const Default = () => <FontSizeField onChange={action("onchange")} />;
export const Selected = () => <FontSizeField value={10} onChange={action("onchange")} />;
