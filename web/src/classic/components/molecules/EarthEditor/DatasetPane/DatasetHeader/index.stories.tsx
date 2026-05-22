import { Meta } from "@storybook/react-vite";

import DatasetHeader from ".";

export default {
  title: "classic/molecules/EarthEditor/DatasetPane/DatasetHeader",
  component: DatasetHeader,
} as Meta;

export const Default = () => <DatasetHeader host="Hoge" />;
