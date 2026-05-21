import { Auth0Provider } from "@auth0/auth0-react";
import { Meta, StoryObj } from "@storybook/react-vite";

import DatasetModal from ".";

const meta: Meta<typeof DatasetModal> = {
  title: "classic/molecules/EarthEditor/DatasetPane/DatasetModal",
  component: DatasetModal,
  decorators: [
    Story => (
      <Auth0Provider
        domain="dummy-domain"
        clientId="dummy-client-id"
        authorizationParams={{
          redirect_uri: window.location.origin,
        }}>
        <Story />
      </Auth0Provider>
    ),
  ],
};

export default meta;

type Story = StoryObj<typeof DatasetModal>;

export const Default: Story = {
  render: () => <DatasetModal isVisible syncLoading={false} />,
};
