import { Meta, StoryObj } from "@storybook/react-vite";

import DataSourceManager from ".";

const meta: Meta<typeof DataSourceManager> = { component: DataSourceManager };
export default meta;
type Story = StoryObj<typeof DataSourceManager>;
export const Default: Story = { args: {} };
