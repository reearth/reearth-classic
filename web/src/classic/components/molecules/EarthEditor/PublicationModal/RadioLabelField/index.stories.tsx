import { Meta } from "@storybook/react-vite";
import { action } from "storybook/actions";

import RadioLabelField, { RadioLabelFieldProps } from ".";

const items: RadioLabelFieldProps["items"] = [
  { label: "デフォルトドメイン", value: "default" },
  { label: "カスタムドメイン", value: "custom" },
];

export default {
  title: "classic/molecules/EarthEditor/PublicationModal/RadioLabelField",
  component: RadioLabelField,
} as Meta;

export const Default = () => (
  <RadioLabelField selectedValue="default" items={items} onChange={action("onchange")} />
);
