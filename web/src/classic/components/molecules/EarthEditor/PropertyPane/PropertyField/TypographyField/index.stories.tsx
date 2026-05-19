import { Meta } from "@storybook/react-vite";
import { action } from "storybook/actions";

import TypographyField from ".";

export default {
  title: "classic/molecules/EarthEditor/PropertyPane/PropertyField/TypographyField",
  component: TypographyField,
} as Meta;

export const Default = () => <TypographyField onChange={action("onchange")} />;
